import React from 'react';

interface ShimmerEffectProps {
  width?: string | number;
  height?: string | number;
  className?: string;
  variant?: 'text' | 'rectangular' | 'circular' | 'rounded';
  lines?: number; // For text variant
  animated?: boolean;
}

export const ShimmerEffect: React.FC<ShimmerEffectProps> = ({
  width = '100%',
  height = '1rem',
  className = '',
  variant = 'rectangular',
  lines = 1,
  animated = true,
}) => {
  const baseClasses = `
    bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200
    ${animated ? 'animate-pulse' : ''}
  `;

  const variantClasses = {
    text: 'rounded',
    rectangular: '',
    circular: 'rounded-full',
    rounded: 'rounded-lg',
  };

  const shimmerStyle = {
    width: typeof width === 'number' ? `${width}px` : width,
    height: typeof height === 'number' ? `${height}px` : height,
  };

  // For text variant, render multiple lines
  if (variant === 'text' && lines > 1) {
    return (
      <div className={`space-y-2 ${className}`}>
        {Array.from({ length: lines }).map((_, index) => (
          <div
            key={index}
            className={`${baseClasses} ${variantClasses[variant]}`}
            style={{
              ...shimmerStyle,
              width: index === lines - 1 ? '75%' : shimmerStyle.width, // Last line shorter
            }}
          />
        ))}
      </div>
    );
  }

  return (
    <div
      className={`${baseClasses} ${variantClasses[variant]} ${className}`}
      style={shimmerStyle}
      role="status"
      aria-label="Loading content"
    />
  );
};

// Predefined shimmer components for common use cases
export const ShimmerText: React.FC<{ lines?: number; className?: string }> = ({ 
  lines = 1, 
  className 
}) => (
  <ShimmerEffect 
    variant="text" 
    height="1rem" 
    lines={lines} 
    className={className} 
  />
);

export const ShimmerButton: React.FC<{ className?: string }> = ({ className }) => (
  <ShimmerEffect 
    variant="rounded" 
    width="120px" 
    height="2.5rem" 
    className={className} 
  />
);

export const ShimmerAvatar: React.FC<{ size?: string; className?: string }> = ({ 
  size = '3rem', 
  className 
}) => (
  <ShimmerEffect 
    variant="circular" 
    width={size} 
    height={size} 
    className={className} 
  />
);

export const ShimmerCard: React.FC<{ className?: string }> = ({ className }) => (
  <div className={`space-y-4 p-4 border border-gray-200 rounded-lg ${className}`}>
    <ShimmerEffect variant="rounded" height="1.5rem" width="60%" />
    <ShimmerText lines={3} />
    <div className="flex space-x-2">
      <ShimmerButton />
      <ShimmerButton />
    </div>
  </div>
);