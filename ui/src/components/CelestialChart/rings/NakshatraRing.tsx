/**
 * NakshatraRing - Middle ring showing 27 lunar mansions
 *
 * Features:
 * - 27 segments of 13°20' each
 * - Highlights current Moon nakshatra
 * - Labels with nakshatra names
 */

import React, { useMemo, useCallback } from 'react';
import type { NakshatraRingProps, Point } from '../types';
import {
  createArcSegmentPath,
  getLabelPosition,
  createNakshatraBoundariesPath
} from '../utils/geometryHelpers';
import { getAllNakshatras, isNakshatraHighlighted } from '../utils/chartCalculations';
import { CHART_COLORS } from '../types';
import type { PanchangamElements } from '../../EclipticBeltVisualization/types/eclipticBelt';

interface ExtendedNakshatraRingProps extends NakshatraRingProps {
  panchangam: PanchangamElements;
}

export const NakshatraRing: React.FC<ExtendedNakshatraRingProps> = ({
  dimensions,
  hoveredId,
  selectedId,
  onHover,
  onSelect,
  currentNakshatra,
  panchangam,
}) => {
  const nakshatras = useMemo(() => getAllNakshatras(), []);
  const { innerRadius, outerRadius, labelOffset } = dimensions.rings.nakshatra;
  const { center } = dimensions;

  // Pre-compute all segment paths
  const segmentPaths = useMemo(() => {
    return nakshatras.map((nakshatra) => ({
      nakshatra,
      path: createArcSegmentPath(
        center,
        innerRadius,
        outerRadius,
        nakshatra.startDegree,
        nakshatra.endDegree
      ),
      labelPos: getLabelPosition(
        center,
        labelOffset || (innerRadius + outerRadius) / 2,
        nakshatra.startDegree,
        nakshatra.endDegree
      ),
    }));
  }, [nakshatras, center, innerRadius, outerRadius, labelOffset]);

  // Boundaries path (single element for all dividers)
  const boundariesPath = useMemo(() =>
    createNakshatraBoundariesPath(center, innerRadius, outerRadius),
    [center, innerRadius, outerRadius]
  );

  const handleMouseEnter = useCallback((id: string, position: Point) => {
    onHover(id, position);
  }, [onHover]);

  const handleMouseLeave = useCallback(() => {
    onHover(null, { x: 0, y: 0 });
  }, [onHover]);

  return (
    <g className="nakshatra-ring">
      {/* Segment fills */}
      {segmentPaths.map(({ nakshatra, path, labelPos }) => {
        const id = `nakshatra-${nakshatra.number}`;
        const isHovered = hoveredId === id;
        const isSelected = selectedId === id;
        const isHighlighted = isNakshatraHighlighted(nakshatra.number - 1, panchangam);

        // Alternate colors for visual distinction
        const baseColor = nakshatra.number % 2 === 0
          ? CHART_COLORS.nakshatraRing
          : '#FFA726'; // Slightly different shade

        let fillColor = baseColor;
        let opacity = 0.5;

        if (isHovered) {
          fillColor = CHART_COLORS.hover;
          opacity = 0.9;
        } else if (isSelected) {
          fillColor = CHART_COLORS.selected;
          opacity = 0.9;
        } else if (isHighlighted) {
          fillColor = CHART_COLORS.moon;
          opacity = 0.8;
        }

        return (
          <g key={id}>
            {/* Segment fill */}
            <path
              d={path}
              fill={fillColor}
              fillOpacity={opacity}
              stroke="none"
              style={{ cursor: 'pointer' }}
              onMouseEnter={() => handleMouseEnter(id, labelPos)}
              onMouseLeave={handleMouseLeave}
              onClick={() => onSelect(id)}
              data-nakshatra={nakshatra.number}
            />

            {/* Current nakshatra highlight border */}
            {isHighlighted && (
              <path
                d={path}
                fill="none"
                stroke={CHART_COLORS.moon}
                strokeWidth={3}
                strokeOpacity={0.9}
                pointerEvents="none"
              />
            )}
          </g>
        );
      })}

      {/* Segment boundaries (single path for performance) */}
      <path
        d={boundariesPath}
        fill="none"
        stroke={CHART_COLORS.text}
        strokeWidth={0.5}
        strokeOpacity={0.3}
        pointerEvents="none"
      />

      {/* Ring borders */}
      <circle
        cx={center.x}
        cy={center.y}
        r={outerRadius}
        fill="none"
        stroke={CHART_COLORS.text}
        strokeWidth={1}
        strokeOpacity={0.4}
        pointerEvents="none"
      />
      <circle
        cx={center.x}
        cy={center.y}
        r={innerRadius}
        fill="none"
        stroke={CHART_COLORS.text}
        strokeWidth={1}
        strokeOpacity={0.3}
        pointerEvents="none"
      />

      {/* Labels - only show for hovered or current nakshatra to avoid crowding */}
      {segmentPaths.map(({ nakshatra, labelPos }) => {
        const id = `nakshatra-${nakshatra.number}`;
        const isHovered = hoveredId === id;
        const isHighlighted = isNakshatraHighlighted(nakshatra.number - 1, panchangam);

        // Only show labels for hovered or highlighted nakshatras
        if (!isHovered && !isHighlighted) return null;

        return (
          <g key={`label-${nakshatra.number}`} pointerEvents="none">
            <text
              x={labelPos.x}
              y={labelPos.y}
              textAnchor="middle"
              dominantBaseline="middle"
              fill={CHART_COLORS.text}
              fontSize={isHovered ? 11 : 9}
              fontWeight={isHovered || isHighlighted ? 'bold' : 'normal'}
            >
              {nakshatra.name}
            </text>
          </g>
        );
      })}
    </g>
  );
};

export default NakshatraRing;
