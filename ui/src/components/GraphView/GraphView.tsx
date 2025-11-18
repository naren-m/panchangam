import React, { useMemo } from 'react';
import { PanchangamData, Settings } from '../../types/panchangam';
import { Calendar, Download, TrendingUp, Moon, Sun, Star, BarChart3 } from 'lucide-react';

interface GraphViewProps {
  year: number;
  month: number;
  panchangamData: Record<string, PanchangamData>;
  settings: Settings;
  onDateClick: (date: Date) => void;
  onExport?: (format: 'csv' | 'json') => void;
}

// Helper function to parse time string to minutes
const parseTimeToMinutes = (time: string): number => {
  const [hours, minutes] = time.split(':').map(Number);
  return hours * 60 + minutes;
};

export const GraphView: React.FC<GraphViewProps> = ({
  year,
  month,
  panchangamData,
  settings,
  onDateClick,
  onExport
}) => {
  // Process data for visualizations
  const chartData = useMemo(() => {
    const entries = Object.entries(panchangamData).sort(([a], [b]) => a.localeCompare(b));

    // Tithi distribution
    const tithiCount: Record<string, number> = {};
    // Nakshatra distribution
    const nakshatraCount: Record<string, number> = {};
    // Daily auspicious count
    const dailyAuspicious: Array<{ date: string, count: number }> = [];
    // Festival days
    const festivalDays: Array<{ date: string, festivals: string[] }> = [];
    // Sunrise/sunset times
    const sunTimes: Array<{ date: string, sunrise: number, sunset: number }> = [];

    entries.forEach(([dateStr, data]) => {
      // Count tithis
      tithiCount[data.tithi] = (tithiCount[data.tithi] || 0) + 1;

      // Count nakshatras
      nakshatraCount[data.nakshatra] = (nakshatraCount[data.nakshatra] || 0) + 1;

      // Count auspicious events
      const auspiciousCount = data.events?.filter(e => e.quality === 'auspicious').length || 0;
      dailyAuspicious.push({ date: dateStr, count: auspiciousCount });

      // Collect festival days
      if (data.festivals && data.festivals.length > 0) {
        festivalDays.push({ date: dateStr, festivals: data.festivals });
      }

      // Parse sun times
      if (data.sunrise_time && data.sunset_time) {
        const sunriseMinutes = parseTimeToMinutes(data.sunrise_time);
        const sunsetMinutes = parseTimeToMinutes(data.sunset_time);
        sunTimes.push({ date: dateStr, sunrise: sunriseMinutes, sunset: sunsetMinutes });
      }
    });

    // Sort and get top items
    const topTithis = Object.entries(tithiCount)
      .sort(([, a], [, b]) => b - a)
      .slice(0, 10);

    const topNakshatras = Object.entries(nakshatraCount)
      .sort(([, a], [, b]) => b - a)
      .slice(0, 10);

    return {
      tithiCount: topTithis,
      nakshatraCount: topNakshatras,
      dailyAuspicious,
      festivalDays,
      sunTimes,
      totalDays: entries.length
    };
  }, [panchangamData]);

  const formatMinutesToTime = (minutes: number): string => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    if (settings.time_format === '12') {
      const period = hours >= 12 ? 'PM' : 'AM';
      const displayHours = hours % 12 || 12;
      return `${displayHours}:${mins.toString().padStart(2, '0')} ${period}`;
    }
    return `${hours.toString().padStart(2, '0')}:${mins.toString().padStart(2, '0')}`;
  };

  if (chartData.totalDays === 0) {
    return (
      <div className="bg-white rounded-lg shadow-lg p-8 text-center">
        <BarChart3 className="w-16 h-16 text-gray-400 mx-auto mb-4" />
        <p className="text-gray-600">No data available for visualization</p>
      </div>
    );
  }

  const maxTithiCount = Math.max(...chartData.tithiCount.map(([, count]) => count));
  const maxNakshatraCount = Math.max(...chartData.nakshatraCount.map(([, count]) => count));
  const maxAuspicious = Math.max(...chartData.dailyAuspicious.map(d => d.count));

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white rounded-lg shadow-lg overflow-hidden border border-orange-200">
        <div className="bg-gradient-to-r from-orange-400 to-orange-500 p-4">
          <div className="flex justify-between items-center">
            <div>
              <h2 className="text-white text-xl font-bold flex items-center gap-2">
                <BarChart3 className="w-6 h-6" />
                Panchangam Analytics
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
      </div>

      {/* Tithi Distribution Chart */}
      <div className="bg-white rounded-lg shadow-lg overflow-hidden border border-orange-200 p-6">
        <h3 className="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
          <Moon className="w-5 h-5 text-orange-500" />
          Tithi Distribution
        </h3>
        <div className="space-y-3">
          {chartData.tithiCount.map(([tithi, count]) => (
            <div key={tithi} className="flex items-center gap-3">
              <div className="w-32 text-sm font-medium text-gray-700 truncate" title={tithi}>
                {tithi}
              </div>
              <div className="flex-1">
                <div className="bg-gray-200 rounded-full h-8 overflow-hidden">
                  <div
                    className="bg-gradient-to-r from-orange-400 to-orange-500 h-full flex items-center justify-end pr-3 transition-all duration-500"
                    style={{ width: `${(count / maxTithiCount) * 100}%` }}
                  >
                    <span className="text-white text-xs font-semibold">
                      {count} {count === 1 ? 'day' : 'days'}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Nakshatra Distribution Chart */}
      <div className="bg-white rounded-lg shadow-lg overflow-hidden border border-orange-200 p-6">
        <h3 className="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
          <Star className="w-5 h-5 text-orange-500" />
          Nakshatra Distribution
        </h3>
        <div className="space-y-3">
          {chartData.nakshatraCount.map(([nakshatra, count]) => (
            <div key={nakshatra} className="flex items-center gap-3">
              <div className="w-32 text-sm font-medium text-gray-700 truncate" title={nakshatra}>
                {nakshatra}
              </div>
              <div className="flex-1">
                <div className="bg-gray-200 rounded-full h-8 overflow-hidden">
                  <div
                    className="bg-gradient-to-r from-blue-400 to-blue-500 h-full flex items-center justify-end pr-3 transition-all duration-500"
                    style={{ width: `${(count / maxNakshatraCount) * 100}%` }}
                  >
                    <span className="text-white text-xs font-semibold">
                      {count} {count === 1 ? 'day' : 'days'}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Auspicious Events Timeline */}
      <div className="bg-white rounded-lg shadow-lg overflow-hidden border border-orange-200 p-6">
        <h3 className="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
          <TrendingUp className="w-5 h-5 text-orange-500" />
          Daily Auspicious Events
        </h3>
        <div className="h-64 flex items-end gap-1 overflow-x-auto pb-2">
          {chartData.dailyAuspicious.map(({ date, count }) => {
            const dateObj = new Date(date);
            const day = dateObj.getDate();
            const isToday = date === new Date().toISOString().split('T')[0];
            const height = maxAuspicious > 0 ? (count / maxAuspicious) * 100 : 0;

            return (
              <div
                key={date}
                className="flex-1 min-w-[30px] flex flex-col items-center gap-1 cursor-pointer group"
                onClick={() => {
                  const [y, m, d] = date.split('-').map(Number);
                  onDateClick(new Date(y, m - 1, d));
                }}
              >
                <div className="relative flex-1 w-full flex items-end justify-center">
                  <div
                    className={`
                      w-full transition-all duration-300
                      ${isToday ? 'bg-orange-500' : 'bg-green-400 group-hover:bg-green-500'}
                      rounded-t
                    `}
                    style={{ height: `${height}%` }}
                    title={`${count} auspicious events`}
                  >
                    {count > 0 && (
                      <div className="text-white text-xs font-semibold text-center pt-1">
                        {count}
                      </div>
                    )}
                  </div>
                </div>
                <div className={`
                  text-xs font-medium
                  ${isToday ? 'text-orange-600 font-bold' : 'text-gray-600'}
                `}>
                  {day}
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Sunrise & Sunset Trend */}
      {chartData.sunTimes.length > 0 && (
        <div className="bg-white rounded-lg shadow-lg overflow-hidden border border-orange-200 p-6">
          <h3 className="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
            <Sun className="w-5 h-5 text-orange-500" />
            Sunrise & Sunset Times
          </h3>
          <div className="h-64 relative">
            <svg className="w-full h-full" viewBox="0 0 800 200" preserveAspectRatio="none">
              {/* Grid lines */}
              {[0, 25, 50, 75, 100].map(percent => (
                <line
                  key={percent}
                  x1="0"
                  y1={200 - (percent * 2)}
                  x2="800"
                  y2={200 - (percent * 2)}
                  stroke="#e5e7eb"
                  strokeWidth="1"
                />
              ))}

              {/* Sunrise line */}
              <polyline
                points={chartData.sunTimes.map((d, i) => {
                  const x = (i / (chartData.sunTimes.length - 1)) * 800;
                  const y = 200 - ((d.sunrise - 240) / 480) * 200; // 4am to 12pm range
                  return `${x},${y}`;
                }).join(' ')}
                fill="none"
                stroke="#f59e0b"
                strokeWidth="3"
                strokeLinecap="round"
                strokeLinejoin="round"
              />

              {/* Sunset line */}
              <polyline
                points={chartData.sunTimes.map((d, i) => {
                  const x = (i / (chartData.sunTimes.length - 1)) * 800;
                  const y = 200 - ((d.sunset - 960) / 480) * 200; // 4pm to 12am range
                  return `${x},${y}`;
                }).join(' ')}
                fill="none"
                stroke="#ef4444"
                strokeWidth="3"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>

            {/* Legend */}
            <div className="mt-4 flex justify-center gap-6">
              <div className="flex items-center gap-2">
                <div className="w-4 h-1 bg-amber-500 rounded"></div>
                <span className="text-sm text-gray-600">Sunrise</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-4 h-1 bg-red-500 rounded"></div>
                <span className="text-sm text-gray-600">Sunset</span>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Festivals Calendar */}
      {chartData.festivalDays.length > 0 && (
        <div className="bg-white rounded-lg shadow-lg overflow-hidden border border-orange-200 p-6">
          <h3 className="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
            <Calendar className="w-5 h-5 text-orange-500" />
            Festival Days ({chartData.festivalDays.length})
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {chartData.festivalDays.map(({ date, festivals }) => {
              const dateObj = new Date(date);
              const isToday = date === new Date().toISOString().split('T')[0];

              return (
                <div
                  key={date}
                  onClick={() => {
                    const [y, m, d] = date.split('-').map(Number);
                    onDateClick(new Date(y, m - 1, d));
                  }}
                  className={`
                    p-4 rounded-lg border-2 cursor-pointer transition-all
                    ${isToday
                      ? 'border-orange-500 bg-orange-50 hover:bg-orange-100'
                      : 'border-purple-200 bg-purple-50 hover:bg-purple-100'
                    }
                  `}
                >
                  <div className={`text-sm font-semibold mb-2 ${isToday ? 'text-orange-700' : 'text-purple-700'}`}>
                    {dateObj.toLocaleDateString(settings.locale, {
                      day: 'numeric',
                      month: 'short',
                      weekday: 'short'
                    })}
                    {isToday && (
                      <span className="ml-2 px-2 py-0.5 text-xs bg-orange-500 text-white rounded-full">
                        Today
                      </span>
                    )}
                  </div>
                  <div className="space-y-1">
                    {festivals.map((festival, i) => (
                      <div key={i} className="text-sm text-gray-700 font-medium">
                        • {festival}
                      </div>
                    ))}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Summary Statistics */}
      <div className="bg-white rounded-lg shadow-lg overflow-hidden border border-orange-200 p-6">
        <h3 className="text-lg font-bold text-gray-800 mb-4">Month Summary</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="text-center p-4 bg-orange-50 rounded-lg">
            <div className="text-3xl font-bold text-orange-600">{chartData.totalDays}</div>
            <div className="text-sm text-gray-600 mt-1">Total Days</div>
          </div>
          <div className="text-center p-4 bg-purple-50 rounded-lg">
            <div className="text-3xl font-bold text-purple-600">{chartData.festivalDays.length}</div>
            <div className="text-sm text-gray-600 mt-1">Festivals</div>
          </div>
          <div className="text-center p-4 bg-blue-50 rounded-lg">
            <div className="text-3xl font-bold text-blue-600">{chartData.tithiCount.length}</div>
            <div className="text-sm text-gray-600 mt-1">Unique Tithis</div>
          </div>
          <div className="text-center p-4 bg-green-50 rounded-lg">
            <div className="text-3xl font-bold text-green-600">
              {chartData.dailyAuspicious.reduce((sum, d) => sum + d.count, 0)}
            </div>
            <div className="text-sm text-gray-600 mt-1">Auspicious Events</div>
          </div>
        </div>
      </div>
    </div>
  );
};
