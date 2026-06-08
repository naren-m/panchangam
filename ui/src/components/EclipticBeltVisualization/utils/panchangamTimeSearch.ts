import type { TithiInfo } from '../types/eclipticBelt';
import {
  DEGREES_PER_KARANA,
  DEGREES_PER_NAKSHATRA,
  DEGREES_PER_TITHI,
  DEGREES_PER_YOGA,
} from '../types/eclipticBelt';
import { calculateAngularDifference, normalizeDegrees } from './panchangamAngles';
import { calculateTithi } from './panchangamElements';

export type PositionCalculator = (date: Date) => { sunLong: number; moonLong: number };

export function calculateTithiStartTime(
  currentDate: Date,
  tithiNumber: number,
  calculatePositions: PositionCalculator,
): Date {
  const maxSearchHours = 48;
  const precisionMs = 60000;
  let low = new Date(currentDate.getTime() - maxSearchHours * 60 * 60 * 1000);
  let high = new Date(currentDate.getTime());

  while (high.getTime() - low.getTime() > precisionMs) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const angle = calculateAngularDifference(moonLong, sunLong);

    if (isAngleInTithiRange(angle, tithiNumber)) {
      high = mid;
    } else {
      low = mid;
    }
  }

  return high;
}

export function calculateTithiEndTime(
  currentDate: Date,
  tithiNumber: number,
  calculatePositions: PositionCalculator,
): Date {
  const maxSearchHours = 48;
  const precisionMs = 60000;
  let low = new Date(currentDate.getTime());
  let high = new Date(currentDate.getTime() + maxSearchHours * 60 * 60 * 1000);

  while (high.getTime() - low.getTime() > precisionMs) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const angle = calculateAngularDifference(moonLong, sunLong);

    if (isAngleInTithiRange(angle, tithiNumber)) {
      low = mid;
    } else {
      high = mid;
    }
  }

  return low;
}

function isAngleInTithiRange(angle: number, tithiNumber: number): boolean {
  const startAngle = ((tithiNumber - 1) * DEGREES_PER_TITHI) % 360;
  const endAngle = (tithiNumber * DEGREES_PER_TITHI) % 360;
  if (startAngle > endAngle) {
    return angle >= startAngle || angle < endAngle;
  }
  return angle >= startAngle && angle < endAngle;
}

export function calculateTithiWithTimes(
  sunLongitude: number,
  moonLongitude: number,
  currentDate: Date,
  calculatePositions: PositionCalculator,
): TithiInfo {
  const basicTithi = calculateTithi(sunLongitude, moonLongitude);
  const startTime = calculateTithiStartTime(currentDate, basicTithi.number, calculatePositions);
  const endTime = calculateTithiEndTime(currentDate, basicTithi.number, calculatePositions);

  return {
    ...basicTithi,
    startTime,
    endTime,
  };
}

function isAngleInYogaRange(combinedAngle: number, yogaNumber: number): boolean {
  const startAngle = ((yogaNumber - 1) * DEGREES_PER_YOGA) % 360;
  const endAngle = (yogaNumber * DEGREES_PER_YOGA) % 360;
  if (startAngle > endAngle) {
    return combinedAngle >= startAngle || combinedAngle < endAngle;
  }
  return combinedAngle >= startAngle && combinedAngle < endAngle;
}

export function calculateYogaStartTime(
  currentDate: Date,
  yogaNumber: number,
  calculatePositions: PositionCalculator,
): Date {
  const maxSearchHours = 48;
  const precisionMs = 60000;
  let low = new Date(currentDate.getTime() - maxSearchHours * 60 * 60 * 1000);
  let high = new Date(currentDate.getTime());

  while (high.getTime() - low.getTime() > precisionMs) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const combinedAngle = normalizeDegrees(sunLong + moonLong);

    if (isAngleInYogaRange(combinedAngle, yogaNumber)) {
      high = mid;
    } else {
      low = mid;
    }
  }

  return high;
}

