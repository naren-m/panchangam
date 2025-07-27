import React from 'react';
import { ErrorMessage } from './ErrorMessage';

interface ApiErrorProps {
  error: string | Error;
  onRetry?: () => void;
  statusCode?: number;
  endpoint?: string;
}

export const ApiError: React.FC<ApiErrorProps> = ({
  error,
  onRetry,
  statusCode,
  endpoint,
}) => {
  const errorMessage = error instanceof Error ? error.message : error;
  
  const getTitle = () => {
    if (statusCode) {
      switch (statusCode) {
        case 400:
          return 'Invalid Request';
        case 401:
          return 'Authentication Required';
        case 403:
          return 'Access Forbidden';
        case 404:
          return 'Data Not Found';
        case 429:
          return 'Too Many Requests';
        case 500:
          return 'Server Error';
        case 502:
          return 'Service Unavailable';
        case 503:
          return 'Service Maintenance';
        default:
          return 'API Error';
      }
    }
    return 'API Error';
  };

  const getMessage = () => {
    if (statusCode) {
      switch (statusCode) {
        case 400:
          return 'The request contains invalid parameters. Please check your input and try again.';
        case 401:
          return 'Authentication is required to access this data.';
        case 403:
          return 'You do not have permission to access this resource.';
        case 404:
          return 'The requested data could not be found. It may have been moved or deleted.';
        case 429:
          return 'Too many requests have been made. Please wait a moment before trying again.';
        case 500:
          return 'An internal server error occurred. Our team has been notified.';
        case 502:
          return 'The service is temporarily unavailable. Please try again in a few moments.';
        case 503:
          return 'The service is currently under maintenance. Please try again later.';
        default:
          return errorMessage;
      }
    }
    return errorMessage;
  };

  const details = process.env.NODE_ENV === 'development' && endpoint 
    ? `Endpoint: ${endpoint}\nStatus: ${statusCode || 'Unknown'}\nError: ${errorMessage}`
    : undefined;

  return (
    <ErrorMessage
      title={getTitle()}
      message={getMessage()}
      type="error"
      onRetry={onRetry}
      showDetails={!!details}
      details={details}
    />
  );
};