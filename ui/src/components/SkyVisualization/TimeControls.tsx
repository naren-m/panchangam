import React, { useState } from 'react';
import { Play, Pause, RotateCcw, Calendar, Clock, Bookmark, ChevronDown } from 'lucide-react';
import type { TimeConfig } from '../../types/skyVisualization';

export interface HistoricalEvent {
  id: string;
  name: string;
  date: Date;
  description: string;
  category: 'eclipse' | 'festival' | 'astronomical' | 'historical';
}

interface TimeControlsProps {
  timeConfig: TimeConfig;
  onTimeConfigChange: (config: Partial<TimeConfig>) => void;
  onJumpToEvent?: (event: HistoricalEvent) => void;
  className?: string;
}

// Predefined historical astronomical events
const HISTORICAL_EVENTS: HistoricalEvent[] = [
  {
    id: 'summer_solstice_2024',
    name: 'Summer Solstice 2024',
    date: new Date('2024-06-20T20:51:00Z'),
    description: 'Longest day of the year in Northern Hemisphere',
    category: 'astronomical',
  },
  {
    id: 'winter_solstice_2024',
    name: 'Winter Solstice 2024',
    date: new Date('2024-12-21T09:20:00Z'),
    description: 'Shortest day of the year in Northern Hemisphere',
    category: 'astronomical',
  },
  {
    id: 'spring_equinox_2024',
    name: 'Vernal Equinox 2024',
    date: new Date('2024-03-20T03:06:00Z'),
    description: 'Day and night are approximately equal',
    category: 'astronomical',
  },
  {
    id: 'autumn_equinox_2024',
    name: 'Autumnal Equinox 2024',
    date: new Date('2024-09-22T12:44:00Z'),
    description: 'Day and night are approximately equal',
    category: 'astronomical',
  },
  {
    id: 'diwali_2024',
    name: 'Diwali 2024 (Amavasya)',
    date: new Date('2024-11-01T00:00:00Z'),
    description: 'Festival of Lights - New Moon',
    category: 'festival',
  },
  {
    id: 'solar_eclipse_2024',
    name: 'Total Solar Eclipse 2024',
    date: new Date('2024-04-08T18:18:00Z'),
    description: 'Total solar eclipse visible from North America',
    category: 'eclipse',
  },
  {
    id: 'lunar_eclipse_2024_march',
    name: 'Penumbral Lunar Eclipse 2024',
    date: new Date('2024-03-25T07:00:00Z'),
    description: 'Penumbral lunar eclipse',
    category: 'eclipse',
  },
  {
    id: 'perseid_meteor_2024',
    name: 'Perseid Meteor Shower Peak 2024',
    date: new Date('2024-08-12T22:00:00Z'),
    description: 'Peak of the annual Perseid meteor shower',
    category: 'astronomical',
  },
  {
    id: 'great_conjunction_2020',
    name: 'Great Conjunction 2020',
    date: new Date('2020-12-21T18:30:00Z'),
    description: 'Jupiter and Saturn closest approach',
    category: 'historical',
  },
];

// Speed presets in minutes per second
const SPEED_PRESETS = [
  { value: 0, label: 'Stopped' },
  { value: 1, label: '1 min/sec' },
  { value: 60, label: '1 hour/sec' },
  { value: 360, label: '6 hours/sec' },
  { value: 1440, label: '1 day/sec' },
  { value: 10080, label: '1 week/sec' },
  { value: 43200, label: '1 month/sec' },
];

