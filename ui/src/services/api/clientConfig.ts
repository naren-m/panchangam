interface RuntimeConfigWindow extends Window {
  __RUNTIME_CONFIG__?: {
    API_ENDPOINT?: string;
  };
}

export function getRuntimeConfig() {
  if (typeof window === 'undefined') {
    return undefined;
  }

  return (window as RuntimeConfigWindow).__RUNTIME_CONFIG__;
}

export function getApiBaseUrl(): string {
  const runtimeConfig = getRuntimeConfig();
  if (runtimeConfig?.API_ENDPOINT) {
    return runtimeConfig.API_ENDPOINT;
  }

  if (import.meta.env.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL;
  }

  if (typeof window !== 'undefined') {
    return window.location.origin;
  }

  return 'http://localhost:8080';
}
