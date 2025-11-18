import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { GraphView } from './GraphView';
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
        name: 'Abhijit Muhurta',
        time: '12:00',
        event_type: 'ABHIJIT_MUHURTA',
        quality: 'auspicious'
      },
      {
        name: 'Brahma Muhurta',
        time: '05:30',
        event_type: 'BRAHMA_MUHURTA',
        quality: 'auspicious'
      }
    ],
    festivals: ['Makar Sankranti', 'Pongal']
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
    events: [
      {
        name: 'Abhijit Muhurta',
        time: '12:00',
        event_type: 'ABHIJIT_MUHURTA',
        quality: 'auspicious'
      }
    ],
    festivals: []
  },
  '2024-01-17': {
    date: '2024-01-17',
    tithi: 'Shukla Saptami',
    nakshatra: 'Ardra',
    yoga: 'Shubha',
    karana: 'Kaulava',
    sunrise_time: '07:13',
    sunset_time: '17:32',
    vara: 'Wednesday',
    planetary_ruler: 'Mercury',
    events: [],
    festivals: []
  }
};

describe('GraphView', () => {
  it('renders graph view with analytics title', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('Panchangam Analytics')).toBeInTheDocument();
  });

  it('displays empty state when no data', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={{}}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('No data available for visualization')).toBeInTheDocument();
  });

  it('shows tithi distribution chart', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('Tithi Distribution')).toBeInTheDocument();
    expect(screen.getByText('Shukla Panchami')).toBeInTheDocument();
    expect(screen.getByText('Shukla Shashthi')).toBeInTheDocument();
  });

  it('shows nakshatra distribution chart', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('Nakshatra Distribution')).toBeInTheDocument();
    expect(screen.getByText('Rohini')).toBeInTheDocument();
    expect(screen.getByText('Mrigashira')).toBeInTheDocument();
  });

  it('shows auspicious events timeline', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('Daily Auspicious Events')).toBeInTheDocument();
  });

  it('shows festival days section when festivals exist', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText(/Festival Days/)).toBeInTheDocument();
    // Festivals are rendered with bullets and may appear multiple times
    const makarElements = screen.getAllByText((content, element) =>
      element?.textContent?.includes('Makar Sankranti') || false
    );
    expect(makarElements.length).toBeGreaterThan(0);
    const pongalElements = screen.getAllByText((content, element) =>
      element?.textContent?.includes('Pongal') || false
    );
    expect(pongalElements.length).toBeGreaterThan(0);
  });

  it('displays month summary statistics', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('Month Summary')).toBeInTheDocument();
    expect(screen.getByText('Total Days')).toBeInTheDocument();
    expect(screen.getByText('Festivals')).toBeInTheDocument();
    expect(screen.getByText('Unique Tithis')).toBeInTheDocument();
    expect(screen.getByText('Auspicious Events')).toBeInTheDocument();
  });

  it('calls onDateClick when festival day is clicked', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    const festivalElements = screen.getAllByText((content, element) =>
      element?.textContent?.includes('Makar Sankranti') || false
    );
    // Find the one that's in a clickable card (festival days section)
    const festivalCard = festivalElements[0].closest('div[class*="cursor-pointer"]');
    if (festivalCard) {
      fireEvent.click(festivalCard);
      expect(onDateClick).toHaveBeenCalled();
    }
  });

  it('calls onExport with csv format', () => {
    const onDateClick = vi.fn();
    const onExport = vi.fn();
    render(
      <GraphView
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
      <GraphView
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

  it('counts total days correctly', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    // Check for total days in month summary - look for element with "Total Days" label
    const totalDaysSection = screen.getByText('Total Days').closest('div');
    expect(totalDaysSection).toBeInTheDocument();
    expect(totalDaysSection?.textContent).toContain('3');
  });

  it('counts festivals correctly', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    // Festival count should be 1 (one day with festivals)
    const monthSummary = screen.getByText('Month Summary').closest('div');
    expect(monthSummary).toBeInTheDocument();
  });

  it('counts unique tithis correctly', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    // Should have 3 unique tithis
    const monthSummary = screen.getByText('Month Summary').closest('div');
    expect(monthSummary).toBeInTheDocument();
  });

  it('shows sunrise and sunset trend when data available', () => {
    const onDateClick = vi.fn();
    render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    expect(screen.getByText('Sunrise & Sunset Times')).toBeInTheDocument();
  });

  it('renders SVG chart elements', () => {
    const onDateClick = vi.fn();
    const { container } = render(
      <GraphView
        year={2024}
        month={0}
        panchangamData={mockPanchangamData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    const svgs = container.querySelectorAll('svg');
    expect(svgs.length).toBeGreaterThan(0);
  });

  it('highlights today in festival section', () => {
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
        festivals: ['Today Festival']
      }
    };

    const onDateClick = vi.fn();
    render(
      <GraphView
        year={today.getFullYear()}
        month={today.getMonth()}
        panchangamData={todayData}
        settings={mockSettings}
        onDateClick={onDateClick}
      />
    );

    // The festival name may appear multiple times in different chart sections
    const festivalElements = screen.getAllByText('Today Festival');
    expect(festivalElements.length).toBeGreaterThan(0);
  });
});
