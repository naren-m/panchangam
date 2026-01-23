import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { FiveAngas } from '../DayDetail/FiveAngas';
import { PanchangamData, Settings } from '../../types/panchangam';

describe('FiveAngas Component', () => {
  const mockData: PanchangamData = {
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
  };

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
  };

  it('renders all five Panchangam elements', () => {
    render(<FiveAngas data={mockData} settings={mockSettings} />);
    
    expect(screen.getByText('Tithi')).toBeInTheDocument();
    expect(screen.getByText('Nakshatra')).toBeInTheDocument();
    expect(screen.getByText('Yoga')).toBeInTheDocument();
    expect(screen.getByText('Karana')).toBeInTheDocument();
    expect(screen.getByText('Vara')).toBeInTheDocument();
  });

  it('displays correct values for each element', () => {
    render(<FiveAngas data={mockData} settings={mockSettings} />);
    
    expect(screen.getByText('Chaturthi (4)')).toBeInTheDocument();
    expect(screen.getByText('Uttara Bhadrapada (26)')).toBeInTheDocument();
    expect(screen.getByText('Siddha (21)')).toBeInTheDocument();
    expect(screen.getByText('Gara (6)')).toBeInTheDocument();
    expect(screen.getByText(/Monday.*Moon/)).toBeInTheDocument();
  });

  it('displays descriptions for each element', () => {
    render(<FiveAngas data={mockData} settings={mockSettings} />);
    
    expect(screen.getByText('Lunar day phase')).toBeInTheDocument();
    expect(screen.getByText('Lunar mansion')).toBeInTheDocument();
    expect(screen.getByText('Sun-Moon combination')).toBeInTheDocument();
    expect(screen.getByText('Half-tithi period')).toBeInTheDocument();
    expect(screen.getByText('Weekday and ruler')).toBeInTheDocument();
  });

  it('renders icons for each element', () => {
    const { container } = render(<FiveAngas data={mockData} settings={mockSettings} />);
    
    // Check for SVG icons (lucide-react renders SVGs)
    const svgs = container.querySelectorAll('svg');
    expect(svgs.length).toBeGreaterThanOrEqual(5);
  });

  it('handles missing optional data gracefully', () => {
    const incompleteData: PanchangamData = {
      ...mockData,
      vara: '',
      planetary_ruler: '',
    };
    
    render(<FiveAngas data={incompleteData} settings={mockSettings} />);
    
    // Should still render without errors
    expect(screen.getByText('Vara')).toBeInTheDocument();
  });
});
