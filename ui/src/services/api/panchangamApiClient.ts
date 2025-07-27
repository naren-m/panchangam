import { apiClient } from './client';
import { PanchangamApiError } from './types';
import type { PanchangamData, GetPanchangamRequest } from '../../types/panchangam';

// API Response interface that matches the actual gRPC response
interface ApiPanchangamEvent {
  name: string;
  time: string;
  event_type: string;
}

interface ApiPanchangamData {
  date: string;
  tithi: string;
  nakshatra: string;
  yoga: string;
  karana: string;
  sunrise_time: string;
  sunset_time: string;
  events: ApiPanchangamEvent[];
}

interface HealthCheckResponse {
  status: string;
  timestamp: string;
  version?: string;
}

/**
 * Validates panchangam request parameters
 */
function validatePanchangamRequest(params: GetPanchangamRequest): void {
  // Date validation
  if (!params.date) {
    throw new PanchangamApiError('Date is required', 'MISSING_DATE');
  }

  const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
  if (!dateRegex.test(params.date)) {
    throw new PanchangamApiError(
      'Invalid date format. Please use YYYY-MM-DD format.',
      'INVALID_DATE_FORMAT'
    );
  }

  // Latitude validation
  if (typeof params.latitude !== 'number') {
    throw new PanchangamApiError('Latitude must be a number', 'INVALID_LATITUDE_TYPE');
  }

  if (params.latitude < -90 || params.latitude > 90) {
    throw new PanchangamApiError(
      'Latitude must be between -90 and 90 degrees',
      'INVALID_LATITUDE_RANGE'
    );
  }

  // Longitude validation
  if (typeof params.longitude !== 'number') {
    throw new PanchangamApiError('Longitude must be a number', 'INVALID_LONGITUDE_TYPE');
  }

  if (params.longitude < -180 || params.longitude > 180) {
    throw new PanchangamApiError(
      'Longitude must be between -180 and 180 degrees',
      'INVALID_LONGITUDE_RANGE'
    );
  }

  // Timezone validation (if provided)
  if (params.timezone) {
    try {
      // Simple timezone validation - check if it's a valid IANA timezone
      Intl.DateTimeFormat(undefined, { timeZone: params.timezone });
    } catch {
      throw new PanchangamApiError(
        'Invalid timezone. Please use IANA timezone format (e.g., "Asia/Kolkata")',
        'INVALID_TIMEZONE'
      );
    }
  }
}

/**
 * Validates panchangam response data
 */
function validatePanchangamResponse(data: any): asserts data is ApiPanchangamData {
  if (!data || typeof data !== 'object') {
    throw new PanchangamApiError('Invalid response format', 'INVALID_RESPONSE_FORMAT');
  }

  const requiredFields = ['date', 'tithi', 'nakshatra', 'yoga', 'karana', 'sunrise_time', 'sunset_time'];
  
  for (const field of requiredFields) {
    if (!(field in data) || typeof data[field] !== 'string') {
      throw new PanchangamApiError(
        `Missing or invalid field: ${field}`,
        'INVALID_RESPONSE_FIELD'
      );
    }
  }

  if (!Array.isArray(data.events)) {
    throw new PanchangamApiError('Events must be an array', 'INVALID_EVENTS_FORMAT');
  }
}

/**
 * Transform API response to match UI types
 */
function transformApiResponse(apiData: ApiPanchangamData, requestDate: string): PanchangamData {
  // Extract day of week for vara calculation - parse date reliably
  const [year, month, day] = requestDate.split('-').map(Number);
  const dateObj = new Date(year, month - 1, day);
  const weekdays = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
  const rulers = ["Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"];
  const dayOfWeek = dateObj.getDay();

  // Extract lunar timing from events
  const moonriseEvent = apiData.events.find(e => e.event_type === 'MOONRISE');
  const moonsetEvent = apiData.events.find(e => e.event_type === 'MOONSET');

  // Determine event quality based on type
  const getEventQuality = (eventType: string): 'auspicious' | 'inauspicious' | 'neutral' => {
    const auspiciousEvents = ['ABHIJIT_MUHURTA', 'BRAHMA_MUHURTA', 'SUNRISE', 'FESTIVAL'];
    const inauspiciousEvents = ['RAHU_KALAM', 'YAMAGANDAM', 'GULIKA_KALAM'];
    
    if (auspiciousEvents.includes(eventType)) return 'auspicious';
    if (inauspiciousEvents.includes(eventType)) return 'inauspicious';
    return 'neutral';
  };

  return {
    date: apiData.date,
    tithi: apiData.tithi,
    nakshatra: apiData.nakshatra,
    yoga: apiData.yoga,
    karana: apiData.karana,
    sunrise_time: apiData.sunrise_time,
    sunset_time: apiData.sunset_time,
    vara: weekdays[dayOfWeek],
    planetary_ruler: rulers[dayOfWeek],
    events: apiData.events.map(event => ({
      name: event.name,
      time: event.time,
      event_type: event.event_type as any,
      quality: getEventQuality(event.event_type)
    })),
    festivals: [], // Not provided by current API
    moonrise_time: moonriseEvent?.time,
    moonset_time: moonsetEvent?.time,
  };
}

