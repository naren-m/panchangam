import { useState, useEffect, useCallback, useRef } from 'react';
import { PanchangamData, Settings } from '../types/panchangam';
import { panchangamApi } from '../services/panchangamApi';
import { formatDateForApi } from '../utils/dateHelpers';
import { requestCache } from '../services/api/requestCache';

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
  const debounceTimerRef = useRef<NodeJS.Timeout | null>(null);
  const lastRequestRef = useRef<string>('');

  const fetchPanchangam = useCallback(async (isRetry = false) => {
    // Create request signature for deduplication
    const requestSignature = `${formatDateForApi(date)}-${JSON.stringify(settings)}`;
    
    // Prevent duplicate requests
    if (requestSignature === lastRequestRef.current && !isRetry) {
      return;
    }
    lastRequestRef.current = requestSignature;

    // Check cache first for non-retry requests
    if (!isRetry) {
      const cacheKey = formatDateForApi(date);
      const cachedData = requestCache.get<PanchangamData>('panchangam', {
        date: cacheKey,
        settings: JSON.stringify(settings)
      });
      
      if (cachedData) {
        setData(cachedData);
        setLoadingState(prev => ({ ...prev, isLoading: false }));
        return;
      }
    }

    // Cancel any existing request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    
    // Clear any pending debounced request
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current);
      debounceTimerRef.current = null;
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
      // Add a small delay to prevent rapid-fire requests
      if (!isRetry) {
        await new Promise(resolve => setTimeout(resolve, 50));
      }

      const response = await panchangamApi.getPanchangam({
        date: formatDateForApi(date),
        latitude: settings.location.latitude,
        longitude: settings.location.longitude,
        timezone: settings.location.timezone,
        region: settings.region,
        calculation_method: settings.calculation_method,
        locale: settings.locale
      });

      // Cache the response for future use
      const cacheKey = formatDateForApi(date);
      requestCache.set('panchangam', {
        date: cacheKey,
        settings: JSON.stringify(settings)
      }, response, 2 * 60 * 1000); // 2 minutes cache

      setData(response);
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        // Request was cancelled, don't update error state
        return;
      }
      
      const isNetworkError = err instanceof Error && (
        err.message.includes('Failed to fetch') ||
        err.message.includes('Network') ||
        err.message.includes('timeout') ||
        err.message.includes('ERR_NAME_NOT_RESOLVED') ||
        err.message.includes('ERR_INSUFFICIENT_RESOURCES')
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

      // Implement exponential backoff for retries
      if (isRetry && loadingState.retryCount < 3) {
        const backoffDelay = Math.min(1000 * Math.pow(2, loadingState.retryCount), 5000);
        setTimeout(() => {
          if (!abortControllerRef.current?.signal.aborted) {
            fetchPanchangam(true);
          }
        }, backoffDelay);
      }
    } finally {
      setLoadingState(prev => ({
        ...prev,
        isLoading: false,
        isRetrying: false,
      }));
    }
  }, [date, settings, loadingState.retryCount]);

  // Debounced fetch function to prevent rapid-fire requests
  const debouncedFetch = useCallback((isRetry = false) => {
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current);
    }
    
    debounceTimerRef.current = setTimeout(() => {
      fetchPanchangam(isRetry);
    }, isRetry ? 0 : 200); // No delay for retries, 200ms delay for new requests
  }, [fetchPanchangam]);

  const retry = useCallback(() => {
    debouncedFetch(true);
  }, [debouncedFetch]);

  useEffect(() => {
    debouncedFetch(false);
    
    // Cleanup function to cancel requests on unmount
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current);
      }
    };
  }, [debouncedFetch]);

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
  const debounceTimerRef = useRef<NodeJS.Timeout | null>(null);
  const lastRequestRef = useRef<string>('');

  const fetchPanchangamRange = useCallback(async (isRetry = false) => {
    // Create request signature for deduplication
    const requestSignature = `${formatDateForApi(startDate)}-${formatDateForApi(endDate)}-${JSON.stringify(settings)}`;
    
    // Prevent duplicate requests
    if (requestSignature === lastRequestRef.current && !isRetry) {
      return;
    }
    lastRequestRef.current = requestSignature;

    // Cancel any existing request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    
    // Clear any pending debounced request
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current);
      debounceTimerRef.current = null;
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
      // Add a small delay to prevent rapid-fire requests
      if (!isRetry) {
        await new Promise(resolve => setTimeout(resolve, 100));
      }

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
        err.message.includes('timeout') ||
        err.message.includes('ERR_NAME_NOT_RESOLVED') ||
        err.message.includes('ERR_INSUFFICIENT_RESOURCES')
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

      // Implement exponential backoff for retries
      if (isRetry && loadingState.retryCount < 3) {
        const backoffDelay = Math.min(1000 * Math.pow(2, loadingState.retryCount), 10000);
        setTimeout(() => {
          if (!abortControllerRef.current?.signal.aborted) {
            fetchPanchangamRange(true);
          }
        }, backoffDelay);
      }
    } finally {
      setLoadingState(prev => ({
        ...prev,
        isLoading: false,
        isRetrying: false,
      }));
    }
  }, [startDate, endDate, settings, loadingState.retryCount]);

  // Debounced fetch function to prevent rapid-fire requests
  const debouncedFetch = useCallback((isRetry = false) => {
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current);
    }
    
    debounceTimerRef.current = setTimeout(() => {
      fetchPanchangamRange(isRetry);
    }, isRetry ? 0 : 300); // No delay for retries, 300ms delay for new requests
  }, [fetchPanchangamRange]);

  const retry = useCallback(() => {
    debouncedFetch(true);
  }, [debouncedFetch]);

  useEffect(() => {
    debouncedFetch(false);
    
    // Cleanup function to cancel requests on unmount
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current);
      }
    };
  }, [debouncedFetch]);

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