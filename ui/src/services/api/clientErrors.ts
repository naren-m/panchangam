import { PanchangamApiError } from './types';

interface ErrorDetails {
  name?: string;
  message?: string;
  status?: number;
  statusText?: string;
  response?: {
    status?: number;
  };
}

export type HttpError = Error & {
  status?: number;
  statusText?: string;
};

const HTTP_ERROR_MESSAGES: Record<number, { code: string; message: string }> = {
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

function getErrorDetails(error: unknown): ErrorDetails {
  if (error instanceof Error) {
    const withStatus = error as HttpError & { response?: { status?: number } };
    return {
      name: error.name,
      message: error.message,
      status: withStatus.status,
      statusText: withStatus.statusText,
      response: withStatus.response,
    };
  }

  if (typeof error === 'object' && error !== null) {
    return error as ErrorDetails;
  }

  return {
    message: String(error),
  };
}

export function transformApiError(error: unknown, requestId: string): PanchangamApiError {
  const details = getErrorDetails(error);

  if (details.name === 'AbortError' || details.name === 'TimeoutError') {
    return new PanchangamApiError(
      'Request timed out. Please check your connection and try again.',
      'REQUEST_TIMEOUT',
      requestId,
      408
    );
  }

  if (details.name === 'TypeError' && details.message?.includes('Failed to fetch')) {
    return new PanchangamApiError(
      'Network error. Please check your internet connection.',
      'NETWORK_ERROR',
      requestId,
      0
    );
  }

  const status = details.status || details.response?.status;

  if (status) {
    const errorInfo = HTTP_ERROR_MESSAGES[status] || {
      code: 'HTTP_ERROR',
      message: `HTTP ${status}: ${details.statusText || 'Unknown error'}`
    };

    return new PanchangamApiError(errorInfo.message, errorInfo.code, requestId, status);
  }

  if (details.message) {
    if (details.message.includes('Failed to fetch') || details.message.includes('NetworkError')) {
      return new PanchangamApiError(
        'Network error. Please check your internet connection.',
        'NETWORK_ERROR',
        requestId,
        0
      );
    }

    if (details.message.includes('timeout') || details.message.includes('Timeout')) {
      return new PanchangamApiError(
        'Request timed out. Please check your connection and try again.',
        'REQUEST_TIMEOUT',
        requestId,
        408
      );
    }
  }

  return new PanchangamApiError(
    details.message || 'An unexpected error occurred.',
    'UNKNOWN_ERROR',
    requestId
  );
}
