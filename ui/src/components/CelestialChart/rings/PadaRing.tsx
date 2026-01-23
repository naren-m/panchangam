/**
 * PadaRing - Innermost ring showing 108 padas
 *
 * Features:
 * - 108 segments of 3°20' each
 * - Optimized rendering: static boundaries + dynamic hover overlay
 * - Highlights current Moon pada
 * - Shows navamsha on hover
 */

import React, { useMemo, useCallback, useState } from 'react';
import type { PadaRingProps, Point, PadaInfo } from '../types';
import {
  createArcSegmentPath,
  createPadaBoundariesPath,
  cartesianToLongitude,
  distanceFromCenter,
  normalizeAngle
} from '../utils/geometryHelpers';
import { getAllPadas, isPadaHighlighted, getPadaAtLongitude } from '../utils/chartCalculations';
import { CHART_COLORS, DEGREES_PER_PADA } from '../types';
import type { PanchangamElements } from '../../EclipticBeltVisualization/types/eclipticBelt';

interface ExtendedPadaRingProps extends PadaRingProps {
  panchangam: PanchangamElements;
}

export const PadaRing: React.FC<ExtendedPadaRingProps> = ({
  dimensions,
  hoveredId,
  selectedId,
  onHover,
  onSelect,
  currentPada,
  panchangam,
}) => {
  const padas = useMemo(() => getAllPadas(), []);
  const { innerRadius, outerRadius } = dimensions.rings.pada;
  const { center } = dimensions;

  // Track internally hovered pada for dynamic overlay
  const [hoveredPadaIndex, setHoveredPadaIndex] = useState<number | null>(null);

  // Static boundaries path (single SVG element for all 108 dividers)
  const boundariesPath = useMemo(() =>
    createPadaBoundariesPath(center, innerRadius, outerRadius),
    [center, innerRadius, outerRadius]
  );

  // Current pada (Moon's position) highlight path
  const currentPadaPath = useMemo(() => {
    const moonPadaIndex = Math.floor(
      normalizeAngle(panchangam.moonPosition.longitude) / DEGREES_PER_PADA
    );
    const pada = padas[moonPadaIndex];
    return createArcSegmentPath(
      center,
      innerRadius,
      outerRadius,
      pada.startDegree,
      pada.endDegree
    );
  }, [center, innerRadius, outerRadius, panchangam.moonPosition.longitude, padas]);

  // Hovered pada path (dynamic overlay)
  const hoveredPadaPath = useMemo(() => {
    if (hoveredPadaIndex === null) return null;
    const pada = padas[hoveredPadaIndex];
    return createArcSegmentPath(
      center,
      innerRadius,
      outerRadius,
      pada.startDegree,
      pada.endDegree
    );
  }, [hoveredPadaIndex, center, innerRadius, outerRadius, padas]);

  // Selected pada path
  const selectedPadaPath = useMemo(() => {
    if (!selectedId || !selectedId.startsWith('pada-')) return null;
    const padaNum = parseInt(selectedId.replace('pada-', ''), 10);
    const pada = padas[padaNum - 1];
    if (!pada) return null;
    return createArcSegmentPath(
      center,
      innerRadius,
      outerRadius,
      pada.startDegree,
      pada.endDegree
    );
  }, [selectedId, center, innerRadius, outerRadius, padas]);

  // Event delegation handler for the entire ring
  const handleMouseMove = useCallback((e: React.MouseEvent<SVGCircleElement>) => {
    const svg = e.currentTarget.ownerSVGElement;
    if (!svg) return;

    const pt = svg.createSVGPoint();
    pt.x = e.clientX;
    pt.y = e.clientY;
    const cursorPt = pt.matrixTransform(svg.getScreenCTM()?.inverse());

    // Check if cursor is within pada ring
    const distance = distanceFromCenter({ x: cursorPt.x, y: cursorPt.y }, center);
    if (distance < innerRadius || distance > outerRadius) {
      setHoveredPadaIndex(null);
      onHover(null, { x: 0, y: 0 });
      return;
    }

    // Calculate which pada the cursor is over
    const longitude = cartesianToLongitude({ x: cursorPt.x, y: cursorPt.y }, center);
    const padaIndex = Math.floor(normalizeAngle(longitude) / DEGREES_PER_PADA);

    if (padaIndex !== hoveredPadaIndex) {
      setHoveredPadaIndex(padaIndex);
      const pada = padas[padaIndex];
      onHover(`pada-${pada.number}`, { x: cursorPt.x, y: cursorPt.y });
    }
  }, [center, innerRadius, outerRadius, hoveredPadaIndex, padas, onHover]);

  const handleMouseLeave = useCallback(() => {
    setHoveredPadaIndex(null);
    onHover(null, { x: 0, y: 0 });
  }, [onHover]);

  const handleClick = useCallback((e: React.MouseEvent<SVGCircleElement>) => {
    if (hoveredPadaIndex !== null) {
      onSelect(`pada-${padas[hoveredPadaIndex].number}`);
    }
  }, [hoveredPadaIndex, padas, onSelect]);

  // Get current pada index for highlight
  const moonPadaIndex = useMemo(() =>
    Math.floor(normalizeAngle(panchangam.moonPosition.longitude) / DEGREES_PER_PADA),
    [panchangam.moonPosition.longitude]
  );

  return (
    <g className="pada-ring">
      {/* Base ring fill */}
      <circle
        cx={center.x}
        cy={center.y}
        r={(innerRadius + outerRadius) / 2}
        fill="none"
        stroke={CHART_COLORS.padaRing}
        strokeWidth={outerRadius - innerRadius}
        strokeOpacity={0.3}
        pointerEvents="none"
      />

      {/* Current pada highlight (Moon's position) */}
      <path
        d={currentPadaPath}
        fill={CHART_COLORS.moon}
        fillOpacity={0.7}
        stroke={CHART_COLORS.moon}
        strokeWidth={2}
        pointerEvents="none"
      />

      {/* Selected pada highlight */}
      {selectedPadaPath && (
        <path
          d={selectedPadaPath}
          fill={CHART_COLORS.selected}
          fillOpacity={0.8}
          stroke={CHART_COLORS.selected}
          strokeWidth={2}
          pointerEvents="none"
        />
      )}

      {/* Hovered pada highlight (dynamic overlay) */}
      {hoveredPadaPath && hoveredPadaIndex !== moonPadaIndex && (
        <path
          d={hoveredPadaPath}
          fill={CHART_COLORS.hover}
          fillOpacity={0.8}
          stroke={CHART_COLORS.hover}
          strokeWidth={1}
          pointerEvents="none"
        />
      )}

      {/* Static boundaries (single path for all 108 dividers) */}
      <path
        d={boundariesPath}
        fill="none"
        stroke={CHART_COLORS.text}
        strokeWidth={0.3}
        strokeOpacity={0.2}
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
        strokeOpacity={0.3}
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

      {/* Invisible interaction layer */}
      <circle
        cx={center.x}
        cy={center.y}
        r={(innerRadius + outerRadius) / 2}
        fill="transparent"
        stroke="transparent"
        strokeWidth={outerRadius - innerRadius}
        style={{ cursor: 'pointer' }}
        onMouseMove={handleMouseMove}
        onMouseLeave={handleMouseLeave}
        onClick={handleClick}
      />

      {/* Current pada label */}
      {moonPadaIndex !== null && (
        <text
          x={center.x}
          y={center.y + innerRadius - 15}
          textAnchor="middle"
          dominantBaseline="middle"
          fill={CHART_COLORS.text}
          fontSize={8}
          fontWeight="bold"
          pointerEvents="none"
        >
          Pada {moonPadaIndex + 1}
        </text>
      )}
    </g>
  );
};

export default PadaRing;
