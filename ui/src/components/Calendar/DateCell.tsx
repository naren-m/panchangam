import React from 'react';
import { PanchangamData, Settings } from '../../types/panchangam';
import { formatTime, isToday } from '../../utils/dateHelpers';
import { Sunrise, Sunset, Star, Moon } from 'lucide-react';

interface DateCellProps {
  date: Date;
  data?: PanchangamData;
  isCurrentMonth: boolean;
  settings: Settings;
  onClick: () => void;
}

export const DateCell: React.FC<DateCellProps> = ({
  date,
  data,
  isCurrentMonth,
  settings,
  onClick
}) => {
  const today = isToday(date);
  const hasFestival = data?.festivals && data.festivals.length > 0;

  return (
    <div
      className={`
        relative min-h-[120px] md:min-h-[140px] p-2 border-r border-b border-orange-100 cursor-pointer
        transition-all duration-200 hover:bg-orange-50 hover:shadow-md
        ${!isCurrentMonth ? 'bg-gray-50 text-gray-400' : 'bg-cream-50'}
        ${today ? 'bg-gradient-to-br from-orange-100 to-yellow-100 border-2 border-orange-400' : ''}
      `}
      onClick={onClick}
    >
      {/* Date number */}
      <div className="flex items-center justify-between mb-1">
        <span className={`
          text-lg font-bold
          ${today ? 'text-orange-600' : isCurrentMonth ? 'text-gray-800' : 'text-gray-400'}
        `}>
          {date.getDate()}
        </span>
        {hasFestival && (
          <div className="w-2 h-2 bg-red-500 rounded-full animate-pulse"></div>
        )}
      </div>

      {/* Panchangam data */}
      {data && isCurrentMonth && (
        <div className="space-y-1 text-xs">
          {/* Tithi */}
          <div className="flex items-center space-x-1">
            <Moon className="w-3 h-3 text-blue-600" />
            <span className="text-blue-700 font-medium truncate">
              {data.tithi.split(' ')[1] || data.tithi}
            </span>
          </div>

          {/* Nakshatra */}
          <div className="flex items-center space-x-1">
            <Star className="w-3 h-3 text-purple-600" />
            <span className="text-purple-700 truncate">
              {data.nakshatra}
            </span>
          </div>

          {/* Sun times */}
          <div className="flex items-center justify-between text-xs">
            <div className="flex items-center space-x-1">
              <Sunrise className="w-3 h-3 text-yellow-600" />
              <span className="text-yellow-700">
                {formatTime(data.sunrise_time, settings.time_format)}
              </span>
            </div>
            <div className="flex items-center space-x-1">
              <Sunset className="w-3 h-3 text-orange-600" />
              <span className="text-orange-700">
                {formatTime(data.sunset_time, settings.time_format)}
              </span>
            </div>
          </div>

          {/* Festival indicator */}
          {hasFestival && (
            <div className="text-xs text-red-600 font-medium truncate">
              {data.festivals![0]}
            </div>
          )}

          {/* Auspicious time indicator */}
          {data.events.some(e => e.quality === 'auspicious') && (
            <div className="absolute bottom-1 right-1">
              <div className="w-2 h-2 bg-green-500 rounded-full"></div>
            </div>
          )}
        </div>
      )}

      {/* Loading state */}
      {!data && isCurrentMonth && (
        <div className="flex items-center justify-center h-full">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-orange-500"></div>
        </div>
      )}
    </div>
  );
};