import React from 'react';
import { Event, Settings } from '../../types/panchangam';
import { formatTimeRange } from '../../utils/dateHelpers';

interface MuhurtaTimelineProps {
  events: Event[];
  settings: Settings;
}

export const MuhurtaTimeline: React.FC<MuhurtaTimelineProps> = ({ events, settings }) => {
  const getQualityColor = (quality: string) => {
    switch (quality) {
      case 'auspicious':
        return 'bg-green-500';
      case 'inauspicious':
        return 'bg-red-500';
      default:
        return 'bg-yellow-500';
    }
  };

  const getQualityBgColor = (quality: string) => {
    switch (quality) {
      case 'auspicious':
        return 'bg-green-50 border-green-200';
      case 'inauspicious':
        return 'bg-red-50 border-red-200';
      default:
        return 'bg-yellow-50 border-yellow-200';
    }
  };

  const sortedEvents = [...events].sort((a, b) => {
    const timeA = a.time.split('-')[0];
    const timeB = b.time.split('-')[0];
    return timeA.localeCompare(timeB);
  });

  return (
    <div className="bg-white rounded-lg p-4 border border-gray-200">
      <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
        <span className="mr-2">⏰</span>
        Muhurta Timeline
      </h3>
      <div className="space-y-3">
        {sortedEvents.map((event, index) => (
          <div
            key={index}
            className={`p-3 rounded-lg border ${getQualityBgColor(event.quality)} transition-all hover:shadow-sm`}
          >
            <div className="flex items-center space-x-3">
              <div className={`w-3 h-3 rounded-full ${getQualityColor(event.quality)}`}></div>
              <div className="flex-1">
                <div className="font-semibold text-gray-800">{event.name}</div>
                <div className="text-sm text-gray-600">
                  {formatTimeRange(event.time, settings.time_format)}
                </div>
              </div>
              <div className={`
                px-2 py-1 rounded-full text-xs font-medium
                ${event.quality === 'auspicious' ? 'bg-green-100 text-green-800' :
                  event.quality === 'inauspicious' ? 'bg-red-100 text-red-800' :
                  'bg-yellow-100 text-yellow-800'}
              `}>
                {event.quality}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};