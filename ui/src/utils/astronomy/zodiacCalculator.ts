// Zodiac (Rashi) calculation utilities

/**
 * Calculate the Zodiac sign (Rashi) from ecliptic longitude
 * @param longitude Ecliptic longitude in degrees
 * @returns Zodiac sign number (1-12) where 1=Aries/Mesha, 2=Taurus/Vrishabha, etc.
 */
export function calculateZodiacFromLongitude(longitude: number): number {
  // Normalize longitude to 0-360 degrees
  let normalizedLong = longitude;
  while (normalizedLong < 0) normalizedLong += 360;
  while (normalizedLong >= 360) normalizedLong -= 360;

  // Each zodiac sign spans 30 degrees
  // There are 12 zodiac signs covering the full 360° ecliptic
  const zodiacSpan = 360.0 / 12.0; // 30 degrees
  
  // Calculate zodiac sign number (1-12)
  const zodiacNumber = Math.floor(normalizedLong / zodiacSpan) + 1;

  // Ensure zodiac number is in valid range (1-12)
  return Math.max(1, Math.min(12, zodiacNumber));
}

/**
 * Get Sanskrit Rashi name from number
 */
export function getRashiName(zodiacNumber: number): string {
  const names = [
    '', // 0 - not used
    'Mesha', 'Vrishabha', 'Mithuna', 'Karka', 'Simha', 'Kanya',
    'Tula', 'Vrishchika', 'Dhanus', 'Makara', 'Kumbha', 'Meena'
  ];
  
  if (zodiacNumber >= 1 && zodiacNumber <= 12) {
    return names[zodiacNumber];
  }
  
  return 'Unknown';
}

/**
 * Get Western zodiac name from number
 */
export function getWesternZodiacName(zodiacNumber: number): string {
  const names = [
    '', // 0 - not used
    'Aries', 'Taurus', 'Gemini', 'Cancer', 'Leo', 'Virgo',
    'Libra', 'Scorpio', 'Sagittarius', 'Capricorn', 'Aquarius', 'Pisces'
  ];
  
  if (zodiacNumber >= 1 && zodiacNumber <= 12) {
    return names[zodiacNumber];
  }
  
  return 'Unknown';
}

/**
 * Get zodiac symbol from number
 */
export function getZodiacSymbol(zodiacNumber: number): string {
  const symbols = [
    '', // 0 - not used
    '♈', '♉', '♊', '♋', '♌', '♍',
    '♎', '♏', '♐', '♑', '♒', '♓'
  ];
  
  if (zodiacNumber >= 1 && zodiacNumber <= 12) {
    return symbols[zodiacNumber];
  }
  
  return '?';
}

/**
 * Get zodiac element from number
 */
export function getZodiacElement(zodiacNumber: number): string {
  const elements = [
    '', // 0 - not used
    'Fire', 'Earth', 'Air', 'Water', 'Fire', 'Earth',
    'Air', 'Water', 'Fire', 'Earth', 'Air', 'Water'
  ];
  
  if (zodiacNumber >= 1 && zodiacNumber <= 12) {
    return elements[zodiacNumber];
  }
  
  return 'Unknown';
}

/**
 * Get zodiac ruling planet from number
 */
export function getZodiacRuler(zodiacNumber: number): string {
  const rulers = [
    '', // 0 - not used
    'Mars', 'Venus', 'Mercury', 'Moon', 'Sun', 'Mercury',
    'Venus', 'Mars', 'Jupiter', 'Saturn', 'Saturn', 'Jupiter'
  ];
  
  if (zodiacNumber >= 1 && zodiacNumber <= 12) {
    return rulers[zodiacNumber];
  }
  
  return 'Unknown';
}

/**
 * Get comprehensive Zodiac information 
 */
export interface ZodiacInfo {
  number: number;
  sanskritName: string;
  westernName: string;
  symbol: string;
  element: string;
  ruler: string;
  longitude: number;
  degreeInSign: number; // 0-29.999... degrees within the sign
}

export function getZodiacInfo(longitude: number): ZodiacInfo {
  const number = calculateZodiacFromLongitude(longitude);
  const zodiacSpan = 30; // degrees per sign
  const startLongitude = (number - 1) * zodiacSpan;
  const degreeInSign = longitude - startLongitude;
  
  return {
    number,
    sanskritName: getRashiName(number),
    westernName: getWesternZodiacName(number),
    symbol: getZodiacSymbol(number),
    element: getZodiacElement(number),
    ruler: getZodiacRuler(number),
    longitude,
    degreeInSign
  };
}

/**
 * Calculate the degree and minute within a zodiac sign
 * @param longitude Ecliptic longitude in degrees
 * @returns Object with degrees and minutes within the current sign
 */
export function getDegreesMinutesInSign(longitude: number): { degrees: number, minutes: number } {
  const zodiacInfo = getZodiacInfo(longitude);
  const totalDegreeInSign = zodiacInfo.degreeInSign;
  
  const degrees = Math.floor(totalDegreeInSign);
  const minutes = Math.floor((totalDegreeInSign - degrees) * 60);
  
  return { degrees, minutes };
}

/**
 * Format zodiac position as a string (e.g., "15°30' Mesha" or "25°45' Leo")
 */
export function formatZodiacPosition(longitude: number, useWestern: boolean = false): string {
  const zodiacInfo = getZodiacInfo(longitude);
  const { degrees, minutes } = getDegreesMinutesInSign(longitude);
  
  const signName = useWestern ? zodiacInfo.westernName : zodiacInfo.sanskritName;
  return `${degrees}°${minutes.toString().padStart(2, '0')}' ${signName}`;
}