/**
 * CelestialChart - Container component for the geocentric celestial visualization
 *
 * Features:
 * - Manages chart state (hover, selection)
 * - Calculates panchangam elements from date/location
 * - Responsive sizing
 * - Delegates rendering to CelestialChartSVG
 */

import React, { useState, useMemo, useCallback, useRef, useEffect } from 'react';
import type {
  CelestialChartProps,
  ChartDimensions,
  ChartInteractionState,
  HoveredElement,
  SelectedElement,
} from './types';
import { createDefaultDimensions } from './types';
import { CelestialChartSVG } from './CelestialChartSVG';
import { calculateChartPanchangamForDate } from './utils/chartCalculations';

// Throttle helper for hover events
const throttle = <Args extends unknown[]>(
  func: (...args: Args) => void,
  limit: number
): ((...args: Args) => void) => {
  let inThrottle = false;
  return (...args: Args) => {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
};

export const CelestialChart: React.FC<CelestialChartProps> = ({
  date,
  className = '',
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [containerSize, setContainerSize] = useState(600);

  // Interaction state
  const [interactionState, setInteractionState] = useState<ChartInteractionState>({
    hoveredElement: null,
    selectedElement: null,
  });

  // Responsive sizing
  useEffect(() => {
    const updateSize = () => {
      if (containerRef.current) {
        const width = containerRef.current.clientWidth;
        // Make chart square, with max size
        const size = Math.min(width, 700);
        setContainerSize(size);
      }
    };

    updateSize();
    window.addEventListener('resize', updateSize);
    return () => window.removeEventListener('resize', updateSize);
  }, []);

  // Chart dimensions based on container size
  const dimensions = useMemo<ChartDimensions>(() =>
    createDefaultDimensions(containerSize),
    [containerSize]
  );

  // Calculate chart-ready panchangam elements from the selected date
  const panchangam = useMemo(() => {
    return calculateChartPanchangamForDate(date);
  }, [date]);

  // Throttled hover handler for performance
  const handleElementHover = useMemo(() =>
    throttle((element: HoveredElement | null) => {
      setInteractionState(prev => ({
        ...prev,
        hoveredElement: element,
      }));
    }, 50), // 20 FPS throttle
    []
  );

  // Selection handler
  const handleElementSelect = useCallback((element: SelectedElement | null) => {
    setInteractionState(prev => ({
      ...prev,
      selectedElement: prev.selectedElement?.id === element?.id ? null : element,
    }));
  }, []);

  // Clear selection when clicking outside
  const handleBackgroundClick = useCallback(() => {
    setInteractionState(prev => ({
      ...prev,
      selectedElement: null,
    }));
  }, []);

  return (
    <div
      ref={containerRef}
      className={`celestial-chart-container ${className}`}
      style={{
        width: '100%',
        maxWidth: '700px',
        margin: '0 auto',
      }}
    >
      {/* Chart Header */}
      <div className="text-center mb-4">
        <h3 className="text-xl font-semibold text-orange-800">
          Celestial Chart
        </h3>
        <p className="text-sm text-orange-600">
          {date.toLocaleDateString('en-US', {
            weekday: 'long',
            year: 'numeric',
            month: 'long',
            day: 'numeric',
          })}
        </p>
      </div>

      {/* Chart SVG */}
      <div
        className="rounded-lg overflow-hidden shadow-lg"
        onClick={handleBackgroundClick}
      >
        <CelestialChartSVG
          dimensions={dimensions}
          panchangam={panchangam}
          interactionState={interactionState}
          onElementHover={handleElementHover}
          onElementSelect={handleElementSelect}
        />
      </div>

      {/* Info Panel for Selected Element */}
      {interactionState.selectedElement && (
        <div className="mt-4 p-4 bg-white rounded-lg shadow-md">
          <h4 className="text-lg font-semibold text-gray-800 mb-2">
            {interactionState.selectedElement.type.charAt(0).toUpperCase() +
              interactionState.selectedElement.type.slice(1)} Details
          </h4>
          <div className="text-sm text-gray-600">
            {(() => {
              const { type, data } = interactionState.selectedElement;
              switch (type) {
                case 'rashi': {
                  const rashi = data as {
                    name: string;
                    westernName: string;
                    element: string;
                    ruler: string;
                    startDegree: number;
                    endDegree: number;
                  };
                  return (
                    <div className="grid grid-cols-2 gap-2">
                      <div><span className="font-medium">Name:</span> {rashi.name}</div>
                      <div><span className="font-medium">Western:</span> {rashi.westernName}</div>
                      <div><span className="font-medium">Element:</span> {rashi.element}</div>
                      <div><span className="font-medium">Ruler:</span> {rashi.ruler}</div>
                      <div><span className="font-medium">Degrees:</span> {rashi.startDegree}° - {rashi.endDegree}°</div>
                    </div>
                  );
                }
                case 'nakshatra': {
                  const nakshatra = data as {
                    name: string;
                    deity: string;
                    symbol: string;
                    startDegree: number;
                    endDegree: number;
                  };
                  return (
                    <div className="grid grid-cols-2 gap-2">
                      <div><span className="font-medium">Name:</span> {nakshatra.name}</div>
                      <div><span className="font-medium">Deity:</span> {nakshatra.deity}</div>
                      <div><span className="font-medium">Symbol:</span> {nakshatra.symbol}</div>
                      <div><span className="font-medium">Degrees:</span> {nakshatra.startDegree.toFixed(2)}° - {nakshatra.endDegree.toFixed(2)}°</div>
                    </div>
                  );
                }
                case 'pada': {
                  const pada = data as {
                    number: number;
                    padaInNakshatra: number;
                    nakshatra: { name: string };
                    navamsha: string;
                    startDegree: number;
                    endDegree: number;
                  };
                  return (
                    <div className="grid grid-cols-2 gap-2">
                      <div><span className="font-medium">Pada:</span> {pada.number} ({pada.nakshatra.name} Pada {pada.padaInNakshatra})</div>
                      <div><span className="font-medium">Navamsha:</span> {pada.navamsha}</div>
                      <div><span className="font-medium">Degrees:</span> {pada.startDegree.toFixed(2)}° - {pada.endDegree.toFixed(2)}°</div>
                    </div>
                  );
                }
                case 'tithi': {
                  const tithi = data as {
                    name: string;
                    paksha: string;
                    number: number;
                    deity: string;
                    angle: number;
                  };
                  return (
                    <div className="grid grid-cols-2 gap-2">
                      <div><span className="font-medium">Tithi:</span> {tithi.name}</div>
                      <div><span className="font-medium">Paksha:</span> {tithi.paksha}</div>
                      <div><span className="font-medium">Number:</span> {tithi.number}/30</div>
                      <div><span className="font-medium">Deity:</span> {tithi.deity}</div>
                      <div><span className="font-medium">Angle:</span> {tithi.angle.toFixed(1)}°</div>
                    </div>
                  );
                }
                case 'sun':
                case 'moon': {
                  const body = data as {
                    longitude: number;
                    label: string;
                  };
                  return (
                    <div className="grid grid-cols-2 gap-2">
                      <div><span className="font-medium">Body:</span> {type.charAt(0).toUpperCase() + type.slice(1)}</div>
                      <div><span className="font-medium">Longitude:</span> {body.longitude.toFixed(2)}°</div>
                      <div className="col-span-2"><span className="font-medium">Position:</span> {body.label}</div>
                    </div>
                  );
                }
                default:
                  return null;
              }
            })()}
          </div>
          <button
            className="mt-3 text-sm text-orange-600 hover:text-orange-800"
            onClick={() => handleElementSelect(null)}
          >
            Clear Selection
          </button>
        </div>
      )}

      {/* Quick Reference */}
      <div className="mt-4 p-3 bg-orange-50 rounded-lg text-sm text-gray-600">
        <p className="font-medium text-gray-700 mb-1">How to Read:</p>
        <ul className="list-disc list-inside space-y-1 text-xs">
          <li><strong>Outer Ring:</strong> 12 Rashis (Zodiac Signs) - 30° each</li>
          <li><strong>Middle Ring:</strong> 27 Nakshatras (Lunar Mansions) - 13°20' each</li>
          <li><strong>Inner Ring:</strong> 108 Padas (Quarters) - 3°20' each</li>
          <li><strong>Arc:</strong> Tithi - Angular separation between Sun and Moon</li>
        </ul>
      </div>
    </div>
  );
};

export default CelestialChart;
