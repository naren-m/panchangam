import React from 'react';
import { X, Calendar, MapPin, Clock } from 'lucide-react';
import { PanchangamData, Settings } from '../../types/panchangam';
import { formatTime, formatTimeRange } from '../../utils/dateHelpers';
import { FiveAngas } from './FiveAngas';
import { MuhurtaTimeline } from './MuhurtaTimeline';
import { EventsList } from './EventsList';

interface DayDetailModalProps {
  date: Date;
  data: PanchangamData | null;
  settings: Settings;
  onClose: () => void;
}

export const DayDetailModal: React.FC<DayDetailModalProps> = ({
  date,
  data,
  settings,
  onClose
}) => {
  if (!data) return null;

  const formatDate = (date: Date) => {
    return date.toLocaleDateString(settings.locale === 'hi' ? 'hi-IN' : 'en-IN', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-hidden">
        {/* Header */}
        <div className="bg-gradient-to-r from-orange-500 to-orange-600 text-white p-6">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-2xl font-bold mb-2">{formatDate(date)}</h2>
              <div className="flex items-center space-x-4 text-orange-100">
                <div className="flex items-center space-x-1">
                  <MapPin className="w-4 h-4" />
                  <span>{settings.location.name}</span>
                </div>
                <div className="flex items-center space-x-1">
                  <Clock className="w-4 h-4" />
                  <span>{settings.calculation_method}</span>
                </div>
              </div>
            </div>
            <button
              onClick={onClose}
              className="p-2 hover:bg-orange-600 rounded-full transition-colors"
            >
              <X className="w-6 h-6" />
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="p-6 overflow-y-auto max-h-[calc(90vh-200px)]">
          <div className="grid md:grid-cols-2 gap-6">
            {/* Left Column */}
            <div className="space-y-6">
              {/* Five Angas */}
              <FiveAngas data={data} settings={settings} />

              {/* Astronomical Times */}
              <div className="bg-blue-50 rounded-lg p-4">
                <h3 className="text-lg font-semibold text-blue-800 mb-3">
                  Astronomical Times
                </h3>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="text-blue-600 font-medium">Sunrise:</span>
                    <div className="text-blue-800">
                      {formatTime(data.sunrise_time, settings.time_format)}
                    </div>
                  </div>
                  <div>
                    <span className="text-blue-600 font-medium">Sunset:</span>
                    <div className="text-blue-800">
                      {formatTime(data.sunset_time, settings.time_format)}
                    </div>
                  </div>
                  {data.moonrise_time && (
                    <div>
                      <span className="text-blue-600 font-medium">Moonrise:</span>
                      <div className="text-blue-800">
                        {formatTime(data.moonrise_time, settings.time_format)}
                      </div>
                    </div>
                  )}
                  {data.moonset_time && (
                    <div>
                      <span className="text-blue-600 font-medium">Moonset:</span>
                      <div className="text-blue-800">
                        {formatTime(data.moonset_time, settings.time_format)}
                      </div>
                    </div>
                  )}
                </div>
              </div>

              {/* Festivals */}
              {data.festivals && data.festivals.length > 0 && (
                <div className="bg-red-50 rounded-lg p-4">
                  <h3 className="text-lg font-semibold text-red-800 mb-3">
                    Festivals & Observances
                  </h3>
                  <div className="space-y-2">
                    {data.festivals.map((festival, index) => (
                      <div key={index} className="flex items-center space-x-2">
                        <div className="w-2 h-2 bg-red-500 rounded-full"></div>
                        <span className="text-red-700 font-medium">{festival}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>

            {/* Right Column */}
            <div className="space-y-6">
              {/* Muhurta Timeline */}
              <MuhurtaTimeline events={data.events} settings={settings} />

              {/* Events List */}
              <EventsList events={data.events} settings={settings} />
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="bg-gray-50 px-6 py-4 border-t">
          <div className="flex items-center justify-between text-sm text-gray-600">
            <div>
              Calculated using {settings.calculation_method} method
            </div>
            <div className="flex items-center space-x-4">
              <button className="text-orange-600 hover:text-orange-700 font-medium">
                Share
              </button>
              <button className="text-orange-600 hover:text-orange-700 font-medium">
                Export
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};