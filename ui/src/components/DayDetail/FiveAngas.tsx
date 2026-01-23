import React from 'react';
import { PanchangamData, Settings } from '../../types/panchangam';
import { Moon, Star, Sun, Zap, Calendar } from 'lucide-react';

interface FiveAngasProps {
  data: PanchangamData;
  settings: Settings;
}

export const FiveAngas: React.FC<FiveAngasProps> = ({ data, settings }) => {
  // Extract Vara from events
  const varaEvent = data.events?.find(event => event.event_type === 'VARA');
  const varaName = varaEvent?.name?.replace('Vara: ', '') || 'Not available';

  const angas = [
    {
      name: 'Tithi',
      value: data.tithi || 'Loading...',
      icon: Moon,
      color: 'blue',
      description: 'Lunar day phase'
    },
    {
      name: 'Nakshatra',
      value: data.nakshatra || 'Loading...',
      icon: Star,
      color: 'purple',
      description: 'Lunar mansion'
    },
    {
      name: 'Yoga',
      value: data.yoga || 'Loading...',
      icon: Sun,
      color: 'yellow',
      description: 'Sun-Moon combination'
    },
    {
      name: 'Karana',
      value: data.karana || 'Loading...',
      icon: Zap,
      color: 'green',
      description: 'Half-tithi period'
    },
    {
      name: 'Vara',
      value: varaName,
      icon: Calendar,
      color: 'red',
      description: 'Weekday and ruler'
    }
  ];

  const getColorClasses = (color: string) => {
    const colors = {
      blue: 'bg-blue-50 text-blue-800 border-blue-200',
      purple: 'bg-purple-50 text-purple-800 border-purple-200',
      yellow: 'bg-yellow-50 text-yellow-800 border-yellow-200',
      green: 'bg-green-50 text-green-800 border-green-200',
      red: 'bg-red-50 text-red-800 border-red-200'
    };
    return colors[color as keyof typeof colors] || colors.blue;
  };

  const getIconColorClasses = (color: string) => {
    const colors = {
      blue: 'text-blue-600',
      purple: 'text-purple-600',
      yellow: 'text-yellow-600',
      green: 'text-green-600',
      red: 'text-red-600'
    };
    return colors[color as keyof typeof colors] || colors.blue;
  };

  return (
    <div className="bg-gradient-to-br from-orange-50 to-yellow-50 rounded-lg p-4 border border-orange-200">
      <h3 className="text-lg font-semibold text-orange-800 mb-4 flex items-center">
        <span className="mr-2">🕉️</span>
        The Five Angas (पञ्चाङ्ग)
      </h3>
      <div className="space-y-3">
        {angas.map((anga, index) => {
          const Icon = anga.icon;
          return (
            <div
              key={index}
              className={`p-3 rounded-lg border ${getColorClasses(anga.color)} transition-all hover:shadow-md`}
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <Icon className={`w-5 h-5 ${getIconColorClasses(anga.color)}`} />
                  <div>
                    <div className="font-semibold">{anga.name}</div>
                    <div className="text-xs opacity-75">{anga.description}</div>
                  </div>
                </div>
                <div className="text-right">
                  <div className="font-bold">{anga.value}</div>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};