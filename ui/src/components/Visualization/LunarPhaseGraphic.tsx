import React from 'react';

interface LunarPhaseGraphicProps {
  phase: number; // 0-1 (0 = new moon, 0.5 = full moon)
  illumination: number; // 0-100
  size?: number;
  showLabel?: boolean;
  showPercentage?: boolean;
}

/**
 * LunarPhaseGraphic - Visualizes the current moon phase
 *
 * @component
 */
export const LunarPhaseGraphic: React.FC<LunarPhaseGraphicProps> = ({
  phase,
  illumination,
  size = 120,
  showLabel = true,
  showPercentage = true,
}) => {
  // Determine phase name
  const getPhaseInfo = (phase: number) => {
    if (phase < 0.03 || phase > 0.97) {
      return { name: 'New Moon', icon: '🌑' };
    } else if (phase < 0.22) {
      return { name: 'Waxing Crescent', icon: '🌒' };
    } else if (phase < 0.28) {
      return { name: 'First Quarter', icon: '🌓' };
    } else if (phase < 0.47) {
      return { name: 'Waxing Gibbous', icon: '🌔' };
    } else if (phase < 0.53) {
      return { name: 'Full Moon', icon: '🌕' };
    } else if (phase < 0.72) {
      return { name: 'Waning Gibbous', icon: '🌖' };
    } else if (phase < 0.78) {
      return { name: 'Last Quarter', icon: '🌗' };
    } else {
      return { name: 'Waning Crescent', icon: '🌘' };
    }
  };

  const phaseInfo = getPhaseInfo(phase);

  // Calculate the terminator position for realistic rendering
  const renderMoonPhase = () => {
    const radius = size / 2;
    const center = size / 2;

    // The phase determines the position of the terminator (light/shadow boundary)
    // phase 0 = new moon (all dark), 0.5 = full moon (all light)
    // We need to draw an ellipse for the terminator

    // For waxing phases (0 to 0.5), the lit portion grows from right to left
    // For waning phases (0.5 to 1), the lit portion shrinks from left to right

    const isWaxing = phase <= 0.5;
    const normalizedPhase = isWaxing ? phase * 2 : (1 - phase) * 2;

    // Calculate terminator curvature (-1 = full dark, 0 = half, 1 = full light)
    const terminator = (normalizedPhase - 0.5) * 2;

    // Create the moon shape with SVG
    return (
      <svg width={size} height={size} className="drop-shadow-lg">
        <defs>
          <radialGradient id="moonGradient">
            <stop offset="0%" stopColor="#f0f0f0" />
            <stop offset="100%" stopColor="#d0d0d0" />
          </radialGradient>
          <radialGradient id="shadowGradient">
            <stop offset="0%" stopColor="#3a3a3a" />
            <stop offset="100%" stopColor="#1a1a1a" />
          </radialGradient>
          <filter id="moonShadow">
            <feGaussianBlur in="SourceAlpha" stdDeviation="2" />
            <feOffset dx="2" dy="2" result="offsetblur" />
            <feMerge>
              <feMergeNode />
              <feMergeNode in="SourceGraphic" />
            </feMerge>
          </filter>
        </defs>

        {/* Outer circle (moon disc) */}
        <circle
          cx={center}
          cy={center}
          r={radius - 2}
          fill="url(#moonGradient)"
          stroke="#999"
          strokeWidth="1"
          filter="url(#moonShadow)"
        />

        {/* Dark portion */}
        <clipPath id="darkClip">
          <circle cx={center} cy={center} r={radius - 2} />
        </clipPath>

        {isWaxing ? (
          // Waxing: darkness on the left
          <ellipse
            cx={center - terminator * radius}
            cy={center}
            rx={Math.abs(terminator) * radius}
            ry={radius}
            fill="url(#shadowGradient)"
            clipPath="url(#darkClip)"
          />
        ) : (
          // Waning: darkness on the right
          <ellipse
            cx={center + terminator * radius}
            cy={center}
            rx={Math.abs(terminator) * radius}
            ry={radius}
            fill="url(#shadowGradient)"
            clipPath="url(#darkClip)"
          />
        )}
      </svg>
    );
  };

  return (
    <div className="flex flex-col items-center gap-2">
      <div
        className="relative"
        role="img"
        aria-label={`Moon phase: ${phaseInfo.name}, ${illumination.toFixed(0)}% illuminated`}
      >
        {renderMoonPhase()}

        {showPercentage && (
          <div className="absolute bottom-0 right-0 bg-black bg-opacity-70 text-white text-xs px-2 py-1 rounded-full">
            {illumination.toFixed(0)}%
          </div>
        )}
      </div>

      {showLabel && (
        <div className="text-center">
          <div className="text-2xl" aria-hidden="true">
            {phaseInfo.icon}
          </div>
          <div className="text-sm font-medium text-gray-700">
            {phaseInfo.name}
          </div>
        </div>
      )}
    </div>
  );
};

interface MuhurtaQualityIndicatorProps {
  quality: 'auspicious' | 'inauspicious' | 'neutral' | 'mixed';
  label: string;
  description?: string;
  compact?: boolean;
}

