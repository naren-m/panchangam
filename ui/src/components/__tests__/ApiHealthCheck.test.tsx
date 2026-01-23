import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { ApiHealthCheck } from '../Settings/ApiHealthCheck';
import { panchangamApi } from '../../services/panchangamApi';

// Mock the API service
vi.mock('../../services/panchangamApi', () => ({
  panchangamApi: {
    healthCheck: vi.fn(),
  },
  apiConfig: {
    baseUrl: 'http://localhost:8080',
    endpoint: 'http://localhost:8080/api/v1/panchangam',
  },
}));

describe('ApiHealthCheck Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows healthy status when API is accessible', async () => {
    vi.mocked(panchangamApi.healthCheck).mockResolvedValue({
      status: 'healthy',
      message: 'API is accessible',
    });

    render(<ApiHealthCheck />);

    // Wait for health check to complete
    await waitFor(() => {
      expect(screen.getByText('Connected')).toBeInTheDocument();
    });

    expect(screen.getByText(/API is accessible/i)).toBeInTheDocument();
  });

  it('shows unhealthy status when API is not accessible', async () => {
    vi.mocked(panchangamApi.healthCheck).mockResolvedValue({
      status: 'unhealthy',
      message: 'Connection failed',
    });

    render(<ApiHealthCheck />);

    await waitFor(() => {
      expect(screen.getByText('Disconnected')).toBeInTheDocument();
    });

    // Check that the error alert is displayed
    expect(screen.getByText(/The app will use fallback data/i)).toBeInTheDocument();
  });

  it('handles health check errors gracefully', async () => {
    vi.mocked(panchangamApi.healthCheck).mockRejectedValue(new Error('Network error'));

    render(<ApiHealthCheck />);

    await waitFor(() => {
      expect(screen.getByText('Disconnected')).toBeInTheDocument();
    });

    expect(screen.getByText(/Network error/i)).toBeInTheDocument();
  });

  it('displays API endpoint information', async () => {
    vi.mocked(panchangamApi.healthCheck).mockResolvedValue({
      status: 'healthy',
      message: 'API is accessible',
    });

    render(<ApiHealthCheck />);

    // Endpoint is displayed regardless of health check status
    expect(screen.getByText(/localhost:8080/i)).toBeInTheDocument();

    // Wait for health check to complete
    await waitFor(() => {
      expect(screen.getByText('Connected')).toBeInTheDocument();
    });
  });
});
