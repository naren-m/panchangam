import React from 'react';

interface ProgressBarProps {
  progress: number; // 0-100
  size?: 'sm' | 'md' | 'lg';
  color?: 'orange' | 'blue' | 'green' | 'red';
  label?: string;
  showPercentage?: boolean;
  animated?: boolean;
}

const sizeClasses = {
  sm: 'h-2',
  md: 'h-3', 
  lg: 'h-4',
};

const colorClasses = {
  orange: 'bg-orange-500',
  blue: 'bg-blue-500',
  green: 'bg-green-500',
  red: 'bg-red-500',
};

const backgroundColorClasses = {
  orange: 'bg-orange-100',
  blue: 'bg-blue-100',
  green: 'bg-green-100',
  red: 'bg-red-100',
};

export const ProgressBar: React.FC<ProgressBarProps> = ({
  progress,
  size = 'md',
  color = 'orange',
  label,
  showPercentage = false,
  animated = true,
}) => {
  // Clamp progress between 0 and 100
  const clampedProgress = Math.max(0, Math.min(100, progress));
  
  return (
    <div className="w-full">
      {/* Label and percentage */}
      {(label || showPercentage) && (
        <div className="flex items-center justify-between mb-2">
          {label && (
            <span className="text-sm font-medium text-gray-700">
              {label}
            </span>
          )}
          {showPercentage && (
            <span className="text-sm text-gray-500">
              {Math.round(clampedProgress)}%
            </span>
          )}
        </div>
      )}
      
      {/* Progress bar */}
      <div 
        className={`
          w-full rounded-full overflow-hidden
          ${sizeClasses[size]} ${backgroundColorClasses[color]}
        `}
        role="progressbar"
        aria-valuenow={clampedProgress}
        aria-valuemin={0}
        aria-valuemax={100}
        aria-label={label || 'Progress'}
      >
        <div
          className={`
            h-full rounded-full transition-all duration-300 ease-out
            ${colorClasses[color]}
            ${animated ? 'animate-pulse' : ''}
          `}
          style={{ width: `${clampedProgress}%` }}
        />
      </div>
    </div>
  );
};