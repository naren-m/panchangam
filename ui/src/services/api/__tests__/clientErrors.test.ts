import { describe, expect, it } from 'vitest';

import { transformApiError } from '../clientErrors';

describe('transformApiError', () => {
  it('maps known HTTP statuses to clear API errors', () => {
    const error = new Error('HTTP 503: Service Unavailable') as Error & {
      status?: number;
      statusText?: string;
    };
    error.status = 503;
    error.statusText = 'Service Unavailable';

    const transformed = transformApiError(error, 'request-id');

    expect(transformed.code).toBe('SERVICE_UNAVAILABLE');
    expect(transformed.message).toBe('Service temporarily unavailable.');
    expect(transformed.status).toBe(503);
    expect(transformed.requestId).toBe('request-id');
  });

  it('keeps unknown HTTP statuses clear', () => {
    const error = new Error('HTTP 418: Teapot') as Error & {
      status?: number;
      statusText?: string;
    };
    error.status = 418;
    error.statusText = 'Teapot';

    const transformed = transformApiError(error, 'request-id');

    expect(transformed.code).toBe('HTTP_ERROR');
    expect(transformed.message).toBe('HTTP 418: Teapot');
    expect(transformed.status).toBe(418);
  });

  it('maps fetch and timeout failures to retryable categories', () => {
    const networkError = transformApiError(new TypeError('Failed to fetch'), 'network-request');
    expect(networkError.code).toBe('NETWORK_ERROR');
    expect(networkError.status).toBe(0);

    const timeout = new Error('Timeout');
    timeout.name = 'AbortError';
    const timeoutError = transformApiError(timeout, 'timeout-request');
    expect(timeoutError.code).toBe('REQUEST_TIMEOUT');
    expect(timeoutError.status).toBe(408);
  });
});