/**
 * MuhurtaQualityIndicator - Visual indicator for muhurta (auspicious timing) quality
 *
 * @component
 */
export const MuhurtaQualityIndicator: React.FC<MuhurtaQualityIndicatorProps> = ({
  quality,
  label,
  description,
  compact = false,
}) => {
  const qualityConfig = {
    auspicious: {
      color: 'bg-green-500',
      borderColor: 'border-green-600',
      textColor: 'text-green-700',
      icon: '✓',
      bgLight: 'bg-green-50',
    },
    inauspicious: {
      color: 'bg-red-500',
      borderColor: 'border-red-600',
      textColor: 'text-red-700',
      icon: '✗',
      bgLight: 'bg-red-50',
    },
    neutral: {
      color: 'bg-gray-400',
      borderColor: 'border-gray-500',
      textColor: 'text-gray-700',
      icon: '−',
      bgLight: 'bg-gray-50',
    },
    mixed: {
      color: 'bg-yellow-500',
      borderColor: 'border-yellow-600',
      textColor: 'text-yellow-700',
      icon: '≈',
      bgLight: 'bg-yellow-50',
    },
  };

  const config = qualityConfig[quality];

  if (compact) {
    return (
      <div
        className={`inline-flex items-center gap-2 px-3 py-1 rounded-full ${config.bgLight} ${config.borderColor} border`}
        role="status"
        aria-label={`${label}: ${quality}`}
      >
        <div
          className={`w-3 h-3 rounded-full ${config.color} flex items-center justify-center text-white text-xs font-bold`}
        >
          {config.icon}
        </div>
        <span className={`text-sm font-medium ${config.textColor}`}>{label}</span>
      </div>
    );
  }

  return (
    <div
      className={`p-4 rounded-lg ${config.bgLight} ${config.borderColor} border-2`}
      role="status"
      aria-label={`${label}: ${quality}${description ? `. ${description}` : ''}`}
    >
      <div className="flex items-start gap-3">
        <div
          className={`w-8 h-8 rounded-full ${config.color} flex items-center justify-center text-white text-lg font-bold flex-shrink-0`}
        >
          {config.icon}
        </div>
        <div className="flex-1">
          <h4 className={`font-semibold ${config.textColor}`}>{label}</h4>
          <div className="text-sm text-gray-600 capitalize mt-1">{quality}</div>
          {description && (
            <p className="text-sm text-gray-600 mt-2">{description}</p>
          )}
        </div>
      </div>
    </div>
  );
};

interface MuhurtaTimelineProps {
  periods: {
    start: string;
    end: string;
    quality: 'auspicious' | 'inauspicious' | 'neutral' | 'mixed';
    label: string;
  }[];
  currentTime?: string;
}

/**
 * MuhurtaTimeline - Timeline visualization of muhurta periods throughout the day
 *
 * @component
 */
export const MuhurtaTimeline: React.FC<MuhurtaTimelineProps> = ({
  periods,
  currentTime,
}) => {
  const qualityColors = {
    auspicious: 'bg-green-500',
    inauspicious: 'bg-red-500',
    neutral: 'bg-gray-400',
    mixed: 'bg-yellow-500',
  };

  return (
    <div className="w-full">
      <div className="relative h-16 bg-gray-100 rounded-lg overflow-hidden">
        {periods.map((period, index) => {
          // Calculate position and width as percentage
          const startHour = parseFloat(period.start.split(':')[0]) + parseFloat(period.start.split(':')[1]) / 60;
          const endHour = parseFloat(period.end.split(':')[0]) + parseFloat(period.end.split(':')[1]) / 60;
          const left = (startHour / 24) * 100;
          const width = ((endHour - startHour) / 24) * 100;

          return (
            <div
              key={index}
              className={`absolute top-0 bottom-0 ${qualityColors[period.quality]} opacity-80 hover:opacity-100 transition-opacity`}
              style={{ left: `${left}%`, width: `${width}%` }}
              title={`${period.label}: ${period.start} - ${period.end}`}
            >
              <div className="h-full flex items-center justify-center text-white text-xs font-semibold px-1 truncate">
                {period.label}
              </div>
            </div>
          );
        })}

        {/* Current time indicator */}
        {currentTime && (
          <>
            {(() => {
              const [hours, minutes] = currentTime.split(':').map(Number);
              const currentHour = hours + minutes / 60;
              const position = (currentHour / 24) * 100;

              return (
                <div
                  className="absolute top-0 bottom-0 w-0.5 bg-blue-600 z-10"
                  style={{ left: `${position}%` }}
                  aria-label={`Current time: ${currentTime}`}
                >
                  <div className="absolute top-0 left-1/2 transform -translate-x-1/2 -translate-y-full">
                    <div className="bg-blue-600 text-white text-xs px-2 py-1 rounded">
                      {currentTime}
                    </div>
                  </div>
                </div>
              );
            })()}
          </>
        )}
      </div>

      {/* Time labels */}
      <div className="flex justify-between mt-2 text-xs text-gray-600">
        <span>00:00</span>
        <span>06:00</span>
        <span>12:00</span>
        <span>18:00</span>
        <span>24:00</span>
      </div>
    </div>
  );
};

export default LunarPhaseGraphic;
