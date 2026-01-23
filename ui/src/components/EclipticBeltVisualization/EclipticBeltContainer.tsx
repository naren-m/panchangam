/**
 * EclipticBeltContainer - Smart Container for the 2D Ecliptic Belt Visualization
 *
 * This component handles:
 * - Calculating Sun and Moon positions from date
 * - Computing Panchangam elements (Tithi, Nakshatra, Yoga, Karana)
 * - Managing interaction state (hover, selection)
 * - Time animation controls
 * - Responsive sizing
 *
 * The actual astronomical calculations use simplified algorithms.
 * For production accuracy, consider integrating Swiss Ephemeris or similar.
 */

import React, { useState, useEffect, useMemo, useCallback, useRef } from 'react';
import {
  EclipticBeltContainerProps,
  EclipticBeltDimensions,
  PanchangamElements,
  TimeControlState,
  LAHIRI_AYANAMSA_2024,
} from './types/eclipticBelt';
import { calculatePanchangamElements, calculateTithiStartTime, calculateTithiEndTime } from './utils/panchangamCalculator';
import {
  generateLayout,
  createDefaultDimensions,
  generateRashiSegments,
  generateNakshatraSegments,
} from './utils/eclipticLayout';
import { generateAllAnnotations, generateVisualizationGuide } from './utils/annotationHelper';
import { EclipticBeltSVG } from './EclipticBeltSVG';

// ============================================================================
// Simplified Astronomical Calculations
// ============================================================================

/**
 * Calculate Julian Day Number from Date
 * Based on Jean Meeus' Astronomical Algorithms
 */
function dateToJulianDay(date: Date): number {
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate() + date.getHours() / 24 + date.getMinutes() / 1440;

  let y = year;
  let m = month;

  if (m <= 2) {
    y -= 1;
    m += 12;
  }

  const a = Math.floor(y / 100);
  const b = 2 - a + Math.floor(a / 4);

  return Math.floor(365.25 * (y + 4716)) +
    Math.floor(30.6001 * (m + 1)) +
    day + b - 1524.5;
}

/**
 * Calculate Sun's ecliptic longitude (simplified)
 *
 * Uses a simplified algorithm that's accurate to within ~1° for most dates.
 * For higher accuracy, use Swiss Ephemeris or VSOP87.
 *
 * @param jd Julian Day
 * @returns Sun's tropical longitude in degrees (0-360)
 */
function calculateSunLongitude(jd: number): number {
  // Days since J2000.0
  const T = (jd - 2451545.0) / 36525.0;

  // Mean longitude of the Sun
  const L0 = 280.46646 + 36000.76983 * T + 0.0003032 * T * T;

  // Mean anomaly of the Sun
  const M = 357.52911 + 35999.05029 * T - 0.0001537 * T * T;

  // Equation of center
  const C = (1.914602 - 0.004817 * T) * Math.sin(M * Math.PI / 180)
    + 0.019993 * Math.sin(2 * M * Math.PI / 180)
    + 0.000289 * Math.sin(3 * M * Math.PI / 180);

  // Sun's true longitude
  let sunLong = L0 + C;

  // Normalize to 0-360
  sunLong = sunLong % 360;
  if (sunLong < 0) sunLong += 360;

  return sunLong;
}

/**
 * Calculate Moon's ecliptic longitude (simplified)
 *
 * Uses a simplified algorithm. For higher accuracy, use Swiss Ephemeris.
 *
 * @param jd Julian Day
 * @returns Moon's tropical longitude in degrees (0-360)
 */
