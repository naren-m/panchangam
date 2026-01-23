/**
 * Panchangam Element Calculator
 *
 * Calculates Tithi, Yoga, and Karana from Sun and Moon positions.
 * These three elements, along with Nakshatra and Vara (weekday),
 * form the five limbs (panch-anga) of the Hindu calendar.
 *
 * Formulas:
 * - Tithi: (Moon longitude - Sun longitude) / 12° → 30 tithis
 * - Yoga: (Sun longitude + Moon longitude) / 13.333° → 27 yogas
 * - Karana: Half-tithi (60 per month) → 11 unique types
 */

import {
  TithiInfo,
  YogaInfo,
  KaranaInfo,
  RashiInfo,
  NakshatraInfo,
  SunPosition,
  MoonPosition,
  PanchangamElements,
  DEGREES_PER_TITHI,
  DEGREES_PER_YOGA,
  DEGREES_PER_KARANA,
  DEGREES_PER_RASHI,
  DEGREES_PER_NAKSHATRA,
} from '../types/eclipticBelt';

// ============================================================================
// Normalization Utilities
// ============================================================================

/**
 * Normalize angle to 0-360 degrees
 */
export function normalizeDegrees(degrees: number): number {
  let normalized = degrees;
  while (normalized < 0) normalized += 360;
  while (normalized >= 360) normalized -= 360;
  return normalized;
}

/**
 * Calculate angular difference (Moon - Sun) handling wrap-around
 */
export function calculateAngularDifference(moonLong: number, sunLong: number): number {
  const diff = normalizeDegrees(moonLong) - normalizeDegrees(sunLong);
  return normalizeDegrees(diff);
}

// ============================================================================
// Tithi Calculation
// ============================================================================

/**
 * 30 Tithi names with their meanings
 */
const TITHI_NAMES = [
  '', // 0 - not used
  'Pratipada', 'Dvitiya', 'Tritiya', 'Chaturthi', 'Panchami',
  'Shashthi', 'Saptami', 'Ashtami', 'Navami', 'Dashami',
  'Ekadashi', 'Dvadashi', 'Trayodashi', 'Chaturdashi', 'Purnima',  // Shukla Paksha (1-15)
  'Pratipada', 'Dvitiya', 'Tritiya', 'Chaturthi', 'Panchami',
  'Shashthi', 'Saptami', 'Ashtami', 'Navami', 'Dashami',
  'Ekadashi', 'Dvadashi', 'Trayodashi', 'Chaturdashi', 'Amavasya'  // Krishna Paksha (16-30)
];

/**
 * Tithi deities (presiding deities for each tithi)
 */
const TITHI_DEITIES = [
  '', // 0 - not used
  'Brahma', 'Vidhata', 'Vishnu', 'Yama', 'Moon',
  'Karttikeya', 'Indra', 'Vasus', 'Sarpa', 'Dharma',
  'Rudra', 'Adityas', 'Kama', 'Shiva', 'Soma',  // Shukla Paksha
  'Agni', 'Brahma', 'Govinda', 'Pitrs', 'Nagas',
  'Vishve Devas', 'Maruts', 'Ashta Vasus', 'Aditi', 'Vishnu',
  'Brahma', 'Vishnu', 'Shiva', 'Chitragupta', 'Pitrs'  // Krishna Paksha
];

/**
 * Calculate Tithi from Sun and Moon longitudes
 *
 * @param sunLongitude Sun's ecliptic longitude in degrees
 * @param moonLongitude Moon's ecliptic longitude in degrees
 * @returns Complete Tithi information
 */
export function calculateTithi(sunLongitude: number, moonLongitude: number): TithiInfo {
  // Angular separation between Moon and Sun
  const angle = calculateAngularDifference(moonLongitude, sunLongitude);

  // Tithi number (1-30), each tithi spans 12°
  const tithiFloat = angle / DEGREES_PER_TITHI;
  const number = Math.floor(tithiFloat) + 1;
  const clampedNumber = Math.max(1, Math.min(30, number));

  // Percentage of current tithi elapsed
  const percentComplete = (tithiFloat - Math.floor(tithiFloat)) * 100;

  // Determine Paksha (fortnight)
  // Shukla Paksha (bright fortnight): Tithi 1-15 (New Moon to Full Moon)
  // Krishna Paksha (dark fortnight): Tithi 16-30 (Full Moon to New Moon)
  const paksha: 'Shukla' | 'Krishna' = clampedNumber <= 15 ? 'Shukla' : 'Krishna';

  // Get the name - use modulo for Krishna Paksha to get proper names
  const nameIndex = clampedNumber <= 15 ? clampedNumber : clampedNumber;

  return {
    number: clampedNumber,
    name: TITHI_NAMES[nameIndex],
    paksha,
    deity: TITHI_DEITIES[nameIndex],
    angle,
    percentComplete
  };
}

