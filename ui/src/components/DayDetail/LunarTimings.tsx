import React from 'react';
import { PanchangamData, Settings } from '../../types/panchangam';
import { formatTime } from '../../utils/dateHelpers';
import { Moon, Star, Sunrise, Sunset } from 'lucide-react';

interface LunarTimingsProps {
  data: PanchangamData;
  settings: Settings;
}

export const LunarTimings: React.FC<LunarTimingsProps> = ({ data, settings }) => {
  // Extract lunar phase information from events
  const moonPhaseEvent = data.events.find(event => event.event_type === 'MOON_PHASE');
  
  // Parse moon phase details from the event name if available
  const getMoonPhaseDetails = () => {
    if (!moonPhaseEvent) return null;
    
    // Extract phase name and illumination from event name
    // Format: "Moon Phase: Full Moon (95.3% illuminated)"
    const phaseMatch = moonPhaseEvent.name.match(/Moon Phase: (.+?) \((.+?)% illuminated\)/);
    if (phaseMatch) {
      return {
        phaseName: phaseMatch[1],
        illumination: parseFloat(phaseMatch[2])
      };
    }
    
    return {
      phaseName: moonPhaseEvent.name.replace('Moon Phase: ', ''),
      illumination: null
    };
  };

  const moonPhaseDetails = getMoonPhaseDetails();

  return (
    <div className="bg-gradient-to-br from-indigo-50 to-purple-50 rounded-lg p-4 border border-indigo-200">
      <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
        <Moon className="w-5 h-5 mr-2 text-indigo-600" />
        Lunar & Solar Timings
      </h3>
      
      <div className="grid grid-cols-2 gap-4">
        {/* Solar Timings */}
        <div className="space-y-3">
          <h4 className="font-medium text-gray-700 flex items-center text-sm">
            <Sunrise className="w-4 h-4 mr-1 text-yellow-500" />
            Solar Events
          </h4>
          
          <div className="space-y-2">
            <div className="flex justify-between items-center text-sm">
              <span className="text-gray-600">Sunrise</span>
              <span className="font-mono text-gray-800">
                {formatTime(data.sunrise_time, settings.time_format)}
              </span>
            </div>
            <div className="flex justify-between items-center text-sm">
              <span className="text-gray-600">Sunset</span>
              <span className="font-mono text-gray-800">
                {formatTime(data.sunset_time, settings.time_format)}
              </span>
            </div>
          </div>
        </div>

        {/* Lunar Timings */}
        <div className="space-y-3">
          <h4 className="font-medium text-gray-700 flex items-center text-sm">
            <Moon className="w-4 h-4 mr-1 text-blue-400" />
            Lunar Events
          </h4>
          
          <div className="space-y-2">
            {data.moonrise_time && (
              <div className="flex justify-between items-center text-sm">
                <span className="text-gray-600">Moonrise</span>
                <span className="font-mono text-gray-800">
                  {formatTime(data.moonrise_time, settings.time_format)}
                </span>
              </div>
            )}
            {data.moonset_time && (
              <div className="flex justify-between items-center text-sm">
                <span className="text-gray-600">Moonset</span>
                <span className="font-mono text-gray-800">
                  {formatTime(data.moonset_time, settings.time_format)}
                </span>
              </div>
            )}
            {!data.moonrise_time && !data.moonset_time && (
              <div className="text-sm text-gray-500 italic">
                Lunar times not available
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Moon Phase Information */}
      {moonPhaseDetails && (
        <div className="mt-4 pt-4 border-t border-indigo-200">
          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <Star className="w-4 h-4 mr-2 text-purple-500" />
              <span className="font-medium text-gray-700">{moonPhaseDetails.phaseName}</span>
            </div>
            {moonPhaseDetails.illumination !== null && (
              <div className="flex items-center">
                <div className="w-16 bg-gray-200 rounded-full h-2 mr-2">
                  <div 
                    className="bg-purple-500 h-2 rounded-full transition-all duration-300"
                    style={{ width: `${moonPhaseDetails.illumination}%` }}
                  ></div>
                </div>
                <span className="text-sm text-gray-600 font-mono">
                  {moonPhaseDetails.illumination.toFixed(1)}%
                </span>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};