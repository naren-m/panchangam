import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { TableView } from './TableView';
import { PanchangamData, Settings } from '../../types/panchangam';

const mockSettings: Settings = {
  calculation_method: 'Drik',
  locale: 'en',
  region: 'California',
  time_format: '12',
  location: {
    name: 'Milpitas, California',
    latitude: 37.4323,
    longitude: -121.9066,
    timezone: 'America/Los_Angeles',
    region: 'California'
  }
};

const mockPanchangamData: Record<string, PanchangamData> = {
  '2024-01-15': {
    date: '2024-01-15',
    tithi: 'Shukla Panchami',
    nakshatra: 'Rohini',
    yoga: 'Siddha',
    karana: 'Bava',
    sunrise_time: '07:15',
    sunset_time: '17:30',
    moonrise_time: '10:20',
    moonset_time: '22:45',
    vara: 'Monday',
    planetary_ruler: 'Moon',
    events: [
      {
        name: 'Sunrise',
        time: '07:15',
        event_type: 'SUNRISE',
        quality: 'neutral'
      },
      {
        name: 'Abhijit Muhurta',
        time: '12:00',
        event_type: 'ABHIJIT_MUHURTA',
        quality: 'auspicious'
      }
    ],
    festivals: ['Makar Sankranti']
  },
  '2024-01-16': {
    date: '2024-01-16',
    tithi: 'Shukla Shashthi',
    nakshatra: 'Mrigashira',
    yoga: 'Sadhya',
    karana: 'Balava',
    sunrise_time: '07:14',
    sunset_time: '17:31',
    vara: 'Tuesday',
    planetary_ruler: 'Mars',
    events: [],
    festivals: []
  }
};

describe('TableView', () => {
  it('renders table view with data', () => {
    const onDateClick = vi.fn();
    render(
      <TableView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('Panchangam Table View')).toBeInTheDocument();
    expect(screen.getByText('Shukla Panchami')).toBeInTheDocument();
    expect(screen.getByText('Rohini')).toBeInTheDocument();
  });

  it('displays empty state when no data', () => {
    const onDateClick = vi.fn();
    render(
      <TableView
        year={2024}
        month={0}
        panchangamData={{}}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('No data available for the selected month')).toBeInTheDocument();
  });

  it('calls onDateClick when row is clicked', () => {
    const onDateClick = vi.fn();
    render(
      <TableView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    const row = screen.getByText('Shukla Panchami').closest('tr');
    fireEvent.click(row!);

    expect(onDateClick).toHaveBeenCalled();
    const calledDate = onDateClick.mock.calls[0][0];
    expect(calledDate instanceof Date).toBe(true);
  });

  it('displays festivals when available', () => {
    const onDateClick = vi.fn();
    render(
      <TableView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('Makar Sankranti')).toBeInTheDocument();
  });

  it('shows auspicious events', () => {
    const onDateClick = vi.fn();
    render(
      <TableView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('Abhijit Muhurta')).toBeInTheDocument();
  });

  it('calls onExport with csv format', () => {
    const onDateClick = vi.fn();
    const onExport = vi.fn();
    render(
      <TableView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
        onExport={onExport}
      />
    );

    const csvButton = screen.getByText('CSV');
    fireEvent.click(csvButton);

    expect(onExport).toHaveBeenCalledWith('csv');
  });

  it('calls onExport with json format', () => {
    const onDateClick = vi.fn();
    const onExport = vi.fn();
    render(
      <TableView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
        onExport={onExport}
      />
    );

    const jsonButton = screen.getByText('JSON');
    fireEvent.click(jsonButton);

    expect(onExport).toHaveBeenCalledWith('json');
  });

  it('highlights today row', () => {
    const today = new Date();
    const todayStr = today.toISOString().split('T')[0];
    const todayData: Record<string, PanchangamData> = {
      [todayStr]: {
        date: todayStr,
        tithi: 'Today Tithi',
        nakshatra: 'Today Nakshatra',
        yoga: 'Siddha',
        karana: 'Bava',
        sunrise_time: '07:00',
        sunset_time: '18:00',
        vara: 'Today',
        planetary_ruler: 'Sun',
        events: [],
        festivals: []
      }
    };

    const onDateClick = vi.fn();
    render(
      <TableView
        year={today.getFullYear()}
        month={today.getMonth()}
        panchangamData={todayData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    // Should have "Today" badge in the table row
    const todayElements = screen.getAllByText('Today');
    expect(todayElements.length).toBeGreaterThan(0);
  });

  it('formats time in 12-hour format when setting is 12', () => {
    const onDateClick = vi.fn();
    render(
      <TableView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText(/7:15 AM/)).toBeInTheDocument();
    expect(screen.getByText(/5:30 PM/)).toBeInTheDocument();
  });

  it('formats time in 24-hour format when setting is 24', () => {
    const settings24h: Settings = {
      ...mockSettings,
      time_format: '24'
    };

    const onDateClick = vi.fn();
    render(
      <TableView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={settings24h}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('07:15')).toBeInTheDocument();
    expect(screen.getByText('17:30')).toBeInTheDocument();
  });

  it('displays summary footer with correct count', () => {
    const onDateClick = vi.fn();
    render(
      <TableView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('Showing 2 days')).toBeInTheDocument();
  });
});
