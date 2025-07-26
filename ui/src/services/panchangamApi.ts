import { PanchangamData, GetPanchangamRequest } from '../types/panchangam';

// Configuration for the API endpoint
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
const API_ENDPOINT = `${API_BASE_URL}/api/v1/panchangam`;

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

// Transform API response to match UI types
const transformApiResponse = (apiData: ApiPanchangamData, requestDate: string): PanchangamData => {
  // Extract day of week for vara calculation
  const dateObj = new Date(requestDate);
  const weekdays = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
  const rulers = ["Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"];
  const dayOfWeek = dateObj.getDay();

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
      event_type: event.event_type as any, // Type assertion for event types
      quality: 'neutral' as const // Default quality since API doesn't provide this
    })),
    festivals: [], // Not provided by current API
    moonrise_time: undefined, // Not provided by current API
    moonset_time: undefined, // Not provided by current API
  };
};

class PanchangamApiService {
  async getPanchangam(request: GetPanchangamRequest): Promise<PanchangamData> {
    try {
      // Build query parameters
      const params = new URLSearchParams({
        date: request.date,
        lat: request.latitude.toString(),
        lng: request.longitude.toString(),
        tz: request.timezone || 'UTC',
        region: request.region || '',
        method: request.calculation_method || 'traditional',
        locale: request.locale || 'en'
      });

      // Make API call to real HTTP gateway
      const response = await fetch(`${API_ENDPOINT}?${params}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
        // Add timeout to prevent hanging requests
        signal: AbortSignal.timeout(30000), // 30 second timeout
      });

      if (!response.ok) {
        // Try to extract error details from response
        let errorMessage = `API request failed: ${response.status} ${response.statusText}`;
        try {
          const errorData = await response.json();
          if (errorData.error && errorData.error.message) {
            errorMessage = errorData.error.message;
          }
        } catch {
          // Ignore JSON parsing errors for error response
        }
        throw new Error(errorMessage);
      }

      const apiData: ApiPanchangamData = await response.json();
      
      // Transform API response to match UI expectations
      return transformApiResponse(apiData, request.date);
      
    } catch (error) {
      console.error('Panchangam API error:', error);
      
      // Provide fallback data if API fails
      if (error instanceof Error && error.message.includes('Failed to fetch')) {
        console.warn('API unavailable, using fallback data');
        return this.getFallbackData(request.date);
      }
      
      throw error;
    }
  }

  async getPanchangamRange(startDate: string, endDate: string, request: Omit<GetPanchangamRequest, 'date'>): Promise<PanchangamData[]> {
    const start = new Date(startDate);
    const end = new Date(endDate);
    const results: PanchangamData[] = [];

    // Process dates in parallel for better performance
    const datePromises: Promise<PanchangamData>[] = [];
    
    for (let d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
      const dateStr = d.toISOString().split('T')[0];
      datePromises.push(this.getPanchangam({ ...request, date: dateStr }));
    }

    // Wait for all requests to complete
    const allResults = await Promise.allSettled(datePromises);
    
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

  // Fallback data generator for when API is unavailable
  private getFallbackData(date: string): PanchangamData {
    const dateObj = new Date(date);
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
          event_type: "MUHURTA" as any,
          quality: "neutral" as const
        }
      ],
      festivals: [],
      moonrise_time: undefined,
      moonset_time: undefined,
    };
  }

  // Health check method to test API availability
  async healthCheck(): Promise<{ status: 'healthy' | 'unhealthy', message: string }> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/health`, {
        method: 'GET',
        signal: AbortSignal.timeout(5000), // 5 second timeout for health check
      });
      
      if (response.ok) {
        return { status: 'healthy', message: 'API is accessible' };
      } else {
        return { status: 'unhealthy', message: `API returned ${response.status}` };
      }
    } catch (error) {
      return { 
        status: 'unhealthy', 
        message: error instanceof Error ? error.message : 'Unknown error' 
      };
    }
  }
}

export const panchangamApi = new PanchangamApiService();

// Export the API configuration for debugging
export const apiConfig = {
  baseUrl: API_BASE_URL,
  endpoint: API_ENDPOINT,
};