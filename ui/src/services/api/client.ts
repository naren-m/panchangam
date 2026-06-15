import {
  ApiClientConfig,
  ApiRequest,
  ApiResponse,
  PanchangamApiError,
  RequestInterceptor,
  ResponseInterceptor,
  ErrorInterceptor
} from './types';
import { getApiBaseUrl, getRuntimeConfig } from './clientConfig';
import { transformApiError, type HttpError } from './clientErrors';
import { delay, generateRequestId } from './clientUtils';

/**
 * Robust HTTP client with retry logic, interceptors, and error handling
 */
export class ApiClient {
  private config: ApiClientConfig;
  private hasCustomBaseURL: boolean;
  private requestInterceptors: RequestInterceptor[] = [];
  private responseInterceptors: ResponseInterceptor[] = [];
  private errorInterceptors: ErrorInterceptor[] = [];

  constructor(config: Partial<ApiClientConfig> = {}) {
    this.hasCustomBaseURL = Boolean(config.baseURL);
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
    let lastError: unknown;
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
    const generatedRequestId = generateRequestId();
    let requestId = generatedRequestId;

    try {
      // Apply request interceptors
      let processedRequest: ApiRequest = {
        ...request,
        headers: {
          'X-Request-ID': generatedRequestId,
          'X-Timestamp': new Date().toISOString(),
          ...request.headers
        }
      };
      for (const interceptor of this.requestInterceptors) {
        processedRequest = await interceptor(processedRequest);
      }
      requestId = processedRequest.headers?.['X-Request-ID'] || generatedRequestId;

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

      // Report HTTP errors before parsing the body so status codes are never lost.
      if (!response.ok) {
        const error = new Error(`HTTP ${response.status}: ${response.statusText}`) as HttpError;
        error.status = response.status;
        error.statusText = response.statusText;
        throw error;
      }

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
        requestId
      };

      // Add performance headers
      apiResponse.headers['x-response-time'] = `${endTime - startTime}ms`;
      apiResponse.headers['x-request-url'] = url.toString();

      // Apply response interceptors
      let processedResponse: ApiResponse<unknown> = apiResponse;
      for (const interceptor of this.responseInterceptors) {
        processedResponse = await interceptor(processedResponse);
      }

      return processedResponse as ApiResponse<T>;

    } catch (error) {
      const transformedError = transformApiError(error, requestId);

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
  async get<T = unknown>(url: string, params?: Record<string, unknown>, options: Partial<ApiRequest> = {}): Promise<ApiResponse<T>> {
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
  async post<T = unknown>(url: string, data?: unknown, options: Partial<ApiRequest> = {}): Promise<ApiResponse<T>> {
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
  async put<T = unknown>(url: string, data?: unknown, options: Partial<ApiRequest> = {}): Promise<ApiResponse<T>> {
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
  async delete<T = unknown>(url: string, options: Partial<ApiRequest> = {}): Promise<ApiResponse<T>> {
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
    if (this.hasCustomBaseURL) {
      return this.config.baseURL;
    }

    // Check for runtime configuration first
    const runtimeConfig = getRuntimeConfig();
    if (runtimeConfig?.API_ENDPOINT) {
      return runtimeConfig.API_ENDPOINT;
    }
    if (import.meta.env.VITE_API_BASE_URL) {
      return import.meta.env.VITE_API_BASE_URL;
    }
    // Use current origin for relative URLs (nginx proxies /api/ to gateway)
    if (typeof window !== 'undefined') {
      return window.location.origin;
    }
    return this.config.baseURL;
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
    if (newConfig.baseURL) {
      this.hasCustomBaseURL = true;
    }
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