/**
 * Get display name for Tithi (includes Paksha)
 */
export function getTithiDisplayName(tithi: TithiInfo): string {
  if (tithi.number === 15) return 'Purnima (Full Moon)';
  if (tithi.number === 30) return 'Amavasya (New Moon)';

  const tithiInPaksha = tithi.paksha === 'Shukla' ? tithi.number : tithi.number - 15;
  return `${tithi.paksha} ${tithi.name} (${tithiInPaksha})`;
}

// ============================================================================
// Yoga Calculation
// ============================================================================

/**
 * 27 Yoga names
 */
const YOGA_NAMES = [
  '', // 0 - not used
  'Vishkambha', 'Priti', 'Ayushman', 'Saubhagya', 'Shobhana',
  'Atiganda', 'Sukarma', 'Dhriti', 'Shula', 'Ganda',
  'Vriddhi', 'Dhruva', 'Vyaghata', 'Harshana', 'Vajra',
  'Siddhi', 'Vyatipata', 'Variyan', 'Parigha', 'Shiva',
  'Siddha', 'Sadhya', 'Shubha', 'Shukla', 'Brahma',
  'Indra', 'Vaidhriti'
];

/**
 * Yoga meanings and nature
 */
const YOGA_INFO: { meaning: string; nature: 'Auspicious' | 'Inauspicious' | 'Mixed' }[] = [
  { meaning: '', nature: 'Mixed' }, // 0 - not used
  { meaning: 'Obstacle remover', nature: 'Auspicious' },
  { meaning: 'Love', nature: 'Auspicious' },
  { meaning: 'Long life', nature: 'Auspicious' },
  { meaning: 'Good fortune', nature: 'Auspicious' },
  { meaning: 'Splendour', nature: 'Auspicious' },
  { meaning: 'Danger', nature: 'Inauspicious' },
  { meaning: 'Good work', nature: 'Auspicious' },
  { meaning: 'Determination', nature: 'Auspicious' },
  { meaning: 'Thorn/Pain', nature: 'Inauspicious' },
  { meaning: 'Danger', nature: 'Inauspicious' },
  { meaning: 'Growth', nature: 'Auspicious' },
  { meaning: 'Steadfast', nature: 'Auspicious' },
  { meaning: 'Destruction', nature: 'Inauspicious' },
  { meaning: 'Joy', nature: 'Auspicious' },
  { meaning: 'Thunder', nature: 'Mixed' },
  { meaning: 'Perfection', nature: 'Auspicious' },
  { meaning: 'Great fall', nature: 'Inauspicious' },
  { meaning: 'Excellence', nature: 'Auspicious' },
  { meaning: 'Obstruction', nature: 'Inauspicious' },
  { meaning: 'Auspiciousness', nature: 'Auspicious' },
  { meaning: 'Accomplishment', nature: 'Auspicious' },
  { meaning: 'Achievement', nature: 'Auspicious' },
  { meaning: 'Auspicious', nature: 'Auspicious' },
  { meaning: 'Bright', nature: 'Auspicious' },
  { meaning: 'Creator', nature: 'Auspicious' },
  { meaning: 'King of gods', nature: 'Auspicious' },
  { meaning: 'Great danger', nature: 'Inauspicious' }
];

/**
 * Calculate Yoga from Sun and Moon longitudes
 *
 * Yoga = (Sun longitude + Moon longitude) / 13.333°
 *
 * @param sunLongitude Sun's ecliptic longitude in degrees
 * @param moonLongitude Moon's ecliptic longitude in degrees
 * @returns Complete Yoga information
 */
export function calculateYoga(sunLongitude: number, moonLongitude: number): YogaInfo {
  // Combined longitude (sum)
  const combinedLongitude = normalizeDegrees(sunLongitude + moonLongitude);

  // Yoga number (1-27), each yoga spans 13.333°
  const yogaFloat = combinedLongitude / DEGREES_PER_YOGA;
  const number = Math.floor(yogaFloat) + 1;
  const clampedNumber = Math.max(1, Math.min(27, number));

  const info = YOGA_INFO[clampedNumber];

  return {
    number: clampedNumber,
    name: YOGA_NAMES[clampedNumber],
    meaning: info.meaning,
    nature: info.nature,
    combinedLongitude
  };
}

// ============================================================================
// Karana Calculation
// ============================================================================

/**
 * 11 Karana types
 *
 * There are 60 karanas in a lunar month (2 per tithi).
 * 4 are fixed (occur once each): Shakuni, Chatushpada, Nagava, Kimstughna
 * 7 are movable (occur 8 times each): Bava, Balava, Kaulava, Taitila, Gara, Vanija, Vishti
 */