function calculateMoonLongitude(jd: number): number {
  // Days since J2000.0
  const T = (jd - 2451545.0) / 36525.0;

  // Moon's mean longitude
  const L = 218.3164477 + 481267.88123421 * T - 0.0015786 * T * T;

  // Moon's mean anomaly
  const M = 134.9633964 + 477198.8675055 * T + 0.0087414 * T * T;

  // Moon's argument of latitude
  const F = 93.2720950 + 483202.0175233 * T - 0.0036539 * T * T;

  // Mean elongation of Moon from Sun
  const D = 297.8501921 + 445267.1114034 * T - 0.0018819 * T * T;

  // Sun's mean anomaly
  const Ms = 357.5291092 + 35999.0502909 * T - 0.0001536 * T * T;

  // Convert to radians
  const toRad = Math.PI / 180;

  // Simplified longitude perturbations
  let longitude = L
    + 6.289 * Math.sin(M * toRad)
    + 1.274 * Math.sin((2 * D - M) * toRad)
    + 0.658 * Math.sin(2 * D * toRad)
    + 0.214 * Math.sin(2 * M * toRad)
    - 0.186 * Math.sin(Ms * toRad)
    - 0.114 * Math.sin(2 * F * toRad);

  // Normalize to 0-360
  longitude = longitude % 360;
  if (longitude < 0) longitude += 360;

  return longitude;
}

/**
 * Convert tropical longitude to sidereal (using Lahiri ayanamsa)
 */
function tropicalToSidereal(tropicalLong: number, date: Date): number {
  // Approximate ayanamsa calculation
  // Ayanamsa increases by about 50.3" per year from the reference epoch
  const year = date.getFullYear();
  const yearsSince2024 = year - 2024 + (date.getMonth() / 12);
  const ayanamsa = LAHIRI_AYANAMSA_2024 + (yearsSince2024 * 50.3 / 3600);

  let sidereal = tropicalLong - ayanamsa;
  if (sidereal < 0) sidereal += 360;

  return sidereal;
}

/**
 * Calculate all celestial positions for a given date
 */
function calculatePositions(date: Date): { sunLong: number; moonLong: number } {
  const jd = dateToJulianDay(date);

  // Get tropical positions
  const sunTropical = calculateSunLongitude(jd);
  const moonTropical = calculateMoonLongitude(jd);

  // Convert to sidereal (Vedic/Indian astrology uses sidereal zodiac)
  const sunLong = tropicalToSidereal(sunTropical, date);
  const moonLong = tropicalToSidereal(moonTropical, date);

  return { sunLong, moonLong };
}

// ============================================================================
// Time Control Component
// ============================================================================

