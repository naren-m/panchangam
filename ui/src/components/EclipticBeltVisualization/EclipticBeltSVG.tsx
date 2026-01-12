/**
 * EclipticBeltSVG - SVG Renderer for the 2D Ecliptic Belt Visualization
 *
 * Renders the complete ecliptic belt visualization as an SVG with:
 * - Rashi (zodiac sign) segments
 * - Nakshatra (lunar mansion) segments
 * - Sun and Moon markers
 * - Tithi arc showing Moon-Sun separation
 * - Interactive hover/click states
 * - Educational annotations
 */

import React, { useCallback, useMemo } from 'react';
import {
  EclipticBeltSVGProps,
  EclipticSegment,
  CelestialMarker,
  RelationshipArc,
} from './types/eclipticBelt';
import { getZoneYRange } from './utils/eclipticLayout';

// ============================================================================
// Sub-components
// ============================================================================

/**
 * Renders the Rashi (zodiac) zone with 12 segments
 */
const RashiZone: React.FC<{
  segments: EclipticSegment[];
  yRange: { top: number; bottom: number };
  hoveredElement: string | null;
  selectedElement: string | null;
  onHover: (id: string | null) => void;
  onClick: (id: string | null) => void;
}> = ({ segments, yRange, hoveredElement, selectedElement, onHover, onClick }) => {
  const height = yRange.bottom - yRange.top;

  return (
    <g className="rashi-zone">
      {segments.map((segment) => {
        const isHovered = hoveredElement === segment.id;
        const isSelected = selectedElement === segment.id;
        const width = segment.endX - segment.startX;

        return (
          <g
            key={segment.id}
            className="rashi-segment"
            onMouseEnter={() => onHover(segment.id)}
            onMouseLeave={() => onHover(null)}
            onClick={() => onClick(segment.id)}
            style={{ cursor: 'pointer' }}
          >
            {/* Background rectangle */}
            <rect
              x={segment.startX}
              y={yRange.top}
              width={width}
              height={height}
              fill={segment.color}
              opacity={isHovered || isSelected ? 0.9 : 0.7}
              stroke={isSelected ? '#333' : 'white'}
              strokeWidth={isSelected ? 2 : 1}
            />
            {/* Label */}
            <text
              x={segment.startX + width / 2}
              y={yRange.top + height / 2}
              textAnchor="middle"
              dominantBaseline="middle"
              fontSize={width > 60 ? 14 : 10}
              fontWeight={isHovered ? 'bold' : 'normal'}
              fill="#333"
            >
              {segment.label}
            </text>
            {/* Degree marker */}
            <text
              x={segment.startX + 2}
              y={yRange.top + 12}
              fontSize={8}
              fill="#666"
            >
              {segment.startDegree}°
            </text>
          </g>
        );
      })}
    </g>
  );
};

/**
 * Renders the Nakshatra zone with 27 segments
 */
const NakshatraZone: React.FC<{
  segments: EclipticSegment[];
  yRange: { top: number; bottom: number };
  hoveredElement: string | null;
  selectedElement: string | null;
  onHover: (id: string | null) => void;
  onClick: (id: string | null) => void;
}> = ({ segments, yRange, hoveredElement, selectedElement, onHover, onClick }) => {
  const height = yRange.bottom - yRange.top;

  return (
    <g className="nakshatra-zone">
      {segments.map((segment) => {
        const isHovered = hoveredElement === segment.id;
        const isSelected = selectedElement === segment.id;
        const width = segment.endX - segment.startX;

        return (
          <g
            key={segment.id}
            className="nakshatra-segment"
            onMouseEnter={() => onHover(segment.id)}
            onMouseLeave={() => onHover(null)}
            onClick={() => onClick(segment.id)}
            style={{ cursor: 'pointer' }}
          >
            {/* Background rectangle */}
            <rect
              x={segment.startX}
              y={yRange.top}
              width={width}
              height={height}
              fill={segment.color}
              opacity={isHovered || isSelected ? 0.85 : 0.6}
              stroke={isSelected ? '#333' : 'white'}
              strokeWidth={isSelected ? 2 : 0.5}
            />
            {/* Label - only show if wide enough */}
            {width > 35 && (
              <text
                x={segment.startX + width / 2}
                y={yRange.top + height / 2}
                textAnchor="middle"
                dominantBaseline="middle"
                fontSize={9}
                fontWeight={isHovered ? 'bold' : 'normal'}
                fill="#333"
                style={{ pointerEvents: 'none' }}
              >
                {segment.label}
              </text>
            )}
          </g>
        );
      })}
    </g>
  );
};