const KARANA_NAMES = [
  'Bava',      // 1 - Movable
  'Balava',    // 2 - Movable
  'Kaulava',   // 3 - Movable
  'Taitila',   // 4 - Movable
  'Gara',      // 5 - Movable (also called Garaja)
  'Vanija',    // 6 - Movable
  'Vishti',    // 7 - Movable (also called Bhadra, inauspicious)
  'Shakuni',   // 8 - Fixed
  'Chatushpada', // 9 - Fixed
  'Nagava',    // 10 - Fixed (also called Naga)
  'Kimstughna' // 11 - Fixed (also called Kinstughna)
];

const KARANA_INFO: { type: 'Movable' | 'Fixed'; nature: 'Auspicious' | 'Inauspicious' | 'Mixed' }[] = [
  { type: 'Movable', nature: 'Auspicious' },    // Bava
  { type: 'Movable', nature: 'Auspicious' },    // Balava
  { type: 'Movable', nature: 'Auspicious' },    // Kaulava
  { type: 'Movable', nature: 'Auspicious' },    // Taitila
  { type: 'Movable', nature: 'Auspicious' },    // Gara
  { type: 'Movable', nature: 'Auspicious' },    // Vanija
  { type: 'Movable', nature: 'Inauspicious' },  // Vishti (Bhadra)
  { type: 'Fixed', nature: 'Mixed' },           // Shakuni
  { type: 'Fixed', nature: 'Mixed' },           // Chatushpada
  { type: 'Fixed', nature: 'Mixed' },           // Nagava
  { type: 'Fixed', nature: 'Auspicious' }       // Kimstughna
];

/**
 * Calculate Karana from Sun and Moon longitudes
 *
 * Each tithi has 2 karanas (first half and second half)
 * Total 60 karanas per lunar month
 *
 * The fixed karanas occur at specific positions:
 * - Kimstughna: First half of Shukla Pratipada (karana 1)
 * - Shakuni: Second half of Krishna Chaturdashi (karana 57)
 * - Chatushpada: First half of Amavasya (karana 58)
 * - Nagava: Second half of Amavasya first half (karana 59)
 * - Kimstughna again: Second half of Amavasya (karana 60)
 *
 * The 7 movable karanas repeat in order for the remaining 56 karanas.
 *
 * @param sunLongitude Sun's ecliptic longitude in degrees
 * @param moonLongitude Moon's ecliptic longitude in degrees
 * @returns Complete Karana information
 */
export function calculateKarana(sunLongitude: number, moonLongitude: number): KaranaInfo {
  // Angular separation between Moon and Sun
  const angle = calculateAngularDifference(moonLongitude, sunLongitude);

  // Karana number (1-60), each karana spans 6°
  const karanaFloat = angle / DEGREES_PER_KARANA;
  const karanaNumber = Math.floor(karanaFloat) + 1;
  const clampedKaranaNumber = Math.max(1, Math.min(60, karanaNumber));

  // Determine which of the 11 karana types this is
  let karanaTypeIndex: number;

  // Fixed karanas at specific positions
  if (clampedKaranaNumber === 1) {
    karanaTypeIndex = 10; // Kimstughna
  } else if (clampedKaranaNumber === 57) {
    karanaTypeIndex = 7;  // Shakuni
  } else if (clampedKaranaNumber === 58) {
    karanaTypeIndex = 8;  // Chatushpada
  } else if (clampedKaranaNumber === 59) {
    karanaTypeIndex = 9;  // Nagava
  } else if (clampedKaranaNumber === 60) {
    karanaTypeIndex = 10; // Kimstughna
  } else {
    // Movable karanas (2-56): cycle through 7 movable types
    // Position 2 starts with Bava (index 0)
    karanaTypeIndex = (clampedKaranaNumber - 2) % 7;
  }

  const info = KARANA_INFO[karanaTypeIndex];

  return {
    number: clampedKaranaNumber,
    name: KARANA_NAMES[karanaTypeIndex],
    type: info.type,
    nature: info.nature
  };
}

// ============================================================================
// Rashi (Zodiac) Calculation
// ============================================================================

/**
 * Rashi (Zodiac Sign) data
 */
