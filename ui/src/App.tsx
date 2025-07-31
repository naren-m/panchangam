import React, { useState, useEffect, useMemo, useRef } from 'react';
import { CalendarGrid } from './components/Calendar/CalendarGrid';
import { CalendarDisplayManager } from './components/Calendar/CalendarDisplayManager';
import { MonthNavigation } from './components/Calendar/MonthNavigation';
import { DayDetailModal } from './components/DayDetail/DayDetailModal';
import { LocationSelector } from './components/LocationPicker/LocationSelector';
import { SettingsPanel } from './components/Settings/SettingsPanel';
import { SkeletonCalendar, LoadingSpinner } from './components/common/Loading';
import { ApiError, NetworkError, ErrorBoundary } from './components/common/Error';
import { useProgressivePanchangam } from './hooks/useProgressivePanchangam';
import { useDayDetail } from './hooks/useDayDetail';
import { Settings, PanchangamData } from './types/panchangam';
import { getCurrentMonthDates } from './utils/dateHelpers';
import { locationService } from './services/locationService';

function App() {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [showLocationSelector, setShowLocationSelector] = useState(false);
  const [showSettings, setShowSettings] = useState(false);
  const [settingsState, setSettingsState] = useState({
    calculation_method: 'Drik',
    locale: 'en',
    region: 'California',
    time_format: '12',
    location: {
      name: 'Milpitas, California',
      latitude: 37.4323,
      longitude: -121.9066,
      timezone: 'America/Los_Angeles',
      region: 'California'
    }
  });

  // Use ref to break infinite loop completely
  const settingsRef = useRef(null);
  
  // Update ref when settings change
  const settings = useMemo(() => {
    const newSettings = {
      calculation_method: settingsState.calculation_method,
      locale: settingsState.locale,
      region: settingsState.region,
      time_format: settingsState.time_format,
      location: {
        name: settingsState.location.name,
        latitude: settingsState.location.latitude,
        longitude: settingsState.location.longitude,
        timezone: settingsState.location.timezone,
        region: settingsState.location.region
      }
    };
    
    // Only update if actually different
    if (!settingsRef.current || JSON.stringify(settingsRef.current) !== JSON.stringify(newSettings)) {
      console.log('📝 Settings actually changed, updating...');
      settingsRef.current = newSettings;
    }
    
    return settingsRef.current;
  }, [
    settingsState.calculation_method,
    settingsState.locale,
    settingsState.region,
    settingsState.time_format,
    settingsState.location.name,
    settingsState.location.latitude,
    settingsState.location.longitude,
    settingsState.location.timezone,
    settingsState.location.region
  ]);

  const year = currentDate.getFullYear();
  const month = currentDate.getMonth();

  // Get the date range for the current month view
  const monthDates = getCurrentMonthDates(year, month);
  const startDate = monthDates[0];
  const endDate = monthDates[monthDates.length - 1];

  // Fetch panchangam data progressively for the visible month
  const { 
    data: panchangamData, 
    loading, 
    isProgressiveLoading,
    progress,
    todayLoaded,
    loadedCount,
    totalCount,
    error, 
    errorState, 
    retry 
  } = useProgressivePanchangam(startDate, endDate, settings);

  // Note: Removed automatic location initialization to prevent infinite re-renders
  // Default location is already set to Milpitas, CA

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
    setSettingsState(prev => ({
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

        {/* Calendar Display Logic - Ensures only ONE calendar renders at a time */}
        <CalendarDisplayManager
          loading={loading}
          hasData={Object.keys(panchangamData).length > 0}
          error={error}
          errorState={errorState}
          isProgressiveLoading={isProgressiveLoading}
          progress={progress}
          loadedCount={loadedCount}
          totalCount={totalCount}
          retry={retry}
          calendarProps={{
            year,
            month,
            panchangamData,
            settings,
            onDateClick: handleDateClick
          }}
        />

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
          onSettingsChange={setSettingsState}
          onClose={() => setShowSettings(false)}
        />
      )}
      </div>
    </ErrorBoundary>
  );
}

export default App;