/**
 * 🎯 Unified API Service Hook
 * 
 * Central hook for all backend integrations with enhanced features:
 * - Unified error handling and retry logic
 * - Optimistic updates and caching
 * - Request deduplication
 * - Background refetching
 * - Offline mode detection
 */

import { useState, useEffect, useCallback, useRef } from 'react';
import { apiService } from '../services/apiService';
import type { PanchangamData, Settings } from '../types/panchangam';
import { formatDateForApi } from '../utils/dateHelpers';

interface ApiState<T> {
  data: T | null;
  isLoading: boolean;
  isRetrying: boolean;
  error: string | null;
  retryCount: number;
  lastFetchTime?: number;
  isOffline: boolean;
}

interface UseApiServiceOptions {
  enableRetry?: boolean;
  maxRetries?: number;
  retryDelay?: number;
  backgroundRefetch?: boolean;
  cacheTime?: number;
}

const defaultOptions: UseApiServiceOptions = {
  enableRetry: true,
  maxRetries: 3,
  retryDelay: 1000,
  backgroundRefetch: true,
  cacheTime: 5 * 60 * 1000 // 5 minutes
};

/**
 * 📊 Hook for fetching single panchangam data
 */
export const usePanchangamData = (
  date: Date, 
  settings: Settings,
  options: UseApiServiceOptions = {}
) => {
  const opts = { ...defaultOptions, ...options };
  const [state, setState] = useState<ApiState<PanchangamData>>({
    data: null,
    isLoading: false,
    isRetrying: false,
    error: null,
    retryCount: 0,
    isOffline: false
  });
  
  const abortControllerRef = useRef<AbortController | null>(null);
  const cacheRef = useRef<Map<string, { data: PanchangamData; timestamp: number }>>(new Map());

  const getCacheKey = useCallback(() => {
    return `${formatDateForApi(date)}-${settings.location.latitude}-${settings.location.longitude}`;
  }, [date, settings]);

  const checkCache = useCallback(() => {
    const key = getCacheKey();
    const cached = cacheRef.current.get(key);
    
    if (cached && Date.now() - cached.timestamp < (opts.cacheTime || 0)) {
      return cached.data;
    }
    return null;
  }, [getCacheKey, opts.cacheTime]);

  const fetchData = useCallback(async (isRetry = false) => {
    // Check cache first
    const cachedData = checkCache();
    if (cachedData && !isRetry) {
      setState(prev => ({ ...prev, data: cachedData, isLoading: false }));
      return;
    }

    // Cancel existing request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    abortControllerRef.current = new AbortController();

    setState(prev => ({
      ...prev,
      isLoading: true,
      isRetrying: isRetry,
      retryCount: isRetry ? prev.retryCount + 1 : 0,
      error: null,
      lastFetchTime: Date.now()
    }));

    try {
      // Check API availability
      const isOnline = await apiService.utils.isApiAvailable();
      
      if (!isOnline) {
        setState(prev => ({ ...prev, isOffline: true }));
      }

      const params = apiService.utils.settingsToApiParams(date, settings);
      const data = await apiService.panchangam.get(params);
      
      // Cache the result
      const key = getCacheKey();
      cacheRef.current.set(key, { data, timestamp: Date.now() });
      
      setState(prev => ({
        ...prev,
        data,
        isLoading: false,
        isRetrying: false,
        error: null,
        isOffline: false
      }));

    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        return;
      }

      const shouldRetry = opts.enableRetry && 
                         state.retryCount < (opts.maxRetries || 0) &&
                         !isRetry;

      setState(prev => ({
        ...prev,
        isLoading: false,
        isRetrying: false,
        error: error instanceof Error ? error.message : 'Unknown error',
        isOffline: error instanceof Error && error.message.includes('Network')
      }));

      // Auto-retry with exponential backoff
      if (shouldRetry) {
        const delay = (opts.retryDelay || 1000) * Math.pow(2, state.retryCount);
        setTimeout(() => fetchData(true), delay);
      }
    }
  }, [date, settings, checkCache, getCacheKey, opts, state.retryCount]);

  const retry = useCallback(() => {
    fetchData(true);
  }, [fetchData]);

  const refetch = useCallback(() => {
    // Clear cache and refetch
    cacheRef.current.clear();
    fetchData(false);
  }, [fetchData]);

  useEffect(() => {
    fetchData(false);

    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, [fetchData]);

  // Background refetch on window focus
  useEffect(() => {
    if (!opts.backgroundRefetch) return;

    const handleFocus = () => {
      if (state.data && state.lastFetchTime) {
        const timeSinceLastFetch = Date.now() - state.lastFetchTime;
        if (timeSinceLastFetch > (opts.cacheTime || 0)) {
          fetchData(false);
        }
      }
    };

    window.addEventListener('focus', handleFocus);
    return () => window.removeEventListener('focus', handleFocus);
  }, [fetchData, state.data, state.lastFetchTime, opts]);

  return {
    ...state,
    retry,
    refetch,
    isCached: !!checkCache()
  };
};

