import React from 'react';
import { WifiOff, Wifi } from 'lucide-react';

interface OfflineIndicatorProps {
  isOffline: boolean;
  className?: string;
  showWhenOnline?: boolean;
  position?: 'fixed' | 'relative';
}

export const OfflineIndicator: React.FC<OfflineIndicatorProps> = ({
  isOffline,
  className = '',
  showWhenOnline = false,
  position = 'fixed',
}) => {
  // Don't render anything if online and showWhenOnline is false
  if (!isOffline && !showWhenOnline) {
    return null;
  }

  const positionClasses = {
    fixed: 'fixed top-4 right-4 z-50',
    relative: 'relative',
  };

  return (
    <div
      className={`
        inline-flex items-center space-x-2 px-3 py-2 rounded-lg shadow-lg
        transition-all duration-300 ease-in-out
        ${isOffline 
          ? 'bg-red-500 text-white' 
          : 'bg-green-500 text-white'
        }
        ${positionClasses[position]}
        ${className}
      `}
      role="status"
      aria-live="polite"
      aria-label={isOffline ? 'You are offline' : 'You are online'}
    >
      {isOffline ? (
        <WifiOff className="w-4 h-4" />
      ) : (
        <Wifi className="w-4 h-4" />
      )}
      <span className="text-sm font-medium">
        {isOffline ? 'Offline' : 'Online'}
      </span>
    </div>
  );
};