const RASHI_DATA: Omit<RashiInfo, 'startDegree' | 'endDegree'>[] = [
  { number: 0, name: '', westernName: '', symbol: '', element: 'Fire', ruler: '' }, // placeholder
  { number: 1, name: 'Mesha', westernName: 'Aries', symbol: '\u2648', element: 'Fire', ruler: 'Mars' },
  { number: 2, name: 'Vrishabha', westernName: 'Taurus', symbol: '\u2649', element: 'Earth', ruler: 'Venus' },
  { number: 3, name: 'Mithuna', westernName: 'Gemini', symbol: '\u264A', element: 'Air', ruler: 'Mercury' },
  { number: 4, name: 'Karka', westernName: 'Cancer', symbol: '\u264B', element: 'Water', ruler: 'Moon' },
  { number: 5, name: 'Simha', westernName: 'Leo', symbol: '\u264C', element: 'Fire', ruler: 'Sun' },
  { number: 6, name: 'Kanya', westernName: 'Virgo', symbol: '\u264D', element: 'Earth', ruler: 'Mercury' },
  { number: 7, name: 'Tula', westernName: 'Libra', symbol: '\u264E', element: 'Air', ruler: 'Venus' },
  { number: 8, name: 'Vrishchika', westernName: 'Scorpio', symbol: '\u264F', element: 'Water', ruler: 'Mars' },
  { number: 9, name: 'Dhanus', westernName: 'Sagittarius', symbol: '\u2650', element: 'Fire', ruler: 'Jupiter' },
  { number: 10, name: 'Makara', westernName: 'Capricorn', symbol: '\u2651', element: 'Earth', ruler: 'Saturn' },
  { number: 11, name: 'Kumbha', westernName: 'Aquarius', symbol: '\u2652', element: 'Air', ruler: 'Saturn' },
  { number: 12, name: 'Meena', westernName: 'Pisces', symbol: '\u2653', element: 'Water', ruler: 'Jupiter' }
];

/**
 * Calculate Rashi (Zodiac Sign) from ecliptic longitude
 */
export function calculateRashi(longitude: number): RashiInfo {
  const normalized = normalizeDegrees(longitude);
  const number = Math.floor(normalized / DEGREES_PER_RASHI) + 1;
  const clampedNumber = Math.max(1, Math.min(12, number));

  const data = RASHI_DATA[clampedNumber];

  return {
    ...data,
    startDegree: (clampedNumber - 1) * DEGREES_PER_RASHI,
    endDegree: clampedNumber * DEGREES_PER_RASHI
  };
}

// ============================================================================
// Nakshatra Calculation (uses existing calculator patterns)
// ============================================================================

/**
 * Nakshatra data
 */
const NAKSHATRA_DATA: Omit<NakshatraInfo, 'startDegree' | 'endDegree' | 'pada'>[] = [
  { number: 0, name: '', deity: '', symbol: '' }, // placeholder
  { number: 1, name: 'Ashwini', deity: 'Ashwini Kumaras', symbol: 'Horse Head' },
  { number: 2, name: 'Bharani', deity: 'Yama', symbol: 'Yoni' },
  { number: 3, name: 'Krittika', deity: 'Agni', symbol: 'Razor' },
  { number: 4, name: 'Rohini', deity: 'Brahma', symbol: 'Cart' },
  { number: 5, name: 'Mrigashira', deity: 'Soma', symbol: 'Deer Head' },
  { number: 6, name: 'Ardra', deity: 'Rudra', symbol: 'Teardrop' },
  { number: 7, name: 'Punarvasu', deity: 'Aditi', symbol: 'Bow' },
  { number: 8, name: 'Pushya', deity: 'Brihaspati', symbol: 'Lotus' },
  { number: 9, name: 'Ashlesha', deity: 'Nagas', symbol: 'Serpent' },
  { number: 10, name: 'Magha', deity: 'Pitrs', symbol: 'Throne' },
  { number: 11, name: 'Purva Phalguni', deity: 'Bhaga', symbol: 'Hammock' },
  { number: 12, name: 'Uttara Phalguni', deity: 'Aryaman', symbol: 'Bed' },
  { number: 13, name: 'Hasta', deity: 'Savitar', symbol: 'Hand' },
  { number: 14, name: 'Chitra', deity: 'Tvashtar', symbol: 'Pearl' },
  { number: 15, name: 'Swati', deity: 'Vayu', symbol: 'Coral' },
  { number: 16, name: 'Vishakha', deity: 'Indragni', symbol: 'Arch' },
  { number: 17, name: 'Anuradha', deity: 'Mitra', symbol: 'Lotus' },
  { number: 18, name: 'Jyeshtha', deity: 'Indra', symbol: 'Earring' },
  { number: 19, name: 'Mula', deity: 'Nirriti', symbol: 'Roots' },
  { number: 20, name: 'Purva Ashadha', deity: 'Apas', symbol: 'Fan' },
  { number: 21, name: 'Uttara Ashadha', deity: 'Vishve Devas', symbol: 'Tusk' },
  { number: 22, name: 'Shravana', deity: 'Vishnu', symbol: 'Ear' },
  { number: 23, name: 'Dhanishta', deity: 'Vasus', symbol: 'Drum' },
  { number: 24, name: 'Shatabhisha', deity: 'Varuna', symbol: 'Circle' },
  { number: 25, name: 'Purva Bhadrapada', deity: 'Aja Ekapada', symbol: 'Sword' },
  { number: 26, name: 'Uttara Bhadrapada', deity: 'Ahir Budhnya', symbol: 'Twin' },
  { number: 27, name: 'Revati', deity: 'Pushan', symbol: 'Fish' }
];

