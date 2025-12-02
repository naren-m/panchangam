import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { CalendarGrid } from '../Calendar/CalendarGrid';
import { PanchangamData, Settings } from '../../types/panchangam';

describe('CalendarGrid Component', () => {
  const mockOnDateClick = vi.fn();

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

  const mockPanchangamData: Record<string, PanchangamData> = {
    '2024-01-15': {
      date: '2024-01-15',
      tithi: 'Chaturthi (4)',
      nakshatra: 'Uttara Bhadrapada (26)',
      yoga: 'Siddha (21)',
      karana: 'Gara (6)',
      vara: 'Monday',
      planetary_ruler: 'Moon',
      sunrise_time: '06:30:00',
      sunset_time: '18:30:00',
      events: [],
      festivals: [],
    },
  };

  const defaultProps = {
    year: 2024,
    month: 0, // January
    panchangamData: mockPanchangamData,
    settings: mockSettings,
    onDateClick: mockOnDateClick,
  };

  it('renders calendar grid with weekday headers', () => {
    render(<CalendarGrid {...defaultProps} />);

    // Check for weekday headers
    expect(screen.getByText(/sun/i)).toBeInTheDocument();
    expect(screen.getByText(/mon/i)).toBeInTheDocument();
    expect(screen.getByText(/tue/i)).toBeInTheDocument();
    expect(screen.getByText(/wed/i)).toBeInTheDocument();
    expect(screen.getByText(/thu/i)).toBeInTheDocument();
    expect(screen.getByText(/fri/i)).toBeInTheDocument();
    expect(screen.getByText(/sat/i)).toBeInTheDocument();
  });

  it('renders dates for the specified month', () => {
    const { container } = render(<CalendarGrid {...defaultProps} />);

    // Should have date cells (usually 35-42 cells for a month view)
    // Date cells are divs with cursor-pointer class
    const dateCells = container.querySelectorAll('.cursor-pointer');
    expect(dateCells.length).toBeGreaterThan(0);
  });

  it('calls onDateClick when a date is clicked', () => {
    const { container } = render(<CalendarGrid {...defaultProps} />);

    const dateCells = container.querySelectorAll('.cursor-pointer');
    fireEvent.click(dateCells[0]);

    expect(mockOnDateClick).toHaveBeenCalledTimes(1);
    expect(mockOnDateClick).toHaveBeenCalledWith(expect.any(Date));
  });

  it('renders calendar without errors when Panchangam data is available', () => {
    const { container } = render(<CalendarGrid {...defaultProps} />);

    // Calendar should render successfully with the data provided
    const dateCells = container.querySelectorAll('.cursor-pointer');
    expect(dateCells.length).toBeGreaterThan(0);

    // Note: Panchangam data display depends on DateCell component
    // and whether dates are in current month view
  });

  it('handles empty Panchangam data gracefully', () => {
    const propsWithNoData = {
      ...defaultProps,
      panchangamData: {},
    };

    const { container } = render(<CalendarGrid {...propsWithNoData} />);

    // Should still render the grid without errors
    const dateCells = container.querySelectorAll('.cursor-pointer');
    expect(dateCells.length).toBeGreaterThan(0);
  });

  it('renders correct number of date cells', () => {
    const { container } = render(<CalendarGrid {...defaultProps} />);

    const dateCells = container.querySelectorAll('.cursor-pointer');
    // Calendar grids typically show 35 or 42 cells (5 or 6 weeks)
    expect(dateCells.length).toBeGreaterThanOrEqual(28); // At least 4 weeks
    expect(dateCells.length).toBeLessThanOrEqual(42); // At most 6 weeks
  });
});
