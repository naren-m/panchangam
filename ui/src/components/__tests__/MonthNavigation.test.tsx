import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MonthNavigation } from '../Calendar/MonthNavigation';
import { Settings } from '../../types/panchangam';

describe('MonthNavigation Component', () => {
  const mockOnPrevMonth = vi.fn();
  const mockOnNextMonth = vi.fn();
  const mockOnToday = vi.fn();
  const mockOnLocationClick = vi.fn();
  const mockOnSettingsClick = vi.fn();

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
  };

  const defaultProps = {
    year: 2024,
    month: 0, // January (0-indexed)
    settings: mockSettings,
    onPrevMonth: mockOnPrevMonth,
    onNextMonth: mockOnNextMonth,
    onToday: mockOnToday,
    onLocationClick: mockOnLocationClick,
    onSettingsClick: mockOnSettingsClick,
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders month and year correctly', () => {
    render(<MonthNavigation {...defaultProps} />);

    expect(screen.getByText(/January/i)).toBeInTheDocument();
    expect(screen.getByText(/2024/i)).toBeInTheDocument();
  });

  it('calls onPrevMonth when previous button is clicked', () => {
    render(<MonthNavigation {...defaultProps} />);

    const prevButton = screen.getByRole('button', { name: /previous/i });
    fireEvent.click(prevButton);

    expect(mockOnPrevMonth).toHaveBeenCalledTimes(1);
  });

  it('calls onNextMonth when next button is clicked', () => {
    render(<MonthNavigation {...defaultProps} />);

    const nextButton = screen.getByRole('button', { name: /next/i });
    fireEvent.click(nextButton);

    expect(mockOnNextMonth).toHaveBeenCalledTimes(1);
  });

  it('calls onToday when today button is clicked', () => {
    render(<MonthNavigation {...defaultProps} />);

    const todayButton = screen.getByRole('button', { name: /today/i });
    fireEvent.click(todayButton);

    expect(mockOnToday).toHaveBeenCalledTimes(1);
  });

  it('displays location information', () => {
    render(<MonthNavigation {...defaultProps} />);

    expect(screen.getByText(/Bangalore/i)).toBeInTheDocument();
  });

  it('handles different months correctly', () => {
    const { rerender } = render(<MonthNavigation {...defaultProps} />);
    expect(screen.getByText(/January/i)).toBeInTheDocument();

    rerender(<MonthNavigation {...defaultProps} month={11} />);
    expect(screen.getByText(/December/i)).toBeInTheDocument();
  });

  it('handles different years correctly', () => {
    const { rerender } = render(<MonthNavigation {...defaultProps} />);
    expect(screen.getByText(/2024/i)).toBeInTheDocument();

    rerender(<MonthNavigation {...defaultProps} year={2025} />);
    expect(screen.getByText(/2025/i)).toBeInTheDocument();
  });
});
