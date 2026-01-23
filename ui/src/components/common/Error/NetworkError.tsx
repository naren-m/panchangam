import React from 'react';
import { AlertTriangle, Server, RefreshCw, Terminal } from 'lucide-react';

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
  const isBackendError = customMessage?.includes('Backend server') || customMessage?.includes('API server');

  return (
    <div className="bg-gradient-to-br from-red-50 to-orange-50 border border-red-200 rounded-xl p-6 shadow-sm">
      <div className="flex items-start gap-4">
        <div className="flex-shrink-0">
          <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
            {isBackendError ? (
              <Server className="w-6 h-6 text-red-600" />
            ) : (
              <AlertTriangle className="w-6 h-6 text-red-600" />
            )}
          </div>
        </div>

        <div className="flex-1">
          <h3 className="text-lg font-semibold text-red-800 mb-2">
            {isBackendError ? 'Backend Server Unavailable' : 'Connection Error'}
          </h3>

          <p className="text-red-700 mb-4">
            {customMessage || 'Unable to connect to the server. Please check your internet connection and try again.'}
          </p>

          {isBackendError && (
            <div className="bg-white/60 border border-red-100 rounded-lg p-4 mb-4">
              <div className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
                <Terminal className="w-4 h-4" />
                <span>To start the backend server:</span>
              </div>
              <div className="bg-gray-900 text-gray-100 rounded-md p-3 font-mono text-sm overflow-x-auto">
                <div className="text-gray-400"># Start the gRPC server first</div>
                <div>go run cmd/grpc-server/main.go</div>
                <div className="mt-2 text-gray-400"># Then start the HTTP gateway</div>
                <div>go run cmd/gateway/main.go</div>
              </div>
            </div>
          )}

          <button
            onClick={onRetry}
            disabled={isRetrying}
            className={`
              inline-flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all
              ${isRetrying
                ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                : 'bg-red-600 text-white hover:bg-red-700 active:scale-95'
              }
            `}
          >
            <RefreshCw className={`w-4 h-4 ${isRetrying ? 'animate-spin' : ''}`} />
            {isRetrying ? 'Retrying...' : 'Try Again'}
          </button>
        </div>
      </div>
    </div>
  );
};