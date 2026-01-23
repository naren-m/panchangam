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

interface LoadingPhase {
  phase: 'today' | 'priority' | 'remaining' | 'complete';
  description: string;
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
  loadingPhase: LoadingPhase;
  todayLoaded: boolean;
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
  loadingPhase,
  todayLoaded,
  retry,
  calendarProps
}) => {
  // Determine current display state
  const isInitialLoading = loading && !hasData && !todayLoaded;
  const hasError = error && !hasData;
  const shouldShowCalendar = hasData || todayLoaded || (!loading && !hasError);
  
  // Get phase-specific loading message
  const getLoadingMessage = () => {
    switch (loadingPhase.phase) {
      case 'today':
        return '🕉️ Loading today\'s tithi...';
      case 'priority':
        return '📅 Loading nearby dates...';
      case 'remaining':
        return '⏳ Loading remaining dates...';
      case 'complete':
        return '✅ All data loaded';
      default:
        return loadingPhase.description;
    }
  };
  
  // Get appropriate progress bar color based on phase
  const getProgressColor = () => {
    switch (loadingPhase.phase) {
      case 'today':
        return 'bg-orange-500';
      case 'priority':
        return 'bg-blue-500';
      case 'remaining':
        return 'bg-green-500';
      case 'complete':
        return 'bg-green-600';
      default:
        return 'bg-blue-500';
    }
  };
  
  const getProgressBgColor = () => {
    switch (loadingPhase.phase) {
      case 'today':
        return 'bg-orange-50 border-orange-200';
      case 'priority':
        return 'bg-blue-50 border-blue-200';
      case 'remaining':
        return 'bg-green-50 border-green-200';
      case 'complete':
        return 'bg-green-50 border-green-200';
      default:
        return 'bg-blue-50 border-blue-200';
    }
  };
  
  const getTextColor = () => {
    switch (loadingPhase.phase) {
      case 'today':
        return 'text-orange-800';
      case 'priority':
        return 'text-blue-800';
      case 'remaining':
        return 'text-green-800';
      case 'complete':
        return 'text-green-800';
      default:
        return 'text-blue-800';
    }
  };

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
          <div className="text-center mb-4">
            <div className="inline-flex items-center px-4 py-2 bg-orange-50 border border-orange-200 rounded-full">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-orange-500 mr-2"></div>
              <span className="text-orange-800 text-sm font-medium">🕉️ Loading today's tithi...</span>
            </div>
          </div>
          <SkeletonCalendar />
        </div>
      )}

      {/* Progressive Loading Indicator - Shows with calendar */}
      {(isProgressiveLoading || (todayLoaded && loading)) && !hasError && (
        <div 
          className={`mb-4 border rounded-lg p-3 transition-all duration-300 ${getProgressBgColor()}`}
          role="progressbar"
          aria-valuenow={Math.round(progress)}
          aria-valuemin={0}
          aria-valuemax={100}
          aria-label={`Loading calendar data: ${Math.round(progress)}% complete`}
        >
          <div className={`flex items-center justify-between text-sm ${getTextColor()}`}>
            <span className="font-medium">{getLoadingMessage()}</span>
            <span aria-live="polite">
              {loadingPhase.phase === 'complete' ? (
                '✅ Complete'
              ) : (
                `${loadedCount}/${totalCount} (${Math.round(progress)}%)`
              )}
            </span>
          </div>
          <div className="mt-2 bg-gray-200 rounded-full h-2 overflow-hidden">
            <div 
              className={`h-2 transition-all duration-500 ease-out ${getProgressColor()}`}
              style={{ width: `${progress}%` }}
              aria-hidden="true"
            />
          </div>
          {loadingPhase.phase === 'today' && todayLoaded && (
            <div className="mt-2 text-xs text-orange-600 flex items-center">
              <span className="inline-block w-2 h-2 bg-orange-500 rounded-full mr-2 animate-pulse"></span>
              Today's tithi is now available!
            </div>
          )}
        </div>
      )}

      {/* Calendar - Lowest Priority but Primary Content */}
      {shouldShowCalendar && !hasError && (
        <div role="main" aria-label="Panchangam calendar">
          {todayLoaded && Object.keys(calendarProps.panchangamData).length === 1 && loading && (
            <div className="mb-4 p-3 bg-gradient-to-r from-orange-50 to-yellow-50 border border-orange-200 rounded-lg">
              <div className="text-center text-orange-800">
                <div className="text-lg font-semibold mb-1">🎉 Today's Tithi Loaded!</div>
                <div className="text-sm opacity-75">Loading additional dates for full month view...</div>
              </div>
            </div>
          )}
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