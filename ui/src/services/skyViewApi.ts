// Sky View API Service
// Fetches real-time planetary positions from the backend

import { CelestialObject } from '../types/skyVisualization';

export interface SkyViewRequest {
  date?: Date;
  time?: string;
  latitude: number;
  longitude: number;
  altitude?: number;
  timezone?: string;
}

export interface SkyViewResponse {
  timestamp: string;
  observer: {
    latitude: number;
    longitude: number;
    altitude: number;
    timezone: string;
  };
  bodies: CelestialBody[];
  visible_bodies: CelestialObject[];
  julian_day: number;
  local_sidereal_time: number;
}

// Get API base URL from environment or default to current origin
const getApiBaseUrl = (): string => {
  // Check for explicit environment variable first
  if (import.meta.env.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL;
  }
  // In browser, use current origin (nginx proxies /api/ to gateway)
  if (typeof window !== 'undefined') {
    return window.location.origin;
  }
  return 'http://localhost:8080';
};

/**
 * Fetches sky view data from the backend API
 * @param request Sky view request parameters
 * @returns Promise with sky view response
 */
export const fetchSkyView = async (request: SkyViewRequest): Promise<SkyViewResponse> => {
  const apiBaseUrl = getApiBaseUrl();

  // Build query parameters
  const params = new URLSearchParams({
    lat: request.latitude.toString(),
    lng: request.longitude.toString(),
  });

  if (request.date) {
    const dateStr = request.date.toISOString().split('T')[0]; // YYYY-MM-DD
    params.append('date', dateStr);

    if (request.time) {
      params.append('time', request.time);
    } else {
      // Add time from Date object
      const hours = request.date.getHours().toString().padStart(2, '0');
      const minutes = request.date.getMinutes().toString().padStart(2, '0');
      const seconds = request.date.getSeconds().toString().padStart(2, '0');
      params.append('time', `${hours}:${minutes}:${seconds}`);
    }
  }

  if (request.altitude !== undefined) {
    params.append('alt', request.altitude.toString());
  }

  if (request.timezone) {
    params.append('tz', request.timezone);
  }

  const url = `${apiBaseUrl}/api/v1/sky-view?${params.toString()}`;

  try {
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(
        errorData.error?.message || `HTTP ${response.status}: ${response.statusText}`
      );
    }

    const data = await response.json();
    return data as SkyViewResponse;
  } catch (error) {
    console.error('Failed to fetch sky view data:', error);
    throw error;
  }
};

/**
 * Converts backend celestial body format to frontend format
 * @param backendBody Backend celestial body
 * @returns Frontend celestial object
 */
export const convertToFrontendFormat = (backendBody: any): CelestialObject => {
  return {
    id: backendBody.id,
    name: backendBody.name,
    type: backendBody.type,
    coordinates: {
      ecliptic: {
        longitude: backendBody.ecliptic_coords.longitude,
        latitude: backendBody.ecliptic_coords.latitude,
        distance: backendBody.ecliptic_coords.distance,
      },
      equatorial: backendBody.equatorial_coords ? {
        rightAscension: backendBody.equatorial_coords.right_ascension,
        declination: backendBody.equatorial_coords.declination,
        distance: backendBody.equatorial_coords.distance,
      } : undefined,
      horizontal: backendBody.horizontal_coords ? {
        azimuth: backendBody.horizontal_coords.azimuth,
        altitude: backendBody.horizontal_coords.altitude,
        distance: backendBody.horizontal_coords.distance,
      } : undefined,
    },
    magnitude: backendBody.magnitude,
    color: backendBody.color,
    size: calculateSizeFromMagnitude(backendBody.magnitude),
    metadata: {
      sanskritName: backendBody.sanskrit_name,
      hindiName: backendBody.hindi_name,
      isVisible: backendBody.is_visible,
      ...backendBody.metadata,
    },
  };
};

/**
 * Calculate visual size from magnitude
 * Brighter objects (lower magnitude) appear larger
 */
function calculateSizeFromMagnitude(magnitude: number): number {
  if (magnitude < -10) {
    // Sun
    return 5;
  } else if (magnitude < 0) {
    // Bright planets and moon
    return 4;
  } else if (magnitude < 2) {
    // Visible planets
    return 2.5;
  } else if (magnitude < 4) {
    return 1.5;
  } else if (magnitude < 6) {
    return 1;
  } else {
    return 0.5;
  }
}

/**
 * Validates latitude value
 */
export const isValidLatitude = (lat: number): boolean => {
  return !isNaN(lat) && lat >= -90 && lat <= 90;
};

/**
 * Validates longitude value
 */
export const isValidLongitude = (lng: number): boolean => {
  return !isNaN(lng) && lng >= -180 && lng <= 180;
};
