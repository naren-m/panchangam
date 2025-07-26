import { vi } from 'vitest';

// Mock environment variables
vi.mock('import.meta.env', () => ({
  VITE_API_BASE_URL: 'http://localhost:8080',
  VITE_API_TIMEOUT: '30000',
  VITE_DEBUG_API: 'true',
  DEV: true
}));

// Mock AbortSignal.timeout for Node.js environment
if (!globalThis.AbortSignal?.timeout) {
  globalThis.AbortSignal = {
    ...globalThis.AbortSignal,
    timeout: (ms: number) => {
      const controller = new AbortController();
      setTimeout(() => controller.abort(), ms);
      return controller.signal;
    }
  } as typeof globalThis.AbortSignal;
}