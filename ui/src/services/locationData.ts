import type { Location } from '../types/panchangam';
import { additionalPopularLocations } from './locationDataAdditional';

export const US_UK_LOCATION_COUNT = 21;

export const popularLocations: Location[] = [
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
  ...additionalPopularLocations,
];
