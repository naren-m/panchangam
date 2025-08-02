import {
  ApiClientConfig,
  ApiRequest,
  ApiResponse,
  ApiError,
  PanchangamApiError,
  RequestInterceptor,
  ResponseInterceptor,
  ErrorInterceptor
} from './types';

/**
 * Generates a unique request ID for correlation
 */
function generateRequestId(): string {
  return `req_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * Delays execution for a specified number of milliseconds
 */
function delay(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}

/**
 * Get API base URL from runtime configuration or environment variables
 */
function getApiBaseUrl(): string {
  // Check for runtime configuration first (for Docker deployments)
  if (typeof window !== 'undefined' && (window as any).__RUNTIME_CONFIG__) {
    const runtimeConfig = (window as any).__RUNTIME_CONFIG__;
    if (runtimeConfig.API_ENDPOINT) {
      return runtimeConfig.API_ENDPOINT;
    }
  }
  
  // Fallback to build-time environment variables
  return import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
}

/**
 * Transforms API errors to standardized format
 */
function transformApiError(error: any, requestId: string, path: string): PanchangamApiError {
  // Handle abort/timeout errors first
  if (error.name === 'AbortError' || error.name === 'TimeoutError') {
    return new PanchangamApiError(
      'Request timed out. Please check your connection and try again.',
      'REQUEST_TIMEOUT',
      requestId,
      408
    );
  }

  // Handle fetch/network errors
  if (error.name === 'TypeError' && error.message.includes('Failed to fetch')) {
    return new PanchangamApiError(
      'Network error. Please check your internet connection.',
      'NETWORK_ERROR',
      requestId,
      0
    );
  }

  // Handle HTTP status errors - check for status on error object directly or nested response
  const status = error.status || (error.response && error.response.status);
  
  if (status) {
    const errorMessages: Record<number, { code: string; message: string }> = {
      400: { code: 'INVALID_REQUEST', message: 'Invalid request parameters.' },
      401: { code: 'UNAUTHORIZED', message: 'Authentication required.' },
      403: { code: 'FORBIDDEN', message: 'Access denied.' },
      404: { code: 'NOT_FOUND', message: 'Resource not found.' },
      429: { code: 'RATE_LIMITED', message: 'Too many requests. Please try again later.' },
      500: { code: 'SERVER_ERROR', message: 'Internal server error. Please try again.' },
      502: { code: 'BAD_GATEWAY', message: 'Service temporarily unavailable.' },
      503: { code: 'SERVICE_UNAVAILABLE', message: 'Service temporarily unavailable.' },
      504: { code: 'GATEWAY_TIMEOUT', message: 'Request timed out at gateway.' }
    };

    const errorInfo = errorMessages[status] || {
      code: 'HTTP_ERROR',
      message: `HTTP ${status}: ${error.statusText || 'Unknown error'}`
    };

    return new PanchangamApiError(errorInfo.message, errorInfo.code, requestId, status);
  }

  // Check for other error patterns in the message
  if (error.message) {
    if (error.message.includes('Failed to fetch') || error.message.includes('NetworkError')) {
      return new PanchangamApiError(
        'Network error. Please check your internet connection.',
        'NETWORK_ERROR',
        requestId,
        0
      );
    }
    
    if (error.message.includes('timeout') || error.message.includes('Timeout')) {
      return new PanchangamApiError(
        'Request timed out. Please check your connection and try again.',
        'REQUEST_TIMEOUT',
        requestId,
        408
      );
    }
  }

  return new PanchangamApiError(
    error.message || 'An unexpected error occurred.',
    'UNKNOWN_ERROR',
    requestId
  );
}

/**
 * Robust HTTP client with retry logic, interceptors, and error handling
 */
export class ApiClient {
  private config: ApiClientConfig;
  private requestInterceptors: RequestInterceptor[] = [];
  private responseInterceptors: ResponseInterceptor[] = [];
  private errorInterceptors: ErrorInterceptor[] = [];

  constructor(config: Partial<ApiClientConfig> = {}) {
    this.config = {
      baseURL: getApiBaseUrl(),
      timeout: parseInt(import.meta.env.VITE_API_TIMEOUT) || 30000,
      retries: parseInt(import.meta.env.VITE_API_RETRIES) || 3,
      retryDelay: parseInt(import.meta.env.VITE_API_RETRY_DELAY) || 1000,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'X-Client-Version': import.meta.env.VITE_APP_VERSION || '1.0.0'
      },
      ...config
    };

    // Add default request interceptor for correlation ID
    this.addRequestInterceptor((request) => {
      const requestId = generateRequestId();
      return {
        ...request,
        headers: {
          ...request.headers,
          'X-Request-ID': requestId,
          'X-Timestamp': new Date().toISOString()
        }
      };
    });

    // Add default response interceptor for logging
    this.addResponseInterceptor((response) => {
      if (import.meta.env.VITE_LOG_LEVEL === 'debug') {
        console.log(`API Response [${response.requestId}]:`, {
          status: response.status,
          url: response.headers['x-request-url'] || 'unknown',
          duration: response.headers['x-response-time'] || 'unknown'
        });
      }
      return response;
    });
  }

  /**
   * Add a request interceptor
   */
  addRequestInterceptor(interceptor: RequestInterceptor): void {
    this.requestInterceptors.push(interceptor);
  }

  /**
   * Add a response interceptor
   */
  addResponseInterceptor(interceptor: ResponseInterceptor): void {
    this.responseInterceptors.push(interceptor);
  }

  /**
   * Add an error interceptor
   */
  addErrorInterceptor(interceptor: ErrorInterceptor): void {
    this.errorInterceptors.push(interceptor);
  }

  /**
   * Execute request with retry logic
   */
  private async executeRequest<T>(request: ApiRequest): Promise<ApiResponse<T>> {
    let lastError: any;
    let actualAttempts = 0;
    
    for (let attempt = 0; attempt <= this.config.retries; attempt++) {
      actualAttempts++;
      try {
        return await this.makeRequest<T>(request);
      } catch (error) {
        lastError = error;
        
        // Don't retry on client errors (4xx) or specific error types
        const shouldNotRetry = error instanceof PanchangamApiError && (
          (error.status && error.status >= 400 && error.status < 500) ||
          error.code === 'INVALID_REQUEST' ||
          error.code === 'UNAUTHORIZED' ||
          error.code === 'FORBIDDEN' ||
          error.code === 'NOT_FOUND'
        );
        
        if (shouldNotRetry) {
          if (import.meta.env.VITE_LOG_LEVEL === 'debug') {
            console.log(`Not retrying client error [${error.code}]:`, error.message);
          }
          break;
        }

        // If this is the last attempt, don't delay
        if (attempt === this.config.retries) {
          break;
        }

        // Exponential backoff with jitter
        const backoffDelay = this.config.retryDelay * Math.pow(2, attempt);
        const jitter = Math.random() * 0.1 * backoffDelay;
        await delay(backoffDelay + jitter);

        if (import.meta.env.VITE_LOG_LEVEL === 'debug') {
          console.warn(`API request failed, retrying (${attempt + 1}/${this.config.retries}):`, error);
        }
      }
    }

    // Enhance error with retry information
    if (lastError instanceof PanchangamApiError) {
      lastError.retryCount = actualAttempts - 1;
    }

    throw lastError;
  }

  /**
   * Make a single HTTP request
   */
  private async makeRequest<T>(request: ApiRequest): Promise<ApiResponse<T>> {
    const requestId = generateRequestId();
    
    try {
      // Apply request interceptors
      let processedRequest = request;
      for (const interceptor of this.requestInterceptors) {
        processedRequest = await interceptor(processedRequest);
      }

      // Build URL with runtime configuration
      const url = new URL(processedRequest.url, this.getRuntimeBaseURL());
      
      // Add query parameters
      if (processedRequest.params) {
        Object.entries(processedRequest.params).forEach(([key, value]) => {
          if (value !== undefined && value !== null) {
            url.searchParams.append(key, String(value));
          }
        });
      }

      // Prepare fetch options
      const fetchOptions: RequestInit = {
        method: processedRequest.method,
        headers: {
          ...this.config.headers,
          ...processedRequest.headers
        },
        signal: AbortSignal.timeout(processedRequest.timeout || this.config.timeout)
      };

      // Add body for non-GET requests
      if (processedRequest.data && processedRequest.method !== 'GET') {
        fetchOptions.body = JSON.stringify(processedRequest.data);
      }

      // Make the request
      const startTime = Date.now();
      const response = await fetch(url.toString(), fetchOptions);
      const endTime = Date.now();

      // Parse response
      let data: T;
      const contentType = response.headers.get('content-type');
      
      if (contentType && contentType.includes('application/json')) {
        data = await response.json();
      } else {
        data = (await response.text()) as unknown as T;
      }

      // Build response object
      const apiResponse: ApiResponse<T> = {
        data,
        status: response.status,
        statusText: response.statusText,
        headers: Object.fromEntries(response.headers.entries()),
        requestId: processedRequest.headers?.['X-Request-ID'] || requestId
      };

      // Add performance headers
      apiResponse.headers['x-response-time'] = `${endTime - startTime}ms`;
      apiResponse.headers['x-request-url'] = url.toString();

      // Check for HTTP errors
      if (!response.ok) {
        const error = new Error(`HTTP ${response.status}: ${response.statusText}`);
        (error as any).status = response.status;
        (error as any).statusText = response.statusText;
        (error as any).response = apiResponse;
        throw error;
      }

      // Apply response interceptors
      let processedResponse = apiResponse;
      for (const interceptor of this.responseInterceptors) {
        processedResponse = await interceptor(processedResponse);
      }

      return processedResponse;

    } catch (error) {
      const transformedError = transformApiError(error, requestId, request.url);
      
      // Apply error interceptors
      let processedError = transformedError;
      for (const interceptor of this.errorInterceptors) {
        processedError = await interceptor(processedError);
      }

      throw processedError;
    }
  }

  /**
   * Make a GET request
   */
  async get<T = any>(url: string, params?: Record<string, any>, options: Partial<ApiRequest> = {}): Promise<ApiResponse<T>> {
    return this.executeRequest<T>({
      method: 'GET',
      url,
      params,
      ...options
    });
  }

  /**
   * Make a POST request
   */
  async post<T = any>(url: string, data?: any, options: Partial<ApiRequest> = {}): Promise<ApiResponse<T>> {
    return this.executeRequest<T>({
      method: 'POST',
      url,
      data,
      ...options
    });
  }

  /**
   * Make a PUT request
   */
  async put<T = any>(url: string, data?: any, options: Partial<ApiRequest> = {}): Promise<ApiResponse<T>> {
    return this.executeRequest<T>({
      method: 'PUT',
      url,
      data,
      ...options
    });
  }

  /**
   * Make a DELETE request
   */
  async delete<T = any>(url: string, options: Partial<ApiRequest> = {}): Promise<ApiResponse<T>> {
    return this.executeRequest<T>({
      method: 'DELETE',
      url,
      ...options
    });
  }

  /**
   * Get current base URL with runtime configuration
   */
  private getRuntimeBaseURL(): string {
    // Check for runtime configuration first
    const runtimeConfig = (window as any).__RUNTIME_CONFIG__;
    return runtimeConfig?.API_ENDPOINT || 
           import.meta.env.VITE_API_BASE_URL || 
           this.config.baseURL;
  }

  /**
   * Get current configuration
   */
  getConfig(): ApiClientConfig {
    return { 
      ...this.config,
      baseURL: this.getRuntimeBaseURL()
    };
  }

  /**
   * Update configuration
   */
  updateConfig(newConfig: Partial<ApiClientConfig>): void {
    this.config = { ...this.config, ...newConfig };
  }
}

// Create and export default client instance
export const apiClient = new ApiClient();

// Export configuration for debugging - make it dynamic
export const apiConfig = {
  get baseUrl() {
    return apiClient.getConfig().baseURL;
  },
  get endpoint() {
    return apiClient.getConfig().baseURL;
  },
  get timeout() {
    return apiClient.getConfig().timeout;
  },
  get retries() {
    return apiClient.getConfig().retries;
  }
};