import { beforeEach, describe, expect, it, vi } from 'vitest';
import { act, renderHook, waitFor } from '@testing-library/react';

import { useProgressivePanchangam } from '../useProgressivePanchangam';
import { panchangamApiClient } from '../../services/api/panchangamApiClient';
import type { PanchangamData, Settings } from '../../types/panchangam';
import { formatDateForApi } from '../../utils/dateHelpers';

vi.mock('../../services/api/panchangamApiClient', () => ({
  panchangamApiClient: {
    getPanchangam: vi.fn(),
  },
}));

const settings: Settings = {
  calculation_method: 'Drik',
  locale: 'en',
  region: 'California',
  time_format: '12',
  location: {
    name: 'Milpitas, California',
    latitude: 37.4323,
    longitude: -121.9066,
    timezone: 'America/Los_Angeles',
    region: 'California',
  },
};

const mockedPanchangamApiClient = vi.mocked(panchangamApiClient);

const createDeferred = <T,>() => {
  let resolve!: (value: T) => void;
  let reject!: (reason?: unknown) => void;
  const promise = new Promise<T>((promiseResolve, promiseReject) => {
    resolve = promiseResolve;
    reject = promiseReject;
  });

  return { promise, resolve, reject };
};

const createPanchangamData = (date: string): PanchangamData => ({
  date,
  tithi: 'Pratipada',
  nakshatra: 'Ashwini',
  yoga: 'Vishkambha',
  karana: 'Bava',
  sunrise_time: '6:30 AM',
  sunset_time: '6:30 PM',
  events: [],
  vara: 'Monday',
  planetary_ruler: 'Moon',
});

describe('useProgressivePanchangam', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('reports network errors without logging batch setup failures', async () => {
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    mockedPanchangamApiClient.getPanchangam.mockImplementation(() => {
      throw new Error('Network setup failed');
    });

    const today = new Date();
    const { result } = renderHook(() => useProgressivePanchangam(today, today, settings));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.errorState).toEqual({
      hasError: true,
      message: 'Backend server is not available. Please ensure the Panchangam API server is running.',
      isNetworkError: true,
    });
    expect(consoleError).not.toHaveBeenCalled();
  });

  it('reports rejected API network errors', async () => {
    mockedPanchangamApiClient.getPanchangam.mockRejectedValue({
      code: 'REQUEST_TIMEOUT',
      message: 'Request timeout',
    });

    const today = new Date();
    const { result } = renderHook(() => useProgressivePanchangam(today, today, settings));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.errorState).toEqual({
      hasError: true,
      message: 'Backend server is not available. Please ensure the Panchangam API server is running.',
      isNetworkError: true,
    });
  });

  it('ignores stale request results after the visible date range changes', async () => {
    const firstDate = new Date();
    const secondDate = new Date(firstDate);
    secondDate.setDate(firstDate.getDate() + 1);

    const firstKey = formatDateForApi(firstDate);
    const secondKey = formatDateForApi(secondDate);
    const firstRequest = createDeferred<PanchangamData>();
    const secondRequest = createDeferred<PanchangamData>();

    mockedPanchangamApiClient.getPanchangam.mockImplementation(({ date }) => {
      if (date === firstKey) {
        return firstRequest.promise;
      }
      if (date === secondKey) {
        return secondRequest.promise;
      }
      return Promise.reject(new Error(`unexpected date ${date}`));
    });

    const { result, rerender } = renderHook(
      ({ startDate, endDate }) => useProgressivePanchangam(startDate, endDate, settings),
      {
        initialProps: {
          startDate: firstDate,
          endDate: firstDate,
        },
      }
    );

    await waitFor(() => {
      expect(mockedPanchangamApiClient.getPanchangam).toHaveBeenCalledWith(expect.objectContaining({ date: firstKey }));
    });

    rerender({
      startDate: secondDate,
      endDate: secondDate,
    });

    await waitFor(() => {
      expect(mockedPanchangamApiClient.getPanchangam).toHaveBeenCalledWith(expect.objectContaining({ date: secondKey }));
    });

    await act(async () => {
      secondRequest.resolve(createPanchangamData(secondKey));
      await secondRequest.promise;
    });

    await waitFor(() => {
      expect(result.current.loadedCount).toBe(1);
      expect(result.current.progress).toBe(100);
      expect(result.current.totalCount).toBe(1);
    });

    await act(async () => {
      firstRequest.resolve(createPanchangamData(firstKey));
      await firstRequest.promise;
    });

    expect(result.current.loadedCount).toBe(1);
    expect(result.current.progress).toBe(100);
    expect(result.current.totalCount).toBe(1);
  });

  it('counts calendar days across daylight-saving changes', async () => {
    mockedPanchangamApiClient.getPanchangam.mockImplementation(({ date }) =>
      Promise.resolve(createPanchangamData(date))
    );
    const startDate = new Date(2024, 10, 1);
    const endDate = new Date(2024, 10, 30);

    const { result } = renderHook(() =>
      useProgressivePanchangam(startDate, endDate, settings)
    );

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.totalCount).toBe(30);
    expect(result.current.loadedCount).toBe(30);
    expect(result.current.progress).toBe(100);
  });
});
