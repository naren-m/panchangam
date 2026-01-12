/**
 * Chart-specific calculations for the Celestial Chart
 *
 * Builds on the panchangam calculator to provide chart-ready data.
 */

import type {
  ChartDimensions,
  CelestialBodyInfo,
  TithiArcInfo,
  PadaInfo,
  Point,
} from '../types';
import type {
  PanchangamElements,
  RashiInfo,
  NakshatraInfo,
  TithiInfo
} from '../../EclipticBeltVisualization/types/eclipticBelt';
import {
  longitudeToCartesian,
  createArcPath,
  normalizeAngle
} from './geometryHelpers';
import {
  CHART_COLORS,
  NAVAMSHA_SEQUENCE,
  DEGREES_PER_PADA,
  PADAS_PER_NAKSHATRA
} from '../types';

// ============================================================================
// Rashi (Zodiac Sign) Data
// ============================================================================

const RASHI_DATA: Array<{
  name: string;
  westernName: string;
  symbol: string;
  element: 'Fire' | 'Earth' | 'Air' | 'Water';
  ruler: string;
}> = [
  { name: 'Mesha', westernName: 'Aries', symbol: '♈', element: 'Fire', ruler: 'Mars' },
  { name: 'Vrishabha', westernName: 'Taurus', symbol: '♉', element: 'Earth', ruler: 'Venus' },
  { name: 'Mithuna', westernName: 'Gemini', symbol: '♊', element: 'Air', ruler: 'Mercury' },
  { name: 'Karka', westernName: 'Cancer', symbol: '♋', element: 'Water', ruler: 'Moon' },
  { name: 'Simha', westernName: 'Leo', symbol: '♌', element: 'Fire', ruler: 'Sun' },
  { name: 'Kanya', westernName: 'Virgo', symbol: '♍', element: 'Earth', ruler: 'Mercury' },
  { name: 'Tula', westernName: 'Libra', symbol: '♎', element: 'Air', ruler: 'Venus' },
  { name: 'Vrishchika', westernName: 'Scorpio', symbol: '♏', element: 'Water', ruler: 'Mars' },
  { name: 'Dhanu', westernName: 'Sagittarius', symbol: '♐', element: 'Fire', ruler: 'Jupiter' },
  { name: 'Makara', westernName: 'Capricorn', symbol: '♑', element: 'Earth', ruler: 'Saturn' },
  { name: 'Kumbha', westernName: 'Aquarius', symbol: '♒', element: 'Air', ruler: 'Saturn' },
  { name: 'Meena', westernName: 'Pisces', symbol: '♓', element: 'Water', ruler: 'Jupiter' },
];

/**
 * Get color for a rashi based on its element
 */
export const getRashiColor = (element: RashiInfo['element']): string => {
  return CHART_COLORS.rashiRing[element.toLowerCase() as keyof typeof CHART_COLORS.rashiRing];
};

/**
 * Get all rashi segments for rendering
 */
export const getAllRashis = (): RashiInfo[] => {
  return RASHI_DATA.map((rashi, index) => ({
    number: index + 1,
    name: rashi.name,
    westernName: rashi.westernName,
    symbol: rashi.symbol,
    element: rashi.element,
    ruler: rashi.ruler,
    startDegree: index * 30,
    endDegree: (index + 1) * 30,
  }));
};

// ============================================================================
// Nakshatra Data
// ============================================================================

const NAKSHATRA_DATA: Array<{
  name: string;
  deity: string;
  symbol: string;
}> = [
  { name: 'Ashwini', deity: 'Ashwini Kumaras', symbol: 'Horse head' },
  { name: 'Bharani', deity: 'Yama', symbol: 'Yoni' },
  { name: 'Krittika', deity: 'Agni', symbol: 'Razor' },
  { name: 'Rohini', deity: 'Brahma', symbol: 'Chariot' },
  { name: 'Mrigashira', deity: 'Soma', symbol: 'Deer head' },
  { name: 'Ardra', deity: 'Rudra', symbol: 'Teardrop' },
  { name: 'Punarvasu', deity: 'Aditi', symbol: 'Bow' },
  { name: 'Pushya', deity: 'Brihaspati', symbol: 'Cow udder' },
  { name: 'Ashlesha', deity: 'Sarpa', symbol: 'Serpent' },
  { name: 'Magha', deity: 'Pitris', symbol: 'Throne' },
  { name: 'Purva Phalguni', deity: 'Bhaga', symbol: 'Hammock' },
  { name: 'Uttara Phalguni', deity: 'Aryaman', symbol: 'Bed' },
  { name: 'Hasta', deity: 'Savitar', symbol: 'Hand' },
  { name: 'Chitra', deity: 'Tvashtar', symbol: 'Pearl' },
  { name: 'Swati', deity: 'Vayu', symbol: 'Coral' },
  { name: 'Vishakha', deity: 'Indragni', symbol: 'Gateway' },
  { name: 'Anuradha', deity: 'Mitra', symbol: 'Lotus' },
  { name: 'Jyeshtha', deity: 'Indra', symbol: 'Earring' },
  { name: 'Mula', deity: 'Nirriti', symbol: 'Roots' },
  { name: 'Purva Ashadha', deity: 'Apas', symbol: 'Fan' },
  { name: 'Uttara Ashadha', deity: 'Vishvadevas', symbol: 'Tusk' },
  { name: 'Shravana', deity: 'Vishnu', symbol: 'Ear' },
  { name: 'Dhanishtha', deity: 'Vasus', symbol: 'Drum' },
  { name: 'Shatabhisha', deity: 'Varuna', symbol: 'Circle' },
  { name: 'Purva Bhadrapada', deity: 'Aja Ekapada', symbol: 'Sword' },
  { name: 'Uttara Bhadrapada', deity: 'Ahir Budhnya', symbol: 'Twins' },
  { name: 'Revati', deity: 'Pushan', symbol: 'Fish' },
];