/**
 * Calculate Nakshatra from Moon's ecliptic longitude
 */
export function calculateNakshatra(moonLongitude: number): NakshatraInfo {
  const normalized = normalizeDegrees(moonLongitude);
  const number = Math.floor(normalized / DEGREES_PER_NAKSHATRA) + 1;
  const clampedNumber = Math.max(1, Math.min(27, number));

  const startDegree = (clampedNumber - 1) * DEGREES_PER_NAKSHATRA;
  const endDegree = clampedNumber * DEGREES_PER_NAKSHATRA;

  // Calculate pada (1-4) - quarter of nakshatra
  const positionInNakshatra = normalized - startDegree;
  const padaSpan = DEGREES_PER_NAKSHATRA / 4;
  const pada = Math.floor(positionInNakshatra / padaSpan) + 1;

  const data = NAKSHATRA_DATA[clampedNumber];

  return {
    ...data,
    startDegree,
    endDegree,
    pada: Math.max(1, Math.min(4, pada))
  };
}

// ============================================================================
// Sun and Moon Position Calculations
// ============================================================================

/**
 * Calculate Sun position with Rashi
 */
export function calculateSunPosition(sunLongitude: number): SunPosition {
  return {
    longitude: normalizeDegrees(sunLongitude),
    rashi: calculateRashi(sunLongitude),
    dailyMotion: 0.9856  // Average daily motion in degrees
  };
}

/**
 * Calculate Moon position with Nakshatra and Rashi
 */
export function calculateMoonPosition(moonLongitude: number, sunLongitude: number): MoonPosition {
  const normalized = normalizeDegrees(moonLongitude);

  // Calculate phase (0 = new moon, 0.5 = full moon, 1 = new moon again)
  const angle = calculateAngularDifference(moonLongitude, sunLongitude);
  const phase = angle / 360;

  return {
    longitude: normalized,
    nakshatra: calculateNakshatra(moonLongitude),
    rashi: calculateRashi(moonLongitude),
    dailyMotion: 13.176,  // Average daily motion in degrees
    phase
  };
}

// ============================================================================
// Complete Panchangam Calculation
// ============================================================================

/**
 * Calculate all Panchangam elements from Sun and Moon positions
 *
 * @param sunLongitude Sun's sidereal ecliptic longitude in degrees
 * @param moonLongitude Moon's sidereal ecliptic longitude in degrees
 * @returns Complete Panchangam elements
 */
export function calculatePanchangamElements(
  sunLongitude: number,
  moonLongitude: number
): PanchangamElements {
  const sunPosition = calculateSunPosition(sunLongitude);
  const moonPosition = calculateMoonPosition(moonLongitude, sunLongitude);

  return {
    tithi: calculateTithi(sunLongitude, moonLongitude),
    nakshatra: moonPosition.nakshatra,
    yoga: calculateYoga(sunLongitude, moonLongitude),
    karana: calculateKarana(sunLongitude, moonLongitude),
    rashi: moonPosition.rashi,
    sunRashi: sunPosition.rashi,
    sunPosition,
    moonPosition
  };
}

// ============================================================================
// Tithi Time Calculations
// ============================================================================

/**
 * Position calculator function type
 * Used for calculating Sun/Moon positions at any given time
 */
export type PositionCalculator = (date: Date) => { sunLong: number; moonLong: number };

/**
 * Calculate the exact time when a tithi started
 *
 * Uses binary search to find when Moon-Sun angular separation
 * crossed into the current tithi's 12° boundary.
 *
 * @param currentDate The current date/time
 * @param tithiNumber The current tithi number (1-30)
 * @param calculatePositions Function to compute Sun/Moon positions for a date
 * @returns The Date when this tithi started
 */
export function calculateTithiStartTime(
  currentDate: Date,
  tithiNumber: number,
  calculatePositions: PositionCalculator
): Date {
  // Calculate the starting angle boundary for this tithi
  // Tithi 1 starts at 0°, Tithi 2 at 12°, Tithi 3 at 24°, etc.
  const startBoundary = ((tithiNumber - 1) * DEGREES_PER_TITHI) % 360;

  // Binary search parameters
  // Tithis average ~24 hours, so search up to 48 hours back
  const MAX_SEARCH_HOURS = 48;
  const PRECISION_MS = 60000; // 1 minute precision

  let low = new Date(currentDate.getTime() - MAX_SEARCH_HOURS * 60 * 60 * 1000);
  let high = new Date(currentDate.getTime());

  // Binary search to find when angle crossed the boundary
  while (high.getTime() - low.getTime() > PRECISION_MS) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const angle = calculateAngularDifference(moonLong, sunLong);

    // Determine if mid is before or after the boundary crossing
    // Handle the wrap-around case (e.g., tithi 30 to tithi 1)
    const currentAngleAboveBoundary = isAngleInTithiRange(angle, tithiNumber);

    if (currentAngleAboveBoundary) {
      // We're still in the current tithi, search earlier
      high = mid;
    } else {
      // We've gone before the tithi start, search later
      low = mid;
    }
  }

  return high;
}

