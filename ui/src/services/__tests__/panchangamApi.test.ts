import { describe, it, expect, vi, beforeEach } from 'vitest';
import { panchangamApi } from '../panchangamApi';

// Mock fetch globally
const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('PanchangamApiService', () => {
  beforeEach(() => {
    mockFetch.mockClear();
  });

  describe('getPanchangam', () => {
    it('should fetch panchangam data successfully', async () => {
      // Mock successful API response
      const mockApiResponse = {
        date: '2024-01-15',
        tithi: 'Chaturthi (4)',
        nakshatra: 'Uttara Bhadrapada (26)',
        yoga: 'Siddha (21)',
        karana: 'Gara (6)',
        sunrise_time: '01:15:32',
        sunset_time: '12:41:47',
        events: [
          {
            name: 'Tithi: Chaturthi',
            time: '01:15:32',
            event_type: 'TITHI'
          }
        ]
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockApiResponse
      });

      const request = {
        date: '2024-01-15',
        latitude: 12.9716,
        longitude: 77.5946,
        timezone: 'Asia/Kolkata',
        region: 'Karnataka',
        calculation_method: 'traditional',
        locale: 'en'
      };

      const result = await panchangamApi.getPanchangam(request);

      // Verify the request was made correctly
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/panchangam'),
        expect.objectContaining({
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
          }
        })
      );

      // Verify the response transformation
      expect(result).toEqual({
        date: '2024-01-15',
        tithi: 'Chaturthi (4)',
        nakshatra: 'Uttara Bhadrapada (26)',
        yoga: 'Siddha (21)',
        karana: 'Gara (6)',
        sunrise_time: '01:15:32',
        sunset_time: '12:41:47',
        vara: 'Monday', // Should be calculated based on date
        planetary_ruler: 'Moon', // Should be calculated based on day of week
        events: [
          {
            name: 'Tithi: Chaturthi',
            time: '01:15:32',
            event_type: 'TITHI',
            quality: 'neutral'
          }
        ],
        festivals: [],
        moonrise_time: undefined,
        moonset_time: undefined
      });
    });

    it('should handle API errors gracefully', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      const request = {
        date: '2024-01-15',
        latitude: 12.9716,
        longitude: 77.5946,
        timezone: 'Asia/Kolkata',
        region: 'Karnataka',
        calculation_method: 'traditional',
        locale: 'en'
      };

      await expect(panchangamApi.getPanchangam(request)).rejects.toThrow('Network error');
    });

    it('should provide fallback data when API is unavailable', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Failed to fetch'));

      const request = {
        date: '2024-01-15',
        latitude: 12.9716,
        longitude: 77.5946,
        timezone: 'Asia/Kolkata',
        region: 'Karnataka',
        calculation_method: 'traditional',
        locale: 'en'
      };

      const result = await panchangamApi.getPanchangam(request);

      // Should return fallback data
      expect(result.date).toBe('2024-01-15');
      expect(result.tithi).toBe('API Unavailable');
      expect(result.nakshatra).toBe('Please check connection');
      expect(result.events).toHaveLength(1);
      expect(result.events[0].name).toBe('API Connection Error');
    });

    it('should handle HTTP error responses', async () => {
      const errorResponse = {
        error: {
          message: 'Invalid latitude parameter',
          code: 'INVALID_PARAMETER'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        statusText: 'Bad Request',
        json: async () => errorResponse
      });

      const request = {
        date: '2024-01-15',
        latitude: 999, // Invalid latitude
        longitude: 77.5946,
        timezone: 'Asia/Kolkata',
        region: 'Karnataka',
        calculation_method: 'traditional',
        locale: 'en'
      };

      await expect(panchangamApi.getPanchangam(request)).rejects.toThrow('Invalid latitude parameter');
    });
  });

  describe('healthCheck', () => {
    it('should return healthy status when API is accessible', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ status: 'healthy' })
      });

      const result = await panchangamApi.healthCheck();

      expect(result).toEqual({
        status: 'healthy',
        message: 'API is accessible'
      });
    });

    it('should return unhealthy status when API is not accessible', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Connection refused'));

      const result = await panchangamApi.healthCheck();

      expect(result).toEqual({
        status: 'unhealthy',
        message: 'Connection refused'
      });
    });
  });

  describe('getPanchangamRange', () => {
    it('should fetch data for multiple dates', async () => {
      const mockApiResponse1 = {
        date: '2024-01-15',
        tithi: 'Chaturthi (4)',
        nakshatra: 'Uttara Bhadrapada (26)',
        yoga: 'Siddha (21)',
        karana: 'Gara (6)',
        sunrise_time: '01:15:32',
        sunset_time: '12:41:47',
        events: []
      };

      const mockApiResponse2 = {
        date: '2024-01-16',
        tithi: 'Panchami (5)',
        nakshatra: 'Revati (27)',
        yoga: 'Sadhya (22)',
        karana: 'Balava (7)',
        sunrise_time: '01:16:00',
        sunset_time: '12:42:00',
        events: []
      };

      mockFetch
        .mockResolvedValueOnce({
          ok: true,
          json: async () => mockApiResponse1
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => mockApiResponse2
        });

      const request = {
        latitude: 12.9716,
        longitude: 77.5946,
        timezone: 'Asia/Kolkata',
        region: 'Karnataka',
        calculation_method: 'traditional',
        locale: 'en'
      };

      const result = await panchangamApi.getPanchangamRange('2024-01-15', '2024-01-16', request);

      expect(result).toHaveLength(2);
      expect(result[0].date).toBe('2024-01-15');
      expect(result[1].date).toBe('2024-01-16');
      expect(mockFetch).toHaveBeenCalledTimes(2);
    });
  });
});