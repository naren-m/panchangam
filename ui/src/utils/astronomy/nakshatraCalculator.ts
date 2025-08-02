// Nakshatra calculation utilities

/**
 * Calculate the Nakshatra number from Moon's ecliptic longitude
 * @param moonLongitude Moon's ecliptic longitude in degrees
 * @returns Nakshatra number (1-27)
 */
export function calculateNakshatraFromLongitude(moonLongitude: number): number {
  // Normalize longitude to 0-360 degrees
  let normalizedLong = moonLongitude;
  while (normalizedLong < 0) normalizedLong += 360;
  while (normalizedLong >= 360) normalizedLong -= 360;

  // Each Nakshatra spans 13°20' (13.333... degrees)
  // There are 27 Nakshatras covering the full 360° zodiac
  const nakshatraSpan = 360.0 / 27.0; // 13.333... degrees
  
  // Calculate Nakshatra number (1-27)
  const nakshatraNumber = Math.floor(normalizedLong / nakshatraSpan) + 1;

  // Ensure Nakshatra number is in valid range (1-27)
  return Math.max(1, Math.min(27, nakshatraNumber));
}

/**
 * Calculate the Pada (1-4) from Moon's ecliptic longitude
 * @param moonLongitude Moon's ecliptic longitude in degrees
 * @returns Pada number (1-4)
 */
export function calculatePadaFromLongitude(moonLongitude: number): number {
  // Normalize longitude to 0-360 degrees
  let normalizedLong = moonLongitude;
  while (normalizedLong < 0) normalizedLong += 360;
  while (normalizedLong >= 360) normalizedLong -= 360;

  const nakshatraSpan = 360.0 / 27.0; // 13.333... degrees
  const padaSpan = nakshatraSpan / 4.0; // 3.333... degrees
  
  // Find position within current nakshatra
  const nakshatraNumber = Math.floor(normalizedLong / nakshatraSpan);
  const positionInNakshatra = normalizedLong - (nakshatraNumber * nakshatraSpan);
  
  // Calculate Pada (1-4)
  const pada = Math.floor(positionInNakshatra / padaSpan) + 1;
  
  // Ensure Pada is in valid range (1-4)
  return Math.max(1, Math.min(4, pada));
}

/**
 * Get Nakshatra name from number
 */
export function getNakshatraName(nakshatraNumber: number): string {
  const names = [
    '', // 0 - not used
    'Ashwini', 'Bharani', 'Krittika', 'Rohini', 'Mrigashira', 'Ardra',
    'Punarvasu', 'Pushya', 'Ashlesha', 'Magha', 'Purva Phalguni', 'Uttara Phalguni',
    'Hasta', 'Chitra', 'Swati', 'Vishakha', 'Anuradha', 'Jyeshtha',
    'Mula', 'Purva Ashadha', 'Uttara Ashadha', 'Shravana', 'Dhanishta', 'Shatabhisha',
    'Purva Bhadrapada', 'Uttara Bhadrapada', 'Revati'
  ];
  
  if (nakshatraNumber >= 1 && nakshatraNumber <= 27) {
    return names[nakshatraNumber];
  }
  
  return 'Unknown';
}

/**
 * Get Nakshatra deity from number
 */
export function getNakshatraDeity(nakshatraNumber: number): string {
  const deities = [
    '', // 0 - not used
    'Ashwini Kumaras', 'Yama', 'Agni', 'Brahma', 'Soma', 'Rudra',
    'Aditi', 'Brihaspati', 'Nagas', 'Pitrs (Ancestors)', 'Bhaga', 'Aryaman',
    'Savitar', 'Tvashtar', 'Vayu', 'Indra-Agni', 'Mitra', 'Indra',
    'Nirriti', 'Apas', 'Vishve Devas', 'Vishnu', 'Vasus', 'Varuna',
    'Aja Ekapada', 'Ahir Budhnya', 'Pushan'
  ];
  
  if (nakshatraNumber >= 1 && nakshatraNumber <= 27) {
    return deities[nakshatraNumber];
  }
  
  return 'Unknown';
}

/**
 * Get Nakshatra symbol from number
 */
export function getNakshatraSymbol(nakshatraNumber: number): string {
  const symbols = [
    '', // 0 - not used
    'Horse Head', 'Yoni', 'Razor/Knife', 'Cart/Chariot', 'Deer Head', 'Teardrop/Diamond',
    'Bow and Quiver', 'Cow Udder', 'Serpent', 'Throne', 'Front Legs of Bed', 'Back Legs of Bed',
    'Hand', 'Bright Jewel', 'Young Shoot of Plant', 'Triumphal Arch', 'Lotus', 'Circular Amulet',
    'Bunch of Roots', 'Elephant Tusk', 'Elephant Tusk', 'Ear/Three Footprints', 'Drum', 'Empty Circle',
    'Front Legs of Funeral Cot', 'Back Legs of Funeral Cot', 'Fish/Pair of Fish'
  ];
  
  if (nakshatraNumber >= 1 && nakshatraNumber <= 27) {
    return symbols[nakshatraNumber];
  }
  
  return 'Unknown';
}

/**
 * Get comprehensive Nakshatra information 
 */
export interface NakshatraInfo {
  number: number;
  name: string;
  deity: string;
  symbol: string;
  pada: number;
  longitude: number;
}

export function getNakshatraInfo(moonLongitude: number): NakshatraInfo {
  const number = calculateNakshatraFromLongitude(moonLongitude);
  const pada = calculatePadaFromLongitude(moonLongitude);
  
  return {
    number,
    name: getNakshatraName(number),
    deity: getNakshatraDeity(number),
    symbol: getNakshatraSymbol(number),
    pada,
    longitude: moonLongitude
  };
}