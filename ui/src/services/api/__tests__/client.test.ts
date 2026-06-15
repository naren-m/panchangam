import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { ApiClient } from '../client';
import { PanchangamApiError } from '../types';

// Mock fetch globally
const mockFetch = vi.fn();
globalThis.fetch = mockFetch;

async function expectApiError(
  request: Promise<unknown>,
  code: string,
  status?: number
): Promise<void> {
  let caughtError: unknown;

  try {
    await request;
  } catch (error) {
    caughtError = error;
  }

  expect(caughtError).toBeInstanceOf(PanchangamApiError);
  expect((caughtError as PanchangamApiError).code).toBe(code);

  if (status !== undefined) {
    expect((caughtError as PanchangamApiError).status).toBe(status);
  }
}

// Mock AbortSignal.timeout
Object.defineProperty(AbortSignal, 'timeout', {
  value: vi.fn((timeout: number) => {
    const controller = new AbortController();
    setTimeout(() => controller.abort(), timeout);
    return controller.signal;
  }),
  writable: true
});

describe('ApiClient', () => {
  let client: ApiClient;

  beforeEach(() => {
    client = new ApiClient({
      baseURL: 'https://api.test.com',
      timeout: 5000,
      retries: 2,
      retryDelay: 100
    });
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.unstubAllEnvs();
    vi.restoreAllMocks();
  });

  describe('Configuration', () => {
    it('should use default configuration', () => {
      const defaultClient = new ApiClient();
      const config = defaultClient.getConfig();

      expect(config.baseURL).toBeDefined();
      expect(config.timeout).toBe(30000);
      expect(config.retries).toBe(3);
      expect(config.headers['Content-Type']).toBe('application/json');
    });

    it('should merge custom configuration', () => {
      const config = client.getConfig();

      expect(config.baseURL).toBe('https://api.test.com');
      expect(config.timeout).toBe(5000);
      expect(config.retries).toBe(2);
    });

    it('should update configuration', () => {
      client.updateConfig({ timeout: 10000 });

      expect(client.getConfig().timeout).toBe(10000);
    });
  });

  describe('GET requests', () => {
    it('should make successful GET request', async () => {
      const mockResponse = { data: 'test data' };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        statusText: 'OK',
        headers: new Headers([['content-type', 'application/json']]),
        json: () => Promise.resolve(mockResponse)
      });

      const response = await client.get('/test');

      expect(mockFetch).toHaveBeenCalledTimes(1);
      expect(response.data).toEqual(mockResponse);
      expect(response.status).toBe(200);
    });

    it('should add query parameters', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        statusText: 'OK',
        headers: new Headers([['content-type', 'application/json']]),
        json: () => Promise.resolve({})
      });

      await client.get('/test', { param1: 'value1', param2: 'value2' });

      const fetchCall = mockFetch.mock.calls[0];
      const url = new URL(fetchCall[0]);

      expect(url.searchParams.get('param1')).toBe('value1');
      expect(url.searchParams.get('param2')).toBe('value2');
    });

    it('should handle non-JSON responses', async () => {
      const textResponse = 'plain text response';
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        statusText: 'OK',
        headers: new Headers([['content-type', 'text/plain']]),
        text: () => Promise.resolve(textResponse)
      });

      const response = await client.get('/test');

      expect(response.data).toBe(textResponse);
    });

    it('should not log successful responses when debug logging is enabled', async () => {
      vi.stubEnv('VITE_LOG_LEVEL', 'debug');
      const debugSpy = vi.spyOn(console, 'debug').mockImplementation(() => undefined);
      const debugClient = new ApiClient({
        baseURL: 'https://api.test.com',
        timeout: 5000,
        retries: 0,
        retryDelay: 1
      });

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        statusText: 'OK',
        headers: new Headers([['content-type', 'application/json']]),
        json: () => Promise.resolve({ ok: true })
      });

      const response = await debugClient.get('/test');

      expect(response.data).toEqual({ ok: true });
      expect(debugSpy).not.toHaveBeenCalled();
    });
  });

  describe('POST requests', () => {
    it('should make successful POST request with data', async () => {
      const requestData = { name: 'test' };
      const responseData = { id: 1, name: 'test' };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        statusText: 'Created',
        headers: new Headers([['content-type', 'application/json']]),
        json: () => Promise.resolve(responseData)
      });

      const response = await client.post('/test', requestData);

      expect(mockFetch).toHaveBeenCalledTimes(1);
      expect(response.data).toEqual(responseData);
      expect(response.status).toBe(201);

      const fetchCall = mockFetch.mock.calls[0];
      const options = fetchCall[1];
      expect(options.method).toBe('POST');
      expect(options.body).toBe(JSON.stringify(requestData));
    });
  });

  describe('Error handling', () => {
    let noRetryClient: ApiClient;

    beforeEach(() => {
      noRetryClient = new ApiClient({
        baseURL: 'https://api.test.com',
        timeout: 5000,
        retries: 0,
        retryDelay: 100
      });
    });

    it('should handle HTTP 404 errors', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        statusText: 'Not Found',
        headers: new Headers(),
        json: () => Promise.resolve({ error: { message: 'Resource not found' } })
      });

      await expectApiError(noRetryClient.get('/test'), 'NOT_FOUND', 404);
    });

    it('should handle HTTP 500 errors', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        headers: new Headers(),
        json: () => Promise.resolve({})
      });

      await expectApiError(noRetryClient.get('/test'), 'SERVER_ERROR', 500);
    });

    it('should handle network errors', async () => {
      const networkError = new TypeError('Failed to fetch');
      mockFetch.mockRejectedValueOnce(networkError);

      await expectApiError(noRetryClient.get('/test'), 'NETWORK_ERROR');
    });

    it('should handle timeout errors', async () => {
      const timeoutError = new Error('Timeout');
      timeoutError.name = 'AbortError';
      mockFetch.mockRejectedValueOnce(timeoutError);

      await expectApiError(noRetryClient.get('/test'), 'REQUEST_TIMEOUT');
    });
  });

  describe('Retry logic', () => {
    it('should retry on server errors', async () => {
      // First two calls fail with 500, third succeeds
      mockFetch
        .mockResolvedValueOnce({
          ok: false,
          status: 500,
          statusText: 'Internal Server Error',
          headers: new Headers(),
          json: () => Promise.resolve({})
        })
        .mockResolvedValueOnce({
          ok: false,
          status: 500,
          statusText: 'Internal Server Error',
          headers: new Headers(),
          json: () => Promise.resolve({})
        })
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          statusText: 'OK',
          headers: new Headers([['content-type', 'application/json']]),
          json: () => Promise.resolve({ success: true })
        });

      const response = await client.get('/test');

      expect(mockFetch).toHaveBeenCalledTimes(3);
      expect(response.data).toEqual({ success: true });
    });

    it('should not retry on 4xx client errors', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        statusText: 'Bad Request',
        headers: new Headers(),
        json: () => Promise.resolve({})
      });

      await expectApiError(client.get('/test'), 'INVALID_REQUEST', 400);

      // Should not retry on client errors
      expect(mockFetch).toHaveBeenCalledTimes(1);
    });

    it('should respect maximum retry attempts', async () => {
      mockFetch.mockResolvedValue({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        headers: new Headers(),
        json: () => Promise.resolve({})
      });

      await expect(client.get('/test')).rejects.toThrow(PanchangamApiError);

      // Should try initial request + 2 retries = 3 total calls
      expect(mockFetch).toHaveBeenCalledTimes(3);
    });

    it('should not log retry attempts when debug logging is enabled', async () => {
      vi.stubEnv('VITE_LOG_LEVEL', 'debug');
      const warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => undefined);
      const retryClient = new ApiClient({
        baseURL: 'https://api.test.com',
        timeout: 5000,
        retries: 1,
        retryDelay: 1
      });

      mockFetch
        .mockResolvedValueOnce({
          ok: false,
          status: 500,
          statusText: 'Internal Server Error',
          headers: new Headers(),
          json: () => Promise.resolve({})
        })
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          statusText: 'OK',
          headers: new Headers([['content-type', 'application/json']]),
          json: () => Promise.resolve({ success: true })
        });

      const response = await retryClient.get('/test');

      expect(response.data).toEqual({ success: true });
      expect(warnSpy).not.toHaveBeenCalled();
    });
  });

  describe('Interceptors', () => {
    it('should apply request interceptors', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        statusText: 'OK',
        headers: new Headers([['content-type', 'application/json']]),
        json: () => Promise.resolve({})
      });

      // Add custom request interceptor
      client.addRequestInterceptor((request) => ({
        ...request,
        headers: {
          ...request.headers,
          'X-Custom-Header': 'test-value'
        }
      }));

      await client.get('/test');

      const fetchCall = mockFetch.mock.calls[0];
      const options = fetchCall[1];
      expect(options.headers['X-Custom-Header']).toBe('test-value');
    });

    it('should apply response interceptors', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        statusText: 'OK',
        headers: new Headers([['content-type', 'application/json']]),
        json: () => Promise.resolve({ original: 'data' })
      });

      // Add custom response interceptor
      client.addResponseInterceptor((response) => ({
        ...response,
        data: {
          ...(response.data as Record<string, unknown>),
          intercepted: true
        }
      }));

      const response = await client.get<Record<string, unknown>>('/test');

      expect(response.data).toEqual({ original: 'data', intercepted: true });
    });
  });

  describe('Request correlation', () => {
    it('should add request ID headers', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        statusText: 'OK',
        headers: new Headers([['content-type', 'application/json']]),
        json: () => Promise.resolve({})
      });

      await client.get('/test');

      const fetchCall = mockFetch.mock.calls[0];
      const options = fetchCall[1];

      expect(options.headers['X-Request-ID']).toMatch(/^req_\d+_/);
      expect(options.headers['X-Timestamp']).toBeDefined();
    });

    it('should use the sent request ID on errors', async () => {
      const noRetryClient = new ApiClient({
        baseURL: 'https://api.test.com',
        timeout: 5000,
        retries: 0,
        retryDelay: 1
      });
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        headers: new Headers(),
        json: () => Promise.resolve({})
      });

      let caughtError: unknown;
      try {
        await noRetryClient.get('/test');
      } catch (error) {
        caughtError = error;
      }

      const fetchCall = mockFetch.mock.calls[0];
      const options = fetchCall[1];

      expect(caughtError).toBeInstanceOf(PanchangamApiError);
      expect((caughtError as PanchangamApiError).requestId).toBe(options.headers['X-Request-ID']);
    });
  });
});
