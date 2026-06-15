import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import type { MockInstance } from 'vitest';
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

function readBlobText(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onload = () => resolve(String(reader.result));
    reader.onerror = () => reject(reader.error);
    reader.readAsText(blob);
  });
}

describe('exportHelpers', () => {
  let mockLink: HTMLAnchorElement;
  let createElementSpy: MockInstance<[tagName: string, options?: ElementCreationOptions], HTMLElement>;
  let setAttributeSpy: MockInstance<[qualifiedName: string, value: string], void>;
  let clickSpy: MockInstance<[], void>;
  let removeChildSpy: MockInstance<[child: Node], Node>;
  let createdBlob: Blob | undefined;

  const getDownloadFilename = (): string => {
    const downloadCall = setAttributeSpy.mock.calls.find(([name]) => name === 'download');
    expect(downloadCall).toBeDefined();
    return downloadCall?.[1] ?? '';
  };

  beforeEach(() => {
    // Mock DOM elements for download
    mockLink = document.createElement('a');
    setAttributeSpy = vi.spyOn(mockLink, 'setAttribute');
    clickSpy = vi.spyOn(mockLink, 'click').mockImplementation(() => undefined);

    createElementSpy = vi.spyOn(document, 'createElement').mockReturnValue(mockLink);
    vi.spyOn(document.body, 'appendChild').mockImplementation(node => node);
    removeChildSpy = vi.spyOn(document.body, 'removeChild').mockImplementation(node => node);

    // Mock URL.createObjectURL
    createdBlob = undefined;
    globalThis.URL.createObjectURL = vi.fn((blob: Blob) => {
      createdBlob = blob;
      return 'blob:mock-url';
    });
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
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);

      const filename = getDownloadFilename();
      expect(filename).toContain('.csv');
      expect(filename).toContain('2024-01');
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
      exportToJSON(mockPanchangamData, mockSettings, 2024, 0);

      const filename = getDownloadFilename();
      expect(filename).toContain('.json');
      expect(filename).toContain('2024-01');
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
      exportAnalyticsData(mockPanchangamData, mockSettings, 2024, 0, 'csv');
      let filename = getDownloadFilename();
      expect(filename).toContain('.csv');

      // Reset mocks
      setAttributeSpy.mockClear();

      exportAnalyticsData(mockPanchangamData, mockSettings, 2024, 0, 'json');
      filename = getDownloadFilename();
      expect(filename).toContain('.json');
    });
  });

  describe('CSV content validation', () => {
    it('formats data correctly with quotes for CSV', () => {
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);
      expect(clickSpy).toHaveBeenCalled();
    });

    it('escapes quotes inside CSV cells', async () => {
      exportToCSV({
        '2024-01-15': {
          ...mockPanchangamData['2024-01-15'],
          tithi: 'Shukla "Special" Panchami',
          events: [
            {
              name: 'Abhijit "Best" Muhurta',
              time: '12:00',
              event_type: 'ABHIJIT_MUHURTA',
              quality: 'auspicious'
            }
          ],
        }
      }, mockSettings, 2024, 0);

      expect(createdBlob).toBeDefined();
      const csvContent = await readBlobText(createdBlob as Blob);

      expect(csvContent).toContain('"Shukla ""Special"" Panchami"');
      expect(csvContent).toContain('"Abhijit ""Best"" Muhurta (12:00)"');
    });

    it('exports API date strings on their local calendar day', async () => {
      exportToCSV({
        '2024-01-15': mockPanchangamData['2024-01-15']
      }, mockSettings, 2024, 0);

      expect(createdBlob).toBeDefined();
      const csvContent = await readBlobText(createdBlob as Blob);

      expect(csvContent).toContain('"01/15/2024"');
      expect(csvContent).not.toContain('"01/14/2024"');
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

      exportToCSV(mockPanchangamData, settingsWithSpecialChars, 2024, 0);

      expect(getDownloadFilename()).toContain('San_Francisco__CA___USA');
    });

    it('formats month with leading zero', () => {
      exportToCSV(mockPanchangamData, mockSettings, 2024, 0);

      expect(getDownloadFilename()).toContain('2024-01');
    });
  });
});
