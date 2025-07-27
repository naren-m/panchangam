import React from 'react';
import { LoadingSpinner } from './LoadingSpinner';

interface LoadingOverlayProps {
  isVisible: boolean;
  message?: string;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  backdrop?: 'light' | 'dark' | 'blur';
}

const backdropClasses = {
  light: 'bg-white bg-opacity-80',
  dark: 'bg-black bg-opacity-50',
  blur: 'bg-white bg-opacity-80 backdrop-blur-sm',
};

export const LoadingOverlay: React.FC<LoadingOverlayProps> = ({
  isVisible,
  message = 'Loading...',
  size = 'lg',
  backdrop = 'blur',
}) => {
  if (!isVisible) return null;

  return (
    <div
      className={`
        fixed inset-0 z-50 flex items-center justify-center
        ${backdropClasses[backdrop]}
        transition-opacity duration-200
      `}
      role="dialog"
      aria-modal="true"
      aria-label="Loading"
    >
      <div className="flex flex-col items-center p-6 bg-white rounded-lg shadow-lg">
        <LoadingSpinner size={size} message={message} />
      </div>
    </div>
  );
};