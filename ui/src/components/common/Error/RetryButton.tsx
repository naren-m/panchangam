import React, { useState, useEffect } from 'react';
import { RefreshCw } from 'lucide-react';

interface RetryButtonProps {
  onRetry: () => void;
  isRetrying?: boolean;
  disabled?: boolean;
  countdown?: number; // Countdown in seconds before retry is available
  variant?: 'primary' | 'secondary' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  className?: string;
  children?: React.ReactNode;
}

const variantClasses = {
  primary: 'bg-orange-500 text-white hover:bg-orange-600 disabled:bg-gray-300',
  secondary: 'bg-white text-orange-600 border border-orange-500 hover:bg-orange-50 disabled:bg-gray-100 disabled:text-gray-400 disabled:border-gray-300',
  ghost: 'bg-transparent text-orange-600 hover:bg-orange-50 disabled:text-gray-400',
};

const sizeClasses = {
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2 text-base',
  lg: 'px-6 py-3 text-lg',
};

export const RetryButton: React.FC<RetryButtonProps> = ({
  onRetry,
  isRetrying = false,
  disabled = false,
  countdown = 0,
  variant = 'primary',
  size = 'md',
  className = '',
  children = 'Retry',
}) => {
  const [remainingTime, setRemainingTime] = useState(countdown);

  useEffect(() => {
    if (countdown > 0) {
      setRemainingTime(countdown);
    }
  }, [countdown]);

  useEffect(() => {
    if (remainingTime > 0) {
      const timer = setTimeout(() => {
        setRemainingTime(prev => prev - 1);
      }, 1000);

      return () => clearTimeout(timer);
    }
  }, [remainingTime]);

  const isDisabled = disabled || isRetrying || remainingTime > 0;

  const handleClick = () => {
    if (!isDisabled) {
      onRetry();
    }
  };

  const getButtonText = () => {
    if (isRetrying) return 'Retrying...';
    if (remainingTime > 0) return `Retry in ${remainingTime}s`;
    return children;
  };

  const getIcon = () => {
    if (isRetrying) {
      return <RefreshCw className="w-4 h-4 animate-spin" />;
    }
    return <RefreshCw className="w-4 h-4" />;
  };

  return (
    <button
      onClick={handleClick}
      disabled={isDisabled}
      className={`
        inline-flex items-center justify-center space-x-2 rounded-lg font-medium 
        transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2
        disabled:cursor-not-allowed disabled:opacity-50
        ${variantClasses[variant]} ${sizeClasses[size]} ${className}
      `}
      type="button"
      aria-label={isRetrying ? 'Retrying...' : 'Retry operation'}
    >
      {getIcon()}
      <span>{getButtonText()}</span>
    </button>
  );
};