const DEGREES_PER_NAKSHATRA = 360 / 27; // 13.333...

/**
 * Get all nakshatra segments for rendering
 */
export const getAllNakshatras = (): NakshatraInfo[] => {
  return NAKSHATRA_DATA.map((nakshatra, index) => ({
    number: index + 1,
    name: nakshatra.name,
    deity: nakshatra.deity,
    symbol: nakshatra.symbol,
    startDegree: index * DEGREES_PER_NAKSHATRA,
    endDegree: (index + 1) * DEGREES_PER_NAKSHATRA,
    pada: 1, // Default, actual pada determined by Moon position
  }));
};

/**
 * Get nakshatra at a specific longitude
 */
export const getNakshatraAtLongitude = (longitude: number): NakshatraInfo => {
  const normalized = normalizeAngle(longitude);
  const index = Math.floor(normalized / DEGREES_PER_NAKSHATRA);
  const nakshatras = getAllNakshatras();
  return nakshatras[index];
};

// ============================================================================
// Pada Calculations
// ============================================================================

/**
 * Get all 108 padas with their details
 */
export const getAllPadas = (): PadaInfo[] => {
  const nakshatras = getAllNakshatras();
  const padas: PadaInfo[] = [];

  for (let globalPada = 0; globalPada < 108; globalPada++) {
    const nakshatraIndex = Math.floor(globalPada / PADAS_PER_NAKSHATRA);
    const padaInNakshatra = (globalPada % PADAS_PER_NAKSHATRA) + 1;
    const navamshaIndex = globalPada % 12;

    padas.push({
      number: globalPada + 1,
      padaInNakshatra,
      nakshatra: nakshatras[nakshatraIndex],
      startDegree: globalPada * DEGREES_PER_PADA,
      endDegree: (globalPada + 1) * DEGREES_PER_PADA,
      navamsha: NAVAMSHA_SEQUENCE[navamshaIndex],
    });
  }

  return padas;
};

/**
 * Get pada at a specific longitude
 */
export const getPadaAtLongitude = (longitude: number): PadaInfo => {
  const normalized = normalizeAngle(longitude);
  const index = Math.floor(normalized / DEGREES_PER_PADA);
  const padas = getAllPadas();
  return padas[index];
};

/**
 * Get current pada from panchangam data
 */
export const getCurrentPada = (panchangam: PanchangamElements): PadaInfo => {
  const moonLongitude = panchangam.moonPosition.longitude;
  return getPadaAtLongitude(moonLongitude);
};

// ============================================================================
// Celestial Body Calculations
// ============================================================================

/**
 * Calculate Sun marker info
 */
export const getSunMarkerInfo = (
  panchangam: PanchangamElements,
  dimensions: ChartDimensions
): CelestialBodyInfo => {
  const longitude = panchangam.sunPosition.longitude;
  const position = longitudeToCartesian(
    dimensions.center,
    dimensions.celestialOrbit,
    longitude
  );

  return {
    type: 'sun',
    longitude,
    position,
    symbol: '☉',
    color: CHART_COLORS.sun,
    size: 24,
    label: `Sun in ${panchangam.sunRashi.name}`,
  };
};

/**
 * Calculate Moon marker info with phase visualization
 */
export const getMoonMarkerInfo = (
  panchangam: PanchangamElements,
  dimensions: ChartDimensions
): CelestialBodyInfo => {
  const longitude = panchangam.moonPosition.longitude;
  const position = longitudeToCartesian(
    dimensions.center,
    dimensions.celestialOrbit,
    longitude
  );

  // Determine moon phase symbol based on tithi
  const tithiNumber = panchangam.tithi.number;
  let symbol = '●'; // Default full moon

  if (tithiNumber === 15) {
    symbol = '●'; // Purnima - Full Moon
  } else if (tithiNumber === 30 || tithiNumber === 1) {
    symbol = '○'; // Amavasya - New Moon
  } else if (tithiNumber < 8) {
    symbol = panchangam.tithi.paksha === 'Shukla' ? '◐' : '◑';
  } else if (tithiNumber < 15) {
    symbol = panchangam.tithi.paksha === 'Shukla' ? '◕' : '◔';
  } else if (tithiNumber < 23) {
    symbol = panchangam.tithi.paksha === 'Krishna' ? '◑' : '◐';
  } else {
    symbol = panchangam.tithi.paksha === 'Krishna' ? '◔' : '◕';
  }

  return {
    type: 'moon',
    longitude,
    position,
    symbol,
    color: CHART_COLORS.moon,
    size: 20,
    label: `Moon in ${panchangam.nakshatra.name} (Pada ${panchangam.nakshatra.pada})`,
  };
};

