/**
 * 🎯 Central Services Export
 * 
 * Single entry point for all service layer exports.
 * This provides a clean, organized way to access all backend integrations.
 */

// 🚀 Primary API Service (Recommended)
export { apiService } from './apiService';

// 🔗 Enhanced Hooks
export { 
  usePanchangamData, 
  usePanchangamRange, 
  useApiHealth 
} from '../hooks/useApiService';

// 🛠️ Low-level clients (for advanced use cases)
export { 
  panchangamApiClient, 
  apiClient, 
  apiConfig 
} from './panchangamApi';

// 📊 Types
export type { 
  PanchangamData, 
  GetPanchangamRequest 
} from '../types/panchangam';

// 🔧 Utilities
export { formatDateForApi } from '../utils/dateHelpers';