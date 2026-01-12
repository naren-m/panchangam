/**
 * Ecliptic Belt Layout Calculator
 *
 * Handles the conversion from celestial coordinates (0-360° longitude)
 * to screen coordinates for the SVG visualization.
 *
 * The belt is rendered as a horizontal strip where:
 * - Left edge = 0° (start of Aries/Mesha)
 * - Right edge = 360° (end of Pisces/Meena, wraps to 0°)
 *
 * The visualization has 5 horizontal zones:
 * 1. Rashi Zone: 12 zodiac signs (30° each)
 * 2. Nakshatra Zone: 27 lunar mansions (13.33° each)
 * 3. Planet Track: Sun and Moon markers
 * 4. Tithi Arc: Arc showing Moon-Sun angular separation
 * 5. Annotation Zone: Educational labels and info cards
 */

import {
  EclipticBeltDimensions,
  EclipticSegment,
  CelestialMarker,
  RelationshipArc,
  ScreenPosition,
  RashiInfo,
  NakshatraInfo,
  PanchangamElements,
  RASHI_COUNT,
  NAKSHATRA_COUNT,
  DEGREES_PER_RASHI,
  DEGREES_PER_NAKSHATRA,
} from '../types/eclipticBelt';

// ============================================================================
// Default Dimensions
// ============================================================================

/**
 * Create default dimensions for the ecliptic belt visualization
 */
export function createDefaultDimensions(
  containerWidth: number = 1200,
  containerHeight: number = 400
): EclipticBeltDimensions {
  const padding = {
    top: 20,
    right: 40,
    bottom: 20,
    left: 40
  };

  // Calculate available drawing area
  const drawableHeight = containerHeight - padding.top - padding.bottom;

  // Allocate height to each zone (proportional)
  const zones = {
    rashiHeight: Math.floor(drawableHeight * 0.18),       // 18% for zodiac signs
    nakshatraHeight: Math.floor(drawableHeight * 0.18),   // 18% for nakshatras
    planetTrackHeight: Math.floor(drawableHeight * 0.28), // 28% for planets
    tithiArcHeight: Math.floor(drawableHeight * 0.16),    // 16% for tithi arc
    annotationHeight: Math.floor(drawableHeight * 0.20)   // 20% for annotations
  };

  return {
    width: containerWidth,
    height: containerHeight,
    padding,
    zones
  };
}

// ============================================================================
// Coordinate Conversion
// ============================================================================

/**
 * Convert ecliptic longitude (0-360°) to X screen coordinate
 */
export function longitudeToX(
  longitude: number,
  dimensions: EclipticBeltDimensions
): number {
  // Normalize longitude to 0-360
  let normalized = longitude % 360;
  if (normalized < 0) normalized += 360;

  // Calculate drawable width
  const drawableWidth = dimensions.width - dimensions.padding.left - dimensions.padding.right;

  // Linear mapping: 0° → left edge, 360° → right edge
  const x = dimensions.padding.left + (normalized / 360) * drawableWidth;

  return x;
}

/**
 * Convert X screen coordinate to ecliptic longitude
 */
export function xToLongitude(
  x: number,
  dimensions: EclipticBeltDimensions
): number {
  const drawableWidth = dimensions.width - dimensions.padding.left - dimensions.padding.right;
  const relativeX = x - dimensions.padding.left;
  const longitude = (relativeX / drawableWidth) * 360;

  return Math.max(0, Math.min(360, longitude));
}

/**
 * Get Y coordinate for the center of each zone
 */
export function getZoneCenterY(
  zone: 'rashi' | 'nakshatra' | 'planets' | 'tithi' | 'annotation',
  dimensions: EclipticBeltDimensions
): number {
  const { padding, zones } = dimensions;

  switch (zone) {
    case 'rashi':
      return padding.top + zones.rashiHeight / 2;

    case 'nakshatra':
      return padding.top + zones.rashiHeight + zones.nakshatraHeight / 2;

    case 'planets':
      return padding.top + zones.rashiHeight + zones.nakshatraHeight + zones.planetTrackHeight / 2;

    case 'tithi':
      return padding.top + zones.rashiHeight + zones.nakshatraHeight +
        zones.planetTrackHeight + zones.tithiArcHeight / 2;

    case 'annotation':
      return padding.top + zones.rashiHeight + zones.nakshatraHeight +
        zones.planetTrackHeight + zones.tithiArcHeight + zones.annotationHeight / 2;

    default:
      return dimensions.height / 2;
  }
}