// ============================================================================
// Tithi Arc Calculations
// ============================================================================

/**
 * Calculate tithi arc connecting Sun and Moon
 */
export const getTithiArcInfo = (
  panchangam: PanchangamElements,
  dimensions: ChartDimensions
): TithiArcInfo => {
  const sunLongitude = panchangam.sunPosition.longitude;
  const moonLongitude = panchangam.moonPosition.longitude;

  // Calculate the shorter arc between Sun and Moon
  let angleDiff = normalizeAngle(moonLongitude - sunLongitude);

  // Generate arc path
  const arcRadius = dimensions.celestialOrbit;
  const arcPath = createArcPath(
    dimensions.center,
    arcRadius,
    sunLongitude,
    moonLongitude
  );

  return {
    sunAngle: sunLongitude,
    moonAngle: moonLongitude,
    radius: arcRadius,
    tithi: panchangam.tithi,
    arcPath,
  };
};

// ============================================================================
// Segment Detection
// ============================================================================

/**
 * Get rashi index from longitude (0-11)
 */
export const getRashiIndex = (longitude: number): number => {
  return Math.floor(normalizeAngle(longitude) / 30);
};

/**
 * Get nakshatra index from longitude (0-26)
 */
export const getNakshatraIndex = (longitude: number): number => {
  return Math.floor(normalizeAngle(longitude) / DEGREES_PER_NAKSHATRA);
};

/**
 * Get pada index from longitude (0-107)
 */
export const getPadaIndex = (longitude: number): number => {
  return Math.floor(normalizeAngle(longitude) / DEGREES_PER_PADA);
};

// ============================================================================
// Highlight Calculations
// ============================================================================

/**
 * Check if a rashi contains the Sun or Moon
 */
export const isRashiHighlighted = (
  rashiIndex: number,
  panchangam: PanchangamElements
): 'sun' | 'moon' | 'both' | null => {
  const sunRashiIndex = getRashiIndex(panchangam.sunPosition.longitude);
  const moonRashiIndex = getRashiIndex(panchangam.moonPosition.longitude);

  const hasSun = rashiIndex === sunRashiIndex;
  const hasMoon = rashiIndex === moonRashiIndex;

  if (hasSun && hasMoon) return 'both';
  if (hasSun) return 'sun';
  if (hasMoon) return 'moon';
  return null;
};

/**
 * Check if a nakshatra is the current one (Moon's position)
 */
export const isNakshatraHighlighted = (
  nakshatraIndex: number,
  panchangam: PanchangamElements
): boolean => {
  const moonNakshatraIndex = getNakshatraIndex(panchangam.moonPosition.longitude);
  return nakshatraIndex === moonNakshatraIndex;
};

/**
 * Check if a pada is the current one (Moon's position)
 */
export const isPadaHighlighted = (
  padaIndex: number,
  panchangam: PanchangamElements
): boolean => {
  const moonPadaIndex = getPadaIndex(panchangam.moonPosition.longitude);
  return padaIndex === moonPadaIndex;
};

// ============================================================================
// Tooltip Content Generation
// ============================================================================

/**
 * Generate tooltip content for a rashi
 */
export const getRashiTooltipContent = (rashi: RashiInfo): string => {
  return `${rashi.symbol} ${rashi.name} (${rashi.westernName})
${rashi.startDegree}° - ${rashi.endDegree}°
Element: ${rashi.element}
Ruler: ${rashi.ruler}`;
};

/**
 * Generate tooltip content for a nakshatra
 */
export const getNakshatraTooltipContent = (nakshatra: NakshatraInfo): string => {
  return `${nakshatra.name}
${nakshatra.startDegree.toFixed(2)}° - ${nakshatra.endDegree.toFixed(2)}°
Deity: ${nakshatra.deity}
Symbol: ${nakshatra.symbol}`;
};

/**
 * Generate tooltip content for a pada
 */
export const getPadaTooltipContent = (pada: PadaInfo): string => {
  return `${pada.nakshatra.name} - Pada ${pada.padaInNakshatra}
${pada.startDegree.toFixed(2)}° - ${pada.endDegree.toFixed(2)}°
Navamsha: ${pada.navamsha}`;
};

/**
 * Generate tooltip content for tithi
 */
export const getTithiTooltipContent = (tithi: TithiInfo): string => {
  return `${tithi.name} (${tithi.paksha} Paksha)
Tithi ${tithi.number}/30
Angular separation: ${tithi.angle.toFixed(1)}°
${tithi.percentComplete.toFixed(0)}% complete`;
};
