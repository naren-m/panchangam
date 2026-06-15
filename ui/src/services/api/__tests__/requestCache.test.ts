import { afterEach, describe, expect, it, vi } from 'vitest';

import { RequestCache } from '../requestCache';

describe('RequestCache', () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  it('uses the same key for the same parameters in any order', () => {
    const cache = new RequestCache({ ttl: 1_000, maxSize: 5 });

    cache.set('/api/v1/panchangam', { lng: -121.9, date: '2024-01-15', lat: 37.4 }, 'cached');

    expect(
      cache.get<string>('/api/v1/panchangam', { date: '2024-01-15', lat: 37.4, lng: -121.9 })
    ).toBe('cached');
  });

  it('removes expired cache entries', () => {
    vi.useFakeTimers();
    const cache = new RequestCache({ ttl: 100, maxSize: 5 });

    cache.set('/endpoint', { date: '2024-01-15' }, 'expired');

    vi.advanceTimersByTime(101);

    expect(cache.get<string>('/endpoint', { date: '2024-01-15' })).toBeNull();
    expect(cache.getStats().size).toBe(0);
  });

  it('removes pending requests after they resolve', async () => {
    const cache = new RequestCache({ ttl: 1_000, maxSize: 5 });
    let resolve!: (value: string) => void;
    const promise = new Promise<string>((promiseResolve) => {
      resolve = promiseResolve;
    });

    cache.setPendingRequest('/endpoint', { b: 2, a: 1 }, promise);

    expect(cache.getPendingRequest('/endpoint', { a: 1, b: 2 })).toBe(promise);

    resolve('done');
    await promise;
    await Promise.resolve();

    expect(cache.getPendingRequest('/endpoint', { a: 1, b: 2 })).toBeNull();
  });

  it('removes the oldest entry when the cache is full', () => {
    const cache = new RequestCache({ ttl: 1_000, maxSize: 2 });

    cache.set('/endpoint', { id: 1 }, 'first');
    cache.set('/endpoint', { id: 2 }, 'second');
    cache.set('/endpoint', { id: 3 }, 'third');

    expect(cache.get<string>('/endpoint', { id: 1 })).toBeNull();
    expect(cache.get<string>('/endpoint', { id: 2 })).toBe('second');
    expect(cache.get<string>('/endpoint', { id: 3 })).toBe('third');
  });
});
