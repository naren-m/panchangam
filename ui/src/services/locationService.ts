import { Location } from '../types/panchangam';

const popularLocations: Location[] = [
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
  }
];

class LocationService {
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
            // Fallback to Chennai if reverse geocoding fails
            resolve(popularLocations[0]);
          }
        },
        (error) => {
          if (error.code === error.PERMISSION_DENIED) {
            console.info('User declined location access, using default location');
          } else {
            console.error('Geolocation error:', error);
          }
          // Fallback to Chennai
          resolve(popularLocations[0]);
        }
      );
    });
  }

  async reverseGeocode(latitude: number, longitude: number): Promise<Location> {
    // In a real implementation, you would use a geocoding service
    // For now, return the closest popular location
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
      timezone: "Asia/Kolkata",
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
    // Simple filtering of popular locations
    const filtered = popularLocations.filter(location =>
      location.name.toLowerCase().includes(query.toLowerCase())
    );
    return filtered;
  }
}

export const locationService = new LocationService();