/**
 * Calculate the exact time when a tithi will end
 *
 * @param currentDate The current date/time
 * @param tithiNumber The current tithi number (1-30)
 * @param calculatePositions Function to compute Sun/Moon positions for a date
 * @returns The Date when this tithi will end
 */
export function calculateTithiEndTime(
  currentDate: Date,
  tithiNumber: number,
  calculatePositions: PositionCalculator
): Date {
  // Calculate the ending angle boundary for this tithi
  const endBoundary = (tithiNumber * DEGREES_PER_TITHI) % 360;

  // Binary search parameters
  const MAX_SEARCH_HOURS = 48;
  const PRECISION_MS = 60000; // 1 minute precision

  let low = new Date(currentDate.getTime());
  let high = new Date(currentDate.getTime() + MAX_SEARCH_HOURS * 60 * 60 * 1000);

  // Binary search to find when angle will cross the boundary
  while (high.getTime() - low.getTime() > PRECISION_MS) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const angle = calculateAngularDifference(moonLong, sunLong);

    // Determine if mid is before or after the boundary crossing
    const stillInCurrentTithi = isAngleInTithiRange(angle, tithiNumber);

    if (stillInCurrentTithi) {
      // We're still in the current tithi, search later
      low = mid;
    } else {
      // We've passed the tithi end, search earlier
      high = mid;
    }
  }

  return low;
}

/**
 * Check if an angle falls within a specific tithi's range
 *
 * @param angle Moon-Sun angular separation (0-360°)
 * @param tithiNumber The tithi number to check (1-30)
 * @returns true if the angle is within this tithi's range
 */
function isAngleInTithiRange(angle: number, tithiNumber: number): boolean {
  const startAngle = ((tithiNumber - 1) * DEGREES_PER_TITHI) % 360;
  const endAngle = (tithiNumber * DEGREES_PER_TITHI) % 360;

  // Handle the wrap-around case for tithi 30 (348° to 360°/0°)
  if (startAngle > endAngle) {
    // Wrap-around: angle must be >= start OR < end
    return angle >= startAngle || angle < endAngle;
  }

  // Normal case: angle must be in [start, end)
  return angle >= startAngle && angle < endAngle;
}

/**
 * Calculate Tithi with start and end times
 *
 * @param sunLongitude Sun's ecliptic longitude in degrees
 * @param moonLongitude Moon's ecliptic longitude in degrees
 * @param currentDate The current date/time (for time calculations)
 * @param calculatePositions Function to compute positions for any date
 * @returns Complete Tithi information including start/end times
 */
export function calculateTithiWithTimes(
  sunLongitude: number,
  moonLongitude: number,
  currentDate: Date,
  calculatePositions: PositionCalculator
): import('../types/eclipticBelt').TithiInfo {
  // First get the basic tithi info
  const basicTithi = calculateTithi(sunLongitude, moonLongitude);

  // Calculate start and end times
  const startTime = calculateTithiStartTime(
    currentDate,
    basicTithi.number,
    calculatePositions
  );

  const endTime = calculateTithiEndTime(
    currentDate,
    basicTithi.number,
    calculatePositions
  );

  return {
    ...basicTithi,
    startTime,
    endTime
  };
}

// ============================================================================
// Educational Helpers
// ============================================================================

/**
 * Get educational explanation for a Tithi
 */
export function getTithiExplanation(tithi: TithiInfo): string {
  const pakshaName = tithi.paksha === 'Shukla' ? 'bright (waxing)' : 'dark (waning)';
  return `${tithi.name} is tithi ${tithi.number} of the ${pakshaName} fortnight. ` +
    `The Moon is ${tithi.angle.toFixed(1)}° ahead of the Sun. ` +
    `Deity: ${tithi.deity}. ${tithi.percentComplete.toFixed(0)}% of this tithi has elapsed.`;
}

/**
 * Get educational explanation for a Yoga
 */
export function getYogaExplanation(yoga: YogaInfo): string {
  return `${yoga.name} yoga (#${yoga.number}) means "${yoga.meaning}". ` +
    `It is ${yoga.nature.toLowerCase()}. ` +
    `Calculated from Sun + Moon = ${yoga.combinedLongitude.toFixed(1)}°.`;
}

