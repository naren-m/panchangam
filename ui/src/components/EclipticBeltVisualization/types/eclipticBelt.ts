/**
 * Type definitions for shared panchangam element calculations.
 */

// ============================================================================
// Core Celestial Types
// ============================================================================

/**
 * Position on the ecliptic belt (0-360 degrees)
 * 0° = Aries/Ashwini start (Vernal Equinox point)
 */
export interface EclipticPosition {
  longitude: number;  // 0-360 degrees
  latitude?: number;  // -90 to +90 degrees (optional for 2D view)
}

/**
 * Sun position with additional metadata
 */
export interface SunPosition extends EclipticPosition {
  rashi: RashiInfo;
  dailyMotion: number;  // degrees per day (~1°)
}

/**
 * Moon position with additional metadata
 */
export interface MoonPosition extends EclipticPosition {
  nakshatra: NakshatraInfo;
  rashi: RashiInfo;
  dailyMotion: number;  // degrees per day (~13°)
  phase: number;        // 0-1 (0 = new, 0.5 = full)
}

// ============================================================================
// Panchangam Element Types
// ============================================================================

/**
 * Tithi (Lunar Day) - Based on angular distance between Sun and Moon
 * Each tithi spans 12° of separation (360° / 30 tithis)
 */
export interface TithiInfo {
  number: number;       // 1-30
  name: string;         // e.g., "Pratipada", "Dvitiya"
  paksha: 'Shukla' | 'Krishna';  // Bright/Dark fortnight
  deity: string;
  angle: number;        // Moon longitude - Sun longitude (0-360°)
  percentComplete: number;  // 0-100% of current tithi elapsed
  startTime?: Date;     // When this tithi started (optional)
  endTime?: Date;       // When this tithi will end (optional)
}

/**
 * Nakshatra (Lunar Mansion) - Moon's position in 27 divisions
 * Each nakshatra spans 13°20' (360° / 27)
 */
export interface NakshatraInfo {
  number: number;       // 1-27
  name: string;         // e.g., "Ashwini", "Bharani"
  deity: string;
  symbol: string;
  startDegree: number;  // 0-360
  endDegree: number;    // 0-360
  pada: number;         // 1-4 (quarter of nakshatra)
}

/**
 * Yoga - Sum of Sun and Moon longitudes divided into 27 parts
 * Each yoga spans 13°20' of combined longitude
 */
export interface YogaInfo {
  number: number;       // 1-27
  name: string;         // e.g., "Vishkambha", "Priti"
  meaning: string;
  nature: 'Auspicious' | 'Inauspicious' | 'Mixed';
  combinedLongitude: number;  // (Sun + Moon) mod 360
}

/**
 * Karana (Half-Tithi) - 11 types, 60 in a lunar month
 * Each karana spans 6° of lunar-solar separation
 */
export interface KaranaInfo {
  number: number;       // 1-60 in lunar month
  name: string;         // e.g., "Bava", "Balava"
  type: 'Movable' | 'Fixed';
  nature: 'Auspicious' | 'Inauspicious' | 'Mixed';
}

/**
 * Rashi (Zodiac Sign) - 12 divisions of 30° each
 */
export interface RashiInfo {
  number: number;       // 1-12
  name: string;         // Sanskrit name: "Mesha", "Vrishabha", etc.
  westernName: string;  // Western name: "Aries", "Taurus", etc.
  symbol: string;       // Unicode symbol: , , etc.
  element: 'Fire' | 'Earth' | 'Air' | 'Water';
  ruler: string;        // Ruling planet
  startDegree: number;  // 0, 30, 60, etc.
  endDegree: number;    // 30, 60, 90, etc.
}

/**
 * Complete Panchangam data for a specific moment
 */
export interface PanchangamElements {
  tithi: TithiInfo;
  nakshatra: NakshatraInfo;
  yoga: YogaInfo;
  karana: KaranaInfo;
  rashi: RashiInfo;  // Moon's rashi
  sunRashi: RashiInfo;
  sunPosition: SunPosition;
  moonPosition: MoonPosition;
}

// ============================================================================
// Constants
// ============================================================================

export const DEGREES_PER_RASHI = 30;  // 360 / 12
export const DEGREES_PER_NAKSHATRA = 13.333333;  // 360 / 27
export const DEGREES_PER_TITHI = 12;  // 360 / 30
export const DEGREES_PER_YOGA = 13.333333;  // 360 / 27
export const DEGREES_PER_KARANA = 6;  // 360 / 60
