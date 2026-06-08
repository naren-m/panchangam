import { vi } from 'vitest';
import '@testing-library/jest-dom/vitest';

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

// Mock IntersectionObserver for components that depend on browser layout APIs.
globalThis.IntersectionObserver = class MockIntersectionObserver implements IntersectionObserver {
  readonly root = null;
  readonly rootMargin = '';
  readonly thresholds = [];

  observe() {
    return null;
  }

  disconnect() {
    return null;
  }

  takeRecords() {
    return [];
  }

  unobserve() {
    return null;
  }
};

// Mock ResizeObserver for components that measure container size.
globalThis.ResizeObserver = class ResizeObserver {
  constructor() {}

  observe() {
    return null;
  }

  disconnect() {
    return null;
  }

  unobserve() {
    return null;
  }
};

// Mock window.matchMedia for responsive design tests.
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});
