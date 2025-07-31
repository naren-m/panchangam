interface CacheEntry<T> {
  data: T;
  timestamp: number;
  expiresAt: number;
}

interface CacheOptions {
  ttl?: number; // Time to live in milliseconds
  maxSize?: number; // Maximum cache size
}

/**
 * Request cache to prevent duplicate API calls and implement request throttling
 */
export class RequestCache {
  private cache = new Map<string, CacheEntry<any>>();
  private pendingRequests = new Map<string, Promise<any>>();
  private options: Required<CacheOptions>;

  constructor(options: CacheOptions = {}) {
    this.options = {
      ttl: options.ttl ?? 5 * 60 * 1000, // Default 5 minutes
      maxSize: options.maxSize ?? 1000
    };
  }

  /**
   * Generate cache key from request parameters
   */
  private generateKey(endpoint: string, params: Record<string, any>): string {
    const sortedParams = Object.keys(params)
      .sort()
      .reduce((result, key) => {
        result[key] = params[key];
        return result;
      }, {} as Record<string, any>);

    return `${endpoint}:${JSON.stringify(sortedParams)}`;
  }

  /**
   * Get cached data if available and not expired
   */
  get<T>(endpoint: string, params: Record<string, any>): T | null {
    const key = this.generateKey(endpoint, params);
    const entry = this.cache.get(key);

    if (!entry) {
      return null;
    }

    if (Date.now() > entry.expiresAt) {
      this.cache.delete(key);
      return null;
    }

    return entry.data;
  }

  /**
   * Store data in cache
   */
  set<T>(endpoint: string, params: Record<string, any>, data: T, customTtl?: number): void {
    const key = this.generateKey(endpoint, params);
    const ttl = customTtl ?? this.options.ttl;
    const now = Date.now();

    // Implement LRU eviction if cache is full
    if (this.cache.size >= this.options.maxSize) {
      const oldestKey = this.cache.keys().next().value;
      this.cache.delete(oldestKey);
    }

    this.cache.set(key, {
      data,
      timestamp: now,
      expiresAt: now + ttl
    });
  }

  /**
   * Get pending request promise or create new one
   */
  getPendingRequest<T>(endpoint: string, params: Record<string, any>): Promise<T> | null {
    const key = this.generateKey(endpoint, params);
    return this.pendingRequests.get(key) || null;
  }

  /**
   * Store pending request promise
   */
  setPendingRequest<T>(endpoint: string, params: Record<string, any>, promise: Promise<T>): void {
    const key = this.generateKey(endpoint, params);
    
    // Clean up when promise resolves or rejects
    const cleanupPromise = promise.finally(() => {
      this.pendingRequests.delete(key);
    });

    this.pendingRequests.set(key, cleanupPromise);
  }

  /**
   * Clear all cached data
   */
  clear(): void {
    this.cache.clear();
    this.pendingRequests.clear();
  }

  /**
   * Remove expired entries
   */
  cleanup(): void {
    const now = Date.now();
    for (const [key, entry] of this.cache.entries()) {
      if (now > entry.expiresAt) {
        this.cache.delete(key);
      }
    }
  }

  /**
   * Get cache statistics
   */
  getStats() {
    return {
      size: this.cache.size,
      pendingRequests: this.pendingRequests.size,
      maxSize: this.options.maxSize,
      ttl: this.options.ttl
    };
  }
}

// Export singleton instance
export const requestCache = new RequestCache({
  ttl: 2 * 60 * 1000, // 2 minutes for panchangam data
  maxSize: 500
});