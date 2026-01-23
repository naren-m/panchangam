/**
 * ChartTooltip - Floating tooltip for chart element details
 *
 * Features:
 * - Positioned near hovered element
 * - Adapts content based on element type
 * - Smooth appearance animation
 */

import React, { useMemo } from 'react';
import type { ChartTooltipProps, PadaInfo } from '../types';
import type { RashiInfo, NakshatraInfo, TithiInfo } from '../../EclipticBeltVisualization/types/eclipticBelt';
import {
  getRashiTooltipContent,
  getNakshatraTooltipContent,
  getPadaTooltipContent,
  getTithiTooltipContent
} from '../utils/chartCalculations';

export const ChartTooltip: React.FC<ChartTooltipProps> = ({
  element,
  chartDimensions,
}) => {
  const { type, data, position } = element;

  // Calculate tooltip position to keep it within bounds
  const tooltipPosition = useMemo(() => {
    const padding = 20;
    const tooltipWidth = 200;
    const tooltipHeight = 100;

    let x = position.x + 15;
    let y = position.y - 15;

    // Keep within chart bounds
    if (x + tooltipWidth > chartDimensions.size - padding) {
      x = position.x - tooltipWidth - 15;
    }
    if (y + tooltipHeight > chartDimensions.size - padding) {
      y = chartDimensions.size - tooltipHeight - padding;
    }
    if (y < padding) {
      y = padding;
    }
    if (x < padding) {
      x = padding;
    }

    return { x, y };
  }, [position, chartDimensions.size]);

  // Generate content based on element type
  const content = useMemo(() => {
    switch (type) {
      case 'rashi':
        return getRashiTooltipContent(data as RashiInfo);
      case 'nakshatra':
        return getNakshatraTooltipContent(data as NakshatraInfo);
      case 'pada':
        return getPadaTooltipContent(data as PadaInfo);
      case 'tithi':
        return getTithiTooltipContent(data as TithiInfo);
      case 'sun':
      case 'moon':
        // Celestial bodies have their own inline tooltips
        return null;
      default:
        return null;
    }
  }, [type, data]);

  if (!content) return null;

  const lines = content.split('\n');

  return (
    <g
      className="chart-tooltip"
      transform={`translate(${tooltipPosition.x}, ${tooltipPosition.y})`}
      pointerEvents="none"
      style={{
        opacity: 1,
        transition: 'opacity 0.15s ease-in-out',
      }}
    >
      {/* Tooltip background */}
      <rect
        x={0}
        y={0}
        width={180}
        height={lines.length * 18 + 16}
        rx={6}
        fill="#37474F"
        fillOpacity={0.95}
      />

      {/* Tooltip content */}
      {lines.map((line, index) => (
        <text
          key={index}
          x={10}
          y={20 + index * 18}
          fill="white"
          fontSize={index === 0 ? 13 : 11}
          fontWeight={index === 0 ? 'bold' : 'normal'}
          opacity={index === 0 ? 1 : 0.85}
        >
          {line}
        </text>
      ))}
    </g>
  );
};

export default ChartTooltip;
