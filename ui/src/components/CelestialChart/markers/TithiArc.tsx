/**
 * TithiArc - Arc connecting Sun and Moon showing angular separation
 *
 * Features:
 * - Curved arc from Sun to Moon position
 * - Animated dash pattern
 * - Tithi information on hover
 * - Color indicates paksha (bright/dark fortnight)
 */

import React, { useMemo } from 'react';
import type { TithiArcProps } from '../types';
import { CHART_COLORS } from '../types';
import { normalizeAngle } from '../utils/geometryHelpers';

export const TithiArc: React.FC<TithiArcProps> = ({
  arcInfo,
  isHovered,
  onHover,
  onClick,
}) => {
  const { sunAngle, moonAngle, radius, tithi, arcPath } = arcInfo;

  // Calculate angular separation for display
  const angularSeparation = useMemo(() => {
    let diff = normalizeAngle(moonAngle - sunAngle);
    return diff;
  }, [sunAngle, moonAngle]);

  // Determine color based on paksha
  const arcColor = useMemo(() => {
    return tithi.paksha === 'Shukla'
      ? '#FFC107' // Golden for bright fortnight (waxing)
      : '#5C6BC0'; // Blue-purple for dark fortnight (waning)
  }, [tithi.paksha]);

  // Arc stroke width based on hover
  const strokeWidth = isHovered ? 4 : 3;

  // Animation ID for dash pattern
  const animationId = 'tithi-arc-dash';

  return (
    <g
      className="tithi-arc"
      onMouseEnter={() => onHover(true, { x: 0, y: 0 })}
      onMouseLeave={() => onHover(false, { x: 0, y: 0 })}
      onClick={onClick}
      style={{ cursor: 'pointer' }}
    >
      {/* Animated dash pattern definition */}
      <defs>
        <linearGradient id="tithi-gradient" x1="0%" y1="0%" x2="100%" y2="0%">
          <stop offset="0%" stopColor={CHART_COLORS.sun} />
          <stop offset="100%" stopColor={CHART_COLORS.moon} />
        </linearGradient>
      </defs>

      {/* Background arc (wider, more visible) */}
      <path
        d={arcPath}
        fill="none"
        stroke={arcColor}
        strokeWidth={strokeWidth + 4}
        strokeOpacity={0.2}
        strokeLinecap="round"
        pointerEvents="none"
      />

      {/* Main arc */}
      <path
        d={arcPath}
        fill="none"
        stroke={arcColor}
        strokeWidth={strokeWidth}
        strokeOpacity={isHovered ? 1 : 0.8}
        strokeLinecap="round"
        strokeDasharray={isHovered ? 'none' : '8 4'}
        style={{
          transition: 'all 0.2s ease',
        }}
      />

      {/* Arc glow on hover */}
      {isHovered && (
        <path
          d={arcPath}
          fill="none"
          stroke={arcColor}
          strokeWidth={strokeWidth + 6}
          strokeOpacity={0.3}
          strokeLinecap="round"
          pointerEvents="none"
        />
      )}

      {/* Angular separation indicator at midpoint */}
      {isHovered && (
        <g pointerEvents="none">
          {/* Calculate midpoint of arc for label placement */}
          {(() => {
            const midAngle = sunAngle + angularSeparation / 2;
            const rad = ((90 - midAngle) * Math.PI) / 180;
            const labelRadius = radius + 20;
            // Approximate center calculation
            const cx = 0; // Will be set by parent
            const cy = 0;

            return (
              <>
                {/* Info box positioned relative to arc */}
                <rect
                  x={-50}
                  y={radius + 25}
                  width={100}
                  height={50}
                  rx={6}
                  fill={CHART_COLORS.text}
                  fillOpacity={0.95}
                  transform={`translate(${radius * 0.7}, ${-radius * 0.3})`}
                />
                <text
                  x={0}
                  y={radius + 42}
                  textAnchor="middle"
                  fill="white"
                  fontSize={12}
                  fontWeight="bold"
                  transform={`translate(${radius * 0.7}, ${-radius * 0.3})`}
                >
                  {tithi.name}
                </text>
                <text
                  x={0}
                  y={radius + 58}
                  textAnchor="middle"
                  fill="white"
                  fontSize={10}
                  transform={`translate(${radius * 0.7}, ${-radius * 0.3})`}
                >
                  {tithi.paksha} Paksha • {angularSeparation.toFixed(1)}°
                </text>
              </>
            );
          })()}
        </g>
      )}
    </g>
  );
};

export default TithiArc;
