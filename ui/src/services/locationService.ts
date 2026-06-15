import type { Location } from '../types/panchangam';
import { popularLocations, US_UK_LOCATION_COUNT } from './locationData';

interface NominatimAddress {
  city?: string;
  town?: string;
  village?: string;
  hamlet?: string;
  state?: string;
  region?: string;
  country?: string;
}

interface NominatimResult {
  lat: string;
  lon: string;
  display_name?: string;
  address?: NominatimAddress;
}

const FAVORITE_COORDINATE_TOLERANCE = 0.001;

class LocationService {
  private readonly FAVORITES_KEY = 'panchangam_favorite_locations';

  // Get user's favorite locations from localStorage
  getFavoriteLocations(): Location[] {
    try {
      const stored = localStorage.getItem(this.FAVORITES_KEY);
      const parsed = stored ? JSON.parse(stored) : [];
      return Array.isArray(parsed) ? parsed : [];
    } catch {
      return [];
    }
  }

  // Add location to favorites
  addToFavorites(location: Location): void {
    try {
      const favorites = this.getFavoriteLocations();

      // Check if already exists (prevent duplicates)
      const exists = favorites.some(fav => this.isSameFavoriteLocation(fav, location));

      if (!exists) {
        favorites.unshift(location); // Add to beginning

        // Limit to 10 favorites
        if (favorites.length > 10) {
          favorites.pop();
        }

        localStorage.setItem(this.FAVORITES_KEY, JSON.stringify(favorites));
      }
    } catch {
      return;
    }
  }

  // Remove location from favorites
  removeFromFavorites(location: Location): void {
    try {
      const favorites = this.getFavoriteLocations();
      const filtered = favorites.filter(fav => !this.isSameFavoriteLocation(fav, location));
      localStorage.setItem(this.FAVORITES_KEY, JSON.stringify(filtered));
    } catch {
      return;
    }
  }

  // Check if location is in favorites
  isFavorite(location: Location): boolean {
    const favorites = this.getFavoriteLocations();
    return favorites.some(fav => this.isSameFavoriteLocation(fav, location));
  }

  private isSameFavoriteLocation(first: Location, second: Location): boolean {
    return Math.abs(first.latitude - second.latitude) < FAVORITE_COORDINATE_TOLERANCE &&
      Math.abs(first.longitude - second.longitude) < FAVORITE_COORDINATE_TOLERANCE;
  }

  // Get locations organized by category for US/UK users
  getLocationsByCategory(): { favorites: Location[], usUk: Location[], popular: Location[] } {
    const favorites = this.getFavoriteLocations();

    // US/UK locations (first 11 + UK locations we added)
    const usUk = popularLocations.slice(0, US_UK_LOCATION_COUNT);

    // Rest as popular (Indian + remaining US cities)
    const popular = popularLocations.slice(US_UK_LOCATION_COUNT);

    return { favorites, usUk, popular };
  }

  async getCurrentLocation(): Promise<Location> {
    return new Promise((resolve, reject) => {
      if (!navigator.geolocation) {
        reject(new Error('Geolocation is not supported'));
        return;
      }

      navigator.geolocation.getCurrentPosition(
        async (position) => {
          const { latitude, longitude } = position.coords;
          try {
            const location = await this.reverseGeocode(latitude, longitude);
            resolve(location);
          } catch {
            // Fallback to Milpitas, CA if reverse geocoding fails
            resolve(popularLocations[0]);
          }
        },
        () => {
          // Fallback to Milpitas, CA
          resolve(popularLocations[0]);
        }
      );
    });
  }

  async reverseGeocode(latitude: number, longitude: number): Promise<Location> {
    try {
      // Use Nominatim reverse geocoding API
      const url = `https://nominatim.openstreetmap.org/reverse?lat=${latitude}&lon=${longitude}&format=json&addressdetails=1`;

      const response = await fetch(url, {
        headers: {
          'User-Agent': 'PanchangamApp/1.0 (https://panchangam.app)'
        }
      });

      if (response.ok) {
        const data = await response.json() as NominatimResult;

        return {
          name: this.formatLocationName(data),
          latitude,
          longitude,
          timezone: this.inferTimezone(latitude, longitude),
          region: this.extractRegion(data)
        };
      }
    } catch {
      return this.createFallbackLocation(latitude, longitude);
    }

    // Fallback: Find closest popular location
    return this.createFallbackLocation(latitude, longitude);
  }

  private createFallbackLocation(latitude: number, longitude: number): Location {
    let closest = popularLocations[0];
    let minDistance = this.calculateDistance(latitude, longitude, closest.latitude, closest.longitude);

    for (const location of popularLocations) {
      const distance = this.calculateDistance(latitude, longitude, location.latitude, location.longitude);
      if (distance < minDistance) {
        minDistance = distance;
        closest = location;
      }
    }

    return {
      name: `Location (${latitude.toFixed(4)}, ${longitude.toFixed(4)})`,
      latitude,
      longitude,
      timezone: this.inferTimezone(latitude, longitude),
      region: closest.region
    };
  }

