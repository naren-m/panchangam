// Re-export the new HTTP client-based API
export { panchangamApiClient as panchangamApi } from './api/panchangamApiClient';
export { apiConfig } from './api/client';

// Export types for backwards compatibility
export type { PanchangamData, GetPanchangamRequest } from '../types/panchangam';