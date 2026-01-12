/**
 * RashiRing - Outermost ring showing 12 zodiac signs
 *
 * Features:
 * - 12 segments of 30° each
 * - Color-coded by element (Fire/Earth/Air/Water)
 * - Highlights Sun and Moon positions
 * - Labels with Sanskrit and Western names
 */

import React, { useMemo, useCallback } from 'react';
import type { RashiRingProps, Point } from '../types';
import { createArcSegmentPath, getLabelPosition, createRashiBoundariesPath } from '../utils/geometryHelpers';
import { getAllRashis, getRashiColor, isRashiHighlighted } from '../utils/chartCalculations';
import { CHART_COLORS } from '../types';
import type { PanchangamElements } from '../../EclipticBeltVisualization/types/eclipticBelt';

interface ExtendedRashiRingProps extends RashiRingProps {
  panchangam: PanchangamElements;
}

export const RashiRing: React.FC<ExtendedRashiRingProps> = ({
  dimensions,
  hoveredId,
  selectedId,
  onHover,
  onSelect,
  panchangam,
}) => {
  const rashis = useMemo(() => getAllRashis(), []);
  const { innerRadius, outerRadius, labelOffset } = dimensions.rings.rashi;
  const { center } = dimensions;

  // Pre-compute all segment paths
  const segmentPaths = useMemo(() => {
    return rashis.map((rashi) => ({
      rashi,
      path: createArcSegmentPath(
        center,
        innerRadius,
        outerRadius,
        rashi.startDegree,
        rashi.endDegree
      ),
      labelPos: getLabelPosition(
        center,
        labelOffset || (innerRadius + outerRadius) / 2,
        rashi.startDegree,
        rashi.endDegree
      ),
    }));
  }, [rashis, center, innerRadius, outerRadius, labelOffset]);

  // Boundaries path (single element for all dividers)
  const boundariesPath = useMemo(() =>
    createRashiBoundariesPath(center, innerRadius, outerRadius),
    [center, innerRadius, outerRadius]
  );

  const handleMouseEnter = useCallback((id: string, position: Point) => {
    onHover(id, position);
  }, [onHover]);

  const handleMouseLeave = useCallback(() => {
    onHover(null, { x: 0, y: 0 });
  }, [onHover]);

  return (
    <g className="rashi-ring">
      {/* Segment fills */}
      {segmentPaths.map(({ rashi, path, labelPos }) => {
        const id = `rashi-${rashi.number}`;
        const isHovered = hoveredId === id;
        const isSelected = selectedId === id;
        const highlight = isRashiHighlighted(rashi.number - 1, panchangam);

        // Determine segment color
        let fillColor = getRashiColor(rashi.element);
        let opacity = 0.6;

        if (isHovered) {
          fillColor = CHART_COLORS.hover;
          opacity = 0.9;
        } else if (isSelected) {
          fillColor = CHART_COLORS.selected;
          opacity = 0.9;
        } else if (highlight) {
          opacity = 0.85;
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
              data-rashi={rashi.number}
            />

            {/* Highlight indicator for Sun/Moon */}
            {highlight && (
              <path
                d={path}
                fill="none"
                stroke={highlight === 'sun' ? CHART_COLORS.sun :
                       highlight === 'moon' ? CHART_COLORS.moon :
                       CHART_COLORS.sun}
                strokeWidth={highlight === 'both' ? 4 : 3}
                strokeOpacity={0.8}
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
        strokeWidth={1}
        strokeOpacity={0.4}
        pointerEvents="none"
      />

      {/* Outer and inner ring borders */}
      <circle
        cx={center.x}
        cy={center.y}
        r={outerRadius}
        fill="none"
        stroke={CHART_COLORS.text}
        strokeWidth={2}
        strokeOpacity={0.6}
        pointerEvents="none"
      />
      <circle
        cx={center.x}
        cy={center.y}
        r={innerRadius}
        fill="none"
        stroke={CHART_COLORS.text}
        strokeWidth={1}
        strokeOpacity={0.4}
        pointerEvents="none"
      />

      {/* Labels */}
      {segmentPaths.map(({ rashi, labelPos }) => {
        const isHovered = hoveredId === `rashi-${rashi.number}`;

        return (
          <g key={`label-${rashi.number}`} pointerEvents="none">
            {/* Symbol */}
            <text
              x={labelPos.x}
              y={labelPos.y - 8}
              textAnchor="middle"
              dominantBaseline="middle"
              fill={isHovered ? CHART_COLORS.text : CHART_COLORS.textLight}
              fontSize={isHovered ? 18 : 16}
              fontWeight={isHovered ? 'bold' : 'normal'}
            >
              {rashi.symbol}
            </text>
            {/* Name (only show on larger displays or hover) */}
            <text
              x={labelPos.x}
              y={labelPos.y + 10}
              textAnchor="middle"
              dominantBaseline="middle"
              fill={CHART_COLORS.textLight}
              fontSize={10}
              opacity={isHovered ? 1 : 0.7}
            >
              {rashi.name}
            </text>
          </g>
        );
      })}
    </g>
  );
};

export default RashiRing;
