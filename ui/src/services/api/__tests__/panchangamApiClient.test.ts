import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { panchangamApiClient } from '../panchangamApiClient';
import { PanchangamApiError } from '../types';
import { requestCache } from '../requestCache';
import * as matchers from '@testing-library/jest-dom/matchers';

expect.extend(matchers);

// Mock the API client
vi.mock('../client', () => ({
  apiClient: {
    get: vi.fn()
  }
}));

import { apiClient } from '../client';

const mockApiClient = apiClient as any;

describe('PanchangamApiClient', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Clear request cache to prevent test pollution
    requestCache.clear();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('getPanchangam', () => {
    const validRequest = {
      date: '2024-01-15',
      latitude: 13.0827,
      longitude: 80.2707,
      timezone: 'Asia/Kolkata',
      region: 'Tamil Nadu',
      calculation_method: 'Drik' as const,
      locale: 'en'
    };

    const mockApiResponse = {
      data: {
        date: '2024-01-15',
        tithi: 'Panchami',
        nakshatra: 'Rohini',
        yoga: 'Vishkumbha',
        karana: 'Bava',
        sunrise_time: '06:30:00',
        sunset_time: '18:15:00',
        events: [
          {
            name: 'Sunrise',
            time: '06:30:00',
            event_type: 'SUNRISE'
          },
          {
            name: 'Moonrise',
            time: '20:30:00',
            event_type: 'MOONRISE'
          }
        ]
      },
      status: 200,
      statusText: 'OK',
      headers: {},
      requestId: 'test-request-id'
    };

    it('should make successful API call with valid parameters', async () => {
      mockApiClient.get.mockResolvedValueOnce(mockApiResponse);

      const result = await panchangamApiClient.getPanchangam(validRequest);

      expect(mockApiClient.get).toHaveBeenCalledWith('/api/v1/panchangam', {
        date: '2024-01-15',
        lat: 13.0827,
        lng: 80.2707,
        tz: 'Asia/Kolkata',
        region: 'Tamil Nadu',
        method: 'Drik',
        locale: 'en'
      });

      expect(result).toEqual({
        date: '2024-01-15',
        tithi: 'Panchami',
        nakshatra: 'Rohini',
        yoga: 'Vishkumbha',
        karana: 'Bava',
        sunrise_time: '06:30:00',
        sunset_time: '18:15:00',
        vara: 'Monday',
        planetary_ruler: 'Moon',
        events: [
          {
            name: 'Sunrise',
            time: '06:30:00',
            event_type: 'SUNRISE',
            quality: 'auspicious'
          },
          {
            name: 'Moonrise',
            time: '20:30:00',
            event_type: 'MOONRISE',
            quality: 'neutral'
          }
        ],
        festivals: [],
        moonrise_time: '20:30:00',
        moonset_time: undefined
      });
    });

    it('should validate request parameters', async () => {
      const invalidRequests = [
        { ...validRequest, date: '' },
        { ...validRequest, date: 'invalid-date' },
        { ...validRequest, latitude: 100 },
        { ...validRequest, latitude: -100 },
        { ...validRequest, longitude: 200 },
        { ...validRequest, longitude: -200 },
        { ...validRequest, timezone: 'Invalid/Timezone' }
      ];

      for (const request of invalidRequests) {
        await expect(panchangamApiClient.getPanchangam(request)).rejects.toThrow(PanchangamApiError);
      }
    });

    it('should handle network errors gracefully', async () => {
      const networkError = new PanchangamApiError(
        'Network error. Please check your internet connection.',
        'NETWORK_ERROR',
        'test-request-id'
      );
      mockApiClient.get.mockRejectedValueOnce(networkError);

      const result = await panchangamApiClient.getPanchangam(validRequest);

      // Should return fallback data
      expect(result.tithi).toBe('API Unavailable');
      expect(result.nakshatra).toBe('Please check connection');
      expect(result.date).toBe('2024-01-15');
      expect(result.vara).toBe('Monday'); // Calculated from date
    });

    it('should handle API errors', async () => {
      const apiError = new PanchangamApiError(
        'Invalid request parameters',
        'INVALID_REQUEST',
        'test-request-id',
        400
      );
      mockApiClient.get.mockRejectedValue(apiError);

      await expect(panchangamApiClient.getPanchangam(validRequest)).rejects.toThrow(PanchangamApiError);
      await expect(panchangamApiClient.getPanchangam(validRequest)).rejects.toThrow('Invalid request parameters');
    });

    it('should validate response data', async () => {
      const invalidResponse = {
        data: {
          // Missing required fields
          date: '2024-01-15'
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        requestId: 'test-request-id'
      };
      mockApiClient.get.mockResolvedValueOnce(invalidResponse);

      await expect(panchangamApiClient.getPanchangam(validRequest)).rejects.toThrow(PanchangamApiError);
    });

    it('should handle timeout errors with fallback', async () => {
      const timeoutError = new PanchangamApiError(
        'Request timed out',
        'REQUEST_TIMEOUT',
        'test-request-id',
        408
      );
      mockApiClient.get.mockRejectedValueOnce(timeoutError);

      const result = await panchangamApiClient.getPanchangam(validRequest);

      // Should return fallback data
      expect(result.tithi).toBe('API Unavailable');
    });
  });

  describe('getPanchangamRange', () => {
    const baseRequest = {
      latitude: 13.0827,
      longitude: 80.2707,
      timezone: 'Asia/Kolkata',
      region: 'Tamil Nadu',
      calculation_method: 'Drik' as const,
      locale: 'en'
    };

    it('should fetch range data and handle parallel requests', async () => {
      const mockResponses = [
        {
          data: { date: '2024-01-01', tithi: 'Pratipada', nakshatra: 'Ashwini', yoga: 'Vishkumbha', karana: 'Bava', sunrise_time: '06:30:00', sunset_time: '18:15:00', events: [] },
          status: 200, statusText: 'OK', headers: {}, requestId: 'req1'
        },
        {
          data: { date: '2024-01-02', tithi: 'Dwitiya', nakshatra: 'Bharani', yoga: 'Priti', karana: 'Balava', sunrise_time: '06:30:00', sunset_time: '18:15:00', events: [] },
          status: 200, statusText: 'OK', headers: {}, requestId: 'req2'
        }
      ];

      mockApiClient.get.mockResolvedValueOnce(mockResponses[0]);
      mockApiClient.get.mockResolvedValueOnce(mockResponses[1]);

      const result = await panchangamApiClient.getPanchangamRange('2024-01-01', '2024-01-02', baseRequest);

      expect(result).toHaveLength(2);
      expect(result[0].date).toBe('2024-01-01');
      expect(result[1].date).toBe('2024-01-02');
      expect(mockApiClient.get).toHaveBeenCalledTimes(2);
    });

    it('should validate date range', async () => {
      // End date before start date
      await expect(
        panchangamApiClient.getPanchangamRange('2024-01-02', '2024-01-01', baseRequest)
      ).rejects.toThrow(PanchangamApiError);

      // Date range too large (> 365 days)
      await expect(
        panchangamApiClient.getPanchangamRange('2024-01-01', '2025-01-02', baseRequest)
      ).rejects.toThrow(PanchangamApiError);
    });

    it('should handle partial failures gracefully', async () => {
      const successResponse = {
        data: { date: '2024-01-01', tithi: 'Pratipada', nakshatra: 'Ashwini', yoga: 'Vishkumbha', karana: 'Bava', sunrise_time: '06:30:00', sunset_time: '18:15:00', events: [] },
        status: 200, statusText: 'OK', headers: {}, requestId: 'req1'
      };

      const apiError = new PanchangamApiError(
        'Server error',
        'SERVER_ERROR',
        'req2',
        500
      );

      mockApiClient.get.mockResolvedValueOnce(successResponse);
      mockApiClient.get.mockRejectedValueOnce(apiError);

      const result = await panchangamApiClient.getPanchangamRange('2024-01-01', '2024-01-02', baseRequest);

      // Should return only successful results (second request failed and should not use fallback for range queries)
      expect(result).toHaveLength(1);
      expect(result[0].date).toBe('2024-01-01');
    });
  });

  describe('healthCheck', () => {
    it('should return healthy status when API responds', async () => {
      const healthResponse = {
        data: {
          status: 'healthy',
          timestamp: '2024-01-15T10:00:00Z',
          version: '1.0.0'
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        requestId: 'health-check-id'
      };

      mockApiClient.get.mockResolvedValueOnce(healthResponse);

      const result = await panchangamApiClient.healthCheck();

      expect(result).toEqual({
        status: 'healthy',
        message: 'API is accessible',
        timestamp: '2024-01-15T10:00:00Z'
      });

      expect(mockApiClient.get).toHaveBeenCalledWith('/api/v1/health', undefined, {
        timeout: 5000
      });
    });

    it('should return unhealthy status when API fails', async () => {
      const apiError = new PanchangamApiError(
        'Service unavailable',
        'SERVICE_UNAVAILABLE',
        'health-check-id',
        503
      );
      mockApiClient.get.mockRejectedValueOnce(apiError);

      const result = await panchangamApiClient.healthCheck();

      expect(result.status).toBe('unhealthy');
      expect(result.message).toContain('Service unavailable');
    });

    it('should handle unknown errors during health check', async () => {
      mockApiClient.get.mockRejectedValueOnce(new Error('Unknown error'));

      const result = await panchangamApiClient.healthCheck();

      expect(result.status).toBe('unhealthy');
      expect(result.message).toBe('Unknown error during health check');
    });
  });

  describe('Edge cases', () => {
    it('should handle missing optional parameters', async () => {
      const minimalRequest = {
        date: '2024-01-15',
        latitude: 13.0827,
        longitude: 80.2707
      };

      const mockResponse = {
        data: {
          date: '2024-01-15',
          tithi: 'Panchami',
          nakshatra: 'Rohini',
          yoga: 'Vishkumbha',
          karana: 'Bava',
          sunrise_time: '06:30:00',
          sunset_time: '18:15:00',
          events: []
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        requestId: 'test-request-id'
      };

      mockApiClient.get.mockResolvedValueOnce(mockResponse);

      const result = await panchangamApiClient.getPanchangam(minimalRequest);

      expect(mockApiClient.get).toHaveBeenCalledWith('/api/v1/panchangam', {
        date: '2024-01-15',
        lat: 13.0827,
        lng: 80.2707,
        tz: 'UTC',
        region: '',
        method: 'traditional',
        locale: 'en'
      });

      expect(result).toBeDefined();
      expect(result.date).toBe('2024-01-15');
    });

    it('should handle boundary latitude and longitude values', async () => {
      const boundaryRequests = [
        { date: '2024-01-15', latitude: 90, longitude: 180 },
        { date: '2024-01-15', latitude: -90, longitude: -180 },
        { date: '2024-01-15', latitude: 0, longitude: 0 }
      ];

      const mockResponse = {
        data: {
          date: '2024-01-15',
          tithi: 'Panchami',
          nakshatra: 'Rohini',
          yoga: 'Vishkumbha',
          karana: 'Bava',
          sunrise_time: '06:30:00',
          sunset_time: '18:15:00',
          events: []
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        requestId: 'test-request-id'
      };

      for (const request of boundaryRequests) {
        mockApiClient.get.mockResolvedValueOnce(mockResponse);
        await expect(panchangamApiClient.getPanchangam(request)).resolves.toBeDefined();
      }
    });
  });
});