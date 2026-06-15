import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { fetchSkyView } from './skyViewApi';

describe('fetchSkyView', () => {
  const mockFetch = vi.fn();

  beforeEach(() => {
    mockFetch.mockReset();
    vi.stubGlobal('fetch', mockFetch);
    vi.spyOn(console, 'error').mockImplementation(() => undefined);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  it('throws the backend error message without logging to console', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 503,
      statusText: 'Service Unavailable',
      json: () => Promise.resolve({
        error: { message: 'Sky view service unavailable' },
      }),
    });

    await expect(fetchSkyView({
      latitude: 37.4323,
      longitude: -121.9066,
    })).rejects.toThrow('Sky view service unavailable');

    expect(console.error).not.toHaveBeenCalled();
  });
});
