import React from 'react';
import { Event, Settings } from '../../types/panchangam';
import { formatTimeRange } from '../../utils/dateHelpers';
import { Clock, AlertTriangle, CheckCircle, Info } from 'lucide-react';

interface EventsListProps {
  events: Event[];
  settings: Settings;
}

export const EventsList: React.FC<EventsListProps> = ({ events, settings }) => {
  const getEventIcon = (eventType: string, quality: string) => {
    switch (quality) {
      case 'auspicious':
        return <CheckCircle className="w-4 h-4 text-green-600" />;
      case 'inauspicious':
        return <AlertTriangle className="w-4 h-4 text-red-600" />;
      default:
        return <Info className="w-4 h-4 text-yellow-600" />;
    }
  };

  const getEventDescription = (eventType: string) => {
    const descriptions = {
      BRAHMA_MUHURTA: 'Ideal time for meditation and spiritual practices',
      RAHU_KALAM: 'Inauspicious period, avoid starting new activities',
      YAMAGANDAM: 'Inauspicious time period ruled by Yama',
      GULIKA_KALAM: 'Period ruled by Gulika, generally avoided',
      ABHIJIT: 'Most auspicious time, good for all activities',
      MUHURTA: 'Auspicious time period',
      GODHULI: 'Twilight period, sacred time for prayers'
    };
    return descriptions[eventType as keyof typeof descriptions] || 'Special time period';
  };

  return (
    <div className="bg-white rounded-lg p-4 border border-gray-200">
      <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
        <Clock className="w-5 h-5 mr-2" />
        Daily Events & Timings
      </h3>
      <div className="space-y-4">
        {events.map((event, index) => (
          <div key={index} className="flex items-start space-x-3 p-3 rounded-lg bg-gray-50 hover:bg-gray-100 transition-colors">
            <div className="flex-shrink-0 mt-1">
              {getEventIcon(event.event_type, event.quality)}
            </div>
            <div className="flex-1 min-w-0">
              <div className="flex items-center justify-between mb-1">
                <h4 className="font-semibold text-gray-800 truncate">{event.name}</h4>
                <span className="text-sm text-gray-600 ml-2">
                  {formatTimeRange(event.time, settings.time_format)}
                </span>
              </div>
              <p className="text-sm text-gray-600 leading-relaxed">
                {getEventDescription(event.event_type)}
              </p>
              <div className="mt-2">
                <span className={`
                  inline-flex items-center px-2 py-1 rounded-full text-xs font-medium
                  ${event.quality === 'auspicious' ? 'bg-green-100 text-green-800' :
                    event.quality === 'inauspicious' ? 'bg-red-100 text-red-800' :
                    'bg-yellow-100 text-yellow-800'}
                `}>
                  {event.event_type.replace(/_/g, ' ').toLowerCase()}
                </span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};