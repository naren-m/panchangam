import React, { useState, useEffect } from 'react';
import { Activity, CheckCircle, XCircle, AlertCircle } from 'lucide-react';
import { panchangamApi, apiConfig } from '../../services/panchangamApi';

interface ApiHealthCheckProps {
  onStatusChange?: (status: 'healthy' | 'unhealthy') => void;
}

export const ApiHealthCheck: React.FC<ApiHealthCheckProps> = ({ onStatusChange }) => {
  const [status, setStatus] = useState<'checking' | 'healthy' | 'unhealthy'>('checking');
  const [message, setMessage] = useState('');
  const [lastChecked, setLastChecked] = useState<Date | null>(null);

  const checkApiHealth = async () => {
    setStatus('checking');
    try {
      const health = await panchangamApi.healthCheck();
      setStatus(health.status);
      setMessage(health.message);
      setLastChecked(new Date());
      onStatusChange?.(health.status);
    } catch (error) {
      setStatus('unhealthy');
      setMessage(error instanceof Error ? error.message : 'Unknown error');
      setLastChecked(new Date());
      onStatusChange?.('unhealthy');
    }
  };

  useEffect(() => {
    checkApiHealth();
    
    // Set up periodic health checks every 30 seconds
    const interval = setInterval(checkApiHealth, 30000);
    return () => clearInterval(interval);
  }, []);

  const getStatusIcon = () => {
    switch (status) {
      case 'checking':
        return <Activity className="w-4 h-4 animate-spin text-blue-500" />;
      case 'healthy':
        return <CheckCircle className="w-4 h-4 text-green-500" />;
      case 'unhealthy':
        return <XCircle className="w-4 h-4 text-red-500" />;
      default:
        return <AlertCircle className="w-4 h-4 text-yellow-500" />;
    }
  };

  const getStatusColor = () => {
    switch (status) {
      case 'checking':
        return 'text-blue-600';
      case 'healthy':
        return 'text-green-600';
      case 'unhealthy':
        return 'text-red-600';
      default:
        return 'text-yellow-600';
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-semibold text-gray-800">API Connection</h3>
        <button
          onClick={checkApiHealth}
          className="text-xs px-2 py-1 bg-gray-100 hover:bg-gray-200 rounded transition-colors"
          disabled={status === 'checking'}
        >
          Refresh
        </button>
      </div>
      
      <div className="space-y-2">
        <div className="flex items-center space-x-2">
          {getStatusIcon()}
          <span className={`text-sm font-medium ${getStatusColor()}`}>
            {status === 'checking' ? 'Checking...' : status === 'healthy' ? 'Connected' : 'Disconnected'}
          </span>
        </div>
        
        <div className="text-xs text-gray-600">
          <div>Endpoint: {apiConfig.endpoint}</div>
          {message && <div>Status: {message}</div>}
          {lastChecked && (
            <div>Last checked: {lastChecked.toLocaleTimeString()}</div>
          )}
        </div>

        {status === 'unhealthy' && (
          <div className="mt-3 p-2 bg-red-50 border border-red-200 rounded text-xs text-red-700">
            <strong>Connection Failed:</strong> The app will use fallback data. 
            Please ensure the API server is running on {apiConfig.baseUrl}.
          </div>
        )}

        {import.meta.env.DEV && (
          <div className="mt-3 p-2 bg-blue-50 border border-blue-200 rounded text-xs text-blue-700">
            <strong>Development Mode:</strong> Make sure to start the gateway server:
            <br />
            <code className="bg-blue-100 px-1 rounded">go run cmd/gateway/main.go</code>
          </div>
        )}
      </div>
    </div>
  );
};