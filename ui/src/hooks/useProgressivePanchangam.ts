import { useState, useEffect } from 'react';
import { PanchangamData, Settings } from '../types/panchangam';
import { panchangamApi } from '../services/panchangamApi';
import { formatDateForApi } from '../utils/dateHelpers';

interface ErrorState {
  hasError: boolean;
  message: string | null;
  statusCode?: number;
  isNetworkError: boolean;
}

export const useProgressivePanchangam = (
  startDate: Date, 
  endDate: Date, 
  settings: Settings
) => {
  const [data, setData] = useState<Record<string, PanchangamData>>({});
  const [loading, setLoading] = useState(true);
  const [errorState, setErrorState] = useState<ErrorState>({
    hasError: false,
    message: null,
    isNetworkError: false,
  });

  // Generate dates for the range
  const generateDateRange = (start: Date, end: Date): string[] => {
    const dates: string[] = [];
    const current = new Date(start);
    
    while (current <= end) {
      dates.push(formatDateForApi(current));
      current.setDate(current.getDate() + 1);
    }
    
    return dates;
  };

  // Load all data at once (simplified approach)
  useEffect(() => {
    console.log('🔥 useProgressivePanchangam useEffect triggered', {
      startDate: startDate.toISOString(),
      endDate: endDate.toISOString(),
      settings_calc_method: settings.calculation_method,
      settings_locale: settings.locale,
      settings_name: settings.location.name
    });
    
    let isCancelled = false;
    
    const loadData = async () => {
      setLoading(true);
      setErrorState({
        hasError: false,
        message: null,
        isNetworkError: false,
      });

      try {
        const dates = generateDateRange(startDate, endDate);
        const today = formatDateForApi(new Date());
        
        // Load today first, then load the rest
        const todayIndex = dates.indexOf(today);
        let orderedDates = dates;
        
        if (todayIndex !== -1) {
          // Put today first
          orderedDates = [today, ...dates.filter(d => d !== today)];
        }

        const results: Record<string, PanchangamData> = {};
        let todayLoaded = false;

        // Load dates sequentially to avoid overloading the server
        for (const dateStr of orderedDates) {
          if (isCancelled) break;

          try {
            const response = await panchangamApi.getPanchangam({
              date: dateStr,
              latitude: settings.location.latitude,
              longitude: settings.location.longitude,
              timezone: settings.location.timezone,
              region: settings.region,
              calculation_method: settings.calculation_method,
              locale: settings.locale
            });

            results[dateStr] = response;

            // Update data incrementally
            if (!isCancelled) {
              setData(prev => ({
                ...prev,
                [dateStr]: response
              }));

              // If today is loaded, we can show the calendar
              if (dateStr === today) {
                todayLoaded = true;
                setLoading(false);
              }
            }
            
            // Add small delay to prevent overwhelming the server
            await new Promise(resolve => setTimeout(resolve, 50));
            
          } catch (err) {
            console.error(`Failed to load date ${dateStr}:`, err);
            
            // Only set error state if today failed
            if (dateStr === today) {
              const isNetworkError = err instanceof Error && (
                err.message.includes('Failed to fetch') ||
                err.message.includes('Network') ||
                err.message.includes('timeout')
              );

              setErrorState({
                hasError: true,
                message: err instanceof Error ? err.message : 'Failed to fetch panchangam data',
                isNetworkError,
              });
              
              setLoading(false);
              return;
            }
          }
        }
        
        // If today wasn't in the range, stop loading anyway
        if (!todayLoaded && !isCancelled) {
          setLoading(false);
        }

      } catch (error) {
        if (!isCancelled) {
          console.error('Failed to load panchangam data:', error);
          
          const isNetworkError = error instanceof Error && (
            error.message.includes('Failed to fetch') ||
            error.message.includes('Network') ||
            error.message.includes('timeout')
          );

          setErrorState({
            hasError: true,
            message: error instanceof Error ? error.message : 'Failed to fetch panchangam data',
            isNetworkError,
          });
          
          setLoading(false);
        }
      }
    };

    setData({});
    loadData();

    return () => {
      isCancelled = true;
    };
  }, [startDate, endDate, settings]);

  // Retry function
  const retry = () => {
    setErrorState({
      hasError: false,
      message: null,
      isNetworkError: false,
    });
    // The useEffect will automatically run again due to dependency changes
  };

  // Calculate metrics
  const totalCount = Math.ceil((endDate.getTime() - startDate.getTime()) / (24 * 60 * 60 * 1000)) + 1;
  const loadedCount = Object.keys(data).length;
  const progress = Math.round((loadedCount / totalCount) * 100);
  const today = formatDateForApi(new Date());
  const todayLoaded = data[today] !== undefined;

  return {
    data,
    loading,
    isProgressiveLoading: loadedCount < totalCount && !loading,
    progress,
    todayLoaded,
    loadedCount,
    totalCount,
    error: errorState.hasError ? errorState.message : null,
    errorState,
    retry,
  };
};