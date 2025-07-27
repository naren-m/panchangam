/**
 * 🎯 Unified API Service Facade
 * 
 * Central point for all backend integrations following frontend best practices:
 * - Single entry point for all API operations
 * - Consistent error handling and transformation
 * - Configuration management
 * - Request/response interceptors
 * - Caching and optimization strategies
 */

import { panchangamApiClient } from './api/panchangamApiClient';
import { apiClient, apiConfig } from './api/client';
import type { 
  PanchangamData, 
  GetPanchangamRequest,
  Settings 
} from '../types/panchangam';

/**
 * 🏗️ Unified API Service Class
 * 
 * Provides a single, consistent interface for all backend operations.
 * This follows the Facade pattern to simplify complex subsystem interactions.
 */
export class ApiService {
  private static instance: ApiService;
  
  // 🔒 Singleton pattern for consistent configuration
  public static getInstance(): ApiService {
    if (!ApiService.instance) {
      ApiService.instance = new ApiService();
    }
    return ApiService.instance;
  }

  private constructor() {
    this.setupInterceptors();
  }

  /**
   * 🔧 Setup global request/response interceptors
   */
  private setupInterceptors(): void {
    // Request interceptor for consistent headers
    apiClient.addRequestInterceptor((request) => ({
      ...request,
      headers: {
        ...request.headers,
        'X-App-Version': import.meta.env.VITE_APP_VERSION || '1.0.0',
        'X-Client-Type': 'web',
        'X-Requested-With': 'PanchangamApp'
      }
    }));

    // Response interceptor for consistent data transformation
    apiClient.addResponseInterceptor((response) => {
      // Add performance metrics
      if (response.headers['x-response-time']) {
        console.debug(`API Response Time: ${response.headers['x-response-time']}`);
      }
      return response;
    });

    // Error interceptor for consistent error handling
    apiClient.addErrorInterceptor(async (error) => {
      // Log errors for monitoring
      console.error('API Error:', {
        code: error.code,
        message: error.message,
        status: error.status,
        requestId: error.requestId
      });

      // Could add error reporting service here
      return error;
    });
  }

  // 📊 Panchangam Operations
  public readonly panchangam = {
    /**
     * Get panchangam data for a specific date
     */
    get: (params: GetPanchangamRequest): Promise<PanchangamData> => {
      return panchangamApiClient.getPanchangam(params);
    },

    /**
     * Get panchangam data for a date range
     */
    getRange: (
      startDate: string, 
      endDate: string, 
      params: Omit<GetPanchangamRequest, 'date'>
    ): Promise<PanchangamData[]> => {
      return panchangamApiClient.getPanchangamRange(startDate, endDate, params);
    },

    /**
     * Check API health status
     */
    healthCheck: () => {
      return panchangamApiClient.healthCheck();
    }
  };

  // 🛠️ System Operations
  public readonly system = {
    /**
     * Get current API configuration
     */
    getConfig: () => {
      return apiClient.getConfig();
    },

    /**
     * Update API configuration
     */
    updateConfig: (config: Partial<typeof apiConfig>) => {
      return apiClient.updateConfig(config);
    },

    /**
     * Get API health and performance metrics
     */
    getMetrics: async () => {
      const config = apiClient.getConfig();
      const health = await panchangamApiClient.healthCheck();
      
      return {
        endpoint: config.baseURL,
        health: health.status,
        message: health.message,
        timestamp: health.timestamp,
        timeout: config.timeout,
        retries: config.retries
      };
    }
  };

  // 🏷️ Utility Methods
  public readonly utils = {
    /**
     * Transform settings to API request parameters
     */
    settingsToApiParams: (
      date: Date, 
      settings: Settings
    ): GetPanchangamRequest => ({
      date: date.toISOString().split('T')[0],
      latitude: settings.location.latitude,
      longitude: settings.location.longitude,
      timezone: settings.location.timezone,
      region: settings.region,
      calculation_method: settings.calculation_method,
      locale: settings.locale
    }),

    /**
     * Check if API is available
     */
    isApiAvailable: async (): Promise<boolean> => {
      try {
        const health = await panchangamApiClient.healthCheck();
        return health.status === 'healthy';
      } catch {
        return false;
      }
    }
  };
}

// 🎯 Export singleton instance
export const apiService = ApiService.getInstance();

// 🔄 Re-export for backwards compatibility
export { panchangamApiClient, apiClient, apiConfig } from './panchangamApi';
export type { PanchangamData, GetPanchangamRequest } from '../types/panchangam';