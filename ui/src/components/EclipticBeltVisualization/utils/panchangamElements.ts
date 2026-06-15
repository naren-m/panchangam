import type {
  KaranaInfo,
  MoonPosition,
  NakshatraInfo,
  PanchangamElements,
  RashiInfo,
  SunPosition,
  TithiInfo,
  YogaInfo,
} from '../types/eclipticBelt';
import {
  DEGREES_PER_KARANA,
  DEGREES_PER_NAKSHATRA,
  DEGREES_PER_RASHI,
  DEGREES_PER_TITHI,
  DEGREES_PER_YOGA,
} from '../types/eclipticBelt';
import { calculateAngularDifference, normalizeDegrees } from './panchangamAngles';
import {
  KARANA_INFO,
  KARANA_NAMES,
  NAKSHATRA_DATA,
  RASHI_DATA,
  TITHI_DEITIES,
  TITHI_NAMES,
  YOGA_INFO,
  YOGA_NAMES,
} from './panchangamElementData';

export function calculateTithi(sunLongitude: number, moonLongitude: number): TithiInfo {
  const angle = calculateAngularDifference(moonLongitude, sunLongitude);
  const tithiFloat = angle / DEGREES_PER_TITHI;
  const number = Math.floor(tithiFloat) + 1;
  const clampedNumber = Math.max(1, Math.min(30, number));
  const percentComplete = (tithiFloat - Math.floor(tithiFloat)) * 100;
  const paksha: 'Shukla' | 'Krishna' = clampedNumber <= 15 ? 'Shukla' : 'Krishna';

  return {
    number: clampedNumber,
    name: TITHI_NAMES[clampedNumber],
    paksha,
    deity: TITHI_DEITIES[clampedNumber],
    angle,
    percentComplete,
  };
}

export function getTithiDisplayName(tithi: TithiInfo): string {
  if (tithi.number === 15) return 'Purnima (Full Moon)';
  if (tithi.number === 30) return 'Amavasya (New Moon)';

  const tithiInPaksha = tithi.paksha === 'Shukla' ? tithi.number : tithi.number - 15;
  return `${tithi.paksha} ${tithi.name} (${tithiInPaksha})`;
}

export function calculateYoga(sunLongitude: number, moonLongitude: number): YogaInfo {
  const combinedLongitude = normalizeDegrees(sunLongitude + moonLongitude);
  const yogaFloat = combinedLongitude / DEGREES_PER_YOGA;
  const number = Math.floor(yogaFloat) + 1;
  const clampedNumber = Math.max(1, Math.min(27, number));
  const info = YOGA_INFO[clampedNumber];

  return {
    number: clampedNumber,
    name: YOGA_NAMES[clampedNumber],
    meaning: info.meaning,
    nature: info.nature,
    combinedLongitude,
  };
}

export function calculateKarana(sunLongitude: number, moonLongitude: number): KaranaInfo {
  const angle = calculateAngularDifference(moonLongitude, sunLongitude);
  const karanaFloat = angle / DEGREES_PER_KARANA;
  const karanaNumber = Math.floor(karanaFloat) + 1;
  const clampedKaranaNumber = Math.max(1, Math.min(60, karanaNumber));

  let karanaTypeIndex: number;
  if (clampedKaranaNumber === 1) {
    karanaTypeIndex = 10;
  } else if (clampedKaranaNumber === 57) {
    karanaTypeIndex = 7;
  } else if (clampedKaranaNumber === 58) {
    karanaTypeIndex = 8;
  } else if (clampedKaranaNumber === 59) {
    karanaTypeIndex = 9;
  } else if (clampedKaranaNumber === 60) {
    karanaTypeIndex = 10;
  } else {
    karanaTypeIndex = (clampedKaranaNumber - 2) % 7;
  }

  const info = KARANA_INFO[karanaTypeIndex];
  return {
    number: clampedKaranaNumber,
    name: KARANA_NAMES[karanaTypeIndex],
    type: info.type,
    nature: info.nature,
  };
}

export function calculateRashi(longitude: number): RashiInfo {
  const normalized = normalizeDegrees(longitude);
  const number = Math.floor(normalized / DEGREES_PER_RASHI) + 1;
  const clampedNumber = Math.max(1, Math.min(12, number));
  const data = RASHI_DATA[clampedNumber];

  return {
    ...data,
    startDegree: (clampedNumber - 1) * DEGREES_PER_RASHI,
    endDegree: clampedNumber * DEGREES_PER_RASHI,
  };
}

export function calculateNakshatra(moonLongitude: number): NakshatraInfo {
  const normalized = normalizeDegrees(moonLongitude);
  const number = Math.floor(normalized / DEGREES_PER_NAKSHATRA) + 1;
  const clampedNumber = Math.max(1, Math.min(27, number));
  const startDegree = (clampedNumber - 1) * DEGREES_PER_NAKSHATRA;
  const endDegree = clampedNumber * DEGREES_PER_NAKSHATRA;
  const positionInNakshatra = normalized - startDegree;
  const padaSpan = DEGREES_PER_NAKSHATRA / 4;
  const pada = Math.floor(positionInNakshatra / padaSpan) + 1;
  const data = NAKSHATRA_DATA[clampedNumber];

  return {
    ...data,
    startDegree,
    endDegree,
    pada: Math.max(1, Math.min(4, pada)),
  };
}

export function calculateSunPosition(sunLongitude: number): SunPosition {
  return {
    longitude: normalizeDegrees(sunLongitude),
    rashi: calculateRashi(sunLongitude),
    dailyMotion: 0.9856,
  };
}

export function calculateMoonPosition(moonLongitude: number, sunLongitude: number): MoonPosition {
  const normalized = normalizeDegrees(moonLongitude);
  const angle = calculateAngularDifference(moonLongitude, sunLongitude);

  return {
    longitude: normalized,
    nakshatra: calculateNakshatra(moonLongitude),
    rashi: calculateRashi(moonLongitude),
    dailyMotion: 13.176,
    phase: angle / 360,
  };
}

export function calculatePanchangamElements(
  sunLongitude: number,
  moonLongitude: number,
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
    moonPosition,
  };
}
