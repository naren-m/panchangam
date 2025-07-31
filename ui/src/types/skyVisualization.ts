// Sky Visualization Types

// Coordinate Systems
export interface EclipticCoordinates {
  longitude: number; // degrees (0-360)
  latitude: number;  // degrees (-90 to +90)
  distance?: number; // AU (astronomical units)
}

export interface EquatorialCoordinates {
  rightAscension: number; // degrees (0-360) or hours (0-24)
  declination: number;    // degrees (-90 to +90)
  distance?: number;      // AU
}

export interface HorizontalCoordinates {
  azimuth: number;   // degrees (0-360, N=0, E=90, S=180, W=270)
  altitude: number;  // degrees (-90 to +90, horizon=0)
  distance?: number; // AU
}

export interface ScreenCoordinates {
  x: number; // pixels
  y: number; // pixels
  z?: number; // depth for 3D rendering
}

// Observer Location
export interface Observer {
  latitude: number;  // degrees
  longitude: number; // degrees
  altitude?: number; // meters above sea level
  timezone?: string;
}

// Sky Sphere Configuration
export interface SkySphereConfig {
  radius: number;
  segments: number;
  rings: number;
  projection: 'stereographic' | 'orthographic' | 'mercator';
  coordinateSystem: 'equatorial' | 'horizontal' | 'ecliptic';
  renderMode: 'webgl' | 'canvas' | 'svg';
}

// Celestial Object
export interface CelestialObject {
  id: string;
  name: string;
  type: 'star' | 'planet' | 'moon' | 'sun' | 'nakshatra' | 'constellation' | 'other';
  coordinates: {
    ecliptic?: EclipticCoordinates;
    equatorial?: EquatorialCoordinates;
    horizontal?: HorizontalCoordinates;
  };
  magnitude?: number;
  color?: string;
  size?: number;
  metadata?: Record<string, any>;
}

// Nakshatra specific
export interface NakshatraVisualization {
  id: number;
  name: string;
  startLongitude: number; // degrees
  endLongitude: number;   // degrees
  stars: CelestialObject[];
  boundaries: ScreenCoordinates[];
  deity?: string;
  symbol?: string;
}

// Camera Configuration
export interface CameraConfig {
  fov: number;        // field of view in degrees
  near: number;       // near clipping plane
  far: number;        // far clipping plane
  position: {
    x: number;
    y: number;
    z: number;
  };
  target: {
    x: number;
    y: number;
    z: number;
  };
}

// Rendering Options
export interface RenderOptions {
  showGrid: boolean;
  showConstellations: boolean;
  showNakshatras: boolean;
  showPlanets: boolean;
  showStars: boolean;
  showLabels: boolean;
  showZodiac: boolean;
  showEcliptic: boolean;
  showEquator: boolean;
  showHorizon: boolean;
  starMagnitudeLimit: number;
  labelMinZoom: number;
}

// Time Configuration
export interface TimeConfig {
  date: Date;
  speed: number; // time multiplier (1 = real-time, 60 = 1 minute per second)
  paused: boolean;
}