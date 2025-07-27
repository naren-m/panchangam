import React from 'react';

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg' | 'xl';
  color?: 'orange' | 'blue' | 'gray' | 'green';
  message?: string;
  fullScreen?: boolean;
}

const sizeClasses = {
  sm: 'h-4 w-4',
  md: 'h-6 w-6',
  lg: 'h-8 w-8',
  xl: 'h-12 w-12',
};

const colorClasses = {
  orange: 'border-orange-500',
  blue: 'border-blue-500',
  gray: 'border-gray-500',
  green: 'border-green-500',
};

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  color = 'orange',
  message,
  fullScreen = false,
}) => {
  const spinnerElement = (
    <div className="flex flex-col items-center justify-center">
      <div
        className={`
          animate-spin rounded-full border-2 border-t-transparent
          ${sizeClasses[size]} ${colorClasses[color]}
        `}
        role="status"
        aria-label="Loading"
      />
      {message && (
        <span className={`mt-2 text-${color}-700 font-medium text-sm`}>
          {message}
        </span>
      )}
    </div>
  );

  if (fullScreen) {
    return (
      <div className="fixed inset-0 bg-white bg-opacity-90 flex items-center justify-center z-50">
        {spinnerElement}
      </div>
    );
  }

  return spinnerElement;
};