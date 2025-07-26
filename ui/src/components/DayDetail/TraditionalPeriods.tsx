import React from 'react';
import { Event, Settings } from '../../types/panchangam';
import { formatTime } from '../../utils/dateHelpers';
import { AlertTriangle, CheckCircle, Clock } from 'lucide-react';

interface TraditionalPeriodsProps {
  events: Event[];
  settings: Settings;
}

export const TraditionalPeriods: React.FC<TraditionalPeriodsProps> = ({ events, settings }) => {
  // Filter traditional period events
  const traditionalEvents = events.filter(event => 
    ['RAHU_KALAM', 'YAMAGANDAM', 'GULIKA_KALAM', 'ABHIJIT_MUHURTA'].includes(event.event_type)
  );

  if (traditionalEvents.length === 0) {
    return null;
  }

  const getEventColor = (eventType: string) => {
    switch (eventType) {
      case 'ABHIJIT_MUHURTA':
        return 'bg-green-50 border-green-200 text-green-800';
      case 'RAHU_KALAM':
        return 'bg-red-50 border-red-200 text-red-800';
      case 'YAMAGANDAM':
        return 'bg-orange-50 border-orange-200 text-orange-800';
      case 'GULIKA_KALAM':
        return 'bg-yellow-50 border-yellow-200 text-yellow-800';
      default:
        return 'bg-gray-50 border-gray-200 text-gray-800';
    }
  };

  const getEventIcon = (eventType: string) => {
    switch (eventType) {
      case 'ABHIJIT_MUHURTA':
        return <CheckCircle className="w-4 h-4 text-green-600" />;
      case 'RAHU_KALAM':
      case 'YAMAGANDAM':
      case 'GULIKA_KALAM':
        return <AlertTriangle className="w-4 h-4 text-red-600" />;
      default:
        return <Clock className="w-4 h-4 text-gray-600" />;
    }
  };

  const getEventTitle = (eventType: string) => {
    const titles = {
      RAHU_KALAM: 'Rahu Kalam',
      YAMAGANDAM: 'Yamagandam',
      GULIKA_KALAM: 'Gulika Kalam',
      ABHIJIT_MUHURTA: 'Abhijit Muhurta'
    };
    return titles[eventType as keyof typeof titles] || eventType;
  };

  const getEventAdvice = (eventType: string) => {
    const advice = {
      RAHU_KALAM: 'Avoid starting new ventures, traveling, or important activities',
      YAMAGANDAM: 'Not suitable for auspicious activities or new beginnings',
      GULIKA_KALAM: 'Generally avoided for important work or ceremonies',
      ABHIJIT_MUHURTA: 'Most favorable time for all activities, ceremonies, and new ventures'
    };
    return advice[eventType as keyof typeof advice] || '';
  };

  return (
    <div className="bg-white rounded-lg p-4 border border-gray-200">
      <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
        <Clock className="w-5 h-5 mr-2 text-indigo-600" />
        Traditional Time Periods
      </h3>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {traditionalEvents.map((event, index) => (
          <div 
            key={index} 
            className={`p-3 rounded-lg border ${getEventColor(event.event_type)}`}
          >
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center">
                {getEventIcon(event.event_type)}
                <h4 className="font-semibold ml-2">
                  {getEventTitle(event.event_type)}
                </h4>
              </div>
              <span className="font-mono text-sm">
                {formatTime(event.time, settings.time_format)}
              </span>
            </div>
            
            <p className="text-sm opacity-90">
              {getEventAdvice(event.event_type)}
            </p>
            
            <div className="mt-2">
              <span className={`
                inline-flex items-center px-2 py-1 rounded-full text-xs font-medium
                ${event.quality === 'auspicious' ? 'bg-green-100 text-green-700' :
                  event.quality === 'inauspicious' ? 'bg-red-100 text-red-700' :
                  'bg-gray-100 text-gray-700'}
              `}>
                {event.quality}
              </span>
            </div>
          </div>
        ))}
      </div>
      
      <div className="mt-4 p-3 bg-blue-50 rounded-lg border border-blue-200">
        <p className="text-sm text-blue-800">
          <strong>Note:</strong> Traditional periods are calculated based on sunrise/sunset times and day of the week. 
          Abhijit Muhurta is considered the most auspicious time, while Rahu Kalam, Yamagandam, and Gulika Kalam 
          are traditionally avoided for important activities.
        </p>
      </div>
    </div>
  );
};