export function calculateYogaEndTime(
  currentDate: Date,
  yogaNumber: number,
  calculatePositions: PositionCalculator,
): Date {
  const maxSearchHours = 48;
  const precisionMs = 60000;
  let low = new Date(currentDate.getTime());
  let high = new Date(currentDate.getTime() + maxSearchHours * 60 * 60 * 1000);

  while (high.getTime() - low.getTime() > precisionMs) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const combinedAngle = normalizeDegrees(sunLong + moonLong);

    if (isAngleInYogaRange(combinedAngle, yogaNumber)) {
      low = mid;
    } else {
      high = mid;
    }
  }

  return low;
}

function isAngleInKaranaRange(angle: number, karanaNumber: number): boolean {
  const startAngle = ((karanaNumber - 1) * DEGREES_PER_KARANA) % 360;
  const endAngle = (karanaNumber * DEGREES_PER_KARANA) % 360;
  if (startAngle > endAngle) {
    return angle >= startAngle || angle < endAngle;
  }
  return angle >= startAngle && angle < endAngle;
}

export function calculateKaranaStartTime(
  currentDate: Date,
  karanaNumber: number,
  calculatePositions: PositionCalculator,
): Date {
  const maxSearchHours = 24;
  const precisionMs = 60000;
  let low = new Date(currentDate.getTime() - maxSearchHours * 60 * 60 * 1000);
  let high = new Date(currentDate.getTime());

  while (high.getTime() - low.getTime() > precisionMs) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const angle = calculateAngularDifference(moonLong, sunLong);

    if (isAngleInKaranaRange(angle, karanaNumber)) {
      high = mid;
    } else {
      low = mid;
    }
  }

  return high;
}

export function calculateKaranaEndTime(
  currentDate: Date,
  karanaNumber: number,
  calculatePositions: PositionCalculator,
): Date {
  const maxSearchHours = 24;
  const precisionMs = 60000;
  let low = new Date(currentDate.getTime());
  let high = new Date(currentDate.getTime() + maxSearchHours * 60 * 60 * 1000);

  while (high.getTime() - low.getTime() > precisionMs) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { sunLong, moonLong } = calculatePositions(mid);
    const angle = calculateAngularDifference(moonLong, sunLong);

    if (isAngleInKaranaRange(angle, karanaNumber)) {
      low = mid;
    } else {
      high = mid;
    }
  }

  return low;
}

function isLongitudeInNakshatraRange(moonLongitude: number, nakshatraNumber: number): boolean {
  const startLongitude = ((nakshatraNumber - 1) * DEGREES_PER_NAKSHATRA) % 360;
  const endLongitude = (nakshatraNumber * DEGREES_PER_NAKSHATRA) % 360;
  if (startLongitude > endLongitude) {
    return moonLongitude >= startLongitude || moonLongitude < endLongitude;
  }
  return moonLongitude >= startLongitude && moonLongitude < endLongitude;
}

export function calculateNakshatraStartTime(
  currentDate: Date,
  nakshatraNumber: number,
  calculatePositions: PositionCalculator,
): Date {
  const maxSearchHours = 36;
  const precisionMs = 60000;
  let low = new Date(currentDate.getTime() - maxSearchHours * 60 * 60 * 1000);
  let high = new Date(currentDate.getTime());

  while (high.getTime() - low.getTime() > precisionMs) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { moonLong } = calculatePositions(mid);
    const normalizedMoon = normalizeDegrees(moonLong);

    if (isLongitudeInNakshatraRange(normalizedMoon, nakshatraNumber)) {
      high = mid;
    } else {
      low = mid;
    }
  }

  return high;
}

export function calculateNakshatraEndTime(
  currentDate: Date,
  nakshatraNumber: number,
  calculatePositions: PositionCalculator,
): Date {
  const maxSearchHours = 36;
  const precisionMs = 60000;
  let low = new Date(currentDate.getTime());
  let high = new Date(currentDate.getTime() + maxSearchHours * 60 * 60 * 1000);

  while (high.getTime() - low.getTime() > precisionMs) {
    const mid = new Date((low.getTime() + high.getTime()) / 2);
    const { moonLong } = calculatePositions(mid);
    const normalizedMoon = normalizeDegrees(moonLong);

    if (isLongitudeInNakshatraRange(normalizedMoon, nakshatraNumber)) {
      low = mid;
    } else {
      high = mid;
    }
  }

  return low;
}
