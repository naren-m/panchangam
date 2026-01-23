/**
 * CelestialChartSVG - Pure SVG renderer for the geocentric celestial chart
 *
 * This is the presenter component that renders the visual elements.
 * All state and logic are handled by the CelestialChart container.
 */

import React, { useCallback, useMemo } from 'react';
import type {
  CelestialChartSVGProps,
  HoveredElement,
  SelectedElement,
  Point,
  PadaInfo
} from './types';
import { CHART_COLORS } from './types';
import { RashiRing, NakshatraRing, PadaRing } from './rings';
import { CelestialBody, TithiArc } from './markers';
import { ChartTooltip } from './tooltips';
import {
  getSunMarkerInfo,
  getMoonMarkerInfo,
  getTithiArcInfo,
  getAllRashis,
  getAllNakshatras,
  getAllPadas,
} from './utils/chartCalculations';

export const CelestialChartSVG: React.FC<CelestialChartSVGProps> = ({
  dimensions,
  panchangam,
  interactionState,
  onElementHover,
  onElementSelect,
}) => {
  const { size, center, earthRadius } = dimensions;
  const { hoveredElement, selectedElement } = interactionState;

  // Calculate celestial body positions
  const sunMarker = useMemo(() =>
    getSunMarkerInfo(panchangam, dimensions),
    [panchangam, dimensions]
  );

  const moonMarker = useMemo(() =>
    getMoonMarkerInfo(panchangam, dimensions),
    [panchangam, dimensions]
  );

  const tithiArc = useMemo(() =>
    getTithiArcInfo(panchangam, dimensions),
    [panchangam, dimensions]
  );

  // Memoized data arrays
  const rashis = useMemo(() => getAllRashis(), []);
  const nakshatras = useMemo(() => getAllNakshatras(), []);
  const padas = useMemo(() => getAllPadas(), []);

  // Ring hover handlers
  const handleRashiHover = useCallback((id: string | null, position: Point) => {
    if (!id) {
      onElementHover(null);
      return;
    }
    const rashiNum = parseInt(id.replace('rashi-', ''), 10);
    const rashi = rashis[rashiNum - 1];
    if (rashi) {
      onElementHover({
        type: 'rashi',
        id,
        data: rashi,
        position,
      });
    }
  }, [rashis, onElementHover]);

  const handleNakshatraHover = useCallback((id: string | null, position: Point) => {
    if (!id) {
      onElementHover(null);
      return;
    }
    const nakshatraNum = parseInt(id.replace('nakshatra-', ''), 10);
    const nakshatra = nakshatras[nakshatraNum - 1];
    if (nakshatra) {
      onElementHover({
        type: 'nakshatra',
        id,
        data: nakshatra,
        position,
      });
    }
  }, [nakshatras, onElementHover]);

  const handlePadaHover = useCallback((id: string | null, position: Point) => {
    if (!id) {
      onElementHover(null);
      return;
    }
    const padaNum = parseInt(id.replace('pada-', ''), 10);
    const pada = padas[padaNum - 1];
    if (pada) {
      onElementHover({
        type: 'pada',
        id,
        data: pada,
        position,
      });
    }
  }, [padas, onElementHover]);

  // Ring select handlers
  const handleRashiSelect = useCallback((id: string) => {
    const rashiNum = parseInt(id.replace('rashi-', ''), 10);
    const rashi = rashis[rashiNum - 1];
    if (rashi) {
      onElementSelect({
        type: 'rashi',
        id,
        data: rashi,
      });
    }
  }, [rashis, onElementSelect]);

  const handleNakshatraSelect = useCallback((id: string) => {
    const nakshatraNum = parseInt(id.replace('nakshatra-', ''), 10);
    const nakshatra = nakshatras[nakshatraNum - 1];
    if (nakshatra) {
      onElementSelect({
        type: 'nakshatra',
        id,
        data: nakshatra,
      });
    }
  }, [nakshatras, onElementSelect]);

  const handlePadaSelect = useCallback((id: string) => {
    const padaNum = parseInt(id.replace('pada-', ''), 10);
    const pada = padas[padaNum - 1];
    if (pada) {
      onElementSelect({
        type: 'pada',
        id,
        data: pada,
      });
    }
  }, [padas, onElementSelect]);

  // Celestial body handlers
  const handleSunHover = useCallback((hovered: boolean, position: Point) => {
    if (hovered) {
      onElementHover({
        type: 'sun',
        id: 'sun',
        data: sunMarker,
        position,
      });
    } else {
      onElementHover(null);
    }
  }, [sunMarker, onElementHover]);

  const handleMoonHover = useCallback((hovered: boolean, position: Point) => {
    if (hovered) {
      onElementHover({
        type: 'moon',
        id: 'moon',
        data: moonMarker,
        position,
      });
    } else {
      onElementHover(null);
    }
  }, [moonMarker, onElementHover]);

  const handleTithiHover = useCallback((hovered: boolean, position: Point) => {
    if (hovered) {
      onElementHover({
        type: 'tithi',
        id: 'tithi',
        data: panchangam.tithi,
        position: { x: center.x, y: center.y },
      });
    } else {
      onElementHover(null);
    }
  }, [panchangam.tithi, center, onElementHover]);

  const handleSunClick = useCallback(() => {
    onElementSelect({
      type: 'sun',
      id: 'sun',
      data: sunMarker,
    });
  }, [sunMarker, onElementSelect]);

  const handleMoonClick = useCallback(() => {
    onElementSelect({
      type: 'moon',
      id: 'moon',
      data: moonMarker,
    });
  }, [moonMarker, onElementSelect]);

  const handleTithiClick = useCallback(() => {
    onElementSelect({
      type: 'tithi',
      id: 'tithi',
      data: panchangam.tithi,
    });
  }, [panchangam.tithi, onElementSelect]);

  return (
    <svg
      width={size}
      height={size}
      viewBox={`0 0 ${size} ${size}`}
      className="celestial-chart-svg"
      style={{ backgroundColor: '#1a1a2e' }}
    >
      {/* Background gradient for space effect */}
      <defs>
        <radialGradient id="space-bg" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#1a1a2e" />
          <stop offset="100%" stopColor="#0d0d1a" />
        </radialGradient>
      </defs>
      <circle
        cx={center.x}
        cy={center.y}
        r={size / 2}
        fill="url(#space-bg)"
      />

      {/* Rings (outer to inner order for proper layering) */}
      <RashiRing
        dimensions={dimensions}
        hoveredId={hoveredElement?.type === 'rashi' ? hoveredElement.id : null}
        selectedId={selectedElement?.type === 'rashi' ? selectedElement.id : null}
        onHover={handleRashiHover}
        onSelect={handleRashiSelect}
        sunRashi={panchangam.sunRashi}
        moonRashi={panchangam.rashi}
        panchangam={panchangam}
      />

      <NakshatraRing
        dimensions={dimensions}
        hoveredId={hoveredElement?.type === 'nakshatra' ? hoveredElement.id : null}
        selectedId={selectedElement?.type === 'nakshatra' ? selectedElement.id : null}
        onHover={handleNakshatraHover}
        onSelect={handleNakshatraSelect}
        currentNakshatra={panchangam.nakshatra}
        panchangam={panchangam}
      />

      <PadaRing
        dimensions={dimensions}
        hoveredId={hoveredElement?.type === 'pada' ? hoveredElement.id : null}
        selectedId={selectedElement?.type === 'pada' ? selectedElement.id : null}
        onHover={handlePadaHover}
        onSelect={handlePadaSelect}
        currentNakshatra={panchangam.nakshatra}
        currentPada={panchangam.nakshatra.pada}
        panchangam={panchangam}
      />

      {/* Earth at center */}
      <g className="earth-marker">
        <defs>
          <radialGradient id="earth-gradient" cx="30%" cy="30%" r="70%">
            <stop offset="0%" stopColor="#4FC3F7" />
            <stop offset="50%" stopColor="#29B6F6" />
            <stop offset="100%" stopColor="#0288D1" />
          </radialGradient>
        </defs>
        <circle
          cx={center.x}
          cy={center.y}
          r={earthRadius}
          fill="url(#earth-gradient)"
          stroke="#01579B"
          strokeWidth={2}
        />
        {/* Simple Earth symbol */}
        <text
          x={center.x}
          y={center.y}
          textAnchor="middle"
          dominantBaseline="central"
          fontSize={earthRadius * 0.9}
          fill="#01579B"
          fontWeight="bold"
        >
          ⊕
        </text>
      </g>

      {/* Tithi arc connecting Sun and Moon */}
      <TithiArc
        arcInfo={tithiArc}
        isHovered={hoveredElement?.type === 'tithi'}
        onHover={handleTithiHover}
        onClick={handleTithiClick}
      />

      {/* Celestial bodies */}
      <CelestialBody
        body={sunMarker}
        isHovered={hoveredElement?.type === 'sun'}
        onHover={handleSunHover}
        onClick={handleSunClick}
      />

      <CelestialBody
        body={moonMarker}
        isHovered={hoveredElement?.type === 'moon'}
        onHover={handleMoonHover}
        onClick={handleMoonClick}
      />

      {/* Tooltip (rendered last for z-index) */}
      {hoveredElement && hoveredElement.type !== 'sun' && hoveredElement.type !== 'moon' && (
        <ChartTooltip
          element={hoveredElement}
          chartDimensions={dimensions}
        />
      )}

      {/* Legend */}
      <g className="chart-legend" transform={`translate(10, ${size - 80})`}>
        <rect
          x={0}
          y={0}
          width={140}
          height={70}
          rx={6}
          fill={CHART_COLORS.text}
          fillOpacity={0.8}
        />
        <text x={10} y={18} fill="white" fontSize={10} fontWeight="bold">
          Current Position
        </text>
        <circle cx={20} cy={35} r={6} fill={CHART_COLORS.sun} />
        <text x={32} y={38} fill="white" fontSize={9}>
          Sun: {panchangam.sunRashi.name}
        </text>
        <circle cx={20} cy={55} r={6} fill={CHART_COLORS.moon} />
        <text x={32} y={58} fill="white" fontSize={9}>
          Moon: {panchangam.nakshatra.name}
        </text>
      </g>
    </svg>
  );
};

export default CelestialChartSVG;
