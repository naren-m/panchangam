import React from 'react';
import { CalendarGrid } from './CalendarGrid';
import { SkeletonCalendar } from '../common/Loading/SkeletonCalendar';
import { ApiError, NetworkError } from '../common/Error';
import { PanchangamData, Settings } from '../../types/panchangam';

interface CalendarProps {
  year: number;
  month: number;
  panchangamData: Record<string, PanchangamData>;
  settings: Settings;
  onDateClick: (date: Date) => void;
}

interface ErrorState {
  hasError: boolean;
  message: string | null;
  statusCode?: number;
  isNetworkError: boolean;
}

interface CalendarDisplayManagerProps {
  loading: boolean;
  hasData: boolean;
  error: string | null;
  errorState: ErrorState;
  isProgressiveLoading: boolean;
  progress: number;
  loadedCount: number;
  totalCount: number;
  retry: () => void;
  calendarProps: CalendarProps;
}

/**
 * CalendarDisplayManager - Ensures only ONE calendar component renders at a time
 * 
 * State Priority (mutually exclusive):
 * 1. Error (highest) - when no data and error exists
 * 2. Loading - when loading and no data exists  
 * 3. Calendar - when data exists or loading is complete
 * 
 * Progressive loading indicator appears alongside calendar when data is loading.
 */
export const CalendarDisplayManager: React.FC<CalendarDisplayManagerProps> = ({
  loading,
  hasData,
  error,
  errorState,
  isProgressiveLoading,
  progress,
  loadedCount,
  totalCount,
  retry,
  calendarProps
}) => {
  // Determine current display state
  const isInitialLoading = loading && !hasData;
  const hasError = error && !hasData;
  const shouldShowCalendar = hasData || (!loading && !hasError);

  return (
    <div className="calendar-display-container">
      {/* Error State - Highest Priority */}
      {hasError && (
        <div className="mb-6" role="alert" aria-live="polite">
          {errorState.isNetworkError ? (
            <NetworkError
              onRetry={retry}
              isRetrying={isProgressiveLoading}
              customMessage={error}
            />
          ) : (
            <ApiError
              error={error}
              onRetry={retry}
              statusCode={errorState.statusCode}
              endpoint="/api/v1/panchangam"
            />
          )}
        </div>
      )}

      {/* Initial Loading State - Medium Priority */}
      {isInitialLoading && !hasError && (
        <div className="mb-6" role="status" aria-live="polite" aria-label="Loading calendar">
          <SkeletonCalendar />
        </div>
      )}

      {/* Progressive Loading Indicator - Shows with calendar */}
      {isProgressiveLoading && hasData && !hasError && (
        <div 
          className="mb-4 bg-blue-50 border border-blue-200 rounded-lg p-3"
          role="progressbar"
          aria-valuenow={progress}
          aria-valuemin={0}
          aria-valuemax={100}
          aria-label={`Loading calendar data: ${progress}% complete`}
        >
          <div className="flex items-center justify-between text-sm text-blue-800">
            <span>Loading calendar data...</span>
            <span aria-live="polite">{loadedCount}/{totalCount} days loaded ({progress}%)</span>
          </div>
          <div className="mt-2 bg-blue-200 rounded-full h-2 overflow-hidden">
            <div 
              className="bg-blue-500 h-2 transition-all duration-300 ease-out"
              style={{ width: `${progress}%` }}
              aria-hidden="true"
            />
          </div>
        </div>
      )}

      {/* Calendar - Lowest Priority but Primary Content */}
      {shouldShowCalendar && !hasError && (
        <div role="main" aria-label="Panchangam calendar">
          <CalendarGrid
            year={calendarProps.year}
            month={calendarProps.month}
            panchangamData={calendarProps.panchangamData}
            settings={calendarProps.settings}
            onDateClick={calendarProps.onDateClick}
          />
        </div>
      )}
    </div>
  );
};