/**
 * Get Y coordinate range for a zone
 */
export function getZoneYRange(
  zone: 'rashi' | 'nakshatra' | 'planets' | 'tithi' | 'annotation',
  dimensions: EclipticBeltDimensions
): { top: number; bottom: number } {
  const { padding, zones } = dimensions;

  let top = padding.top;

  switch (zone) {
    case 'rashi':
      return { top, bottom: top + zones.rashiHeight };

    case 'nakshatra':
      top += zones.rashiHeight;
      return { top, bottom: top + zones.nakshatraHeight };

    case 'planets':
      top += zones.rashiHeight + zones.nakshatraHeight;
      return { top, bottom: top + zones.planetTrackHeight };

    case 'tithi':
      top += zones.rashiHeight + zones.nakshatraHeight + zones.planetTrackHeight;
      return { top, bottom: top + zones.tithiArcHeight };

    case 'annotation':
      top += zones.rashiHeight + zones.nakshatraHeight + zones.planetTrackHeight + zones.tithiArcHeight;
      return { top, bottom: top + zones.annotationHeight };

    default:
      return { top: 0, bottom: dimensions.height };
  }
}

// ============================================================================
// Rashi (Zodiac) Segments
// ============================================================================

/**
 * Color palette for the 12 Rashis based on elements
 */
const RASHI_COLORS = {
  Fire: '#FF6B6B',   // Warm red for Aries, Leo, Sagittarius
  Earth: '#4ECDC4',  // Teal for Taurus, Virgo, Capricorn
  Air: '#95E1D3',    // Light green for Gemini, Libra, Aquarius
  Water: '#74B9FF'   // Blue for Cancer, Scorpio, Pisces
};

/**
 * Generate all 12 Rashi segments for the visualization
 */
export function generateRashiSegments(dimensions: EclipticBeltDimensions): EclipticSegment[] {
  const segments: EclipticSegment[] = [];

  const rashiData = [
    { name: 'Mesha', symbol: '\u2648', element: 'Fire' },
    { name: 'Vrishabha', symbol: '\u2649', element: 'Earth' },
    { name: 'Mithuna', symbol: '\u264A', element: 'Air' },
    { name: 'Karka', symbol: '\u264B', element: 'Water' },
    { name: 'Simha', symbol: '\u264C', element: 'Fire' },
    { name: 'Kanya', symbol: '\u264D', element: 'Earth' },
    { name: 'Tula', symbol: '\u264E', element: 'Air' },
    { name: 'Vrishchika', symbol: '\u264F', element: 'Water' },
    { name: 'Dhanus', symbol: '\u2650', element: 'Fire' },
    { name: 'Makara', symbol: '\u2651', element: 'Earth' },
    { name: 'Kumbha', symbol: '\u2652', element: 'Air' },
    { name: 'Meena', symbol: '\u2653', element: 'Water' }
  ];

  for (let i = 0; i < RASHI_COUNT; i++) {
    const startDegree = i * DEGREES_PER_RASHI;
    const endDegree = (i + 1) * DEGREES_PER_RASHI;

    const startX = longitudeToX(startDegree, dimensions);
    const endX = longitudeToX(endDegree, dimensions);

    const data = rashiData[i];

    segments.push({
      id: `rashi-${i + 1}`,
      startX,
      endX,
      startDegree,
      endDegree,
      label: `${data.symbol} ${data.name}`,
      color: RASHI_COLORS[data.element as keyof typeof RASHI_COLORS]
    });
  }

  return segments;
}

// ============================================================================
// Nakshatra Segments
// ============================================================================

/**
 * Color gradient for 27 Nakshatras (subtle variations)
 */
function getNakshatraColor(index: number): string {
  // Use a gradient from orange to purple across the 27 nakshatras
  const hue = 30 + (index * 270 / 27); // 30° (orange) to 300° (purple)
  const saturation = 60;
  const lightness = 75;
  return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
}

/**
 * Generate all 27 Nakshatra segments for the visualization
 */
