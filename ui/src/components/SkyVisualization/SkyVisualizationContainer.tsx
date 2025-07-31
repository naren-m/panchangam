import React, { useState, useEffect } from 'react';
import { Eye, EyeOff, Settings, Clock, Maximize2, Minimize2 } from 'lucide-react';
import SkySphere from './SkySphere';
import { 
  Observer, 
  TimeConfig, 
  RenderOptions, 
  CelestialObject,
  SkySphereConfig 
} from '../../types/skyVisualization';
import { getNakshatraInfo, calculateNakshatraFromLongitude } from '../../utils/astronomy/nakshatraCalculator';
import { getZodiacInfo, calculateZodiacFromLongitude } from '../../utils/astronomy/zodiacCalculator';

interface SkyVisualizationContainerProps {
  latitude: number;
  longitude: number;
  date?: Date;
  className?: string;
}

export const SkyVisualizationContainer: React.FC<SkyVisualizationContainerProps> = ({
  latitude,
  longitude,
  date = new Date(),
  className
}) => {
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [showControls, setShowControls] = useState(true);
  const [celestialObjects, setCelestialObjects] = useState<CelestialObject[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentNakshatra, setCurrentNakshatra] = useState<number | undefined>();
  const [currentRashi, setCurrentRashi] = useState<number | undefined>();
  
  const [timeConfig, setTimeConfig] = useState<TimeConfig>({
    date: date,
    speed: 1,
    paused: false
  });

  const [renderOptions, setRenderOptions] = useState<RenderOptions>({
    showGrid: true,
    showConstellations: true,
    showNakshatras: true,
    showPlanets: true,
    showStars: true,
    showLabels: false,
    showZodiac: true,
    showEcliptic: true,
    showEquator: true,
    showHorizon: true,
    starMagnitudeLimit: 6.0,
    labelMinZoom: 2.0
  });

  const observer: Observer = {
    latitude,
    longitude,
    altitude: 0
  };

  const skyConfig: Partial<SkySphereConfig> = {
    projection: 'stereographic',
    coordinateSystem: 'equatorial'
  };

  // Fetch planetary positions
  useEffect(() => {
    const fetchPlanetaryData = async () => {
      try {
        setLoading(true);
        setError(null);

        // Calculate approximate Moon position based on date
        // This is a simplified calculation - in a real implementation, 
        // this would use precise ephemeris data from the backend
        const daysSinceJ2000 = (timeConfig.date.getTime() - new Date('2000-01-01T12:00:00Z').getTime()) / (1000 * 60 * 60 * 24);
        
        // Rough approximation of lunar longitude (Moon moves ~13.2 degrees per day)
        const moonLongitude = (218.316 + 13.176396 * daysSinceJ2000) % 360;
        
        // Calculate current nakshatra from moon position
        const nakshatraNumber = calculateNakshatraFromLongitude(moonLongitude);
        setCurrentNakshatra(nakshatraNumber);
        
        // Rough approximation of solar longitude (Sun moves ~1 degree per day)
        const sunLongitude = (280.459 + 0.98564736 * daysSinceJ2000) % 360;
        
        // Calculate current rashi (zodiac sign) from sun position
        const rashiNumber = calculateZodiacFromLongitude(sunLongitude);
        setCurrentRashi(rashiNumber);
        
        const sampleObjects: CelestialObject[] = [
          {
            id: 'sun',
            name: 'Sun',
            type: 'sun',
            coordinates: {
              ecliptic: {
                longitude: sunLongitude,
                latitude: 0,
                distance: 1
              }
            },
            magnitude: -26.7,
            color: '#ffee00',
            size: 5,
            metadata: {
              rashi: rashiNumber,
              rashiInfo: getZodiacInfo(sunLongitude)
            }
          },
          {
            id: 'moon',
            name: 'Moon',
            type: 'moon',
            coordinates: {
              ecliptic: {
                longitude: moonLongitude,
                latitude: (Math.sin(daysSinceJ2000 * 0.1) * 5), // Approximate lunar latitude variation
                distance: 0.00257
              }
            },
            magnitude: -12.6,
            color: '#ffffff',
            size: 4,
            metadata: {
              nakshatra: nakshatraNumber,
              nakshatraInfo: getNakshatraInfo(moonLongitude)
            }
          },
          // Add some bright stars
          {
            id: 'sirius',
            name: 'Sirius',
            type: 'star',
            coordinates: {
              ecliptic: {
                longitude: 101.287,
                latitude: -39.608,
                distance: 8.6
              }
            },
            magnitude: -1.46,
            color: '#aaccff'
          },
          {
            id: 'canopus',
            name: 'Canopus',
            type: 'star',
            coordinates: {
              ecliptic: {
                longitude: 95.988,
                latitude: -52.696,
                distance: 310
              }
            },
            magnitude: -0.74,
            color: '#fff4e6'
          }
        ];

        setCelestialObjects(sampleObjects);
        setLoading(false);
      } catch (err) {
        setError('Failed to load celestial data');
        setLoading(false);
      }
    };

    fetchPlanetaryData();
  }, [date]);

  // Time animation
  useEffect(() => {
    if (timeConfig.paused || timeConfig.speed === 0) return;

    const interval = setInterval(() => {
      setTimeConfig(prev => ({
        ...prev,
        date: new Date(prev.date.getTime() + 60000 * prev.speed) // Advance by speed minutes
      }));
    }, 1000); // Update every second

    return () => clearInterval(interval);
  }, [timeConfig.paused, timeConfig.speed]);

  const toggleFullscreen = () => {
    setIsFullscreen(!isFullscreen);
  };

  const handleRenderOptionToggle = (option: keyof RenderOptions) => {
    setRenderOptions(prev => ({
      ...prev,
      [option]: !prev[option]
    }));
  };

  const handleTimeSpeedChange = (speed: number) => {
    setTimeConfig(prev => ({ ...prev, speed }));
  };

  const handlePauseToggle = () => {
    setTimeConfig(prev => ({ ...prev, paused: !prev.paused }));
  };

  const containerClass = isFullscreen 
    ? 'fixed inset-0 z-50 bg-black' 
    : `relative ${className || ''}`;

  if (loading) {
    return (
      <div className={`${containerClass} flex items-center justify-center bg-gray-900`}>
        <div className="text-white">Loading sky visualization...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`${containerClass} flex items-center justify-center bg-gray-900`}>
        <div className="text-red-500">{error}</div>
      </div>
    );
  }

  return (
    <div className={containerClass}>
      <SkySphere
        config={skyConfig}
        celestialObjects={celestialObjects}
        observer={observer}
        timeConfig={timeConfig}
        renderOptions={renderOptions}
        currentNakshatra={currentNakshatra}
        currentRashi={currentRashi}
        className="w-full h-full"
      />

      {/* Control Panel */}
      {showControls && (
        <div className="absolute top-4 left-4 bg-gray-800 bg-opacity-90 rounded-lg p-4 text-white max-w-xs">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Sky View Controls</h3>
            <button
              onClick={() => setShowControls(false)}
              className="text-gray-400 hover:text-white"
            >
              <EyeOff size={20} />
            </button>
          </div>

          {/* Time Controls */}
          <div className="mb-4">
            <h4 className="text-sm font-medium mb-2">Time Controls</h4>
            <div className="flex items-center gap-2 mb-2">
              <button
                onClick={handlePauseToggle}
                className={`px-3 py-1 rounded ${
                  timeConfig.paused ? 'bg-blue-600' : 'bg-gray-600'
                } hover:bg-opacity-80`}
              >
                {timeConfig.paused ? 'Play' : 'Pause'}
              </button>
              <select
                value={timeConfig.speed}
                onChange={(e) => handleTimeSpeedChange(Number(e.target.value))}
                className="bg-gray-700 rounded px-2 py-1"
              >
                <option value={0}>Stopped</option>
                <option value={1}>Real-time</option>
                <option value={60}>1 hour/sec</option>
                <option value={1440}>1 day/sec</option>
                <option value={10080}>1 week/sec</option>
              </select>
            </div>
            <div className="text-xs text-gray-400">
              {timeConfig.date.toLocaleString()}
            </div>
          </div>

          {/* Display Options */}
          <div className="space-y-2">
            <h4 className="text-sm font-medium mb-2">Display Options</h4>
            {Object.entries({
              showGrid: 'Coordinate Grid',
              showNakshatras: 'Nakshatras',
              showPlanets: 'Planets',
              showStars: 'Stars',
              showEcliptic: 'Ecliptic',
              showEquator: 'Celestial Equator',
              showHorizon: 'Horizon'
            }).map(([key, label]) => (
              <label key={key} className="flex items-center gap-2 text-sm">
                <input
                  type="checkbox"
                  checked={renderOptions[key as keyof RenderOptions] as boolean}
                  onChange={() => handleRenderOptionToggle(key as keyof RenderOptions)}
                  className="rounded"
                />
                {label}
              </label>
            ))}
          </div>
        </div>
      )}

      {/* Toggle Controls Button */}
      {!showControls && (
        <button
          onClick={() => setShowControls(true)}
          className="absolute top-4 left-4 bg-gray-800 bg-opacity-90 rounded-lg p-2 text-white hover:bg-opacity-100"
        >
          <Settings size={20} />
        </button>
      )}

      {/* Fullscreen Toggle */}
      <button
        onClick={toggleFullscreen}
        className="absolute top-4 right-4 bg-gray-800 bg-opacity-90 rounded-lg p-2 text-white hover:bg-opacity-100"
      >
        {isFullscreen ? <Minimize2 size={20} /> : <Maximize2 size={20} />}
      </button>

      {/* Info Panel */}
      <div className="absolute bottom-4 left-4 bg-gray-800 bg-opacity-90 rounded-lg p-3 text-white text-sm max-w-xs">
        <div>Location: {latitude.toFixed(2)}°, {longitude.toFixed(2)}°</div>
        
        {currentRashi && (
          <div className="mt-2 border-t border-gray-600 pt-2">
            <div className="font-medium">Current Rashi: {getZodiacInfo(
              celestialObjects.find(obj => obj.id === 'sun')?.coordinates.ecliptic?.longitude || 0
            ).sanskritName}</div>
            <div className="text-xs text-gray-400">
              {getZodiacInfo(
                celestialObjects.find(obj => obj.id === 'sun')?.coordinates.ecliptic?.longitude || 0
              ).symbol} {getZodiacInfo(
                celestialObjects.find(obj => obj.id === 'sun')?.coordinates.ecliptic?.longitude || 0
              ).westernName} • {getZodiacInfo(
                celestialObjects.find(obj => obj.id === 'sun')?.coordinates.ecliptic?.longitude || 0
              ).element}
            </div>
          </div>
        )}
        
        {currentNakshatra && (
          <div className="mt-2 border-t border-gray-600 pt-2">
            <div className="font-medium">Current Nakshatra: {getNakshatraInfo(
              celestialObjects.find(obj => obj.id === 'moon')?.coordinates.ecliptic?.longitude || 0
            ).name}</div>
            <div className="text-xs text-gray-400">#{currentNakshatra} • {
              getNakshatraInfo(
                celestialObjects.find(obj => obj.id === 'moon')?.coordinates.ecliptic?.longitude || 0
              ).deity
            }</div>
          </div>
        )}
        
        <div className="text-xs text-gray-400 mt-2">
          Drag to rotate • Scroll to zoom
        </div>
      </div>
    </div>
  );
};

export default SkyVisualizationContainer;