/**
 * Renders a celestial body marker (Sun or Moon)
 */
const CelestialMarkerComponent: React.FC<{
  marker: CelestialMarker;
  isHovered: boolean;
  isSelected: boolean;
  onHover: (id: string | null) => void;
  onClick: (id: string | null) => void;
}> = ({ marker, isHovered, isSelected, onHover, onClick }) => {
  const markerId = `marker-${marker.type}`;
  const scale = isHovered || isSelected ? 1.2 : 1;

  return (
    <g
      className={`celestial-marker ${marker.type}`}
      transform={`translate(${marker.position.x}, ${marker.position.y})`}
      onMouseEnter={() => onHover(markerId)}
      onMouseLeave={() => onHover(null)}
      onClick={() => onClick(markerId)}
      style={{ cursor: 'pointer' }}
    >
      {/* Glow effect */}
      {(isHovered || isSelected) && (
        <circle
          r={marker.size * 1.5}
          fill={marker.color}
          opacity={0.3}
        />
      )}

      {/* Main circle */}
      <circle
        r={marker.size * scale}
        fill={marker.color}
        stroke={isSelected ? '#333' : 'white'}
        strokeWidth={isSelected ? 3 : 2}
      />

      {/* Symbol */}
      <text
        textAnchor="middle"
        dominantBaseline="middle"
        fontSize={marker.size * 1.2}
        fill={marker.type === 'sun' ? '#000' : '#444'}
        style={{ pointerEvents: 'none' }}
      >
        {marker.type === 'sun' ? '\u2609' : '\u263D'}
      </text>

      {/* Label below */}
      <text
        y={marker.size + 15}
        textAnchor="middle"
        fontSize={12}
        fontWeight="bold"
        fill="#333"
        style={{ pointerEvents: 'none' }}
      >
        {marker.type === 'sun' ? 'Sun' : 'Moon'}
      </text>

      {/* Longitude */}
      <text
        y={marker.size + 28}
        textAnchor="middle"
        fontSize={10}
        fill="#666"
        style={{ pointerEvents: 'none' }}
      >
        {marker.longitude.toFixed(1)}°
      </text>
    </g>
  );
};

/**
 * Renders the Tithi arc connecting Sun and Moon
 */
const TithiArcComponent: React.FC<{
  arc: RelationshipArc;
  height: number;
  isHovered: boolean;
  isSelected: boolean;
  onHover: (id: string | null) => void;
  onClick: (id: string | null) => void;
}> = ({ arc, height, isHovered, isSelected, onHover, onClick }) => {
  // Calculate arc path
  // We draw a curved path from Sun to Moon
  const startX = Math.min(arc.startX, arc.endX);
  const endX = Math.max(arc.startX, arc.endX);
  const midX = (startX + endX) / 2;

  // Arc height based on angle (larger angle = higher arc)
  const arcHeight = Math.min(height * 0.4, (arc.angle / 360) * height * 0.8);

  // SVG path for the arc
  const pathD = `M ${arc.startX} ${arc.y} Q ${midX} ${arc.y - arcHeight} ${arc.endX} ${arc.y}`;

  return (
    <g
      className="tithi-arc"
      onMouseEnter={() => onHover('tithi-arc')}
      onMouseLeave={() => onHover(null)}
      onClick={() => onClick('tithi-arc')}
      style={{ cursor: 'pointer' }}
    >
      {/* Arc path */}
      <path
        d={pathD}
        fill="none"
        stroke={arc.color}
        strokeWidth={isHovered || isSelected ? 4 : 2}
        strokeDasharray={isHovered ? 'none' : '5,3'}
        opacity={isHovered || isSelected ? 1 : 0.7}
      />

      {/* Arrow at Moon end */}
      <circle
        cx={arc.endX}
        cy={arc.y}
        r={4}
        fill={arc.color}
      />

      {/* Label at arc peak */}
      <text
        x={midX}
        y={arc.y - arcHeight - 8}
        textAnchor="middle"
        fontSize={11}
        fontWeight={isHovered ? 'bold' : 'normal'}
        fill="#333"
        style={{ pointerEvents: 'none' }}
      >
        {arc.label}
      </text>

      {/* Angle indicator */}
      <text
        x={midX}
        y={arc.y - arcHeight - 22}
        textAnchor="middle"
        fontSize={10}
        fill="#666"
        style={{ pointerEvents: 'none' }}
      >
        Moon-Sun: {arc.angle.toFixed(1)}°
      </text>
    </g>
  );
};

