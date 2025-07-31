// Coordinate Transformation Utilities
// Based on astronomical algorithms from Jean Meeus "Astronomical Algorithms"

import { 
  EclipticCoordinates, 
  EquatorialCoordinates, 
  HorizontalCoordinates, 
  ScreenCoordinates,
  Observer 
} from '../../types/skyVisualization';

// Constants
const DEG_TO_RAD = Math.PI / 180;
const RAD_TO_DEG = 180 / Math.PI;
const J2000_OBLIQUITY = 23.43929111; // Earth's obliquity at J2000 epoch in degrees

/**
 * Converts ecliptic coordinates to equatorial coordinates
 * @param ecliptic Ecliptic coordinates (longitude, latitude in degrees)
 * @param jd Julian day number (for precession calculation)
 * @returns Equatorial coordinates (RA, Dec in degrees)
 */
export function eclipticToEquatorial(
  ecliptic: EclipticCoordinates, 
  jd: number = 2451545.0 // J2000.0 default
): EquatorialCoordinates {
  // Calculate obliquity of ecliptic for the given epoch
  const T = (jd - 2451545.0) / 36525.0; // Julian centuries from J2000
  const epsilon = J2000_OBLIQUITY - 
    (46.8150 * T + 0.00059 * T * T - 0.001813 * T * T * T) / 3600;
  
  const lambda = ecliptic.longitude * DEG_TO_RAD;
  const beta = ecliptic.latitude * DEG_TO_RAD;
  const eps = epsilon * DEG_TO_RAD;
  
  // Convert to equatorial coordinates
  const sinLambda = Math.sin(lambda);
  const cosLambda = Math.cos(lambda);
  const sinBeta = Math.sin(beta);
  const cosBeta = Math.cos(beta);
  const sinEps = Math.sin(eps);
  const cosEps = Math.cos(eps);
  
  // Right Ascension
  const ra = Math.atan2(
    sinLambda * cosEps - Math.tan(beta) * sinEps,
    cosLambda
  );
  
  // Declination
  const dec = Math.asin(
    sinBeta * cosEps + cosBeta * sinEps * sinLambda
  );
  
  // Normalize RA to 0-360 degrees
  let raDeg = ra * RAD_TO_DEG;
  if (raDeg < 0) raDeg += 360;
  
  return {
    rightAscension: raDeg,
    declination: dec * RAD_TO_DEG,
    distance: ecliptic.distance
  };
}

/**
 * Converts equatorial coordinates to horizontal coordinates
 * @param equatorial Equatorial coordinates (RA, Dec in degrees)
 * @param observer Observer location (lat, lon in degrees)
 * @param time Date/time of observation
 * @returns Horizontal coordinates (azimuth, altitude in degrees)
 */
export function equatorialToHorizontal(
  equatorial: EquatorialCoordinates,
  observer: Observer,
  time: Date
): HorizontalCoordinates {
  // Calculate Local Sidereal Time (LST)
  const lst = getLocalSiderealTime(observer.longitude, time);
  
  // Hour angle
  const ha = (lst - equatorial.rightAscension) * DEG_TO_RAD;
  
  const dec = equatorial.declination * DEG_TO_RAD;
  const lat = observer.latitude * DEG_TO_RAD;
  
  const sinDec = Math.sin(dec);
  const cosDec = Math.cos(dec);
  const sinLat = Math.sin(lat);
  const cosLat = Math.cos(lat);
  const cosHA = Math.cos(ha);
  
  // Altitude
  const alt = Math.asin(
    sinDec * sinLat + cosDec * cosLat * cosHA
  );
  
  // Azimuth
  const az = Math.atan2(
    -cosDec * Math.sin(ha),
    sinDec * cosLat - cosDec * sinLat * cosHA
  );
  
  // Convert to degrees and normalize azimuth to 0-360
  let azDeg = az * RAD_TO_DEG;
  if (azDeg < 0) azDeg += 360;
  
  return {
    azimuth: azDeg,
    altitude: alt * RAD_TO_DEG,
    distance: equatorial.distance
  };
}

/**
 * Converts horizontal coordinates to screen coordinates
 * @param horizontal Horizontal coordinates (azimuth, altitude)
 * @param width Canvas width
 * @param height Canvas height
 * @param projection Projection type
 * @param fov Field of view in degrees (for perspective projection)
 * @returns Screen coordinates (x, y in pixels)
 */
