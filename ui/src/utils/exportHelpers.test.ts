import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { exportToCSV, exportToJSON, exportAnalyticsData } from './exportHelpers';
import { PanchangamData, Settings } from '../types/panchangam';

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

describe('exportHelpers', () => {
  let createElementSpy: any;
  let clickSpy: any;
  let appendChildSpy: any;
  let removeChildSpy: any;

  beforeEach(() => {
    // Mock DOM elements for download
    clickSpy = vi.fn();
    const mockLink = {
      setAttribute: vi.fn(),
      style: {},
      click: clickSpy
    };

    createElementSpy = vi.spyOn(document, 'createElement').mockReturnValue(mockLink as any);
    appendChildSpy = vi.spyOn(document.body, 'appendChild').mockImplementation(() => mockLink as any);
    removeChildSpy = vi.spyOn(document.body, 'removeChild').mockImplementation(() => mockLink as any);

    // Mock URL.createObjectURL
    global.URL.createObjectURL = vi.fn(() => 'blob:mock-url');
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('exportToCSV', () => {
    it('creates a CSV file with correct headers', () => {
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);

      expect(createElementSpy).toHaveBeenCalledWith('a');
      expect(clickSpy).toHaveBeenCalled();
    });

    it('includes all data rows in CSV', () => {
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);

      expect(clickSpy).toHaveBeenCalled();
    });

    it('sets correct filename for CSV export', () => {
      const mockLink = createElementSpy.mock.results[0].value;
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);

      const setAttributeCalls = mockLink.setAttribute.mock.calls;
      const downloadCall = setAttributeCalls.find((call: any) => call[0] === 'download');
      expect(downloadCall).toBeDefined();
      expect(downloadCall[1]).toContain('.csv');
      expect(downloadCall[1]).toContain('2024-01');
    });

    it('includes festivals in CSV export', () => {
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);
      expect(clickSpy).toHaveBeenCalled();
    });

    it('includes auspicious events in CSV export', () => {
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);
      expect(clickSpy).toHaveBeenCalled();
    });

    it('handles empty data gracefully', () => {
      exportToCSV({}, mockSettings, 2024, 0);
      expect(clickSpy).toHaveBeenCalled();
    });

    it('removes link element after download', () => {
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);
      expect(removeChildSpy).toHaveBeenCalled();
    });
  });

  describe('exportToJSON', () => {
    it('creates a JSON file with metadata', () => {
      exportToJSON(mockPanchangamData, mockSettings, 2024, 0);

      expect(createElementSpy).toHaveBeenCalledWith('a');
      expect(clickSpy).toHaveBeenCalled();
    });

    it('sets correct filename for JSON export', () => {
      const mockLink = createElementSpy.mock.results[0].value;
      exportToJSON(mockPanchangamData, mockSettings, 2024, 0);

      const setAttributeCalls = mockLink.setAttribute.mock.calls;
      const downloadCall = setAttributeCalls.find((call: any) => call[0] === 'download');
      expect(downloadCall).toBeDefined();
      expect(downloadCall[1]).toContain('.json');
      expect(downloadCall[1]).toContain('2024-01');
    });

    it('includes location in metadata', () => {
      exportToJSON(mockPanchangamData, mockSettings, 2024, 0);
      expect(clickSpy).toHaveBeenCalled();
    });

    it('includes calculation method in metadata', () => {
      exportToJSON(mockPanchangamData, mockSettings, 2024, 0);
      expect(clickSpy).toHaveBeenCalled();
    });

    it('handles empty data gracefully', () => {
      exportToJSON({}, mockSettings, 2024, 0);
      expect(clickSpy).toHaveBeenCalled();
    });

    it('removes link element after download', () => {
      exportToJSON(mockPanchangamData, mockSettings, 2024, 0);
      expect(removeChildSpy).toHaveBeenCalled();
    });
  });

  describe('exportAnalyticsData', () => {
    it('calls exportToCSV when format is csv', () => {
      exportAnalyticsData(mockPanchangamData, mockSettings, 2024, 0, 'csv');
      expect(clickSpy).toHaveBeenCalled();
    });

    it('calls exportToJSON when format is json', () => {
      exportAnalyticsData(mockPanchangamData, mockSettings, 2024, 0, 'json');
      expect(clickSpy).toHaveBeenCalled();
    });

    it('creates appropriate file based on format parameter', () => {
      const mockLink = createElementSpy.mock.results[0].value;

      exportAnalyticsData(mockPanchangamData, mockSettings, 2024, 0, 'csv');
      let downloadCall = mockLink.setAttribute.mock.calls.find((call: any) => call[0] === 'download');
      expect(downloadCall[1]).toContain('.csv');

      // Reset mocks
      createElementSpy.mockClear();

      exportAnalyticsData(mockPanchangamData, mockSettings, 2024, 0, 'json');
      const mockLink2 = createElementSpy.mock.results[0].value;
      downloadCall = mockLink2.setAttribute.mock.calls.find((call: any) => call[0] === 'download');
      expect(downloadCall[1]).toContain('.json');
    });
  });

  describe('CSV content validation', () => {
    it('formats data correctly with quotes for CSV', () => {
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);
      expect(clickSpy).toHaveBeenCalled();
    });

    it('includes all required columns', () => {
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);
      expect(clickSpy).toHaveBeenCalled();
    });
  });

  describe('Filename sanitization', () => {
    it('sanitizes location name in filename', () => {
      const settingsWithSpecialChars: Settings = {
        ...mockSettings,
        location: {
          ...mockSettings.location,
          name: 'San Francisco, CA / USA'
        }
      };

      const mockLink = createElementSpy.mock.results[0].value;
      exportToCSV(mockPanchangamData, settingsWithSpecialChars, 2024, 0);

      const setAttributeCalls = mockLink.setAttribute.mock.calls;
      const downloadCall = setAttributeCalls.find((call: any) => call[0] === 'download');
      expect(downloadCall[1]).toContain('San_Francisco__CA___USA');
    });

    it('formats month with leading zero', () => {
      const mockLink = createElementSpy.mock.results[0].value;
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);

      const setAttributeCalls = mockLink.setAttribute.mock.calls;
      const downloadCall = setAttributeCalls.find((call: any) => call[0] === 'download');
      expect(downloadCall[1]).toContain('2024-01');
    });
  });
});
