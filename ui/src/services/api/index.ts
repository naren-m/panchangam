// Export main API client
export { apiClient, apiConfig } from './client';
export { panchangamApiClient } from './panchangamApiClient';

// Export types
export type {
  ApiClientConfig,
  ApiRequest,
  ApiResponse,
  ApiError,
  RequestInterceptor,
  ResponseInterceptor,
  ErrorInterceptor
} from './types';

export { PanchangamApiError } from './types';

// Re-export for backwards compatibility with existing code
export { panchangamApiClient as panchangamApi } from './panchangamApiClient';