export const TimeControls: React.FC<TimeControlsProps> = ({
  timeConfig,
  onTimeConfigChange,
  onJumpToEvent,
  className = '',
}) => {
  const [showDatePicker, setShowDatePicker] = useState(false);
  const [showEventsList, setShowEventsList] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');

  const handlePlayPause = () => {
    onTimeConfigChange({ paused: !timeConfig.paused });
  };

  const handleReset = () => {
    onTimeConfigChange({
      date: new Date(),
      speed: 1,
      paused: false,
    });
  };

  const handleSpeedChange = (speed: number) => {
    onTimeConfigChange({ speed });
  };

  const handleDateChange = (dateString: string) => {
    const newDate = new Date(dateString);
    if (!isNaN(newDate.getTime())) {
      const currentDate = timeConfig.date;
      // Preserve time of day
      newDate.setHours(currentDate.getHours());
      newDate.setMinutes(currentDate.getMinutes());
      newDate.setSeconds(currentDate.getSeconds());
      onTimeConfigChange({ date: newDate });
    }
  };

  const handleTimeChange = (timeString: string) => {
    const [hours, minutes] = timeString.split(':').map(Number);
    const newDate = new Date(timeConfig.date);
    newDate.setHours(hours);
    newDate.setMinutes(minutes);
    onTimeConfigChange({ date: newDate });
  };

  const handleEventSelect = (event: HistoricalEvent) => {
    onTimeConfigChange({ date: new Date(event.date), paused: true });
    setShowEventsList(false);
    if (onJumpToEvent) {
      onJumpToEvent(event);
    }
  };

  const filteredEvents = selectedCategory === 'all'
    ? HISTORICAL_EVENTS
    : HISTORICAL_EVENTS.filter(e => e.category === selectedCategory);

  const formatDate = (date: Date) => {
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: true,
    });
  };

  return (
    <div className={`bg-gray-800 bg-opacity-90 rounded-lg p-4 ${className}`}>
      <div className="space-y-4">
        {/* Current Time Display */}
        <div className="text-white text-center">
          <div className="text-lg font-semibold">{formatDate(timeConfig.date)}</div>
          <div className="text-sm text-gray-300">{formatTime(timeConfig.date)}</div>
        </div>

        {/* Play/Pause and Reset Controls */}
        <div className="flex items-center justify-center gap-2">
          <button
            onClick={handlePlayPause}
            className={`p-2 rounded-lg ${
              timeConfig.paused
                ? 'bg-blue-600 hover:bg-blue-700'
                : 'bg-red-600 hover:bg-red-700'
            } text-white transition-colors`}
            title={timeConfig.paused ? 'Play' : 'Pause'}
          >
            {timeConfig.paused ? <Play size={20} /> : <Pause size={20} />}
          </button>

          <button
            onClick={handleReset}
            className="p-2 rounded-lg bg-gray-700 hover:bg-gray-600 text-white transition-colors"
            title="Reset to current time"
          >
            <RotateCcw size={20} />
          </button>

          <button
            onClick={() => setShowDatePicker(!showDatePicker)}
            className="p-2 rounded-lg bg-gray-700 hover:bg-gray-600 text-white transition-colors"
            title="Pick date and time"
          >
            <Calendar size={20} />
          </button>

          <button
            onClick={() => setShowEventsList(!showEventsList)}
            className="p-2 rounded-lg bg-gray-700 hover:bg-gray-600 text-white transition-colors"
            title="Historical events"
          >
            <Bookmark size={20} />
          </button>
        </div>

        {/* Speed Control */}
        <div className="space-y-2">
          <label className="text-sm text-gray-300 block">Playback Speed</label>
          <select
            value={timeConfig.speed}
            onChange={(e) => handleSpeedChange(Number(e.target.value))}
            className="w-full bg-gray-700 text-white rounded px-3 py-2 text-sm"
          >
            {SPEED_PRESETS.map((preset) => (
              <option key={preset.value} value={preset.value}>
                {preset.label}
              </option>
            ))}
          </select>
        </div>

        {/* Date/Time Picker */}
        {showDatePicker && (
          <div className="space-y-2 border-t border-gray-700 pt-4">
            <div>
              <label className="text-xs text-gray-400 block mb-1">Date</label>
              <input
                type="date"
                value={timeConfig.date.toISOString().split('T')[0]}
                onChange={(e) => handleDateChange(e.target.value)}
                className="w-full bg-gray-700 text-white rounded px-2 py-1 text-sm"
              />
            </div>
            <div>
              <label className="text-xs text-gray-400 block mb-1">Time</label>
              <input
                type="time"
                value={`${String(timeConfig.date.getHours()).padStart(2, '0')}:${String(
                  timeConfig.date.getMinutes()
                ).padStart(2, '0')}`}
                onChange={(e) => handleTimeChange(e.target.value)}
                className="w-full bg-gray-700 text-white rounded px-2 py-1 text-sm"
              />
            </div>
          </div>
        )}

        {/* Historical Events List */}
        {showEventsList && (
          <div className="border-t border-gray-700 pt-4 space-y-2">
            <div className="flex items-center justify-between mb-2">
              <label className="text-sm text-gray-300">Historical Events</label>
              <select
                value={selectedCategory}
                onChange={(e) => setSelectedCategory(e.target.value)}
                className="bg-gray-700 text-white rounded px-2 py-1 text-xs"
              >
                <option value="all">All</option>
                <option value="astronomical">Astronomical</option>
                <option value="eclipse">Eclipses</option>
                <option value="festival">Festivals</option>
                <option value="historical">Historical</option>
              </select>
            </div>

            <div className="max-h-60 overflow-y-auto space-y-1">
              {filteredEvents.map((event) => (
                <button
                  key={event.id}
                  onClick={() => handleEventSelect(event)}
                  className="w-full text-left p-2 bg-gray-700 hover:bg-gray-600 rounded text-sm text-white transition-colors"
                >
                  <div className="font-medium">{event.name}</div>
                  <div className="text-xs text-gray-400">
                    {formatDate(event.date)} - {event.description}
                  </div>
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Time Range Slider (for future enhancement) */}
        <div className="space-y-1">
          <label className="text-xs text-gray-400 block">Quick Navigation</label>
          <div className="flex gap-1">
            <button
              onClick={() => {
                const newDate = new Date(timeConfig.date);
                newDate.setHours(newDate.getHours() - 6);
                onTimeConfigChange({ date: newDate });
              }}
              className="flex-1 bg-gray-700 hover:bg-gray-600 text-white text-xs py-1 rounded"
            >
              -6h
            </button>
            <button
              onClick={() => {
                const newDate = new Date(timeConfig.date);
                newDate.setHours(newDate.getHours() - 1);
                onTimeConfigChange({ date: newDate });
              }}
              className="flex-1 bg-gray-700 hover:bg-gray-600 text-white text-xs py-1 rounded"
            >
              -1h
            </button>
            <button
              onClick={() => {
                const newDate = new Date(timeConfig.date);
                newDate.setHours(newDate.getHours() + 1);
                onTimeConfigChange({ date: newDate });
              }}
              className="flex-1 bg-gray-700 hover:bg-gray-600 text-white text-xs py-1 rounded"
            >
              +1h
            </button>
            <button
              onClick={() => {
                const newDate = new Date(timeConfig.date);
                newDate.setHours(newDate.getHours() + 6);
                onTimeConfigChange({ date: newDate });
              }}
              className="flex-1 bg-gray-700 hover:bg-gray-600 text-white text-xs py-1 rounded"
            >
              +6h
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TimeControls;
