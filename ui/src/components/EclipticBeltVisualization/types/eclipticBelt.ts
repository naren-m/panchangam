/**
 * Type definitions for the 2D Ecliptic Belt Visualization
 *
 * The ecliptic belt shows the Sun's apparent path through the sky,
 * divided into 12 Rashis (zodiac signs) and 27 Nakshatras (lunar mansions).
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
  symbol: string;       // Unicode symbol: ♈, ♉, etc.
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
// Visualization Layout Types
// ============================================================================

/**
 * Dimensions and positioning for the SVG visualization
 */
export interface EclipticBeltDimensions {
  width: number;
  height: number;
  padding: {
    top: number;
    right: number;
    bottom: number;
    left: number;
  };
  zones: {
    rashiHeight: number;
    nakshatraHeight: number;
    planetTrackHeight: number;
    tithiArcHeight: number;
    annotationHeight: number;
  };
}

/**
 * Screen coordinates for rendering
 */
export interface ScreenPosition {
  x: number;
  y: number;
}

/**
 * A segment on the ecliptic belt (for Rashis or Nakshatras)
 */
export interface EclipticSegment {
  id: string;
  startX: number;
  endX: number;
  startDegree: number;
  endDegree: number;
  label: string;
  color?: string;
}

/**
 * Marker for celestial bodies on the belt
 */
export interface CelestialMarker {
  type: 'sun' | 'moon';
  position: ScreenPosition;
  longitude: number;
  label: string;
  color: string;
  size: number;
}

/**
 * Arc showing angular relationship (e.g., tithi arc between Sun and Moon)
 */
export interface RelationshipArc {
  type: 'tithi';
  startX: number;
  endX: number;
  y: number;
  angle: number;  // Angular separation in degrees
  label: string;
  color: string;
}

/**
 * Educational annotation for the visualization
 */
export interface Annotation {
  id: string;
  type: 'tithi' | 'nakshatra' | 'yoga' | 'karana' | 'rashi';
  content: string;
  detail: string;
  position: ScreenPosition;
  highlight?: {
    startX: number;
    endX: number;
  };
}

// ============================================================================
// Component Props Types
// ============================================================================

/**
 * Props for the main EclipticBeltContainer
 */
export interface EclipticBeltContainerProps {
  date: Date;
  latitude: number;
  longitude: number;
  timezone?: string;
  panchangamData?: Record<string, any>;  // From useProgressivePanchangam
  onClose?: () => void;
  className?: string;
}

/**
 * Props for the SVG renderer
 */
export interface EclipticBeltSVGProps {
  dimensions: EclipticBeltDimensions;
  panchangam: PanchangamElements;
  rashiSegments: EclipticSegment[];
  nakshatraSegments: EclipticSegment[];
  sunMarker: CelestialMarker;
  moonMarker: CelestialMarker;
  tithiArc: RelationshipArc;
  annotations: Annotation[];
  selectedElement: string | null;
  onElementSelect: (elementId: string | null) => void;
  onElementHover: (elementId: string | null) => void;
  hoveredElement: string | null;
  showLabels?: boolean;
  animationEnabled?: boolean;
}

/**
 * Time control state for animation
 */
export interface TimeControlState {
  currentDate: Date;
  isPlaying: boolean;
  playbackSpeed: number;  // 1 = real-time, 60 = 1 min/sec, etc.
  direction: 'forward' | 'backward';
}

/**
 * Interaction state for the visualization
 */
export interface InteractionState {
  selectedElement: string | null;
  hoveredElement: string | null;
  zoomLevel: number;
  panOffset: { x: number; y: number };
}

// ============================================================================
// Calculation Input/Output Types
// ============================================================================

/**
 * Input for panchangam calculations
 */
export interface PanchangamCalculationInput {
  date: Date;
  latitude: number;
  longitude: number;
  timezone: string;
}

/**
 * Output from panchangam calculations
 */
export interface PanchangamCalculationResult {
  elements: PanchangamElements;
  julianDay: number;
  localSiderealTime: number;
  ayanamsa: number;  // Precession correction for sidereal positions
}

// ============================================================================
// Constants
// ============================================================================

export const RASHI_COUNT = 12;
export const NAKSHATRA_COUNT = 27;
export const TITHI_COUNT = 30;
export const YOGA_COUNT = 27;
export const KARANA_COUNT = 11;  // Unique types (60 in month)

export const DEGREES_PER_RASHI = 30;  // 360 / 12
export const DEGREES_PER_NAKSHATRA = 13.333333;  // 360 / 27
export const DEGREES_PER_TITHI = 12;  // 360 / 30
export const DEGREES_PER_YOGA = 13.333333;  // 360 / 27
export const DEGREES_PER_KARANA = 6;  // 360 / 60

/**
 * Ayanamsa value for sidereal calculations (Lahiri)
 * This shifts tropical positions to sidereal (as used in Indian astronomy)
 * Approximate value for 2024, should be calculated precisely for accuracy
 */
export const LAHIRI_AYANAMSA_2024 = 24.17;  // degrees