  private calculateDistance(lat1: number, lon1: number, lat2: number, lon2: number): number {
    const R = 6371; // Earth's radius in km
    const dLat = (lat2 - lat1) * Math.PI / 180;
    const dLon = (lon2 - lon1) * Math.PI / 180;
    const a = Math.sin(dLat/2) * Math.sin(dLat/2) +
              Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
              Math.sin(dLon/2) * Math.sin(dLon/2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
    return R * c;
  }

  getPopularLocations(): Location[] {
    return popularLocations;
  }

  async searchLocations(query: string): Promise<Location[]> {
    // First, search through popular locations for quick results
    const popularFiltered = popularLocations.filter(location =>
      location.name.toLowerCase().includes(query.toLowerCase())
    );

    // If we have good matches from popular locations, return them
    if (popularFiltered.length > 0 && query.length < 5) {
      return popularFiltered.slice(0, 10);
    }

    // For longer queries, try geocoding API for broader search
    try {
      const geocodedResults = await this.geocodeSearch(query);

      // Combine popular results with geocoded results, removing duplicates
      const allResults = [...popularFiltered];

      for (const geocoded of geocodedResults) {
        // Check if this location is already in our results (within 1km)
        const isDuplicate = allResults.some(existing =>
          this.calculateDistance(
            existing.latitude, existing.longitude,
            geocoded.latitude, geocoded.longitude
          ) < 1 // 1km threshold
        );

        if (!isDuplicate) {
          allResults.push(geocoded);
        }
      }

      return allResults.slice(0, 15); // Limit to 15 results
    } catch {
      return popularFiltered.slice(0, 10);
    }
  }

  private async geocodeSearch(query: string): Promise<Location[]> {
    // Using OpenStreetMap Nominatim API (free, no API key required)
    const encodedQuery = encodeURIComponent(query);
    const url = `https://nominatim.openstreetmap.org/search?q=${encodedQuery}&format=json&limit=10&addressdetails=1`;

    const response = await fetch(url, {
      headers: {
        'User-Agent': 'PanchangamApp/1.0 (https://panchangam.app)'
      }
    });

    if (!response.ok) {
      throw new Error(`Geocoding API error: ${response.status}`);
    }

    const data = await response.json() as NominatimResult[];

    return data.map((item) => ({
      name: this.formatLocationName(item),
      latitude: parseFloat(item.lat),
      longitude: parseFloat(item.lon),
      timezone: this.inferTimezone(parseFloat(item.lat), parseFloat(item.lon)),
      region: this.extractRegion(item)
    }));
  }

  private formatLocationName(item: NominatimResult): string {
    const address = item.address || {};
    const displayName = item.display_name || '';

    // Try to build a nice name from address components
    const city = address.city || address.town || address.village || address.hamlet;
    const state = address.state || address.region;
    const country = address.country;

    if (city && state && country !== 'United States') {
      return `${city}, ${state}, ${country}`;
    } else if (city && state) {
      return `${city}, ${state}`;
    } else if (city && country) {
      return `${city}, ${country}`;
    } else {
      // Fallback to shortened display name
      const parts = displayName.split(',');
      if (parts.length > 2) {
        return `${parts[0]}, ${parts[1]}`.trim();
      }
      return parts[0]?.trim() || displayName;
    }
  }

  private extractRegion(item: NominatimResult): string {
    const address = item.address || {};
    return address.state || address.region || address.country || 'Unknown';
  }

  private inferTimezone(latitude: number, longitude: number): string {
    // Simple timezone inference based on longitude
    // This is a basic implementation - for production, you'd want a proper timezone API

    // US timezones
    if (latitude > 25 && latitude < 49 && longitude > -125 && longitude < -66) {
      if (longitude > -90) return "America/New_York";
      if (longitude > -105) return "America/Chicago";
      if (longitude > -120) return "America/Denver";
      return "America/Los_Angeles";
    }

    // Alaska
    if (latitude > 60 && longitude > -170 && longitude < -140) {
      return "America/Anchorage";
    }

    // Hawaii
    if (latitude > 18 && latitude < 23 && longitude > -162 && longitude < -154) {
      return "Pacific/Honolulu";
    }

    // India
    if (latitude > 6 && latitude < 38 && longitude > 68 && longitude < 98) {
      return "Asia/Kolkata";
    }

    // Europe (rough approximation)
    if (latitude > 35 && latitude < 70 && longitude > -10 && longitude < 40) {
      if (longitude < 15) return "Europe/London";
      return "Europe/Berlin";
    }

    // Australia (rough approximation)
    if (latitude > -45 && latitude < -10 && longitude > 110 && longitude < 155) {
      return "Australia/Sydney";
    }

    // Default fallback based on longitude
    const utcOffset = Math.round(longitude / 15);
    if (utcOffset >= -12 && utcOffset <= 12) {
      return `Etc/GMT${utcOffset <= 0 ? '+' : '-'}${Math.abs(utcOffset)}`;
    }

    return "UTC";
  }
}

export const locationService = new LocationService();