const TimeControls: React.FC<{
  timeState: TimeControlState;
  onDateChange: (date: Date) => void;
  onPlayPause: () => void;
  onSpeedChange: (speed: number) => void;
  onToday: () => void;
}> = ({ timeState, onDateChange, onPlayPause, onSpeedChange, onToday }) => {
  const formatDate = (date: Date) => {
    return date.toLocaleDateString('en-US', {
      weekday: 'short',
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const adjustTime = (hours: number) => {
    const newDate = new Date(timeState.currentDate);
    newDate.setHours(newDate.getHours() + hours);
    onDateChange(newDate);
  };

  return (
    <div className="flex flex-wrap items-center justify-center gap-3 p-3 bg-orange-50 rounded-lg border border-orange-200">
      {/* Navigation buttons */}
      <button
        onClick={() => adjustTime(-24)}
        className="px-3 py-1.5 bg-orange-100 hover:bg-orange-200 rounded text-sm font-medium text-orange-800 transition-colors"
        title="Previous day"
      >
        ◀◀ Day
      </button>
      <button
        onClick={() => adjustTime(-1)}
        className="px-3 py-1.5 bg-orange-100 hover:bg-orange-200 rounded text-sm font-medium text-orange-800 transition-colors"
        title="Previous hour"
      >
        ◀ Hour
      </button>

      {/* Current date display */}
      <div className="px-4 py-1.5 bg-white rounded-lg border border-orange-200 min-w-[280px] text-center">
        <span className="text-sm font-medium text-gray-800">
          {formatDate(timeState.currentDate)}
        </span>
      </div>

      <button
        onClick={() => adjustTime(1)}
        className="px-3 py-1.5 bg-orange-100 hover:bg-orange-200 rounded text-sm font-medium text-orange-800 transition-colors"
        title="Next hour"
      >
        Hour ▶
      </button>
      <button
        onClick={() => adjustTime(24)}
        className="px-3 py-1.5 bg-orange-100 hover:bg-orange-200 rounded text-sm font-medium text-orange-800 transition-colors"
        title="Next day"
      >
        Day ▶▶
      </button>

      {/* Today button */}
      <button
        onClick={onToday}
        className="px-3 py-1.5 bg-orange-500 hover:bg-orange-600 text-white rounded text-sm font-medium transition-colors"
      >
        Today
      </button>

      {/* Play/Pause (animation) */}
      <button
        onClick={onPlayPause}
        className={`px-3 py-1.5 rounded text-sm font-medium transition-colors ${timeState.isPlaying
            ? 'bg-red-100 hover:bg-red-200 text-red-800'
            : 'bg-green-100 hover:bg-green-200 text-green-800'
          }`}
      >
        {timeState.isPlaying ? '⏸ Pause' : '▶ Play'}
      </button>

      {/* Speed selector */}
      {timeState.isPlaying && (
        <select
          value={timeState.playbackSpeed}
          onChange={(e) => onSpeedChange(Number(e.target.value))}
          className="px-2 py-1.5 bg-white border border-orange-200 rounded text-sm"
        >
          <option value={1}>1x (1 min/sec)</option>
          <option value={60}>60x (1 hr/sec)</option>
          <option value={1440}>1440x (1 day/sec)</option>
        </select>
      )}
    </div>
  );
};

// ============================================================================
// Info Panel Component
// ============================================================================

const InfoPanel: React.FC<{
  panchangam: PanchangamElements;
  selectedElement: string | null;
  onClose: () => void;
}> = ({ panchangam, selectedElement, onClose }) => {
  const { tithi, nakshatra, yoga, karana, sunPosition, moonPosition, rashi, sunRashi } = panchangam;

  return (
    <div className="bg-white rounded-lg shadow-lg p-4 max-w-md">
      <div className="flex justify-between items-start mb-3">
        <h3 className="text-lg font-semibold text-orange-800">Panchangam Details</h3>
        <button
          onClick={onClose}
          className="text-gray-400 hover:text-gray-600"
        >
          ✕
        </button>
      </div>

      <div className="space-y-3 text-sm">
        {/* Tithi */}
        <div className={`p-2 rounded ${selectedElement?.includes('tithi') ? 'bg-orange-100' : 'bg-gray-50'}`}>
          <div className="font-medium text-gray-700">🌙 Tithi (Lunar Day)</div>
          <div className="text-gray-600">
            {tithi.paksha} {tithi.name} ({tithi.percentComplete.toFixed(0)}% complete)
          </div>
          <div className="text-xs text-gray-500">
            Moon-Sun: {tithi.angle.toFixed(1)}° | Deity: {tithi.deity}
          </div>
          {tithi.startTime && tithi.endTime && (
            <div className="text-xs text-blue-600 mt-1 pt-1 border-t border-gray-200">
              ⏱ Started: {tithi.startTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
              {' '} | Ends: {tithi.endTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
              {tithi.endTime.getDate() !== tithi.startTime.getDate() && (
                <span className="text-gray-500"> (+1 day)</span>
              )}
            </div>
          )}
        </div>

        {/* Nakshatra */}
        <div className={`p-2 rounded ${selectedElement?.includes('nakshatra') ? 'bg-orange-100' : 'bg-gray-50'}`}>
          <div className="font-medium text-gray-700">⭐ Nakshatra (Lunar Mansion)</div>
          <div className="text-gray-600">
            {nakshatra.name} (Pada {nakshatra.pada})
          </div>
          <div className="text-xs text-gray-500">
            Deity: {nakshatra.deity} | Symbol: {nakshatra.symbol}
          </div>
        </div>

        {/* Yoga */}
        <div className={`p-2 rounded ${selectedElement?.includes('yoga') ? 'bg-orange-100' : 'bg-gray-50'}`}>
          <div className="font-medium text-gray-700">🔮 Yoga</div>
          <div className="text-gray-600">
            {yoga.name} - "{yoga.meaning}"
          </div>
          <div className="text-xs text-gray-500">
            Nature: {yoga.nature} | Combined: {yoga.combinedLongitude.toFixed(1)}°
          </div>
        </div>

        {/* Karana */}
        <div className={`p-2 rounded ${selectedElement?.includes('karana') ? 'bg-orange-100' : 'bg-gray-50'}`}>
          <div className="font-medium text-gray-700">⌚ Karana (Half-Tithi)</div>
          <div className="text-gray-600">{karana.name}</div>
          <div className="text-xs text-gray-500">
            Type: {karana.type} | Nature: {karana.nature}
          </div>
        </div>

        {/* Celestial Positions */}
        <div className="p-2 bg-yellow-50 rounded">
          <div className="font-medium text-gray-700">☉ Sun &amp; 🌑 Moon</div>
          <div className="grid grid-cols-2 gap-2 mt-1 text-xs text-gray-600">
            <div>
              Sun: {sunRashi.symbol} {sunRashi.name}
              <br />
              <span className="text-gray-500">{sunPosition.longitude.toFixed(2)}°</span>
            </div>
            <div>
              Moon: {rashi.symbol} {rashi.name}
              <br />
              <span className="text-gray-500">{moonPosition.longitude.toFixed(2)}°</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

// ============================================================================
// Main Container Component
// ============================================================================

export const EclipticBeltContainer: React.FC<EclipticBeltContainerProps> = ({
  date: initialDate,
  latitude,
  longitude,
  timezone = 'America/Los_Angeles',
  panchangamData,
  onClose,
  className = '',
}) => {
  // Container ref for responsive sizing
  const containerRef = useRef<HTMLDivElement>(null);

  // State
  const [dimensions, setDimensions] = useState<EclipticBeltDimensions>(
    createDefaultDimensions(1000, 400)
  );
  const [timeState, setTimeState] = useState<TimeControlState>({
    currentDate: initialDate,
    isPlaying: false,
    playbackSpeed: 60,
    direction: 'forward',
  });
  const [selectedElement, setSelectedElement] = useState<string | null>(null);
  const [hoveredElement, setHoveredElement] = useState<string | null>(null);
  const [showGuide, setShowGuide] = useState(false);

  // Handle resize
  useEffect(() => {
    const updateDimensions = () => {
      if (containerRef.current) {
        const width = containerRef.current.offsetWidth;
        const height = Math.min(450, Math.max(350, width * 0.35));
        setDimensions(createDefaultDimensions(width, height));
      }
    };

    updateDimensions();
    window.addEventListener('resize', updateDimensions);
    return () => window.removeEventListener('resize', updateDimensions);
  }, []);

  // Animation loop
  useEffect(() => {
    if (!timeState.isPlaying) return;

    const intervalMs = 1000; // Update every second
    const interval = setInterval(() => {
      setTimeState((prev) => {
        const newDate = new Date(prev.currentDate);
        // Add minutes based on playback speed
        newDate.setMinutes(newDate.getMinutes() + prev.playbackSpeed);
        return { ...prev, currentDate: newDate };
      });
    }, intervalMs);

    return () => clearInterval(interval);
  }, [timeState.isPlaying, timeState.playbackSpeed]);

  // Calculate positions and panchangam
  const { sunLong, moonLong } = useMemo(
    () => calculatePositions(timeState.currentDate),
    [timeState.currentDate]
  );

  const panchangam = useMemo(() => {
    const elements = calculatePanchangamElements(sunLong, moonLong);

    // Calculate tithi start and end times
    const tithiStartTime = calculateTithiStartTime(
      timeState.currentDate,
      elements.tithi.number,
      calculatePositions
    );
    const tithiEndTime = calculateTithiEndTime(
      timeState.currentDate,
      elements.tithi.number,
      calculatePositions
    );

    // Enhance tithi with times
    return {
      ...elements,
      tithi: {
        ...elements.tithi,
        startTime: tithiStartTime,
        endTime: tithiEndTime,
      },
    };
  }, [sunLong, moonLong, timeState.currentDate]);

  // Generate layout
  const layout = useMemo(
    () => generateLayout(panchangam, dimensions.width, dimensions.height),
    [panchangam, dimensions]
  );

  // Generate annotations
  const annotations = useMemo(
    () => generateAllAnnotations(panchangam, dimensions),
    [panchangam, dimensions]
  );

  // Handlers
  const handleDateChange = useCallback((newDate: Date) => {
    setTimeState((prev) => ({ ...prev, currentDate: newDate, isPlaying: false }));
  }, []);

  const handlePlayPause = useCallback(() => {
    setTimeState((prev) => ({ ...prev, isPlaying: !prev.isPlaying }));
  }, []);

  const handleSpeedChange = useCallback((speed: number) => {
    setTimeState((prev) => ({ ...prev, playbackSpeed: speed }));
  }, []);

  const handleToday = useCallback(() => {
    setTimeState((prev) => ({ ...prev, currentDate: new Date(), isPlaying: false }));
  }, []);

  const handleElementSelect = useCallback((id: string | null) => {
    setSelectedElement(id);
  }, []);

  const handleElementHover = useCallback((id: string | null) => {
    setHoveredElement(id);
  }, []);

  return (
    <div
      ref={containerRef}
      className={`ecliptic-belt-container bg-white rounded-xl shadow-xl overflow-hidden ${className}`}
    >
      {/* Header */}
      <div className="flex justify-between items-center px-4 py-3 bg-gradient-to-r from-orange-500 to-amber-500">
        <div className="flex items-center gap-3">
          <h2 className="text-xl font-bold text-white">
            🌌 Ecliptic Belt Visualization
          </h2>
          <button
            onClick={() => setShowGuide(!showGuide)}
            className="px-2 py-1 bg-white/20 hover:bg-white/30 rounded text-white text-sm transition-colors"
          >
            {showGuide ? 'Hide' : 'Show'} Guide
          </button>
        </div>
        {onClose && (
          <button
            onClick={onClose}
            className="text-white hover:text-orange-100 transition-colors"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        )}
      </div>

      {/* Guide (collapsible) */}
      {showGuide && (
        <div className="px-4 py-3 bg-blue-50 border-b border-blue-100 text-sm text-gray-700">
          <p className="mb-2"><strong>Understanding the Visualization:</strong></p>
          <ul className="list-disc list-inside space-y-1 text-xs">
            <li><strong>Top row:</strong> 12 Rashis (zodiac signs), each 30° wide</li>
            <li><strong>Second row:</strong> 27 Nakshatras (lunar mansions), each 13.3° wide</li>
            <li><strong>☉ Sun &amp; 🌑 Moon:</strong> Current celestial positions</li>
            <li><strong>Arc:</strong> Moon-Sun angular separation (determines Tithi)</li>
            <li>Click any element for details. Use time controls to animate.</li>
          </ul>
        </div>
      )}

      {/* Time Controls */}
      <TimeControls
        timeState={timeState}
        onDateChange={handleDateChange}
        onPlayPause={handlePlayPause}
        onSpeedChange={handleSpeedChange}
        onToday={handleToday}
      />

      {/* Main Visualization */}
      <div className="p-4">
        <EclipticBeltSVG
          dimensions={dimensions}
          panchangam={panchangam}
          rashiSegments={layout.rashiSegments}
          nakshatraSegments={layout.nakshatraSegments}
          sunMarker={layout.sunMarker}
          moonMarker={layout.moonMarker}
          tithiArc={layout.tithiArc}
          annotations={annotations}
          selectedElement={selectedElement}
          onElementSelect={handleElementSelect}
          onElementHover={handleElementHover}
          hoveredElement={hoveredElement}
          showLabels={true}
          animationEnabled={timeState.isPlaying}
        />
      </div>

      {/* Info Panel (shown when element is selected) */}
      {selectedElement && (
        <div className="px-4 pb-4">
          <InfoPanel
            panchangam={panchangam}
            selectedElement={selectedElement}
            onClose={() => setSelectedElement(null)}
          />
        </div>
      )}

      {/* Footer with calculation note */}
      <div className="px-4 py-2 bg-gray-50 text-xs text-gray-500 text-center border-t">
        Calculations use Lahiri Ayanamsa for sidereal positions.
        For educational purposes - accuracy may vary.
      </div>
    </div>
  );
};

export default EclipticBeltContainer;
