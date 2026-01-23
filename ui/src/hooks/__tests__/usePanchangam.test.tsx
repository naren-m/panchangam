import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { usePanchangam, usePanchangamRange } from '../usePanchangam';

// Mock the API service
vi.mock('../../services/panchangamApi', () => ({
  panchangamApi: {
    getPanchangam: vi.fn(),
    getPanchangamRange: vi.fn(),
  },
}));

// Import the mocked panchangamApi
import { panchangamApi } from '../../services/panchangamApi';
const mockPanchangamApi = panchangamApi as any;

describe('usePanchangam', () => {
  const mockDate = new Date('2024-01-15');
  const mockSettings = {
    calculation_method: 'Drik' as const,
    locale: 'en',
    region: 'Tamil Nadu',
    time_format: '12' as const,
    location: {
      name: 'Chennai, Tamil Nadu',
      latitude: 13.0827,
      longitude: 80.2707,
      timezone: 'Asia/Kolkata',
      region: 'Tamil Nadu'
    }
  };

  const mockPanchangamData = {
    date: '2024-01-15',
    tithi: 'Panchami',
    nakshatra: 'Rohini',
    yoga: 'Vishkumbha',
    karana: 'Bava',
    sunrise_time: '06:30:00',
    sunset_time: '18:15:00',
    vara: 'Monday',
    planetary_ruler: 'Moon',
    events: [],
    festivals: [],
    moonrise_time: '20:30:00',
    moonset_time: '08:45:00',
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('successful data fetch', () => {
    it('should fetch and return panchangam data', async () => {
      mockPanchangamApi.getPanchangam.mockResolvedValue(mockPanchangamData);

      const { result } = renderHook(() => usePanchangam(mockDate, mockSettings));

      // Initially loading
      expect(result.current.loading).toBe(true);
      expect(result.current.data).toBe(null);
      expect(result.current.error).toBe(null);

      // Wait for data to load
      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.data).toEqual(mockPanchangamData);
      expect(result.current.error).toBe(null);
      expect(mockPanchangamApi.getPanchangam).toHaveBeenCalledWith({
        date: '2024-01-15',
        latitude: 13.0827,
        longitude: 80.2707,
        timezone: 'Asia/Kolkata',
        region: 'Tamil Nadu',
        calculation_method: 'Drik',
        locale: 'en'
      });
    });
  });

  describe('error handling', () => {
    it('should handle API errors', async () => {
      const errorMessage = 'API request failed';
      mockPanchangamApi.getPanchangam.mockRejectedValue(new Error(errorMessage));

      const { result } = renderHook(() => usePanchangam(mockDate, mockSettings));

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.data).toBe(null);
      expect(result.current.error).toBe(errorMessage);
      expect(result.current.errorState.hasError).toBe(true);
      expect(result.current.errorState.message).toBe(errorMessage);
    });

    it('should detect network errors', async () => {
      const networkError = new Error('Failed to fetch');
      mockPanchangamApi.getPanchangam.mockRejectedValue(networkError);

      const { result } = renderHook(() => usePanchangam(mockDate, mockSettings));

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.errorState.isNetworkError).toBe(true);
    });

    it('should parse status codes from API errors', async () => {
      const apiError = new Error('API request failed: 404 Not Found');
      mockPanchangamApi.getPanchangam.mockRejectedValue(apiError);

      const { result } = renderHook(() => usePanchangam(mockDate, mockSettings));

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.errorState.statusCode).toBe(404);
    });
  });

  describe('retry functionality', () => {
    it('should retry when retry function is called', async () => {
      mockPanchangamApi.getPanchangam
        .mockRejectedValueOnce(new Error('First call failed'))
        .mockResolvedValueOnce(mockPanchangamData);

      const { result } = renderHook(() => usePanchangam(mockDate, mockSettings));

      // Wait for first call to fail
      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.error).toBeTruthy();
      expect(result.current.retryCount).toBe(0);

      // Retry and wait for completion
      await act(async () => {
        result.current.retry();
      });

      // Wait for retry to complete
      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.data).toEqual(mockPanchangamData);
      expect(result.current.error).toBe(null);
      expect(result.current.retryCount).toBe(1);
      expect(mockPanchangamApi.getPanchangam).toHaveBeenCalledTimes(2);
    });
  });

  describe('abort controller', () => {
    it('should cancel previous requests when new ones are made', async () => {
      let resolveFirst: (value: any) => void;
      let resolveSecond: (value: any) => void;

      const firstPromise = new Promise(resolve => { resolveFirst = resolve; });
      const secondPromise = new Promise(resolve => { resolveSecond = resolve; });

      mockPanchangamApi.getPanchangam
        .mockReturnValueOnce(firstPromise)
        .mockReturnValueOnce(secondPromise);

      const { result, rerender } = renderHook(
        ({ date }) => usePanchangam(date, mockSettings),
        { initialProps: { date: mockDate } }
      );

      // Start second request before first completes
      const newDate = new Date('2024-01-16');
      rerender({ date: newDate });

      // Resolve second request
      resolveSecond!(mockPanchangamData);

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      // Only second request should complete successfully
      expect(mockPanchangamApi.getPanchangam).toHaveBeenCalledTimes(2);
      expect(result.current.data).toEqual(mockPanchangamData);
    });
  });
});

describe('usePanchangamRange', () => {
  const startDate = new Date('2024-01-01');
  const endDate = new Date('2024-01-03');
  const mockSettings = {
    calculation_method: 'Drik' as const,
    locale: 'en',
    region: 'Tamil Nadu',
    time_format: '12' as const,
    location: {
      name: 'Chennai, Tamil Nadu',
      latitude: 13.0827,
      longitude: 80.2707,
      timezone: 'Asia/Kolkata',
      region: 'Tamil Nadu'
    }
  };

  const mockPanchangamData = {
    date: '2024-01-15',
    tithi: 'Panchami',
    nakshatra: 'Rohini',
    yoga: 'Vishkumbha',
    karana: 'Bava',
    sunrise_time: '06:30:00',
    sunset_time: '18:15:00',
    vara: 'Monday',
    planetary_ruler: 'Moon',
    events: [],
    festivals: [],
    moonrise_time: '20:30:00',
    moonset_time: '08:45:00',
  };

  const mockRangeData = [
    { ...mockPanchangamData, date: '2024-01-01' },
    { ...mockPanchangamData, date: '2024-01-02' },
    { ...mockPanchangamData, date: '2024-01-03' },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch range data and organize by date', async () => {
    mockPanchangamApi.getPanchangamRange.mockResolvedValue(mockRangeData);

    const { result } = renderHook(() => usePanchangamRange(startDate, endDate, mockSettings));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.data).toEqual({
      '2024-01-01': mockRangeData[0],
      '2024-01-02': mockRangeData[1],
      '2024-01-03': mockRangeData[2],
    });

    expect(mockPanchangamApi.getPanchangamRange).toHaveBeenCalledWith(
      '2024-01-01',
      '2024-01-03',
      {
        latitude: 13.0827,
        longitude: 80.2707,
        timezone: 'Asia/Kolkata',
        region: 'Tamil Nadu',
        calculation_method: 'Drik',
        locale: 'en'
      }
    );
  });

  it('should handle range data fetch errors', async () => {
    const errorMessage = 'Range fetch failed';
    mockPanchangamApi.getPanchangamRange.mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => usePanchangamRange(startDate, endDate, mockSettings));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.data).toEqual({});
    expect(result.current.error).toBe(errorMessage);
  });
});