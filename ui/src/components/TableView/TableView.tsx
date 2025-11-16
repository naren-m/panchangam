import React, { useMemo } from 'react';
import { PanchangamData, Settings } from '../../types/panchangam';
import { Download, Calendar, Sun, Moon, Star, Clock } from 'lucide-react';

interface TableViewProps {
  year: number;
  month: number;
  panchangamData: Record<string, PanchangamData>;
  settings: Settings;
  onDateClick: (date: Date) => void;
  onExport?: (format: 'csv' | 'json') => void;
}

export const TableView: React.FC<TableViewProps> = ({
  year,
  month,
  panchangamData,
  settings,
  onDateClick,
  onExport
}) => {
  // Sort and prepare data for display
  const sortedData = useMemo(() => {
    const entries = Object.entries(panchangamData);
    return entries.sort(([dateA], [dateB]) => dateA.localeCompare(dateB));
  }, [panchangamData]);

  const handleRowClick = (dateStr: string) => {
    const [year, month, day] = dateStr.split('-').map(Number);
    onDateClick(new Date(year, month - 1, day));
  };

  const getEventQualityColor = (quality: string) => {
    switch (quality) {
      case 'auspicious':
        return 'text-green-600 bg-green-50';
      case 'inauspicious':
        return 'text-red-600 bg-red-50';
      default:
        return 'text-gray-600 bg-gray-50';
    }
  };

  const formatTime = (time: string) => {
    if (!time) return '-';
    if (settings.time_format === '12') {
      const [hours, minutes] = time.split(':').map(Number);
      const period = hours >= 12 ? 'PM' : 'AM';
      const displayHours = hours % 12 || 12;
      return `${displayHours}:${minutes.toString().padStart(2, '0')} ${period}`;
    }
    return time;
  };

  if (sortedData.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow-lg p-8 text-center">
        <Calendar className="w-16 h-16 text-gray-400 mx-auto mb-4" />
        <p className="text-gray-600">No data available for the selected month</p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden border border-orange-200">
      {/* Header with export options */}
      <div className="bg-gradient-to-r from-orange-400 to-orange-500 p-4">
        <div className="flex justify-between items-center">
          <div>
            <h2 className="text-white text-xl font-bold flex items-center gap-2">
              <Calendar className="w-6 h-6" />
              Panchangam Table View
            </h2>
            <p className="text-orange-100 text-sm mt-1">
              {new Date(year, month, 1).toLocaleDateString(settings.locale, {
                month: 'long',
                year: 'numeric'
              })}
            </p>
          </div>
          {onExport && (
            <div className="flex gap-2">
              <button
                onClick={() => onExport('csv')}
                className="bg-white text-orange-600 px-4 py-2 rounded-lg hover:bg-orange-50 transition-colors flex items-center gap-2 text-sm font-medium"
              >
                <Download className="w-4 h-4" />
                CSV
              </button>
              <button
                onClick={() => onExport('json')}
                className="bg-white text-orange-600 px-4 py-2 rounded-lg hover:bg-orange-50 transition-colors flex items-center gap-2 text-sm font-medium"
              >
                <Download className="w-4 h-4" />
                JSON
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Scrollable table container */}
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-orange-100 sticky top-0 z-10">
            <tr>
              <th className="px-4 py-3 text-left text-xs font-semibold text-orange-900 uppercase tracking-wider">
                <div className="flex items-center gap-2">
                  <Calendar className="w-4 h-4" />
                  Date
                </div>
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-orange-900 uppercase tracking-wider">
                Day
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-orange-900 uppercase tracking-wider">
                <div className="flex items-center gap-2">
                  <Moon className="w-4 h-4" />
                  Tithi
                </div>
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-orange-900 uppercase tracking-wider">
                <div className="flex items-center gap-2">
                  <Star className="w-4 h-4" />
                  Nakshatra
                </div>
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-orange-900 uppercase tracking-wider">
                Yoga
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-orange-900 uppercase tracking-wider">
                Karana
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-orange-900 uppercase tracking-wider">
                <div className="flex items-center gap-2">
                  <Sun className="w-4 h-4" />
                  Sunrise
                </div>
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-orange-900 uppercase tracking-wider">
                <div className="flex items-center gap-2">
                  <Sun className="w-4 h-4" />
                  Sunset
                </div>
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-orange-900 uppercase tracking-wider">
                Festivals
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-orange-900 uppercase tracking-wider">
                Auspicious Times
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {sortedData.map(([dateStr, data], index) => {
              const date = new Date(dateStr);
              const isToday = dateStr === new Date().toISOString().split('T')[0];
              const isWeekend = date.getDay() === 0 || date.getDay() === 6;

              const auspiciousEvents = data.events?.filter(e =>
                e.quality === 'auspicious' &&
                ['MUHURTA', 'ABHIJIT_MUHURTA', 'ABHIJIT', 'BRAHMA_MUHURTA'].includes(e.event_type)
              ) || [];

              return (
                <tr
                  key={dateStr}
                  onClick={() => handleRowClick(dateStr)}
                  className={`
                    transition-colors cursor-pointer
                    ${isToday ? 'bg-orange-50 hover:bg-orange-100' : 'hover:bg-gray-50'}
                    ${isWeekend ? 'bg-blue-50 hover:bg-blue-100' : ''}
                    ${index % 2 === 0 && !isToday && !isWeekend ? 'bg-white' : ''}
                  `}
                >
                  <td className="px-4 py-3 whitespace-nowrap">
                    <div className="flex items-center">
                      <div className={`
                        text-sm font-medium
                        ${isToday ? 'text-orange-900 font-bold' : 'text-gray-900'}
                      `}>
                        {date.toLocaleDateString(settings.locale, {
                          day: '2-digit',
                          month: 'short',
                          year: 'numeric'
                        })}
                      </div>
                      {isToday && (
                        <span className="ml-2 px-2 py-1 text-xs font-semibold bg-orange-500 text-white rounded-full">
                          Today
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <span className={`text-sm ${isWeekend ? 'font-semibold text-blue-700' : 'text-gray-900'}`}>
                      {data.vara}
                    </span>
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <span className="text-sm text-gray-900">{data.tithi}</span>
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <span className="text-sm text-gray-900">{data.nakshatra}</span>
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <span className="text-sm text-gray-600">{data.yoga}</span>
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <span className="text-sm text-gray-600">{data.karana}</span>
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <div className="flex items-center gap-1 text-sm text-gray-900">
                      <Clock className="w-3 h-3 text-orange-500" />
                      {formatTime(data.sunrise_time)}
                    </div>
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <div className="flex items-center gap-1 text-sm text-gray-900">
                      <Clock className="w-3 h-3 text-orange-500" />
                      {formatTime(data.sunset_time)}
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    {data.festivals && data.festivals.length > 0 ? (
                      <div className="flex flex-wrap gap-1">
                        {data.festivals.map((festival, i) => (
                          <span
                            key={i}
                            className="inline-flex items-center px-2 py-1 text-xs font-medium bg-purple-100 text-purple-800 rounded-full"
                          >
                            {festival}
                          </span>
                        ))}
                      </div>
                    ) : (
                      <span className="text-sm text-gray-400">-</span>
                    )}
                  </td>
                  <td className="px-4 py-3">
                    {auspiciousEvents.length > 0 ? (
                      <div className="flex flex-wrap gap-1">
                        {auspiciousEvents.slice(0, 2).map((event, i) => (
                          <span
                            key={i}
                            className={`inline-flex items-center px-2 py-1 text-xs font-medium rounded-full ${getEventQualityColor(event.quality)}`}
                            title={`${event.name}: ${formatTime(event.time)}`}
                          >
                            {event.name}
                          </span>
                        ))}
                        {auspiciousEvents.length > 2 && (
                          <span className="inline-flex items-center px-2 py-1 text-xs font-medium bg-gray-100 text-gray-600 rounded-full">
                            +{auspiciousEvents.length - 2} more
                          </span>
                        )}
                      </div>
                    ) : (
                      <span className="text-sm text-gray-400">-</span>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/* Summary footer */}
      <div className="bg-gray-50 px-4 py-3 border-t border-gray-200">
        <div className="flex justify-between items-center text-sm text-gray-600">
          <div>
            Showing {sortedData.length} {sortedData.length === 1 ? 'day' : 'days'}
          </div>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 bg-orange-200 rounded"></div>
              <span>Today</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 bg-blue-100 rounded"></div>
              <span>Weekend</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
