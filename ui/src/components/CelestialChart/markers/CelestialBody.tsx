/**
 * CelestialBody - Sun and Moon markers on the celestial orbit
 *
 * Features:
 * - Animated glow effect
 * - Moon phase visualization
 * - Hover and click interactions
 */

import React, { useMemo } from 'react';
import type { CelestialBodyProps } from '../types';
import { CHART_COLORS } from '../types';

export const CelestialBody: React.FC<CelestialBodyProps> = ({
  body,
  isHovered,
  onHover,
  onClick,
}) => {
  const { type, position, symbol, color, size, label } = body;

  // Generate unique IDs for gradients/filters
  const gradientId = `${type}-gradient`;
  const glowId = `${type}-glow`;
  const shadowId = `${type}-shadow`;

  // Sun-specific rendering
  const renderSun = useMemo(() => (
    <g>
      {/* Glow effect */}
      <defs>
        <radialGradient id={gradientId} cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#FFEB3B" stopOpacity="1" />
          <stop offset="60%" stopColor="#FFC107" stopOpacity="0.8" />
          <stop offset="100%" stopColor="#FF9800" stopOpacity="0.6" />
        </radialGradient>
        <filter id={glowId} x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="3" result="blur" />
          <feComposite in="SourceGraphic" in2="blur" operator="over" />
        </filter>
      </defs>

      {/* Outer glow */}
      <circle
        cx={position.x}
        cy={position.y}
        r={size * 1.3}
        fill="#FFC107"
        fillOpacity={0.3}
        pointerEvents="none"
      />

      {/* Sun body */}
      <circle
        cx={position.x}
        cy={position.y}
        r={size / 2}
        fill={`url(#${gradientId})`}
        stroke="#FF9800"
        strokeWidth={2}
        filter={isHovered ? `url(#${glowId})` : undefined}
        style={{
          cursor: 'pointer',
          transition: 'all 0.2s ease',
          transform: isHovered ? 'scale(1.1)' : 'scale(1)',
          transformOrigin: `${position.x}px ${position.y}px`,
        }}
      />

      {/* Sun rays (decorative) */}
      {[0, 45, 90, 135, 180, 225, 270, 315].map((angle) => {
        const rad = (angle * Math.PI) / 180;
        const innerR = size / 2 + 2;
        const outerR = size / 2 + 8;
        const x1 = position.x + innerR * Math.cos(rad);
        const y1 = position.y + innerR * Math.sin(rad);
        const x2 = position.x + outerR * Math.cos(rad);
        const y2 = position.y + outerR * Math.sin(rad);

        return (
          <line
            key={angle}
            x1={x1}
            y1={y1}
            x2={x2}
            y2={y2}
            stroke="#FFC107"
            strokeWidth={2}
            strokeLinecap="round"
            opacity={0.7}
            pointerEvents="none"
          />
        );
      })}

      {/* Sun symbol */}
      <text
        x={position.x}
        y={position.y}
        textAnchor="middle"
        dominantBaseline="central"
        fontSize={size * 0.6}
        fill="#795548"
        fontWeight="bold"
        pointerEvents="none"
      >
        ☉
      </text>
    </g>
  ), [position, size, isHovered, gradientId, glowId]);

  // Moon-specific rendering with phase
  const renderMoon = useMemo(() => (
    <g>
      <defs>
        <radialGradient id={gradientId} cx="30%" cy="30%" r="70%">
          <stop offset="0%" stopColor="#FFFFFF" stopOpacity="1" />
          <stop offset="70%" stopColor="#E0E0E0" stopOpacity="0.9" />
          <stop offset="100%" stopColor="#BDBDBD" stopOpacity="0.8" />
        </radialGradient>
        <filter id={shadowId} x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="1" dy="1" stdDeviation="2" floodColor="#000" floodOpacity="0.3" />
        </filter>
      </defs>

      {/* Outer glow */}
      <circle
        cx={position.x}
        cy={position.y}
        r={size * 1.2}
        fill="#B0BEC5"
        fillOpacity={0.2}
        pointerEvents="none"
      />

      {/* Moon body */}
      <circle
        cx={position.x}
        cy={position.y}
        r={size / 2}
        fill={`url(#${gradientId})`}
        stroke="#90A4AE"
        strokeWidth={1.5}
        filter={`url(#${shadowId})`}
        style={{
          cursor: 'pointer',
          transition: 'all 0.2s ease',
          transform: isHovered ? 'scale(1.1)' : 'scale(1)',
          transformOrigin: `${position.x}px ${position.y}px`,
        }}
      />

      {/* Moon phase overlay - simplified crescent */}
      <text
        x={position.x}
        y={position.y}
        textAnchor="middle"
        dominantBaseline="central"
        fontSize={size * 0.8}
        fill="#546E7A"
        pointerEvents="none"
      >
        {symbol}
      </text>
    </g>
  ), [position, size, symbol, isHovered, gradientId, shadowId]);

  return (
    <g
      className={`celestial-body celestial-body-${type}`}
      onMouseEnter={() => onHover(true, position)}
      onMouseLeave={() => onHover(false, position)}
      onClick={onClick}
      role="button"
      aria-label={label}
    >
      {type === 'sun' ? renderSun : renderMoon}

      {/* Hover label */}
      {isHovered && (
        <g pointerEvents="none">
          <rect
            x={position.x - 60}
            y={position.y - size - 30}
            width={120}
            height={22}
            rx={4}
            fill={CHART_COLORS.text}
            fillOpacity={0.9}
          />
          <text
            x={position.x}
            y={position.y - size - 15}
            textAnchor="middle"
            dominantBaseline="middle"
            fontSize={11}
            fill="white"
            fontWeight="500"
          >
            {label}
          </text>
        </g>
      )}
    </g>
  );
};

export default CelestialBody;