/**
 * Get educational explanation for a Karana
 */
export function getKaranaExplanation(karana: KaranaInfo): string {
  return `${karana.name} is a ${karana.type.toLowerCase()} karana (#${karana.number} in this lunar month). ` +
    `It is considered ${karana.nature.toLowerCase()} for activities.`;
}

// ============================================================================
// Yoga Time Calculations
// ============================================================================

/**
 * Check if an angle falls within a specific yoga's range
 *
 * @param combinedAngle Sun + Moon combined longitude (0-360°)
 * @param yogaNumber The yoga number to check (1-27)
 * @returns true if the angle is within this yoga's range
 */
function isAngleInYogaRange(combinedAngle: number, yogaNumber: number): boolean {
  const startAngle = ((yogaNumber - 1) * DEGREES_PER_YOGA) % 360;
  const endAngle = (yogaNumber * DEGREES_PER_YOGA) % 360;

  // Handle the wrap-around case
  if (startAngle > endAngle) {
    return combinedAngle >= startAngle || combinedAngle < endAngle;
  }

  return combinedAngle >= startAngle && combinedAngle < endAngle;
}

/**
 * Calculate the exact time when a yoga started
 *
 * Uses binary search to find when Sun + Moon combined longitude
 * crossed into the current yoga's 13.333° boundary.
 *
 * @param currentDate The current date/time
 * @param yogaNumber The current yoga number (1-27)
 * @param calculatePositions Function to compute Sun/Moon positions for a date
 * @returns The Date when this yoga started
 */
export function calculateYogaStartTime(
  currentDate: Date,
  yogaNumber: number,
  calculatePositions: PositionCalculator
): Date {
  const MAX_SEARCH_HOURS = 48;
  const PRECISION_MS = 60000; // 1 minute precision

  let low = new Date(currentDate.getTime() - MAX_SEARCH_HOURS * 60 * 60 * 1000);
  let high = new Date(currentDate.getTime());

  while (high.getTime() - low.getTime() > PRECISION_MS) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const combinedAngle = normalizeDegrees(sunLong + moonLong);

    const currentAngleInRange = isAngleInYogaRange(combinedAngle, yogaNumber);

    if (currentAngleInRange) {
      high = mid;
    } else {
      low = mid;
    }
  }

  return high;
}

/**
 * Calculate the exact time when a yoga will end
 *
 * @param currentDate The current date/time
 * @param yogaNumber The current yoga number (1-27)
 * @param calculatePositions Function to compute Sun/Moon positions for a date
 * @returns The Date when this yoga will end
 */
export function calculateYogaEndTime(
  currentDate: Date,
  yogaNumber: number,
  calculatePositions: PositionCalculator
): Date {
  const MAX_SEARCH_HOURS = 48;
  const PRECISION_MS = 60000;

  let low = new Date(currentDate.getTime());
  let high = new Date(currentDate.getTime() + MAX_SEARCH_HOURS * 60 * 60 * 1000);

  while (high.getTime() - low.getTime() > PRECISION_MS) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const combinedAngle = normalizeDegrees(sunLong + moonLong);

    const stillInCurrentYoga = isAngleInYogaRange(combinedAngle, yogaNumber);

    if (stillInCurrentYoga) {
      low = mid;
    } else {
      high = mid;
    }
  }

  return low;
}

// ============================================================================
// Karana Time Calculations
// ============================================================================

/**
 * Check if an angle falls within a specific karana's range
 *
 * @param angle Moon-Sun angular separation (0-360°)
 * @param karanaNumber The karana number to check (1-60)
 * @returns true if the angle is within this karana's range
 */
function isAngleInKaranaRange(angle: number, karanaNumber: number): boolean {
  const startAngle = ((karanaNumber - 1) * DEGREES_PER_KARANA) % 360;
  const endAngle = (karanaNumber * DEGREES_PER_KARANA) % 360;

  // Handle the wrap-around case
  if (startAngle > endAngle) {
    return angle >= startAngle || angle < endAngle;
  }

  return angle >= startAngle && angle < endAngle;
}

/**
 * Calculate the exact time when a karana started
 *
 * Uses binary search to find when Moon-Sun angular separation
 * crossed into the current karana's 6° boundary.
 *
 * @param currentDate The current date/time
 * @param karanaNumber The current karana number (1-60)
 * @param calculatePositions Function to compute Sun/Moon positions for a date
 * @returns The Date when this karana started
 */
