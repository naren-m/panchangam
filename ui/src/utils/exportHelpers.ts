import { PanchangamData, Settings } from '../types/panchangam';
import { parseApiDate } from './dateHelpers';

function csvCell(value: string): string {
  return `"${value.replace(/"/g, '""')}"`;
}

function exportFilename(year: number, month: number, locationName: string, extension: string): string {
  const monthNumber = String(month + 1).padStart(2, '0');
  const safeLocationName = locationName.replace(/[^a-z0-9]/gi, '_');

  return `panchangam_${year}-${monthNumber}_${safeLocationName}.${extension}`;
}

function downloadTextFile(content: string, type: string, filename: string): void {
  const blob = new Blob([content], { type });
  const link = document.createElement('a');
  const url = URL.createObjectURL(blob);

  link.setAttribute('href', url);
  link.setAttribute('download', filename);
  link.style.visibility = 'hidden';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

/**
 * Export panchangam data to CSV format
 */
export function exportToCSV(
  panchangamData: Record<string, PanchangamData>,
  settings: Settings,
  year: number,
  month: number
): void {
  // Prepare CSV header
  const headers = [
    'Date',
    'Day',
    'Tithi',
    'Nakshatra',
    'Yoga',
    'Karana',
    'Sunrise',
    'Sunset',
    'Moonrise',
    'Moonset',
    'Planetary Ruler',
    'Festivals',
    'Auspicious Events'
  ];

  // Sort data by date
  const sortedEntries = Object.entries(panchangamData).sort(([a], [b]) => a.localeCompare(b));

  // Convert data to CSV rows
  const rows = sortedEntries.map(([dateStr, data]) => {
    const date = parseApiDate(dateStr);
    const formattedDate = date.toLocaleDateString(settings.locale, {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    });

    const festivals = data.festivals?.join('; ') || '';
    const auspiciousEvents = data.events
      ?.filter(e => e.quality === 'auspicious')
      .map(e => `${e.name} (${e.time})`)
      .join('; ') || '';

    return [
      formattedDate,
      data.vara,
      data.tithi,
      data.nakshatra,
      data.yoga,
      data.karana,
      data.sunrise_time,
      data.sunset_time,
      data.moonrise_time || '-',
      data.moonset_time || '-',
      data.planetary_ruler,
      festivals,
      auspiciousEvents
    ];
  });

  // Combine headers and rows
  const csvContent = [
    headers.join(','),
    ...rows.map(row => row.map(csvCell).join(','))
  ].join('\n');

  downloadTextFile(
    csvContent,
    'text/csv;charset=utf-8;',
    exportFilename(year, month, settings.location.name, 'csv')
  );
}

/**
 * Export panchangam data to JSON format
 */
export function exportToJSON(
  panchangamData: Record<string, PanchangamData>,
  settings: Settings,
  year: number,
  month: number
): void {
  // Prepare export data with metadata
  const exportData = {
    metadata: {
      generated_at: new Date().toISOString(),
      location: settings.location,
      calculation_method: settings.calculation_method,
      locale: settings.locale,
      region: settings.region,
      year,
      month: month + 1,
      month_name: new Date(year, month, 1).toLocaleDateString(settings.locale, { month: 'long' })
    },
    data: panchangamData
  };

  // Convert to JSON with pretty printing
  const jsonContent = JSON.stringify(exportData, null, 2);

  downloadTextFile(
    jsonContent,
    'application/json;charset=utf-8;',
    exportFilename(year, month, settings.location.name, 'json')
  );
}

/**
 * Export specific view data (used by graph view analytics)
 */
export function exportAnalyticsData(
  panchangamData: Record<string, PanchangamData>,
  settings: Settings,
  year: number,
  month: number,
  format: 'csv' | 'json'
): void {
  if (format === 'csv') {
    exportToCSV(panchangamData, settings, year, month);
  } else {
    exportToJSON(panchangamData, settings, year, month);
  }
}
