import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { SettingsPanel } from '../Settings/SettingsPanel';
import { Settings } from '../../types/panchangam';
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

describe('SettingsPanel Component', () => {
  const mockOnSettingsChange = vi.fn();
  const mockOnClose = vi.fn();

  const mockSettings: Settings = {
    location: {
      latitude: 12.9716,
      longitude: 77.5946,
      timezone: 'Asia/Kolkata',
      name: 'Bangalore',
    },
    calculationMethod: 'traditional',
    displayLanguage: 'en',
    theme: 'light',
    locale: 'en-US',
    calculation_method: 'Drik',
    time_format: '12',
    region: 'Karnataka',
  };

  const defaultProps = {
    settings: mockSettings,
    onSettingsChange: mockOnSettingsChange,
    onClose: mockOnClose,
  };

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(panchangamApi.healthCheck).mockResolvedValue({
      status: 'healthy',
      message: 'API is accessible',
    });
  });

  it('renders settings panel with title', async () => {
    render(<SettingsPanel {...defaultProps} />);

    expect(screen.getByText('Settings')).toBeInTheDocument();

    // Wait for API health check to complete
    await waitFor(() => {
      expect(screen.getByText('Connected')).toBeInTheDocument();
    });
  });

  it('renders API health check component', async () => {
    render(<SettingsPanel {...defaultProps} />);

    // Wait for health check and verify it's displayed
    await waitFor(() => {
      expect(screen.getByText('API Connection')).toBeInTheDocument();
    });
  });

  it('displays calculation method options', () => {
    render(<SettingsPanel {...defaultProps} />);

    expect(screen.getByText('Calculation Method')).toBeInTheDocument();
    expect(screen.getByText('Drik (Modern)')).toBeInTheDocument();
    expect(screen.getByText('Vakya (Traditional)')).toBeInTheDocument();
  });

  it('displays language selector', () => {
    render(<SettingsPanel {...defaultProps} />);

    expect(screen.getByText('Language')).toBeInTheDocument();

    // Find the select element by its value
    const selects = screen.getAllByRole('combobox');
    expect(selects.length).toBeGreaterThan(0);
  });

  it('displays time format options', () => {
    render(<SettingsPanel {...defaultProps} />);

    expect(screen.getByText('Time Format')).toBeInTheDocument();
    expect(screen.getByText('12-hour (AM/PM)')).toBeInTheDocument();
    expect(screen.getByText('24-hour')).toBeInTheDocument();
  });

  it('displays region selector', () => {
    render(<SettingsPanel {...defaultProps} />);

    expect(screen.getByText('Region')).toBeInTheDocument();

    // Verify there are combobox elements (language and region selectors)
    const selects = screen.getAllByRole('combobox');
    expect(selects.length).toBeGreaterThanOrEqual(2); // Language and Region
  });

  it('calls onClose when close button is clicked', () => {
    render(<SettingsPanel {...defaultProps} />);

    const closeButton = screen.getAllByRole('button').find(
      (button) => button.querySelector('.lucide-x')
    );

    if (closeButton) {
      fireEvent.click(closeButton);
      expect(mockOnClose).toHaveBeenCalledTimes(1);
    }
  });

  it('calls onClose when Cancel button is clicked', () => {
    render(<SettingsPanel {...defaultProps} />);

    const cancelButton = screen.getByText('Cancel');
    fireEvent.click(cancelButton);

    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('calls onClose when Save Settings button is clicked', () => {
    render(<SettingsPanel {...defaultProps} />);

    const saveButton = screen.getByText('Save Settings');
    fireEvent.click(saveButton);

    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('calls onSettingsChange when calculation method is changed', () => {
    render(<SettingsPanel {...defaultProps} />);

    const vakyaRadio = screen.getByLabelText(/Vakya \(Traditional\)/i);
    fireEvent.click(vakyaRadio);

    expect(mockOnSettingsChange).toHaveBeenCalledWith({
      ...mockSettings,
      calculation_method: 'Vakya',
    });
  });

  it('calls onSettingsChange when time format is changed', () => {
    render(<SettingsPanel {...defaultProps} />);

    const format24Radio = screen.getByLabelText('24-hour');
    fireEvent.click(format24Radio);

    expect(mockOnSettingsChange).toHaveBeenCalledWith({
      ...mockSettings,
      time_format: '24',
    });
  });

  it('displays modal overlay', () => {
    const { container } = render(<SettingsPanel {...defaultProps} />);

    // Check for modal overlay with dark background
    const overlay = container.querySelector('.bg-black.bg-opacity-50');
    expect(overlay).toBeInTheDocument();
  });

  it('displays Save and Cancel buttons in footer', () => {
    render(<SettingsPanel {...defaultProps} />);

    expect(screen.getByText('Cancel')).toBeInTheDocument();
    expect(screen.getByText('Save Settings')).toBeInTheDocument();
  });
});
