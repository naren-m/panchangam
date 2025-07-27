import React from 'react';
import { ErrorMessage } from './ErrorMessage';

interface NetworkErrorProps {
  onRetry?: () => void;
  isRetrying?: boolean;
  customMessage?: string;
}

export const NetworkError: React.FC<NetworkErrorProps> = ({
  onRetry,
  isRetrying = false,
  customMessage,
}) => {
  const defaultMessage = 'Unable to connect to the server. Please check your internet connection and try again.';
  
  return (
    <ErrorMessage
      title="Connection Error"
      message={customMessage || defaultMessage}
      type="error"
      onRetry={isRetrying ? undefined : onRetry}
      details={isRetrying ? 'Retrying...' : undefined}
    />
  );
};