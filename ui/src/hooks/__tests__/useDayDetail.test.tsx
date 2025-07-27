import { describe, it, expect, vi } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useDayDetail } from '../useDayDetail';

describe('useDayDetail', () => {
  const mockRetry = vi.fn();
  
  const mockPanchangamData = {
    '2024-01-15': {
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
    },
    '2024-01-16': {
      date: '2024-01-16',
      tithi: 'Shashthi',
      nakshatra: 'Mrigashirsha',
      yoga: 'Priti',
      karana: 'Balava',
      sunrise_time: '06:30:00',
      sunset_time: '18:16:00',
      vara: 'Tuesday',
      planetary_ruler: 'Mars',
      events: [],
      festivals: [],
      moonrise_time: '21:15:00',
      moonset_time: '09:30:00',
    }
  };

  describe('data retrieval', () => {
    it('should return null when no date is selected', () => {
      const { result } = renderHook(() => useDayDetail({
        selectedDate: null,
        panchangamData: mockPanchangamData,
        loading: false,
        error: null,
        retry: mockRetry,
      }));

      expect(result.current.data).toBe(null);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBe(null);
    });

    it('should return data for selected date when available', () => {
      const selectedDate = new Date('2024-01-15');
      
      const { result } = renderHook(() => useDayDetail({
        selectedDate,
        panchangamData: mockPanchangamData,
        loading: false,
        error: null,
        retry: mockRetry,
      }));

      expect(result.current.data).toEqual(mockPanchangamData['2024-01-15']);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBe(null);
    });

    it('should return null when data is not available for selected date', () => {
      const selectedDate = new Date('2024-01-17'); // Not in mock data
      
      const { result } = renderHook(() => useDayDetail({
        selectedDate,
        panchangamData: mockPanchangamData,
        loading: false,
        error: null,
        retry: mockRetry,
      }));

      expect(result.current.data).toBe(null);
    });
  });

  describe('loading states', () => {
    it('should show loading when data is not available and currently loading', () => {
      const selectedDate = new Date('2024-01-17'); // Not in mock data
      
      const { result } = renderHook(() => useDayDetail({
        selectedDate,
        panchangamData: mockPanchangamData,
        loading: true,
        error: null,
        retry: mockRetry,
      }));

      expect(result.current.isLoading).toBe(true);
      expect(result.current.data).toBe(null);
    });

    it('should not show loading when data is available even if loading is true', () => {
      const selectedDate = new Date('2024-01-15'); // Available in mock data
      
      const { result } = renderHook(() => useDayDetail({
        selectedDate,
        panchangamData: mockPanchangamData,
        loading: true,
        error: null,
        retry: mockRetry,
      }));

      expect(result.current.isLoading).toBe(false);
      expect(result.current.data).toEqual(mockPanchangamData['2024-01-15']);
    });

    it('should not show loading when not currently loading', () => {
      const selectedDate = new Date('2024-01-17'); // Not in mock data
      
      const { result } = renderHook(() => useDayDetail({
        selectedDate,
        panchangamData: mockPanchangamData,
        loading: false,
        error: null,
        retry: mockRetry,
      }));

      expect(result.current.isLoading).toBe(false);
    });
  });

  describe('error handling', () => {
    it('should show error when data is not available and there is an error', () => {
      const selectedDate = new Date('2024-01-17'); // Not in mock data
      const error = 'Failed to fetch data';
      
      const { result } = renderHook(() => useDayDetail({
        selectedDate,
        panchangamData: mockPanchangamData,
        loading: false,
        error,
        retry: mockRetry,
      }));

      expect(result.current.error).toBe(error);
      expect(result.current.data).toBe(null);
    });

    it('should not show error when data is available even if there is an error', () => {
      const selectedDate = new Date('2024-01-15'); // Available in mock data
      const error = 'Failed to fetch data';
      
      const { result } = renderHook(() => useDayDetail({
        selectedDate,
        panchangamData: mockPanchangamData,
        loading: false,
        error,
        retry: mockRetry,
      }));

      expect(result.current.error).toBe(null);
      expect(result.current.data).toEqual(mockPanchangamData['2024-01-15']);
    });

    it('should not show error when there is no error', () => {
      const selectedDate = new Date('2024-01-17'); // Not in mock data
      
      const { result } = renderHook(() => useDayDetail({
        selectedDate,
        panchangamData: mockPanchangamData,
        loading: false,
        error: null,
        retry: mockRetry,
      }));

      expect(result.current.error).toBe(null);
    });
  });

  describe('retry functionality', () => {
    it('should pass through retry function', () => {
      const { result } = renderHook(() => useDayDetail({
        selectedDate: new Date('2024-01-15'),
        panchangamData: mockPanchangamData,
        loading: false,
        error: null,
        retry: mockRetry,
      }));

      expect(result.current.retry).toBe(mockRetry);
    });
  });

  describe('date string conversion', () => {
    it('should handle different date formats correctly', () => {
      const selectedDate = new Date('2024-01-15T10:30:00Z'); // With time
      
      const { result } = renderHook(() => useDayDetail({
        selectedDate,
        panchangamData: mockPanchangamData,
        loading: false,
        error: null,
        retry: mockRetry,
      }));

      expect(result.current.data).toEqual(mockPanchangamData['2024-01-15']);
    });

    it('should update data when selected date changes', () => {
      const { result, rerender } = renderHook(
        ({ selectedDate }) => useDayDetail({
          selectedDate,
          panchangamData: mockPanchangamData,
          loading: false,
          error: null,
          retry: mockRetry,
        }),
        { initialProps: { selectedDate: new Date('2024-01-15') } }
      );

      expect(result.current.data).toEqual(mockPanchangamData['2024-01-15']);

      // Change selected date
      rerender({ selectedDate: new Date('2024-01-16') });

      expect(result.current.data).toEqual(mockPanchangamData['2024-01-16']);
    });
  });
});