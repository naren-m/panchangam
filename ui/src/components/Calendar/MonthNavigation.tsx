import React from 'react';
import { ChevronLeft, ChevronRight, Calendar, MapPin, Globe } from 'lucide-react';
import { getMonthName } from '../../utils/dateHelpers';
import { Settings } from '../../types/panchangam';

interface MonthNavigationProps {
  year: number;
  month: number;
  settings: Settings;
  onPrevMonth: () => void;
  onNextMonth: () => void;
  onToday: () => void;
  onLocationClick: () => void;
  onSettingsClick: () => void;
  onSkyViewClick?: () => void;
}

export const MonthNavigation: React.FC<MonthNavigationProps> = ({
  year,
  month,
  settings,
  onPrevMonth,
  onNextMonth,
  onToday,
  onLocationClick,
  onSettingsClick,
  onSkyViewClick
}) => {
  const monthName = getMonthName(month, settings.locale);

  return (
    <div className="bg-white rounded-lg shadow-md p-4 mb-6 border border-orange-200">
      <div className="flex items-center justify-between flex-wrap gap-4">
        {/* Month navigation */}
        <div className="flex items-center space-x-4">
          <button
            onClick={onPrevMonth}
            className="p-2 hover:bg-orange-100 rounded-full transition-colors"
            aria-label="Previous month"
          >
            <ChevronLeft className="w-5 h-5 text-orange-600" />
          </button>
          
          <div className="text-center">
            <h1 className="text-2xl md:text-3xl font-bold text-gray-800">
              {monthName} {year}
            </h1>
            <p className="text-sm text-gray-600 mt-1">
              पञ्चाङ्गम् • Panchangam
            </p>
          </div>

          <button
            onClick={onNextMonth}
            className="p-2 hover:bg-orange-100 rounded-full transition-colors"
            aria-label="Next month"
          >
            <ChevronRight className="w-5 h-5 text-orange-600" />
          </button>
        </div>

        {/* Action buttons */}
        <div className="flex items-center space-x-3">
          <button
            onClick={onToday}
            className="flex items-center space-x-2 px-4 py-2 bg-orange-500 text-white rounded-lg hover:bg-orange-600 transition-colors"
          >
            <Calendar className="w-4 h-4" />
            <span className="hidden sm:inline">Today</span>
          </button>

          <button
            onClick={onLocationClick}
            className="flex items-center space-x-2 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
            title={settings.location.name}
          >
            <MapPin className="w-4 h-4" />
            <span className="hidden lg:inline max-w-32 truncate">
              {settings.location.name.split(',')[0]}
            </span>
          </button>

          {onSkyViewClick && (
            <button
              onClick={onSkyViewClick}
              className="flex items-center space-x-2 px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition-colors"
              title="Sky View"
            >
              <Globe className="w-4 h-4" />
              <span className="hidden sm:inline">Sky View</span>
            </button>
          )}

          <button
            onClick={onSettingsClick}
            className="p-2 text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
            aria-label="Settings"
          >
            <div className="w-5 h-5 text-gray-600">⚙️</div>
          </button>
        </div>
      </div>
    </div>
  );
};