import { useState, useEffect, useCallback, useRef } from 'react';
import { PanchangamData, Settings } from '../types/panchangam';
import { panchangamApi } from '../services/panchangamApi';
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

export const usePanchangam = (date: Date, settings: Settings) => {
  const [data, setData] = useState<PanchangamData | null>(null);
  const [loadingState, setLoadingState] = useState<LoadingState>({
    isLoading: false,
    isRetrying: false,
    retryCount: 0,
  });
  const [errorState, setErrorState] = useState<ErrorState>({
    hasError: false,
    message: null,
    isNetworkError: false,
  });
  const abortControllerRef = useRef<AbortController | null>(null);

  const fetchPanchangam = useCallback(async (isRetry = false) => {
    // Cancel any existing request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    
    // Create new abort controller
    abortControllerRef.current = new AbortController();
    
    setLoadingState(prev => ({
      ...prev,
      isLoading: true,
      isRetrying: isRetry,
      retryCount: isRetry ? prev.retryCount + 1 : 0,
      lastFetchTime: Date.now(),
    }));
    
    setErrorState({
      hasError: false,
      message: null,
      isNetworkError: false,
    });

    try {
      const response = await panchangamApi.getPanchangam({
        date: formatDateForApi(date),
        latitude: settings.location.latitude,
        longitude: settings.location.longitude,
        timezone: settings.location.timezone,
        region: settings.region,
        calculation_method: settings.calculation_method,
        locale: settings.locale
      });

      setData(response);
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        // Request was cancelled, don't update error state
        return;
      }
      
      const isNetworkError = err instanceof Error && (
        err.message.includes('Failed to fetch') ||
        err.message.includes('Network') ||
        err.message.includes('timeout')
      );
      
      const statusCode = err instanceof Error && err.message.includes('API request failed:') 
        ? parseInt(err.message.match(/\d{3}/)?.[0] || '0')
        : undefined;
      
      setErrorState({
        hasError: true,
        message: err instanceof Error ? err.message : 'Failed to fetch panchangam data',
        statusCode,
        isNetworkError,
      });
    } finally {
      setLoadingState(prev => ({
        ...prev,
        isLoading: false,
        isRetrying: false,
      }));
    }
  }, [date, settings]);

  const retry = useCallback(() => {
    fetchPanchangam(true);
  }, [fetchPanchangam]);

  useEffect(() => {
    fetchPanchangam(false);
    
    // Cleanup function to cancel requests on unmount
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
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
    isLoading: false,
    isRetrying: false,
    retryCount: 0,
  });
  const [errorState, setErrorState] = useState<ErrorState>({
    hasError: false,
    message: null,
    isNetworkError: false,
  });
  const abortControllerRef = useRef<AbortController | null>(null);

  const fetchPanchangamRange = useCallback(async (isRetry = false) => {
    // Cancel any existing request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    
    // Create new abort controller
    abortControllerRef.current = new AbortController();
    
    setLoadingState(prev => ({
      ...prev,
      isLoading: true,
      isRetrying: isRetry,
      retryCount: isRetry ? prev.retryCount + 1 : 0,
      lastFetchTime: Date.now(),
    }));
    
    setErrorState({
      hasError: false,
      message: null,
      isNetworkError: false,
    });

    try {
      const response = await panchangamApi.getPanchangamRange(
        formatDateForApi(startDate),
        formatDateForApi(endDate),
        {
          latitude: settings.location.latitude,
          longitude: settings.location.longitude,
          timezone: settings.location.timezone,
          region: settings.region,
          calculation_method: settings.calculation_method,
          locale: settings.locale
        }
      );

      const dataMap: Record<string, PanchangamData> = {};
      response.forEach(item => {
        dataMap[item.date] = item;
      });

      setData(dataMap);
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        // Request was cancelled, don't update error state
        return;
      }
      
      const isNetworkError = err instanceof Error && (
        err.message.includes('Failed to fetch') ||
        err.message.includes('Network') ||
        err.message.includes('timeout')
      );
      
      const statusCode = err instanceof Error && err.message.includes('API request failed:') 
        ? parseInt(err.message.match(/\d{3}/)?.[0] || '0')
        : undefined;
      
      setErrorState({
        hasError: true,
        message: err instanceof Error ? err.message : 'Failed to fetch panchangam data',
        statusCode,
        isNetworkError,
      });
    } finally {
      setLoadingState(prev => ({
        ...prev,
        isLoading: false,
        isRetrying: false,
      }));
    }
  }, [startDate, endDate, settings]);

  const retry = useCallback(() => {
    fetchPanchangamRange(true);
  }, [fetchPanchangamRange]);

  useEffect(() => {
    fetchPanchangamRange(false);
    
    // Cleanup function to cancel requests on unmount
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
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