export function generateNakshatraSegments(dimensions: EclipticBeltDimensions): EclipticSegment[] {
  const segments: EclipticSegment[] = [];

  const nakshatraNames = [
    'Ashwini', 'Bharani', 'Krittika', 'Rohini', 'Mrigashira', 'Ardra',
    'Punarvasu', 'Pushya', 'Ashlesha', 'Magha', 'P.Phal', 'U.Phal',
    'Hasta', 'Chitra', 'Swati', 'Vishakha', 'Anuradha', 'Jyeshtha',
    'Mula', 'P.Ash', 'U.Ash', 'Shravana', 'Dhanishta', 'Shatabhisha',
    'P.Bha', 'U.Bha', 'Revati'
  ];

  for (let i = 0; i < NAKSHATRA_COUNT; i++) {
    const startDegree = i * DEGREES_PER_NAKSHATRA;
    const endDegree = (i + 1) * DEGREES_PER_NAKSHATRA;

    const startX = longitudeToX(startDegree, dimensions);
    const endX = longitudeToX(endDegree, dimensions);

    segments.push({
      id: `nakshatra-${i + 1}`,
      startX,
      endX,
      startDegree,
      endDegree,
      label: nakshatraNames[i],
      color: getNakshatraColor(i)
    });
  }

  return segments;
}

// ============================================================================
// Celestial Markers (Sun and Moon)
// ============================================================================

/**
 * Create a Sun marker at the given longitude
 */
export function createSunMarker(
  longitude: number,
  dimensions: EclipticBeltDimensions
): CelestialMarker {
  const x = longitudeToX(longitude, dimensions);
  const y = getZoneCenterY('planets', dimensions);

  return {
    type: 'sun',
    position: { x, y },
    longitude,
    label: '\u2609 Sun',  // Unicode sun symbol
    color: '#FFB300',     // Warm yellow-orange
    size: 24
  };
}

/**
 * Create a Moon marker at the given longitude
 */
export function createMoonMarker(
  longitude: number,
  phase: number,
  dimensions: EclipticBeltDimensions
): CelestialMarker {
  const x = longitudeToX(longitude, dimensions);
  const y = getZoneCenterY('planets', dimensions);

  // Choose moon symbol based on phase
  let moonSymbol = '☽';  // First quarter moon
  if (phase < 0.125) moonSymbol = '🌑';      // New moon
  else if (phase < 0.375) moonSymbol = '🌒'; // Waxing crescent
  else if (phase < 0.625) moonSymbol = '🌕'; // Full moon
  else if (phase < 0.875) moonSymbol = '🌘'; // Waning crescent
  else moonSymbol = '🌑';                    // New moon

  return {
    type: 'moon',
    position: { x, y },
    longitude,
    label: `${moonSymbol} Moon`,
    color: '#C0C0C0',  // Silver
    size: 20
  };
}

// ============================================================================
// Tithi Arc (Moon-Sun Relationship)
// ============================================================================

/**
 * Create a Tithi arc showing the angular separation between Sun and Moon
 *
 * The arc visually represents how far the Moon has traveled from the Sun,
 * which determines the current Tithi (lunar day).
 */
export function createTithiArc(
  sunLongitude: number,
  moonLongitude: number,
  dimensions: EclipticBeltDimensions
): RelationshipArc {
  const sunX = longitudeToX(sunLongitude, dimensions);
  const moonX = longitudeToX(moonLongitude, dimensions);
  const y = getZoneCenterY('tithi', dimensions);

  // Calculate angular separation
  let angle = moonLongitude - sunLongitude;
  if (angle < 0) angle += 360;

  // Determine tithi number for label
  const tithiNumber = Math.floor(angle / 12) + 1;
  const tithiNames = [
    '', 'Pratipada', 'Dvitiya', 'Tritiya', 'Chaturthi', 'Panchami',
    'Shashthi', 'Saptami', 'Ashtami', 'Navami', 'Dashami',
    'Ekadashi', 'Dvadashi', 'Trayodashi', 'Chaturdashi', 'Purnima',
    'Pratipada', 'Dvitiya', 'Tritiya', 'Chaturthi', 'Panchami',
    'Shashthi', 'Saptami', 'Ashtami', 'Navami', 'Dashami',
    'Ekadashi', 'Dvadashi', 'Trayodashi', 'Chaturdashi', 'Amavasya'
  ];

  const tithiName = tithiNames[Math.min(tithiNumber, 30)];
  const paksha = tithiNumber <= 15 ? 'Shukla' : 'Krishna';

  return {
    type: 'tithi',
    startX: sunX,
    endX: moonX,
    y,
    angle,
    label: `${paksha} ${tithiName} (${angle.toFixed(1)}°)`,
    color: tithiNumber <= 15 ? '#FFD700' : '#4A4A4A'  // Gold for Shukla, dark for Krishna
  };
}

