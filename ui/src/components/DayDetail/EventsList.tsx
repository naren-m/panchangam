import React from 'react';
import { Event, Settings } from '../../types/panchangam';
import { formatTimeRange } from '../../utils/dateHelpers';
import { Clock, AlertTriangle, CheckCircle, Info, Sun, Moon, Star, Sunrise, Sunset, Calendar } from 'lucide-react';

interface EventsListProps {
  events: Event[];
  settings: Settings;
}

export const EventsList: React.FC<EventsListProps> = ({ events, settings }) => {
  const getEventIcon = (eventType: string, quality: string) => {
    // Specific icons for different event types
    switch (eventType) {
      case 'SUNRISE':
        return <Sunrise className="w-4 h-4 text-yellow-500" />;
      case 'SUNSET':
        return <Sunset className="w-4 h-4 text-orange-500" />;
      case 'MOONRISE':
      case 'MOONSET':
        return <Moon className="w-4 h-4 text-blue-400" />;
      case 'MOON_PHASE':
        return <Star className="w-4 h-4 text-purple-400" />;
      case 'RAHU_KALAM':
      case 'YAMAGANDAM':
      case 'GULIKA_KALAM':
        return <AlertTriangle className="w-4 h-4 text-red-600" />;
      case 'ABHIJIT_MUHURTA':
      case 'BRAHMA_MUHURTA':
        return <CheckCircle className="w-4 h-4 text-green-600" />;
      case 'FESTIVAL':
        return <Calendar className="w-4 h-4 text-purple-600" />;
      default:
        // Fall back to quality-based icons
        switch (quality) {
          case 'auspicious':
            return <CheckCircle className="w-4 h-4 text-green-600" />;
          case 'inauspicious':
            return <AlertTriangle className="w-4 h-4 text-red-600" />;
          default:
            return <Info className="w-4 h-4 text-blue-600" />;
        }
    }
  };

  const getEventDescription = (eventType: string) => {
    const descriptions = {
      // Solar events
      SUNRISE: 'Beginning of the day, time for prayers and new beginnings',
      SUNSET: 'End of the day, time for reflection and evening prayers',
      
      // Lunar events
      MOONRISE: 'Moon rises above the horizon, influences tides and emotions',
      MOONSET: 'Moon sets below the horizon, time for rest and introspection',
      MOON_PHASE: 'Current lunar phase, affects spiritual and material activities',
      
      // Traditional inauspicious periods
      RAHU_KALAM: 'Inauspicious period ruled by Rahu, avoid starting new activities',
      YAMAGANDAM: 'Inauspicious time period ruled by Yama, Lord of Death',
      GULIKA_KALAM: 'Period ruled by Gulika (son of Saturn), generally avoided for new ventures',
      
      // Auspicious periods
      ABHIJIT_MUHURTA: 'Most auspicious period of the day, excellent for all activities',
      BRAHMA_MUHURTA: 'Time of Brahma, ideal for meditation and spiritual practices',
      
      // Panchangam elements
      TITHI: 'Lunar day, influences emotional and spiritual activities',
      NAKSHATRA: 'Lunar mansion, affects personal characteristics and timing',
      YOGA: 'Auspicious combination of Sun and Moon positions',
      KARANA: 'Half of a Tithi, influences daily activities',
      VARA: 'Day of the week, ruled by specific planetary energy',
      
      // General
      MUHURTA: 'Auspicious time period for specific activities',
      GODHULI: 'Twilight period, sacred time for prayers and rituals',
      
      // Festivals
      FESTIVAL: 'Traditional festival or observance day, significant in Hindu calendar'
    };
    return descriptions[eventType as keyof typeof descriptions] || 'Special time period with traditional significance';
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