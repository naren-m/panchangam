import React, { useState, useEffect } from 'react';
import { CalendarGrid } from './components/Calendar/CalendarGrid';
import { MonthNavigation } from './components/Calendar/MonthNavigation';
import { DayDetailModal } from './components/DayDetail/DayDetailModal';
import { LocationSelector } from './components/LocationPicker/LocationSelector';
import { SettingsPanel } from './components/Settings/SettingsPanel';
import { SkeletonCalendar, LoadingSpinner } from './components/common/Loading';
import { ApiError, NetworkError, ErrorBoundary, OfflineIndicator } from './components/common/Error';
import { usePanchangamRange } from './hooks/usePanchangam';
import { useDayDetail } from './hooks/useDayDetail';
import { Settings, PanchangamData } from './types/panchangam';
import { getCurrentMonthDates } from './utils/dateHelpers';
import { locationService } from './services/locationService';

function App() {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [showLocationSelector, setShowLocationSelector] = useState(false);
  const [showSettings, setShowSettings] = useState(false);
  const [settings, setSettings] = useState<Settings>({
    calculation_method: 'Drik',
    locale: 'en',
    region: 'Tamil Nadu',
    time_format: '12',
    location: {
      name: 'Chennai, Tamil Nadu',
      latitude: 13.0827,
      longitude: 80.2707,
      timezone: 'Asia/Kolkata',
      region: 'Tamil Nadu'
    }
  });

  const year = currentDate.getFullYear();
  const month = currentDate.getMonth();

  // Get the date range for the current month view
  const monthDates = getCurrentMonthDates(year, month);
  const startDate = monthDates[0];
  const endDate = monthDates[monthDates.length - 1];

  // Fetch panchangam data for the visible month
  const { 
    data: panchangamData, 
    loading, 
    isRetrying, 
    error, 
    errorState, 
    retry,
    isOffline,
    offlineState 
  } = usePanchangamRange(startDate, endDate, settings);

  // Initialize location on first load
  useEffect(() => {
    const initializeLocation = async () => {
      try {
        const location = await locationService.getCurrentLocation();
        setSettings(prev => ({
          ...prev,
          location,
          region: location.region
        }));
      } catch (error) {
        console.error('Failed to get initial location:', error);
        // Keep default location (Chennai)
      }
    };

    initializeLocation();
  }, []);

  const handlePrevMonth = () => {
    setCurrentDate(new Date(year, month - 1, 1));
  };

  const handleNextMonth = () => {
    setCurrentDate(new Date(year, month + 1, 1));
  };

  const handleToday = () => {
    setCurrentDate(new Date());
  };

  const handleDateClick = (date: Date) => {
    setSelectedDate(date);
  };

  const handleLocationSelect = (location: any) => {
    setSettings(prev => ({
      ...prev,
      location,
      region: location.region
    }));
  };

  // Get day detail data and loading states
  const dayDetail = useDayDetail({
    selectedDate,
    panchangamData,
    loading,
    error,
    retry,
  });

  return (
    <ErrorBoundary>
      <div className="min-h-screen bg-gradient-to-br from-orange-50 to-yellow-50">
        <div className="container mx-auto px-4 py-6 max-w-7xl">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-4xl md:text-5xl font-bold text-orange-800 mb-2">
            🕉️ Panchangam
          </h1>
          <p className="text-orange-600 text-lg">
            Hindu Calendar & Astronomical Almanac
          </p>
        </div>

        {/* Offline Indicator */}
        <OfflineIndicator isOffline={isOffline} />

        {/* Navigation */}
        <MonthNavigation
          year={year}
          month={month}
          settings={settings}
          onPrevMonth={handlePrevMonth}
          onNextMonth={handleNextMonth}
          onToday={handleToday}
          onLocationClick={() => setShowLocationSelector(true)}
          onSettingsClick={() => setShowSettings(true)}
        />

        {/* Loading State */}
        {loading && !isRetrying && (
          <div className="mb-6">
            <SkeletonCalendar />
          </div>
        )}

        {/* Error State */}
        {error && (
          <div className="mb-6">
            {errorState.isNetworkError ? (
              <NetworkError
                onRetry={retry}
                isRetrying={isRetrying}
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

        {/* Retry Loading Indicator */}
        {isRetrying && (
          <div className="mb-6 text-center">
            <LoadingSpinner
              size="md"
              color="orange"
              message="Retrying..."
            />
          </div>
        )}

        {/* Calendar */}
        {!loading || isRetrying ? (
          <CalendarGrid
            year={year}
            month={month}
            panchangamData={panchangamData}
            settings={settings}
            onDateClick={handleDateClick}
          />
        ) : null}

        {/* Footer */}
        <div className="text-center mt-8 text-gray-600">
          <p className="text-sm">
            Calculated using {settings.calculation_method} method for {settings.location.name}
          </p>
          <p className="text-xs mt-2">
            May the divine blessings guide you through auspicious times 🙏
          </p>
        </div>
      </div>

      {/* Modals */}
      {selectedDate && (
        <DayDetailModal
          date={selectedDate}
          data={dayDetail.data}
          settings={settings}
          isLoading={dayDetail.isLoading}
          error={dayDetail.error}
          onRetry={dayDetail.retry}
          onClose={() => setSelectedDate(null)}
        />
      )}

      {showLocationSelector && (
        <LocationSelector
          currentLocation={settings.location}
          onLocationSelect={handleLocationSelect}
          onClose={() => setShowLocationSelector(false)}
        />
      )}

      {showSettings && (
        <SettingsPanel
          settings={settings}
          onSettingsChange={setSettings}
          onClose={() => setShowSettings(false)}
        />
      )}
        </div>
      </div>
    </ErrorBoundary>
  );
}

export default App;