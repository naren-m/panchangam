import { Location } from '../types/panchangam';

const popularLocations: Location[] = [
  // US/UK Focused Locations - Default Location First
  {
    name: "Milpitas, California",
    latitude: 37.4323,
    longitude: -121.9066,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  
  // UK Locations - Major Cities
  {
    name: "London, England",
    latitude: 51.5074,
    longitude: -0.1278,
    timezone: "Europe/London",
    region: "England"
  },
  {
    name: "Birmingham, England",
    latitude: 52.4862,
    longitude: -1.8904,
    timezone: "Europe/London",
    region: "England"
  },
  {
    name: "Manchester, England",
    latitude: 53.4808,
    longitude: -2.2426,
    timezone: "Europe/London",
    region: "England"
  },
  {
    name: "Leeds, England",
    latitude: 53.8008,
    longitude: -1.5491,
    timezone: "Europe/London",
    region: "England"
  },
  {
    name: "Glasgow, Scotland",
    latitude: 55.8642,
    longitude: -4.2518,
    timezone: "Europe/London",
    region: "Scotland"
  },
  {
    name: "Liverpool, England",
    latitude: 53.4084,
    longitude: -2.9916,
    timezone: "Europe/London",
    region: "England"
  },
  {
    name: "Newcastle, England",
    latitude: 54.9783,
    longitude: -1.6178,
    timezone: "Europe/London",
    region: "England"
  },
  {
    name: "Sheffield, England",
    latitude: 53.3811,
    longitude: -1.4701,
    timezone: "Europe/London",
    region: "England"
  },
  {
    name: "Bristol, England",
    latitude: 51.4545,
    longitude: -2.5879,
    timezone: "Europe/London",
    region: "England"
  },
  {
    name: "Edinburgh, Scotland",
    latitude: 55.9533,
    longitude: -3.1883,
    timezone: "Europe/London",
    region: "Scotland"
  },
  
  // Indian Locations
  {
    name: "Chennai, Tamil Nadu",
    latitude: 13.0827,
    longitude: 80.2707,
    timezone: "Asia/Kolkata",
    region: "Tamil Nadu"
  },
  {
    name: "Mumbai, Maharashtra", 
    latitude: 19.0760,
    longitude: 72.8777,
    timezone: "Asia/Kolkata",
    region: "Maharashtra"
  },
  {
    name: "Delhi, India",
    latitude: 28.6139,
    longitude: 77.2090,
    timezone: "Asia/Kolkata",
    region: "Delhi"
  },
  {
    name: "Bangalore, Karnataka",
    latitude: 12.9716,
    longitude: 77.5946,
    timezone: "Asia/Kolkata",
    region: "Karnataka"
  },
  {
    name: "Varanasi, Uttar Pradesh",
    latitude: 25.3176,
    longitude: 82.9739,
    timezone: "Asia/Kolkata",
    region: "Uttar Pradesh"
  },
  {
    name: "Tirupati, Andhra Pradesh",
    latitude: 13.6288,
    longitude: 79.4192,
    timezone: "Asia/Kolkata",
    region: "Andhra Pradesh"
  },
  
  // US Locations - Major Cities
  {
    name: "New York, NY",
    latitude: 40.7128,
    longitude: -74.0060,
    timezone: "America/New_York",
    region: "New York"
  },
  {
    name: "Los Angeles, CA",
    latitude: 34.0522,
    longitude: -118.2437,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Chicago, IL",
    latitude: 41.8781,
    longitude: -87.6298,
    timezone: "America/Chicago",
    region: "Illinois"
  },
  {
    name: "Houston, TX",
    latitude: 29.7604,
    longitude: -95.3698,
    timezone: "America/Chicago",
    region: "Texas"
  },
  {
    name: "Phoenix, AZ",
    latitude: 33.4484,
    longitude: -112.0740,
    timezone: "America/Phoenix",
    region: "Arizona"
  },
  {
    name: "Philadelphia, PA",
    latitude: 39.9526,
    longitude: -75.1652,
    timezone: "America/New_York",
    region: "Pennsylvania"
  },
  {
    name: "San Antonio, TX",
    latitude: 29.4241,
    longitude: -98.4936,
    timezone: "America/Chicago",
    region: "Texas"
  },
  {
    name: "San Diego, CA",
    latitude: 32.7157,
    longitude: -117.1611,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Dallas, TX",
    latitude: 32.7767,
    longitude: -96.7970,
    timezone: "America/Chicago",
    region: "Texas"
  },
  {
    name: "San Jose, CA",
    latitude: 37.3382,
    longitude: -121.8863,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Austin, TX",
    latitude: 30.2672,
    longitude: -97.7431,
    timezone: "America/Chicago",
    region: "Texas"
  },
  {
    name: "Jacksonville, FL",
    latitude: 30.3322,
    longitude: -81.6557,
    timezone: "America/New_York",
    region: "Florida"
  },
  {
    name: "San Francisco, CA",
    latitude: 37.7749,
    longitude: -122.4194,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Columbus, OH",
    latitude: 39.9612,
    longitude: -82.9988,
    timezone: "America/New_York",
    region: "Ohio"
  },
  {
    name: "Charlotte, NC",
    latitude: 35.2271,
    longitude: -80.8431,
    timezone: "America/New_York",
    region: "North Carolina"
  },
  {
    name: "Fort Worth, TX",
    latitude: 32.7555,
    longitude: -97.3308,
    timezone: "America/Chicago",
    region: "Texas"
  },
  {
    name: "Indianapolis, IN",
    latitude: 39.7684,
    longitude: -86.1581,
    timezone: "America/Indiana/Indianapolis",
    region: "Indiana"
  },
  {
    name: "Seattle, WA",
    latitude: 47.6062,
    longitude: -122.3321,
    timezone: "America/Los_Angeles",
    region: "Washington"
  },
  {
    name: "Denver, CO",
    latitude: 39.7392,
    longitude: -104.9903,
    timezone: "America/Denver",
    region: "Colorado"
  },
  {
    name: "Washington, DC",
    latitude: 38.9072,
    longitude: -77.0369,
    timezone: "America/New_York",
    region: "District of Columbia"
  },
  {
    name: "Boston, MA",
    latitude: 42.3601,
    longitude: -71.0589,
    timezone: "America/New_York",
    region: "Massachusetts"
  },
  {
    name: "Nashville, TN",
    latitude: 36.1627,
    longitude: -86.7816,
    timezone: "America/Chicago",
    region: "Tennessee"
  },
  {
    name: "El Paso, TX",
    latitude: 31.7619,
    longitude: -106.4850,
    timezone: "America/Denver",
    region: "Texas"
  },
  {
    name: "Detroit, MI",
    latitude: 42.3314,
    longitude: -83.0458,
    timezone: "America/Detroit",
    region: "Michigan"
  },
  {
    name: "Portland, OR",
    latitude: 45.5152,
    longitude: -122.6784,
    timezone: "America/Los_Angeles",
    region: "Oregon"
  },
  {
    name: "Las Vegas, NV",
    latitude: 36.1699,
    longitude: -115.1398,
    timezone: "America/Los_Angeles",
    region: "Nevada"
  },
  {
    name: "Memphis, TN",
    latitude: 35.1495,
    longitude: -90.0490,
    timezone: "America/Chicago",
    region: "Tennessee"
  },
  {
    name: "Louisville, KY",
    latitude: 38.2527,
    longitude: -85.7585,
    timezone: "America/New_York",
    region: "Kentucky"
  },
  {
    name: "Baltimore, MD",
    latitude: 39.2904,
    longitude: -76.6122,
    timezone: "America/New_York",
    region: "Maryland"
  },
  {
    name: "Milwaukee, WI",
    latitude: 43.0389,
    longitude: -87.9065,
    timezone: "America/Chicago",
    region: "Wisconsin"
  },
  {
    name: "Albuquerque, NM",
    latitude: 35.0844,
    longitude: -106.6504,
    timezone: "America/Denver",
    region: "New Mexico"
  },
  {
    name: "Tucson, AZ",
    latitude: 32.2226,
    longitude: -110.9747,
    timezone: "America/Phoenix",
    region: "Arizona"
  },
  {
    name: "Fresno, CA",
    latitude: 36.7378,
    longitude: -119.7871,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Mesa, AZ",
    latitude: 33.4152,
    longitude: -111.8315,
    timezone: "America/Phoenix",
    region: "Arizona"
  },
  {
    name: "Sacramento, CA",
    latitude: 38.5816,
    longitude: -121.4944,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Atlanta, GA",
    latitude: 33.7490,
    longitude: -84.3880,
    timezone: "America/New_York",
    region: "Georgia"
  },
  {
    name: "Kansas City, MO",
    latitude: 39.0997,
    longitude: -94.5786,
    timezone: "America/Chicago",
    region: "Missouri"
  },
  {
    name: "Colorado Springs, CO",
    latitude: 38.8339,
    longitude: -104.8214,
    timezone: "America/Denver",
    region: "Colorado"
  },
  {
    name: "Miami, FL",
    latitude: 25.7617,
    longitude: -80.1918,
    timezone: "America/New_York",
    region: "Florida"
  },
  {
    name: "Raleigh, NC",
    latitude: 35.7796,
    longitude: -78.6382,
    timezone: "America/New_York",
    region: "North Carolina"
  },
  {
    name: "Omaha, NE",
    latitude: 41.2565,
    longitude: -95.9345,
    timezone: "America/Chicago",
    region: "Nebraska"
  },
  {
    name: "Long Beach, CA",
    latitude: 33.7701,
    longitude: -118.1937,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Virginia Beach, VA",
    latitude: 36.8529,
    longitude: -75.9780,
    timezone: "America/New_York",
    region: "Virginia"
  },
  {
    name: "Oakland, CA",
    latitude: 37.8044,
    longitude: -122.2711,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Minneapolis, MN",
    latitude: 44.9778,
    longitude: -93.2650,
    timezone: "America/Chicago",
    region: "Minnesota"
  },
  {
    name: "Tulsa, OK",
    latitude: 36.1540,
    longitude: -95.9928,
    timezone: "America/Chicago",
    region: "Oklahoma"
  },
  {
    name: "Arlington, TX",
    latitude: 32.7357,
    longitude: -97.1081,
    timezone: "America/Chicago",
    region: "Texas"
  },
  {
    name: "Tampa, FL",
    latitude: 27.9506,
    longitude: -82.4572,
    timezone: "America/New_York",
    region: "Florida"
  },
  {
    name: "New Orleans, LA",
    latitude: 29.9511,
    longitude: -90.0715,
    timezone: "America/Chicago",
    region: "Louisiana"
  },
  {
    name: "Wichita, KS",
    latitude: 37.6872,
    longitude: -97.3301,
    timezone: "America/Chicago",
    region: "Kansas"
  },
  {
    name: "Cleveland, OH",
    latitude: 41.4993,
    longitude: -81.6944,
    timezone: "America/New_York",
    region: "Ohio"
  },
  {
    name: "Bakersfield, CA",
    latitude: 35.3733,
    longitude: -119.0187,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Aurora, CO",
    latitude: 39.7294,
    longitude: -104.8319,
    timezone: "America/Denver",
    region: "Colorado"
  },
  {
    name: "Honolulu, HI",
    latitude: 21.3099,
    longitude: -157.8581,
    timezone: "Pacific/Honolulu",
    region: "Hawaii"
  },
  {
    name: "Anaheim, CA",
    latitude: 33.8366,
    longitude: -117.9143,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Santa Ana, CA",
    latitude: 33.7455,
    longitude: -117.8677,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Corpus Christi, TX",
    latitude: 27.8006,
    longitude: -97.3964,
    timezone: "America/Chicago",
    region: "Texas"
  },
  {
    name: "Riverside, CA",
    latitude: 33.9533,
    longitude: -117.3962,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Lexington, KY",
    latitude: 38.0406,
    longitude: -84.5037,
    timezone: "America/New_York",
    region: "Kentucky"
  },
  {
    name: "Stockton, CA",
    latitude: 37.9577,
    longitude: -121.2908,
    timezone: "America/Los_Angeles",
    region: "California"
  },
  {
    name: "Henderson, NV",
    latitude: 36.0397,
    longitude: -114.9817,
    timezone: "America/Los_Angeles",
    region: "Nevada"
  },
  {
    name: "Saint Paul, MN",
    latitude: 44.9537,
    longitude: -93.0900,
    timezone: "America/Chicago",
    region: "Minnesota"
  },
  {
    name: "St. Louis, MO",
    latitude: 38.6270,
    longitude: -90.1994,
    timezone: "America/Chicago",
    region: "Missouri"
  },
  {
    name: "Cincinnati, OH",
    latitude: 39.1031,
    longitude: -84.5120,
    timezone: "America/New_York",
    region: "Ohio"
  },
  {
    name: "Pittsburgh, PA",
    latitude: 40.4406,
    longitude: -79.9959,
    timezone: "America/New_York",
    region: "Pennsylvania"
  },
  {
    name: "Greensboro, NC",
    latitude: 36.0726,
    longitude: -79.7920,
    timezone: "America/New_York",
    region: "North Carolina"
  },
  {
    name: "Anchorage, AK",
    latitude: 61.2181,
    longitude: -149.9003,
    timezone: "America/Anchorage",
    region: "Alaska"
  }
];

class LocationService {
  private readonly FAVORITES_KEY = 'panchangam_favorite_locations';
  
  // Get user's favorite locations from localStorage
  getFavoriteLocations(): Location[] {
    try {
      const stored = localStorage.getItem(this.FAVORITES_KEY);
      return stored ? JSON.parse(stored) : [];
    } catch (error) {
      console.warn('Failed to load favorite locations:', error);
      return [];
    }
  }
  
  // Add location to favorites
  addToFavorites(location: Location): void {
    try {
      const favorites = this.getFavoriteLocations();
      
      // Check if already exists (prevent duplicates)
      const exists = favorites.some(fav => 
        Math.abs(fav.latitude - location.latitude) < 0.001 &&
        Math.abs(fav.longitude - location.longitude) < 0.001
      );
      
      if (!exists) {
        favorites.unshift(location); // Add to beginning
        
        // Limit to 10 favorites
        if (favorites.length > 10) {
          favorites.pop();
        }
        
        localStorage.setItem(this.FAVORITES_KEY, JSON.stringify(favorites));
      }
    } catch (error) {
      console.warn('Failed to save favorite location:', error);
    }
  }
  
  // Remove location from favorites
  removeFromFavorites(location: Location): void {
    try {
      const favorites = this.getFavoriteLocations();
      const filtered = favorites.filter(fav => 
        Math.abs(fav.latitude - location.latitude) >= 0.001 ||
        Math.abs(fav.longitude - location.longitude) >= 0.001
      );
      localStorage.setItem(this.FAVORITES_KEY, JSON.stringify(filtered));
    } catch (error) {
      console.warn('Failed to remove favorite location:', error);
    }
  }
  
  // Check if location is in favorites
  isFavorite(location: Location): boolean {
    const favorites = this.getFavoriteLocations();
    return favorites.some(fav => 
      Math.abs(fav.latitude - location.latitude) < 0.001 &&
      Math.abs(fav.longitude - location.longitude) < 0.001
    );
  }
  
  // Get locations organized by category for US/UK users
  getLocationsByCategory(): { favorites: Location[], usUk: Location[], popular: Location[] } {
    const favorites = this.getFavoriteLocations();
    
    // US/UK locations (first 11 + UK locations we added)
    const usUkCount = 11 + 10; // US default + UK locations
    const usUk = popularLocations.slice(0, usUkCount);
    
    // Rest as popular (Indian + remaining US cities)
    const popular = popularLocations.slice(usUkCount);
    
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
          } catch (error) {
            // Fallback to Milpitas, CA if reverse geocoding fails
            resolve(popularLocations[0]);
          }
        },
        (error) => {
          if (error.code === error.PERMISSION_DENIED) {
            console.info('User declined location access, using default location');
          } else {
            console.error('Geolocation error:', error);
          }
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
        const data = await response.json();
        
        return {
          name: this.formatLocationName(data),
          latitude,
          longitude,
          timezone: this.inferTimezone(latitude, longitude),
          region: this.extractRegion(data)
        };
      }
    } catch (error) {
      console.warn('Reverse geocoding failed, using fallback:', error);
    }

    // Fallback: Find closest popular location
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
    } catch (error) {
      console.warn('Geocoding failed, falling back to popular locations:', error);
      return popularFiltered.slice(0, 10);
    }
  }

  private async geocodeSearch(query: string): Promise<Location[]> {
    // Using OpenStreetMap Nominatim API (free, no API key required)
    const encodedQuery = encodeURIComponent(query);
    const url = `https://nominatim.openstreetmap.org/search?q=${encodedQuery}&format=json&limit=10&addressdetails=1`;
    
    try {
      const response = await fetch(url, {
        headers: {
          'User-Agent': 'PanchangamApp/1.0 (https://panchangam.app)'
        }
      });
      
      if (!response.ok) {
        throw new Error(`Geocoding API error: ${response.status}`);
      }
      
      const data = await response.json();
      
      return data.map((item: any) => ({
        name: this.formatLocationName(item),
        latitude: parseFloat(item.lat),
        longitude: parseFloat(item.lon),
        timezone: this.inferTimezone(parseFloat(item.lat), parseFloat(item.lon)),
        region: this.extractRegion(item)
      }));
    } catch (error) {
      console.error('Geocoding search failed:', error);
      throw error;
    }
  }

  private formatLocationName(item: any): string {
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

  private extractRegion(item: any): string {
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