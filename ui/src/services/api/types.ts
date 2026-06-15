// HTTP Client Types
export interface ApiClientConfig {
  baseURL: string;
  timeout: number;
  retries: number;
  retryDelay: number;
  headers: Record<string, string>;
}

export interface ApiRequest<T = unknown> {
  url: string;
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  data?: T;
  params?: Record<string, unknown>;
  headers?: Record<string, string>;
  timeout?: number;
}

export interface ApiResponse<T = unknown> {
  data: T;
  status: number;
  statusText: string;
  headers: Record<string, string>;
  requestId: string;
}

export interface ApiError {
  code: string;
  message: string;
  details?: unknown;
  requestId: string;
  timestamp: string;
  path: string;
  status: number;
}

export class PanchangamApiError extends Error {
  public readonly code: string;
  public readonly requestId?: string;
  public readonly status?: number;
  public retryCount?: number;

  constructor(message: string, code: string = 'UNKNOWN_ERROR', requestId?: string, status?: number) {
    super(message);
    this.name = 'PanchangamApiError';
    this.code = code;
    this.requestId = requestId;
    this.status = status;
  }
}

// Request/Response interceptor types
export type RequestInterceptor = (request: ApiRequest) => ApiRequest | Promise<ApiRequest>;
export type ResponseInterceptor = (response: ApiResponse<unknown>) => ApiResponse<unknown> | Promise<ApiResponse<unknown>>;
export type ErrorInterceptor = (error: PanchangamApiError) => PanchangamApiError | Promise<PanchangamApiError>;
