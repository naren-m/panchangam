import { useState, useEffect } from 'react';

interface UseOfflineOptions {
  onOnline?: () => void;
  onOffline?: () => void;
}

export interface OfflineState {
  isOffline: boolean;
  wasOffline: boolean;
  lastOnlineTime: Date | null;
  lastOfflineTime: Date | null;
}

export const useOffline = (options: UseOfflineOptions = {}): OfflineState => {
  const { onOnline, onOffline } = options;

  const [offlineState, setOfflineState] = useState<OfflineState>({
    isOffline: !navigator.onLine,
    wasOffline: false,
    lastOnlineTime: navigator.onLine ? new Date() : null,
    lastOfflineTime: !navigator.onLine ? new Date() : null,
  });

  useEffect(() => {
    const handleOnline = () => {
      const now = new Date();
      setOfflineState(prev => ({
        ...prev,
        isOffline: false,
        wasOffline: prev.isOffline,
        lastOnlineTime: now,
      }));
      
      if (onOnline) {
        onOnline();
      }
    };

    const handleOffline = () => {
      const now = new Date();
      setOfflineState(prev => ({
        ...prev,
        isOffline: true,
        wasOffline: false, // Reset wasOffline when going offline
        lastOfflineTime: now,
      }));
      
      if (onOffline) {
        onOffline();
      }
    };

    // Add event listeners
    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    // Cleanup
    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, [onOnline, onOffline]);

  return offlineState;
};