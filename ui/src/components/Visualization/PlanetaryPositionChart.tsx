import React, { useMemo } from 'react';

interface PlanetaryPosition {
  name: string;
  longitude: number;
  latitude: number;
  distance: number;
  speed: number;
  isRetrograde?: boolean;
}

interface PlanetaryPositionChartProps {
  positions: PlanetaryPosition[];
  showRetrograde?: boolean;
  size?: number;
  interactive?: boolean;
}

/**
 * PlanetaryPositionChart - Visualizes planetary positions in the zodiac wheel
 *
 * @component
 */
export const PlanetaryPositionChart: React.FC<PlanetaryPositionChartProps> = ({
  positions,
  showRetrograde = true,
  size = 400,
  interactive = true,
}) => {
  const [hoveredPlanet, setHoveredPlanet] = React.useState<string | null>(null);

  // Zodiac signs
  const zodiacSigns = [
    { name: 'Aries', symbol: '♈', startDeg: 0 },
    { name: 'Taurus', symbol: '♉', startDeg: 30 },
    { name: 'Gemini', symbol: '♊', startDeg: 60 },
    { name: 'Cancer', symbol: '♋', startDeg: 90 },
    { name: 'Leo', symbol: '♌', startDeg: 120 },
    { name: 'Virgo', symbol: '♍', startDeg: 150 },
    { name: 'Libra', symbol: '♎', startDeg: 180 },
    { name: 'Scorpio', symbol: '♏', startDeg: 210 },
    { name: 'Sagittarius', symbol: '♐', startDeg: 240 },
    { name: 'Capricorn', symbol: '♑', startDeg: 270 },
    { name: 'Aquarius', symbol: '♒', startDeg: 300 },
    { name: 'Pisces', symbol: '♓', startDeg: 330 },
  ];

  // Planet symbols and colors
  const planetConfig: Record<string, { symbol: string; color: string }> = {
    sun: { symbol: '☉', color: '#FDB813' },
    moon: { symbol: '☽', color: '#C0C0C0' },
    mercury: { symbol: '☿', color: '#87CEEB' },
    venus: { symbol: '♀', color: '#FFB6C1' },
    mars: { symbol: '♂', color: '#FF4500' },
    jupiter: { symbol: '♃', color: '#FFA500' },
    saturn: { symbol: '♄', color: '#8B7355' },
    uranus: { symbol: '♅', color: '#4FD5D6' },
    neptune: { symbol: '♆', color: '#4169E1' },
    pluto: { symbol: '♇', color: '#8B4513' },
  };

  const center = size / 2;
  const radius = size * 0.35;

  // Convert longitude to SVG coordinates
  const polarToCartesian = (degrees: number, r: number) => {
    // Adjust for SVG coordinate system (0 degrees at top, clockwise)
    const radians = ((degrees - 90) * Math.PI) / 180;
    return {
      x: center + r * Math.cos(radians),
      y: center + r * Math.sin(radians),
    };
  };

  // Calculate planet positions
  const planetPositions = useMemo(() => {
    return positions.map((planet) => {
      const pos = polarToCartesian(planet.longitude, radius);
      const config = planetConfig[planet.name.toLowerCase()] || {
        symbol: '●',
        color: '#666',
      };
      return {
        ...planet,
        ...pos,
        ...config,
      };
    });
  }, [positions, radius, center]);

  return (
    <div className="relative">
      <svg
        width={size}
        height={size}
        className="border border-gray-300 rounded-lg bg-gradient-to-br from-indigo-50 to-purple-50"
        role="img"
        aria-label="Planetary position chart showing zodiac wheel with planet positions"
      >
        {/* Zodiac wheel background */}
        <circle
          cx={center}
          cy={center}
          r={radius}
          fill="none"
          stroke="#ddd"
          strokeWidth="2"
        />

        {/* Zodiac sign divisions */}
        {zodiacSigns.map((sign, index) => {
          const startPos = polarToCartesian(sign.startDeg, radius * 0.9);
          const endPos = polarToCartesian(sign.startDeg, radius * 1.1);
          const labelPos = polarToCartesian(sign.startDeg + 15, radius * 1.25);

          return (
            <g key={sign.name}>
              {/* Division line */}
              <line
                x1={startPos.x}
                y1={startPos.y}
                x2={endPos.x}
                y2={endPos.y}
                stroke="#ccc"
                strokeWidth="1"
              />
              {/* Sign symbol */}
              <text
                x={labelPos.x}
                y={labelPos.y}
                textAnchor="middle"
                dominantBaseline="middle"
                fontSize="20"
                fill="#666"
                aria-label={sign.name}
              >
                {sign.symbol}
              </text>
            </g>
          );
        })}

        {/* Planet positions */}
        {planetPositions.map((planet) => (
          <g
            key={planet.name}
            onMouseEnter={() => interactive && setHoveredPlanet(planet.name)}
            onMouseLeave={() => interactive && setHoveredPlanet(null)}
            style={{ cursor: interactive ? 'pointer' : 'default' }}
          >
            {/* Planet marker */}
            <circle
              cx={planet.x}
              cy={planet.y}
              r={hoveredPlanet === planet.name ? 16 : 12}
              fill={planet.color}
              stroke="#fff"
              strokeWidth="2"
              className="transition-all duration-200"
            />
            {/* Planet symbol */}
            <text
              x={planet.x}
              y={planet.y}
              textAnchor="middle"
              dominantBaseline="middle"
              fontSize="16"
              fill="#fff"
              fontWeight="bold"
              pointerEvents="none"
              aria-label={`${planet.name} at ${planet.longitude.toFixed(2)} degrees`}
            >
              {planet.symbol}
            </text>

            {/* Retrograde indicator */}
            {showRetrograde && planet.isRetrograde && (
              <text
                x={planet.x + 18}
                y={planet.y - 18}
                fontSize="12"
                fill="#ff0000"
                fontWeight="bold"
                aria-label={`${planet.name} is retrograde`}
              >
                ℞
              </text>
            )}
          </g>
        ))}
      </svg>

      {/* Legend / Hover info */}
      {hoveredPlanet && (
        <div className="absolute top-2 left-2 bg-white p-3 rounded-lg shadow-lg border border-gray-200 z-10">
          {(() => {
            const planet = planetPositions.find((p) => p.name === hoveredPlanet);
            if (!planet) return null;

            const signIndex = Math.floor(planet.longitude / 30);
            const signDegree = planet.longitude % 30;

            return (
              <div className="text-sm">
                <div className="font-bold text-gray-900">{planet.name}</div>
                <div className="text-gray-600">
                  {zodiacSigns[signIndex].name} {signDegree.toFixed(2)}°
                </div>
                <div className="text-gray-500 text-xs">
                  Longitude: {planet.longitude.toFixed(2)}°
                </div>
                {planet.isRetrograde && (
                  <div className="text-red-600 text-xs font-semibold mt-1">
                    Retrograde
                  </div>
                )}
              </div>
            );
          })()}
        </div>
      )}
    </div>
  );
};

export default PlanetaryPositionChart;