/**
 * Renders the annotation panel at the bottom
 */
const AnnotationPanel: React.FC<{
  panchangamSummary: string;
  yRange: { top: number; bottom: number };
  width: number;
  padding: { left: number; right: number };
}> = ({ panchangamSummary, yRange, width, padding }) => {
  const contentWidth = width - padding.left - padding.right;
  const height = yRange.bottom - yRange.top;

  // Parse the summary into lines
  const lines = panchangamSummary.split('\n').filter(line => line.trim());

  return (
    <g className="annotation-panel">
      {/* Background */}
      <rect
        x={padding.left}
        y={yRange.top}
        width={contentWidth}
        height={height}
        fill="#FFF9E6"
        rx={8}
        ry={8}
        stroke="#E0D5B5"
        strokeWidth={1}
      />

      {/* Content - simple text rendering */}
      {lines.slice(0, 4).map((line, index) => (
        <text
          key={index}
          x={padding.left + 15}
          y={yRange.top + 18 + index * 16}
          fontSize={12}
          fill="#333"
        >
          {line.replace(/\*\*/g, '').replace(/\*/g, '')}
        </text>
      ))}
    </g>
  );
};

// ============================================================================
// Main Component
// ============================================================================

export const EclipticBeltSVG: React.FC<EclipticBeltSVGProps> = ({
  dimensions,
  panchangam,
  rashiSegments,
  nakshatraSegments,
  sunMarker,
  moonMarker,
  tithiArc,
  annotations,
  selectedElement,
  onElementSelect,
  onElementHover,
  hoveredElement,
  showLabels = true,
  animationEnabled = true,
}) => {
  // Memoize zone Y ranges
  const rashiYRange = useMemo(
    () => getZoneYRange('rashi', dimensions),
    [dimensions]
  );
  const nakshatraYRange = useMemo(
    () => getZoneYRange('nakshatra', dimensions),
    [dimensions]
  );
  const planetsYRange = useMemo(
    () => getZoneYRange('planets', dimensions),
    [dimensions]
  );
  const tithiYRange = useMemo(
    () => getZoneYRange('tithi', dimensions),
    [dimensions]
  );
  const annotationYRange = useMemo(
    () => getZoneYRange('annotation', dimensions),
    [dimensions]
  );

  // Handlers
  const handleHover = useCallback((id: string | null) => {
    onElementHover(id);
  }, [onElementHover]);

  const handleClick = useCallback((id: string | null) => {
    onElementSelect(selectedElement === id ? null : id);
  }, [onElementSelect, selectedElement]);

  // Generate summary for annotation panel
  const summaryText = useMemo(() => {
    const { tithi, nakshatra, yoga, karana } = panchangam;
    return [
      `Tithi: ${tithi.paksha} ${tithi.name} (${tithi.percentComplete.toFixed(0)}%)`,
      `Nakshatra: ${nakshatra.name} (Pada ${nakshatra.pada})`,
      `Yoga: ${yoga.name} - ${yoga.meaning}`,
      `Karana: ${karana.name}`
    ].join('\n');
  }, [panchangam]);

  return (
    <svg
      width={dimensions.width}
      height={dimensions.height}
      viewBox={`0 0 ${dimensions.width} ${dimensions.height}`}
      className="ecliptic-belt-svg"
      style={{
        background: 'linear-gradient(to bottom, #f8f4e8 0%, #fff9e6 100%)',
        borderRadius: '12px',
      }}
    >
      {/* Definitions for gradients and filters */}
      <defs>
        {/* Glow filter for celestial bodies */}
        <filter id="glow" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="3" result="coloredBlur" />
          <feMerge>
            <feMergeNode in="coloredBlur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>

        {/* Sun gradient */}
        <radialGradient id="sunGradient" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#FFE066" />
          <stop offset="100%" stopColor="#FFB300" />
        </radialGradient>

        {/* Moon gradient */}
        <radialGradient id="moonGradient" cx="30%" cy="30%" r="70%">
          <stop offset="0%" stopColor="#F5F5F5" />
          <stop offset="100%" stopColor="#A0A0A0" />
        </radialGradient>
      </defs>

      {/* Zone labels on the left */}
      <g className="zone-labels">
        <text
          x={dimensions.padding.left - 5}
          y={rashiYRange.top + (rashiYRange.bottom - rashiYRange.top) / 2}
          textAnchor="end"
          dominantBaseline="middle"
          fontSize={10}
          fill="#666"
          transform={`rotate(-90, ${dimensions.padding.left - 5}, ${rashiYRange.top + (rashiYRange.bottom - rashiYRange.top) / 2})`}
        >
          Rashis
        </text>
        <text
          x={dimensions.padding.left - 5}
          y={nakshatraYRange.top + (nakshatraYRange.bottom - nakshatraYRange.top) / 2}
          textAnchor="end"
          dominantBaseline="middle"
          fontSize={10}
          fill="#666"
          transform={`rotate(-90, ${dimensions.padding.left - 5}, ${nakshatraYRange.top + (nakshatraYRange.bottom - nakshatraYRange.top) / 2})`}
        >
          Nakshatras
        </text>
      </g>

      {/* Rashi (Zodiac) Zone */}
      <RashiZone
        segments={rashiSegments}
        yRange={rashiYRange}
        hoveredElement={hoveredElement}
        selectedElement={selectedElement}
        onHover={handleHover}
        onClick={handleClick}
      />

      {/* Nakshatra Zone */}
      <NakshatraZone
        segments={nakshatraSegments}
        yRange={nakshatraYRange}
        hoveredElement={hoveredElement}
        selectedElement={selectedElement}
        onHover={handleHover}
        onClick={handleClick}
      />

      {/* Planet Track Background */}
      <rect
        x={dimensions.padding.left}
        y={planetsYRange.top}
        width={dimensions.width - dimensions.padding.left - dimensions.padding.right}
        height={planetsYRange.bottom - planetsYRange.top}
        fill="#FFF8DC"
        opacity={0.5}
      />

      {/* Center line for planet track */}
      <line
        x1={dimensions.padding.left}
        y1={(planetsYRange.top + planetsYRange.bottom) / 2}
        x2={dimensions.width - dimensions.padding.right}
        y2={(planetsYRange.top + planetsYRange.bottom) / 2}
        stroke="#DDD"
        strokeWidth={1}
        strokeDasharray="5,5"
      />

      {/* Tithi Arc */}
      <TithiArcComponent
        arc={tithiArc}
        height={tithiYRange.bottom - tithiYRange.top}
        isHovered={hoveredElement === 'tithi-arc'}
        isSelected={selectedElement === 'tithi-arc'}
        onHover={handleHover}
        onClick={handleClick}
      />

      {/* Celestial Markers */}
      <CelestialMarkerComponent
        marker={sunMarker}
        isHovered={hoveredElement === 'marker-sun'}
        isSelected={selectedElement === 'marker-sun'}
        onHover={handleHover}
        onClick={handleClick}
      />
      <CelestialMarkerComponent
        marker={moonMarker}
        isHovered={hoveredElement === 'marker-moon'}
        isSelected={selectedElement === 'marker-moon'}
        onHover={handleHover}
        onClick={handleClick}
      />

      {/* Annotation Panel */}
      {showLabels && (
        <AnnotationPanel
          panchangamSummary={summaryText}
          yRange={annotationYRange}
          width={dimensions.width}
          padding={dimensions.padding}
        />
      )}

      {/* Scale/Legend at bottom */}
      <g className="legend">
        <text
          x={dimensions.width / 2}
          y={dimensions.height - 5}
          textAnchor="middle"
          fontSize={10}
          fill="#999"
        >
          0° (Aries) ← Ecliptic Longitude → 360° (Pisces)
        </text>
      </g>
    </svg>
  );
};

export default EclipticBeltSVG;