export function horizontalToScreen(
  horizontal: HorizontalCoordinates,
  width: number,
  height: number,
  projection: 'stereographic' | 'orthographic' | 'mercator' = 'stereographic',
  fov: number = 90
): ScreenCoordinates {
  const az = horizontal.azimuth * DEG_TO_RAD;
  const alt = horizontal.altitude * DEG_TO_RAD;
  
  let x: number, y: number;
  
  switch (projection) {
    case 'stereographic': {
      // Stereographic projection (conformal, preserves angles)
      const k = 2 / (1 + Math.sin(alt));
      const r = k * Math.cos(alt);
      x = r * Math.sin(az);
      y = r * Math.cos(az);
      break;
    }
    
    case 'orthographic': {
      // Orthographic projection (like viewing a globe from far away)
      if (alt < 0) {
        // Below horizon, don't render
        return { x: -1, y: -1, z: -1 };
      }
      x = Math.cos(alt) * Math.sin(az);
      y = Math.cos(alt) * Math.cos(az);
      break;
    }
    
    case 'mercator': {
      // Mercator projection (cylindrical)
      x = az / Math.PI - 1;
      y = Math.log(Math.tan(Math.PI / 4 + alt / 2)) / Math.PI;
      break;
    }
  }
  
  // Scale to canvas dimensions
  const scale = Math.min(width, height) / (2 * Math.tan(fov * DEG_TO_RAD / 2));
  
  return {
    x: width / 2 + x * scale,
    y: height / 2 - y * scale, // Flip y-axis for screen coordinates
    z: horizontal.distance
  };
}

/**
 * Calculate Local Sidereal Time
 * @param longitude Observer longitude in degrees
 * @param time Date/time
 * @returns LST in degrees
 */
function getLocalSiderealTime(longitude: number, time: Date): number {
  // Convert to Julian Day
  const jd = dateToJulianDay(time);
  
  // Calculate centuries from J2000
  const T = (jd - 2451545.0) / 36525.0;
  
  // Greenwich mean sidereal time at 0h UT
  let gmst0 = 280.46061837 + 
    360.98564736629 * (jd - 2451545.0) +
    0.000387933 * T * T -
    T * T * T / 38710000.0;
  
  // Normalize to 0-360
  gmst0 = gmst0 % 360;
  if (gmst0 < 0) gmst0 += 360;
  
  // Add hour angle
  const ut = time.getUTCHours() + 
    time.getUTCMinutes() / 60 + 
    time.getUTCSeconds() / 3600;
  
  let lst = gmst0 + ut * 15 + longitude;
  
  // Normalize to 0-360
  lst = lst % 360;
  if (lst < 0) lst += 360;
  
  return lst;
}

/**
 * Convert Date to Julian Day number
 * @param date Date object
 * @returns Julian Day number
 */
export function dateToJulianDay(date: Date): number {
  const year = date.getUTCFullYear();
  const month = date.getUTCMonth() + 1;
  const day = date.getUTCDate();
  const hour = date.getUTCHours();
  const minute = date.getUTCMinutes();
  const second = date.getUTCSeconds();
  
  let a = Math.floor((14 - month) / 12);
  let y = year + 4800 - a;
  let m = month + 12 * a - 3;
  
  let jdn = day + Math.floor((153 * m + 2) / 5) + 
    365 * y + Math.floor(y / 4) - 
    Math.floor(y / 100) + Math.floor(y / 400) - 32045;
  
  // Add time of day
  jdn += (hour - 12) / 24 + minute / 1440 + second / 86400;
  
  return jdn;
}

/**
 * Chain transformation from ecliptic to screen coordinates
 * @param ecliptic Ecliptic coordinates
 * @param observer Observer location
 * @param time Observation time
 * @param canvasWidth Canvas width
 * @param canvasHeight Canvas height
 * @param projection Projection type
 * @returns Screen coordinates
 */
export function eclipticToScreen(
  ecliptic: EclipticCoordinates,
  observer: Observer,
  time: Date,
  canvasWidth: number,
  canvasHeight: number,
  projection: 'stereographic' | 'orthographic' | 'mercator' = 'stereographic'
): ScreenCoordinates {
  const jd = dateToJulianDay(time);
  const equatorial = eclipticToEquatorial(ecliptic, jd);
  const horizontal = equatorialToHorizontal(equatorial, observer, time);
  return horizontalToScreen(horizontal, canvasWidth, canvasHeight, projection);
}