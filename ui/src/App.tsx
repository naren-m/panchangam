import React, { useState, useEffect } from 'react';
import { CalendarGrid } from './components/Calendar/CalendarGrid';
import { MonthNavigation } from './components/Calendar/MonthNavigation';
import { DayDetailModal } from './components/DayDetail/DayDetailModal';
import { LocationSelector } from './components/LocationPicker/LocationSelector';
import { SettingsPanel } from './components/Settings/SettingsPanel';
import { usePanchangamRange } from './hooks/usePanchangam';
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
  const { data: panchangamData, loading, error } = usePanchangamRange(startDate, endDate, settings);

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

  const getSelectedDateData = (): PanchangamData | null => {
    if (!selectedDate) return null;
    const dateStr = selectedDate.toISOString().split('T')[0];
    return panchangamData[dateStr] || null;
  };

  return (
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

        {/* Loading State */}
        {loading && (
          <div className="text-center py-12">
            <div className="inline-flex items-center space-x-2">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
              <span className="text-orange-700 font-medium">Loading Panchangam data...</span>
            </div>
          </div>
        )}

        {/* Error State */}
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
            <div className="flex items-center space-x-2">
              <div className="text-red-600">⚠️</div>
              <div>
                <h3 className="font-semibold text-red-800">Error Loading Data</h3>
                <p className="text-red-700 text-sm">{error}</p>
              </div>
            </div>
          </div>
        )}

        {/* Calendar */}
        <CalendarGrid
          year={year}
          month={month}
          panchangamData={panchangamData}
          settings={settings}
          onDateClick={handleDateClick}
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
          data={getSelectedDateData()}
          settings={settings}
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
  );
}

export default App;