// ============================================================================
// Complete Layout Generation
// ============================================================================

/**
 * Generate all layout elements for the visualization
 */
export interface EclipticBeltLayout {
  dimensions: EclipticBeltDimensions;
  rashiSegments: EclipticSegment[];
  nakshatraSegments: EclipticSegment[];
  sunMarker: CelestialMarker;
  moonMarker: CelestialMarker;
  tithiArc: RelationshipArc;
}

export function generateLayout(
  panchangam: PanchangamElements,
  containerWidth: number = 1200,
  containerHeight: number = 400
): EclipticBeltLayout {
  const dimensions = createDefaultDimensions(containerWidth, containerHeight);

  return {
    dimensions,
    rashiSegments: generateRashiSegments(dimensions),
    nakshatraSegments: generateNakshatraSegments(dimensions),
    sunMarker: createSunMarker(panchangam.sunPosition.longitude, dimensions),
    moonMarker: createMoonMarker(
      panchangam.moonPosition.longitude,
      panchangam.moonPosition.phase,
      dimensions
    ),
    tithiArc: createTithiArc(
      panchangam.sunPosition.longitude,
      panchangam.moonPosition.longitude,
      dimensions
    )
  };
}

// ============================================================================
// Hover/Selection Detection
// ============================================================================

/**
 * Find which segment (Rashi or Nakshatra) contains a given X coordinate
 */
export function findSegmentAtX(
  x: number,
  segments: EclipticSegment[]
): EclipticSegment | null {
  return segments.find(seg => x >= seg.startX && x < seg.endX) || null;
}

/**
 * Find which zone contains a given Y coordinate
 */
export function findZoneAtY(
  y: number,
  dimensions: EclipticBeltDimensions
): 'rashi' | 'nakshatra' | 'planets' | 'tithi' | 'annotation' | null {
  const zones: Array<'rashi' | 'nakshatra' | 'planets' | 'tithi' | 'annotation'> =
    ['rashi', 'nakshatra', 'planets', 'tithi', 'annotation'];

  for (const zone of zones) {
    const range = getZoneYRange(zone, dimensions);
    if (y >= range.top && y < range.bottom) {
      return zone;
    }
  }

  return null;
}

/**
 * Check if a point is near a celestial marker
 */
export function isNearMarker(
  point: ScreenPosition,
  marker: CelestialMarker,
  threshold: number = 20
): boolean {
  const dx = point.x - marker.position.x;
  const dy = point.y - marker.position.y;
  const distance = Math.sqrt(dx * dx + dy * dy);
  return distance <= threshold;
}

// ============================================================================
// Animation Helpers
// ============================================================================

/**
 * Calculate interpolated longitude between two positions
 * (useful for smooth animation)
 */
export function interpolateLongitude(
  fromLong: number,
  toLong: number,
  progress: number
): number {
  // Handle wrap-around at 360°
  let diff = toLong - fromLong;

  // Take the shorter path
  if (diff > 180) diff -= 360;
  if (diff < -180) diff += 360;

  let result = fromLong + diff * progress;

  // Normalize to 0-360
  if (result < 0) result += 360;
  if (result >= 360) result -= 360;

  return result;
}

/**
 * Calculate positions for time-based animation
 * Moon moves ~13°/day, Sun moves ~1°/day
 */
export function calculateAnimatedPositions(
  baseSunLong: number,
  baseMoonLong: number,
  hoursOffset: number
): { sunLong: number; moonLong: number } {
  // Degrees per hour
  const sunSpeed = 1 / 24;      // ~1° per day
  const moonSpeed = 13.2 / 24;  // ~13.2° per day

  const sunLong = baseSunLong + sunSpeed * hoursOffset;
  const moonLong = baseMoonLong + moonSpeed * hoursOffset;

  return {
    sunLong: sunLong % 360,
    moonLong: moonLong % 360
  };
}
