import { useState, useEffect, useCallback, useRef } from 'react';
import { PanchangamData, Settings } from '../types/panchangam';
import { panchangamApiClient } from '../services/api/panchangamApiClient';
import { formatDateForApi } from '../utils/dateHelpers';

interface LoadingPhase {
  phase: 'today' | 'priority' | 'remaining' | 'complete';
  description: string;
}

interface UseProgressivePanchangamReturn {
  data: Record<string, PanchangamData>;
  loading: boolean;
  isProgressiveLoading: boolean;
  progress: number;
  todayLoaded: boolean;
  loadedCount: number;
  totalCount: number;
  error: string | null;
  errorState: {
    hasError: boolean;
    message: string | null;
    isNetworkError: boolean;
  };
  loadingPhase: LoadingPhase;
  retry: () => void;
}

/**
 * Hook for truly progressive loading of panchangam data
 * Phase 1: Load today's data first for immediate display
 * Phase 2: Load ±5 days around today for quick navigation
 * Phase 3: Load remaining dates in the month
 */
export function useProgressivePanchangam(
  startDate: Date,
  endDate: Date,
  settings: Settings
): UseProgressivePanchangamReturn {
  const [allData, setAllData] = useState<Record<string, PanchangamData>>({});
  const [loadedCount, setLoadedCount] = useState(0);
  const [progress, setProgress] = useState(0);
  const [todayLoaded, setTodayLoaded] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [errorState, setErrorState] = useState({
    hasError: false,
    message: null as string | null,
    isNetworkError: false
  });
  const [loadingPhase, setLoadingPhase] = useState<LoadingPhase>({
    phase: 'today',
    description: 'Loading today\'s tithi...'
  });
  const abortControllerRef = useRef<AbortController | null>(null);
  
  // Calculate total days
  const totalDays = Math.ceil((endDate.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24)) + 1;
  
  // Generate date arrays for progressive loading
  const getDatesForProgressiveLoading = useCallback(() => {
    const today = new Date();
    const todayStr = today.toISOString().split('T')[0];
    
    // All dates in the range
    const allDates: Date[] = [];
    for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 1)) {
      allDates.push(new Date(d));
    }
    
    // Phase 1: Today (if in range)
    const todayDates = allDates.filter(d => 
      d.toISOString().split('T')[0] === todayStr
    );
    
    // Phase 2: Priority dates (±5 days from today, excluding today)
    const priorityDates = allDates.filter(d => {
      const dateStr = d.toISOString().split('T')[0];
      if (dateStr === todayStr) return false;
      
      const daysDiff = Math.abs((d.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
      return daysDiff <= 5;
    });
    
    // Phase 3: Remaining dates
    const remainingDates = allDates.filter(d => {
      const dateStr = d.toISOString().split('T')[0];
      if (dateStr === todayStr) return false;
      
      const daysDiff = Math.abs((d.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
      return daysDiff > 5;
    });
    
    return { todayDates, priorityDates, remainingDates, allDates };
  }, [startDate, endDate]);
  
  // Fetch data for a specific set of dates
  const fetchDatesData = useCallback(async (dates: Date[], phase: LoadingPhase['phase']) => {
    if (dates.length === 0) return;
    
    try {
      // Fetch data for all dates in parallel (but limited batch size)
      const batchSize = 5;
      const batches: Date[][] = [];
      
      for (let i = 0; i < dates.length; i += batchSize) {
        batches.push(dates.slice(i, i + batchSize));
      }
      
      for (const batch of batches) {
        if (abortControllerRef.current?.signal.aborted) return;
        
        const promises = batch.map(date => 
          panchangamApiClient.getPanchangam({
            date: formatDateForApi(date),
            latitude: settings.location.latitude,
            longitude: settings.location.longitude,
            timezone: settings.location.timezone,
            region: settings.region,
            calculation_method: settings.calculation_method,
            locale: settings.locale
          })
        );
        
        const results = await Promise.allSettled(promises);
        
        // Process successful results
        const newData: Record<string, PanchangamData> = {};
        results.forEach((result, index) => {
          if (result.status === 'fulfilled') {
            const dateStr = formatDateForApi(batch[index]);
            newData[dateStr] = result.value;
          }
        });
        
        // Update state with new data
        setAllData(prev => ({ ...prev, ...newData }));
        setLoadedCount(prev => prev + Object.keys(newData).length);
        
        // Check if today is loaded
        const todayStr = new Date().toISOString().split('T')[0];
        if (newData[todayStr] && !todayLoaded) {
          setTodayLoaded(true);
        }
      }
    } catch (err) {
      console.error(`Error fetching ${phase} data:`, err);
    }
  }, [settings, todayLoaded]);
  
  // Main progressive loading function
  const loadProgressively = useCallback(async () => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    
    abortControllerRef.current = new AbortController();
    
    setLoading(true);
    setError(null);
    setErrorState({ hasError: false, message: null, isNetworkError: false });
    setAllData({});
    setLoadedCount(0);
    setProgress(0);
    setTodayLoaded(false);
    
    try {
      const { todayDates, priorityDates, remainingDates } = getDatesForProgressiveLoading();
      
      // Phase 1: Load today's data first
      if (todayDates.length > 0) {
        setLoadingPhase({ phase: 'today', description: 'Loading today\'s tithi...' });
        await fetchDatesData(todayDates, 'today');
        
        if (abortControllerRef.current?.signal.aborted) return;
      }
      
      // Phase 2: Load priority dates (±5 days)
      if (priorityDates.length > 0) {
        setLoadingPhase({ phase: 'priority', description: 'Loading nearby dates...' });
        await fetchDatesData(priorityDates, 'priority');
        
        if (abortControllerRef.current?.signal.aborted) return;
      }
      
      // Phase 3: Load remaining dates
      if (remainingDates.length > 0) {
        setLoadingPhase({ phase: 'remaining', description: 'Loading remaining dates...' });
        await fetchDatesData(remainingDates, 'remaining');
        
        if (abortControllerRef.current?.signal.aborted) return;
      }
      
      // Complete
      setLoadingPhase({ phase: 'complete', description: 'All data loaded' });
      
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        return;
      }
      
      const isNetworkError = err instanceof Error && (
        err.message.includes('Failed to fetch') ||
        err.message.includes('Network') ||
        err.message.includes('timeout')
      );
      
      setError(err instanceof Error ? err.message : 'Failed to fetch panchangam data');
      setErrorState({
        hasError: true,
        message: err instanceof Error ? err.message : 'Failed to fetch panchangam data',
        isNetworkError
      });
    } finally {
      setLoading(false);
    }
  }, [getDatesForProgressiveLoading, fetchDatesData]);
  
  // Update progress based on loaded count
  useEffect(() => {
    if (totalDays > 0) {
      setProgress((loadedCount / totalDays) * 100);
    }
  }, [loadedCount, totalDays]);
  
  // Retry function
  const retry = useCallback(() => {
    loadProgressively();
  }, [loadProgressively]);
  
  // Load data when dependencies change
  useEffect(() => {
    loadProgressively();
    
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, [loadProgressively]);
  
  return {
    data: allData,
    loading,
    isProgressiveLoading: loading && loadedCount > 0,
    progress,
    todayLoaded,
    loadedCount,
    totalCount: totalDays,
    error,
    errorState,
    loadingPhase,
    retry
  };
}