export function calculateKaranaStartTime(
  currentDate: Date,
  karanaNumber: number,
  calculatePositions: PositionCalculator
): Date {
  const MAX_SEARCH_HOURS = 24; // Karanas are shorter (~12 hours)
  const PRECISION_MS = 60000;

  let low = new Date(currentDate.getTime() - MAX_SEARCH_HOURS * 60 * 60 * 1000);
  let high = new Date(currentDate.getTime());

  while (high.getTime() - low.getTime() > PRECISION_MS) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const angle = calculateAngularDifference(moonLong, sunLong);

    const currentAngleInRange = isAngleInKaranaRange(angle, karanaNumber);

    if (currentAngleInRange) {
      high = mid;
    } else {
      low = mid;
    }
  }

  return high;
}

/**
 * Calculate the exact time when a karana will end
 *
 * @param currentDate The current date/time
 * @param karanaNumber The current karana number (1-60)
 * @param calculatePositions Function to compute Sun/Moon positions for a date
 * @returns The Date when this karana will end
 */
export function calculateKaranaEndTime(
  currentDate: Date,
  karanaNumber: number,
  calculatePositions: PositionCalculator
): Date {
  const MAX_SEARCH_HOURS = 24;
  const PRECISION_MS = 60000;

  let low = new Date(currentDate.getTime());
  let high = new Date(currentDate.getTime() + MAX_SEARCH_HOURS * 60 * 60 * 1000);

  while (high.getTime() - low.getTime() > PRECISION_MS) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const angle = calculateAngularDifference(moonLong, sunLong);

    const stillInCurrentKarana = isAngleInKaranaRange(angle, karanaNumber);

    if (stillInCurrentKarana) {
      low = mid;
    } else {
      high = mid;
    }
  }

  return low;
}

// ============================================================================
// Nakshatra Time Calculations
// ============================================================================

/**
 * Check if a Moon longitude falls within a specific nakshatra's range
 *
 * @param moonLong Moon's ecliptic longitude (0-360°)
 * @param nakshatraNumber The nakshatra number to check (1-27)
 * @returns true if the longitude is within this nakshatra's range
 */
function isLongitudeInNakshatraRange(moonLong: number, nakshatraNumber: number): boolean {
  const startLong = ((nakshatraNumber - 1) * DEGREES_PER_NAKSHATRA) % 360;
  const endLong = (nakshatraNumber * DEGREES_PER_NAKSHATRA) % 360;

  // Handle the wrap-around case for nakshatra 27 (Revati)
  if (startLong > endLong) {
    return moonLong >= startLong || moonLong < endLong;
  }

  return moonLong >= startLong && moonLong < endLong;
}

/**
 * Calculate the exact time when a nakshatra started
 *
 * Uses binary search to find when Moon's longitude
 * crossed into the current nakshatra's 13.333° boundary.
 *
 * @param currentDate The current date/time
 * @param nakshatraNumber The current nakshatra number (1-27)
 * @param calculatePositions Function to compute Sun/Moon positions for a date
 * @returns The Date when this nakshatra started
 */
export function calculateNakshatraStartTime(
  currentDate: Date,
  nakshatraNumber: number,
  calculatePositions: PositionCalculator
): Date {
  const MAX_SEARCH_HOURS = 36; // Nakshatras last ~1 day
  const PRECISION_MS = 60000;

  let low = new Date(currentDate.getTime() - MAX_SEARCH_HOURS * 60 * 60 * 1000);
  let high = new Date(currentDate.getTime());

  while (high.getTime() - low.getTime() > PRECISION_MS) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { moonLong } = calculatePositions(mid);
    const normalizedMoon = normalizeDegrees(moonLong);

    const currentInRange = isLongitudeInNakshatraRange(normalizedMoon, nakshatraNumber);

    if (currentInRange) {
      high = mid;
    } else {
      low = mid;
    }
  }

  return high;
}

/**
 * Calculate the exact time when a nakshatra will end
 *
 * @param currentDate The current date/time
 * @param nakshatraNumber The current nakshatra number (1-27)
 * @param calculatePositions Function to compute Sun/Moon positions for a date
 * @returns The Date when this nakshatra will end
 */
export function calculateNakshatraEndTime(
  currentDate: Date,
  nakshatraNumber: number,
  calculatePositions: PositionCalculator
): Date {
  const MAX_SEARCH_HOURS = 36;
  const PRECISION_MS = 60000;

  let low = new Date(currentDate.getTime());
  let high = new Date(currentDate.getTime() + MAX_SEARCH_HOURS * 60 * 60 * 1000);

  while (high.getTime() - low.getTime() > PRECISION_MS) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { moonLong } = calculatePositions(mid);
    const normalizedMoon = normalizeDegrees(moonLong);

    const stillInCurrentNakshatra = isLongitudeInNakshatraRange(normalizedMoon, nakshatraNumber);

    if (stillInCurrentNakshatra) {
      low = mid;
    } else {
      high = mid;
    }
  }

  return low;
}
