import { useState, useEffect, useCallback, useRef } from 'react';
import { PanchangamData, Settings } from '../types/panchangam';
import { panchangamApiClient } from '../services/api/panchangamApiClient';
import { formatDateForApi } from '../utils/dateHelpers';

interface LoadingState {
  isLoading: boolean;
  isRetrying: boolean;
  retryCount: number;
  lastFetchTime?: number;
}

interface ErrorState {
  hasError: boolean;
  message: string | null;
  statusCode?: number;
  isNetworkError: boolean;
}

function clearError(): ErrorState {
  return {
    hasError: false,
    message: null,
    isNetworkError: false,
  };
}

function getErrorState(error: unknown): ErrorState {
  const message = error instanceof Error ? error.message : 'Failed to fetch panchangam data';
  const statusFromError = typeof error === 'object' && error !== null && 'status' in error
    ? Number((error as { status?: unknown }).status)
    : undefined;
  const statusFromMessage = message.match(/\b[1-5]\d{2}\b/)?.[0];
  const statusCode = Number.isFinite(statusFromError)
    ? statusFromError
    : statusFromMessage
      ? Number(statusFromMessage)
      : undefined;
  const lowerMessage = message.toLowerCase();

  return {
    hasError: true,
    message,
    statusCode,
    isNetworkError:
      lowerMessage.includes('failed to fetch') ||
      lowerMessage.includes('network') ||
      lowerMessage.includes('timeout') ||
      lowerMessage.includes('err_name_not_resolved') ||
      lowerMessage.includes('err_insufficient_resources'),
  };
}

function buildSingleRequest(date: Date, settings: Settings) {
  return {
    date: formatDateForApi(date),
    latitude: settings.location.latitude,
    longitude: settings.location.longitude,
    timezone: settings.location.timezone,
    region: settings.region,
    calculation_method: settings.calculation_method,
    locale: settings.locale
  };
}

function buildRangeRequest(settings: Settings) {
  return {
    latitude: settings.location.latitude,
    longitude: settings.location.longitude,
    timezone: settings.location.timezone,
    region: settings.region,
    calculation_method: settings.calculation_method,
    locale: settings.locale
  };
}

export const usePanchangam = (date: Date, settings: Settings) => {
  const [data, setData] = useState<PanchangamData | null>(null);
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: true,
    isRetrying: false,
    retryCount: 0,
  });
  const [errorState, setErrorState] = useState<ErrorState>(clearError);
  const abortControllerRef = useRef<AbortController | null>(null);
  const requestIdRef = useRef(0);

  const fetchPanchangam = useCallback(async (isRetry = false) => {
    abortControllerRef.current?.abort();

    const requestId = requestIdRef.current + 1;
    requestIdRef.current = requestId;
    const controller = new AbortController();
    abortControllerRef.current = controller;

    setLoadingState(prev => ({
      isLoading: true,
      isRetrying: isRetry,
      retryCount: isRetry ? prev.retryCount + 1 : 0,
      lastFetchTime: Date.now(),
    }));
    setErrorState(clearError());

    try {
      const response = await panchangamApiClient.getPanchangam(buildSingleRequest(date, settings));

      if (controller.signal.aborted || requestId !== requestIdRef.current) {
        return;
      }

      setData(response);
      setErrorState(clearError());
    } catch (error) {
      if (controller.signal.aborted || requestId !== requestIdRef.current) {
        return;
      }

      setErrorState(getErrorState(error));
    } finally {
      if (!controller.signal.aborted && requestId === requestIdRef.current) {
        setLoadingState(prev => ({
          ...prev,
          isLoading: false,
          isRetrying: false,
        }));
      }
    }
  }, [date, settings]);

  const retry = useCallback(() => {
    void fetchPanchangam(true);
  }, [fetchPanchangam]);

  useEffect(() => {
    void fetchPanchangam(false);

    return () => {
      requestIdRef.current += 1;
      abortControllerRef.current?.abort();
    };
  }, [fetchPanchangam]);

  return {
    data,
    loading: loadingState.isLoading,
    isRetrying: loadingState.isRetrying,
    retryCount: loadingState.retryCount,
    error: errorState.hasError ? errorState.message : null,
    errorState,
    retry,
  };
};

export const usePanchangamRange = (startDate: Date, endDate: Date, settings: Settings) => {
  const [data, setData] = useState<Record<string, PanchangamData>>({});
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: true,
    isRetrying: false,
    retryCount: 0,
  });
  const [errorState, setErrorState] = useState<ErrorState>(clearError);
  const abortControllerRef = useRef<AbortController | null>(null);
  const requestIdRef = useRef(0);

  const fetchPanchangamRange = useCallback(async (isRetry = false) => {
    abortControllerRef.current?.abort();

    const requestId = requestIdRef.current + 1;
    requestIdRef.current = requestId;
    const controller = new AbortController();
    abortControllerRef.current = controller;

    setLoadingState(prev => ({
      isLoading: true,
      isRetrying: isRetry,
      retryCount: isRetry ? prev.retryCount + 1 : 0,
      lastFetchTime: Date.now(),
    }));
    setErrorState(clearError());

    try {
      const response = await panchangamApiClient.getPanchangamRange(
        formatDateForApi(startDate),
        formatDateForApi(endDate),
        buildRangeRequest(settings)
      );

      if (controller.signal.aborted || requestId !== requestIdRef.current) {
        return;
      }

      const dataMap: Record<string, PanchangamData> = {};
      response.forEach(item => {
        dataMap[item.date] = item;
      });

      setData(dataMap);
      setErrorState(clearError());
    } catch (error) {
      if (controller.signal.aborted || requestId !== requestIdRef.current) {
        return;
      }

      setErrorState(getErrorState(error));
    } finally {
      if (!controller.signal.aborted && requestId === requestIdRef.current) {
        setLoadingState(prev => ({
          ...prev,
          isLoading: false,
          isRetrying: false,
        }));
      }
    }
  }, [startDate, endDate, settings]);

  const retry = useCallback(() => {
    void fetchPanchangamRange(true);
  }, [fetchPanchangamRange]);

  useEffect(() => {
    void fetchPanchangamRange(false);

    return () => {
      requestIdRef.current += 1;
      abortControllerRef.current?.abort();
    };
  }, [fetchPanchangamRange]);

  return {
    data,
    loading: loadingState.isLoading,
    isRetrying: loadingState.isRetrying,
    retryCount: loadingState.retryCount,
    error: errorState.hasError ? errorState.message : null,
    errorState,
    retry,
  };
};
