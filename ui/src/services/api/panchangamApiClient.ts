import { apiClient } from './client';
import { PanchangamApiError } from './types';
import { requestCache } from './requestCache';
import { formatDateForApi, getCalendarDayDifference, parseApiDate } from '../../utils/dateHelpers';
import type { PanchangamData, GetPanchangamRequest, Event as PanchangamEvent } from '../../types/panchangam';

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

const WEEKDAYS = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
const PLANETARY_RULERS = ['Sun', 'Moon', 'Mars', 'Mercury', 'Jupiter', 'Venus', 'Saturn'];
const AUSPICIOUS_EVENT_TYPES = new Set(['ABHIJIT_MUHURTA', 'BRAHMA_MUHURTA', 'SUNRISE', 'FESTIVAL']);
const INAUSPICIOUS_EVENT_TYPES = new Set(['RAHU_KALAM', 'YAMAGANDAM', 'GULIKA_KALAM']);

function getVaraDetails(date: string): Pick<PanchangamData, 'vara' | 'planetary_ruler'> {
  const dayOfWeek = parseApiDate(date).getDay();

  return {
    vara: WEEKDAYS[dayOfWeek],
    planetary_ruler: PLANETARY_RULERS[dayOfWeek]
  };
}

function getEventQuality(eventType: string): PanchangamEvent['quality'] {
  if (AUSPICIOUS_EVENT_TYPES.has(eventType)) return 'auspicious';
  if (INAUSPICIOUS_EVENT_TYPES.has(eventType)) return 'inauspicious';
  return 'neutral';
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
function validatePanchangamResponse(data: unknown): asserts data is ApiPanchangamData {
  if (!data || typeof data !== 'object') {
    throw new PanchangamApiError('Invalid response format', 'INVALID_RESPONSE_FORMAT');
  }

  const response = data as Record<string, unknown>;
  const requiredFields = ['date', 'tithi', 'nakshatra', 'yoga', 'karana', 'sunrise_time', 'sunset_time'];

  for (const field of requiredFields) {
    if (!(field in response) || typeof response[field] !== 'string') {
      throw new PanchangamApiError(
        `Missing or invalid field: ${field}`,
        'INVALID_RESPONSE_FIELD'
      );
    }
  }

  if (!Array.isArray(response.events)) {
    throw new PanchangamApiError('Events must be an array', 'INVALID_EVENTS_FORMAT');
  }
}

/**
 * Transform API response to match UI types
 */
function transformApiResponse(apiData: ApiPanchangamData, requestDate: string): PanchangamData {
  const varaDetails = getVaraDetails(requestDate);

  // Extract lunar timing from events
  const moonriseEvent = apiData.events.find(e => e.event_type === 'MOONRISE');
  const moonsetEvent = apiData.events.find(e => e.event_type === 'MOONSET');

  // Extract tithi start time from events
  const tithiEvent = apiData.events.find(e => e.event_type === 'TITHI');
  const tithiStartTime = tithiEvent?.time;

  return {
    date: apiData.date,
    tithi: apiData.tithi,
    tithi_start_time: tithiStartTime,
    nakshatra: apiData.nakshatra,
    yoga: apiData.yoga,
    karana: apiData.karana,
    sunrise_time: apiData.sunrise_time,
    sunset_time: apiData.sunset_time,
    ...varaDetails,
    events: apiData.events.map(event => ({
      name: event.name,
      time: event.time,
      event_type: event.event_type as PanchangamEvent['event_type'],
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
  const varaDetails = getVaraDetails(date);

  return {
    date,
    tithi: "API Unavailable",
    nakshatra: "Please check connection",
    yoga: "Offline Mode",
    karana: "No Data",
    sunrise_time: "06:30:00",
    sunset_time: "18:30:00",
    ...varaDetails,
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

      const requestParams = {
        date: params.date,
        lat: params.latitude,
        lng: params.longitude,
        tz: params.timezone || 'UTC',
        region: params.region || '',
        method: params.calculation_method || 'traditional',
        locale: params.locale || 'en'
      };

      const endpoint = '/api/v1/panchangam';

      // Check cache first
      const cachedData = requestCache.get<PanchangamData>(endpoint, requestParams);
      if (cachedData) {
        return cachedData;
      }

      // Check for pending request to prevent duplicate calls
      const pendingRequest = requestCache.getPendingRequest<ApiPanchangamData>(endpoint, requestParams);
      if (pendingRequest) {
        const response = await pendingRequest;
        validatePanchangamResponse(response);
        return transformApiResponse(response, params.date);
      }

      // Create new request
      const apiPromise = apiClient.get<ApiPanchangamData>(endpoint, requestParams);
      requestCache.setPendingRequest(endpoint, requestParams, apiPromise.then(r => r.data));

      const response = await apiPromise;

      // Validate response data
      validatePanchangamResponse(response.data);

      // Transform data
      const transformedData = transformApiResponse(response.data, params.date);

      // Cache the transformed result
      requestCache.set(endpoint, requestParams, transformedData);

      return transformedData;

    } catch (error) {
      // Handle specific error types
      if (error instanceof PanchangamApiError) {
        // For network errors, provide fallback data
        if (error.code === 'NETWORK_ERROR' || error.code === 'REQUEST_TIMEOUT') {
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
    const start = parseApiDate(startDate);
    const end = parseApiDate(endDate);

    // Validate date range
    if (start > end) {
      throw new PanchangamApiError('Start date must be before end date', 'INVALID_DATE_RANGE');
    }

    const daysDiff = getCalendarDayDifference(start, end);
    if (daysDiff > 365) {
      throw new PanchangamApiError('Date range cannot exceed 365 days', 'DATE_RANGE_TOO_LARGE');
    }

    // Process dates with controlled concurrency to prevent flooding
    const results: PanchangamData[] = [];
    const batchSize = 5; // Process 5 requests at a time
    const dates: string[] = [];

    // Collect all dates first
    for (let d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
      dates.push(formatDateForApi(d));
    }

    // Process in batches with delay between batches
    for (let i = 0; i < dates.length; i += batchSize) {
      const batch = dates.slice(i, i + batchSize);
      const batchPromises = batch.map((dateStr, index) => {
        // Add small delay to stagger requests within batch
        return new Promise<PanchangamData>((resolve, reject) => {
          setTimeout(async () => {
            try {
              const result = await this.getPanchangam({ ...params, date: dateStr });
              resolve(result);
            } catch (error) {
              reject(error);
            }
          }, index * 100); // 100ms delay between requests in batch
        });
      });

      const batchResults = await Promise.allSettled(batchPromises);

      // Extract successful results from this batch
      batchResults.forEach((result) => {
        if (result.status === 'fulfilled') {
          results.push(result.value);
        }
      });

      // Delay between batches (except for last batch)
      if (i + batchSize < dates.length) {
        await new Promise(resolve => setTimeout(resolve, 200));
      }
    }

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