/**
 * 📅 Hook for fetching panchangam range data
 */
export const usePanchangamRange = (
  startDate: Date,
  endDate: Date,
  settings: Settings,
  options: UseApiServiceOptions = {}
) => {
  const opts = { ...defaultOptions, ...options };
  const [state, setState] = useState<ApiState<Record<string, PanchangamData>>>({
    data: null,
    isLoading: false,
    isRetrying: false,
    error: null,
    retryCount: 0,
    isOffline: false
  });

  const abortControllerRef = useRef<AbortController | null>(null);

  const fetchRange = useCallback(async (isRetry = false) => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    abortControllerRef.current = new AbortController();

    setState(prev => ({
      ...prev,
      isLoading: true,
      isRetrying: isRetry,
      retryCount: isRetry ? prev.retryCount + 1 : 0,
      error: null,
      lastFetchTime: Date.now()
    }));

    try {
      const startDateStr = formatDateForApi(startDate);
      const endDateStr = formatDateForApi(endDate);
      const baseParams = apiService.utils.settingsToApiParams(startDate, settings);
      
      const { date, ...params } = baseParams; // Remove date from params
      const dataArray = await apiService.panchangam.getRange(startDateStr, endDateStr, params);
      
      // Convert array to map
      const dataMap: Record<string, PanchangamData> = {};
      dataArray.forEach(item => {
        dataMap[item.date] = item;
      });

      setState(prev => ({
        ...prev,
        data: dataMap,
        isLoading: false,
        isRetrying: false,
        error: null,
        isOffline: false
      }));

    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        return;
      }

      setState(prev => ({
        ...prev,
        isLoading: false,
        isRetrying: false,
        error: error instanceof Error ? error.message : 'Unknown error',
        isOffline: error instanceof Error && error.message.includes('Network')
      }));
    }
  }, [startDate, endDate, settings]);

  const retry = useCallback(() => {
    fetchRange(true);
  }, [fetchRange]);

  useEffect(() => {
    fetchRange(false);

    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, [fetchRange]);

  return {
    ...state,
    retry
  };
};

/**
 * 🏥 Hook for API health monitoring
 */
export const useApiHealth = () => {
  const [health, setHealth] = useState<{
    status: 'healthy' | 'unhealthy' | 'checking';
    message: string;
    metrics?: any;
  }>({ status: 'checking', message: 'Checking API status...' });

  const checkHealth = useCallback(async () => {
    try {
      setHealth(prev => ({ ...prev, status: 'checking' }));
      
      const [healthResult, metrics] = await Promise.all([
        apiService.panchangam.healthCheck(),
        apiService.system.getMetrics()
      ]);

      setHealth({
        status: healthResult.status,
        message: healthResult.message,
        metrics
      });
    } catch (error) {
      setHealth({
        status: 'unhealthy',
        message: error instanceof Error ? error.message : 'Health check failed'
      });
    }
  }, []);

  useEffect(() => {
    checkHealth();
    
    // Periodic health checks
    const interval = setInterval(checkHealth, 30000);
    return () => clearInterval(interval);
  }, [checkHealth]);

  return {
    ...health,
    refresh: checkHealth
  };
};