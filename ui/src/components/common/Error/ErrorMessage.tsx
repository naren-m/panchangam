import { useState, type FC } from 'react';
import { AlertCircle, AlertTriangle, Info, type LucideIcon } from 'lucide-react';

type ErrorMessageType = 'error' | 'warning' | 'info';

interface ErrorMessageProps {
  title?: string;
  message: string;
  type?: ErrorMessageType;
  onRetry?: () => void;
  onDismiss?: () => void;
  showDetails?: boolean;
  details?: string;
}

interface ErrorMessageStyle {
  container: string;
  Icon: LucideIcon;
  iconColor: string;
  titleColor: string;
  messageColor: string;
  buttonColor: string;
}

const typeStyles: Record<ErrorMessageType, ErrorMessageStyle> = {
  error: {
    container: 'bg-red-50 border-red-200',
    Icon: AlertCircle,
    iconColor: 'text-red-600',
    titleColor: 'text-red-800',
    messageColor: 'text-red-700',
    buttonColor: 'bg-red-600 hover:bg-red-700',
  },
  warning: {
    container: 'bg-yellow-50 border-yellow-200',
    Icon: AlertTriangle,
    iconColor: 'text-yellow-600',
    titleColor: 'text-yellow-800',
    messageColor: 'text-yellow-700',
    buttonColor: 'bg-yellow-600 hover:bg-yellow-700',
  },
  info: {
    container: 'bg-blue-50 border-blue-200',
    Icon: Info,
    iconColor: 'text-blue-600',
    titleColor: 'text-blue-800',
    messageColor: 'text-blue-700',
    buttonColor: 'bg-blue-600 hover:bg-blue-700',
  },
};

export const ErrorMessage: FC<ErrorMessageProps> = ({
  title,
  message,
  type = 'error',
  onRetry,
  onDismiss,
  showDetails = false,
  details,
}) => {
  const [showDetailText, setShowDetailText] = useState(false);
  const styles = typeStyles[type];
  const StatusIcon = styles.Icon;

  return (
    <div className={`${styles.container} border rounded-lg p-4 mb-4`} role="alert">
      <div className="flex items-start space-x-3">
        <StatusIcon
          className={`${styles.iconColor} h-5 w-5 flex-shrink-0`}
          aria-hidden="true"
          data-testid={`${type}-status-icon`}
        />
        
        <div className="flex-1 min-w-0">
          {title && (
            <h3 className={`font-semibold ${styles.titleColor} mb-1`}>
              {title}
            </h3>
          )}
          
          <p className={`${styles.messageColor} text-sm`}>
            {message}
          </p>
          
          {showDetails && details && (
            <div className="mt-2">
              <button
                onClick={() => setShowDetailText(!showDetailText)}
                className={`text-xs ${styles.messageColor} underline hover:no-underline`}
              >
                {showDetailText ? 'Hide details' : 'Show details'}
              </button>
              
              {showDetailText && (
                <div className="mt-2 p-2 bg-white bg-opacity-50 rounded text-xs font-mono">
                  {details}
                </div>
              )}
            </div>
          )}
        </div>
        
        <div className="flex flex-col space-y-2">
          {onDismiss && (
            <button
              onClick={onDismiss}
              className={`${styles.iconColor} hover:opacity-70 text-lg leading-none`}
              aria-label="Dismiss"
            >
              ×
            </button>
          )}
        </div>
      </div>
      
      {onRetry && (
        <div className="mt-3 flex justify-end">
          <button
            onClick={onRetry}
            className={`
              ${styles.buttonColor} text-white px-3 py-1 rounded text-sm
              transition-colors duration-200
            `}
          >
            Try Again
          </button>
        </div>
      )}
    </div>
  );
};