/**
 * Generate fallback data when API is unavailable
 */
function generateFallbackData(date: string): PanchangamData {
  // Parse date more reliably to avoid timezone issues
  const [year, month, day] = date.split('-').map(Number);
  const dateObj = new Date(year, month - 1, day);
  const weekdays = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
  const rulers = ["Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"];
  const dayOfWeek = dateObj.getDay();

  return {
    date,
    tithi: "API Unavailable",
    nakshatra: "Please check connection",
    yoga: "Offline Mode",
    karana: "No Data",
    sunrise_time: "06:30:00",
    sunset_time: "18:30:00",
    vara: weekdays[dayOfWeek],
    planetary_ruler: rulers[dayOfWeek],
    events: [
      {
        name: "API Connection Error",
        time: "00:00:00",
        event_type: "MUHURTA" as const,
        quality: "neutral" as const
      }
    ],
    festivals: [],
    moonrise_time: undefined,
    moonset_time: undefined,
  };
}

/**
 * Panchangam API client with validation and error handling
 */
export class PanchangamApiClient {
  /**
   * Get panchangam data for a specific date
   */
  async getPanchangam(params: GetPanchangamRequest): Promise<PanchangamData> {
    try {
      // Validate input parameters
      validatePanchangamRequest(params);

      // Make API call
      const response = await apiClient.get<ApiPanchangamData>('/api/v1/panchangam', {
        date: params.date,
        lat: params.latitude,
        lng: params.longitude,
        tz: params.timezone || 'UTC',
        region: params.region || '',
        method: params.calculation_method || 'traditional',
        locale: params.locale || 'en'
      });

      // Validate response data
      validatePanchangamResponse(response.data);

      // Transform and return
      return transformApiResponse(response.data, params.date);

    } catch (error) {
      console.error('Panchangam API error:', error);

      // Handle specific error types
      if (error instanceof PanchangamApiError) {
        // For network errors, provide fallback data
        if (error.code === 'NETWORK_ERROR' || error.code === 'REQUEST_TIMEOUT') {
          console.warn('API unavailable, using fallback data');
          return generateFallbackData(params.date);
        }
        // Re-throw other API errors (like validation errors)
        throw error;
      }

      // Re-throw unexpected errors as they are
      throw error;
    }
  }

  /**
   * Get panchangam data for a date range
   */
  async getPanchangamRange(
    startDate: string, 
    endDate: string, 
    params: Omit<GetPanchangamRequest, 'date'>
  ): Promise<PanchangamData[]> {
    const start = new Date(startDate);
    const end = new Date(endDate);
    
    // Validate date range
    if (start > end) {
      throw new PanchangamApiError('Start date must be before end date', 'INVALID_DATE_RANGE');
    }

    const daysDiff = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24));
    if (daysDiff > 365) {
      throw new PanchangamApiError('Date range cannot exceed 365 days', 'DATE_RANGE_TOO_LARGE');
    }

    // Process dates in parallel for better performance
    const datePromises: Promise<PanchangamData>[] = [];
    
    for (let d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
      const dateStr = d.toISOString().split('T')[0];
      datePromises.push(this.getPanchangam({ ...params, date: dateStr }));
    }

    // Wait for all requests to complete
    const allResults = await Promise.allSettled(datePromises);
    const results: PanchangamData[] = [];
    
    // Extract successful results and log failures
    allResults.forEach((result, index) => {
      if (result.status === 'fulfilled') {
        results.push(result.value);
      } else {
        const date = new Date(start);
        date.setDate(date.getDate() + index);
        console.error(`Failed to fetch data for ${date.toISOString().split('T')[0]}:`, result.reason);
      }
    });

    return results;
  }

  /**
   * Health check to test API availability
   */
  async healthCheck(): Promise<{ status: 'healthy' | 'unhealthy'; message: string; timestamp?: string }> {
    try {
      const response = await apiClient.get<HealthCheckResponse>('/api/v1/health', undefined, {
        timeout: 5000 // 5 second timeout for health check
      });

      return {
        status: 'healthy',
        message: 'API is accessible',
        timestamp: response.data.timestamp
      };

    } catch (error) {
      if (error instanceof PanchangamApiError) {
        return {
          status: 'unhealthy',
          message: `API health check failed: ${error.message}`
        };
      }

      return {
        status: 'unhealthy',
        message: 'Unknown error during health check'
      };
    }
  }
}

// Create and export singleton instance
export const panchangamApiClient = new PanchangamApiClient();