import { describe, expect, it } from 'vitest';

import { normalizeDegrees } from './panchangamAngles';
import { calculatePanchangamElements } from './panchangamElements';
import {
  calculateTithiEndTime,
  calculateTithiStartTime,
  type PositionCalculator,
} from './panchangamTimeSearch';

describe('panchangam calculation modules', () => {
  it('normalizes negative and overflow angles', () => {
    expect(normalizeDegrees(-10)).toBe(350);
    expect(normalizeDegrees(725)).toBe(5);
  });

  it('calculates the core panchangam elements from sun and moon longitude', () => {
    const elements = calculatePanchangamElements(10, 25);

    expect(elements.tithi.number).toBe(2);
    expect(elements.tithi.paksha).toBe('Shukla');
    expect(elements.karana.number).toBe(3);
    expect(elements.yoga.number).toBe(3);
    expect(elements.nakshatra.number).toBe(2);
    expect(elements.rashi.name).toBe('Mesha');
    expect(elements.sunRashi.name).toBe('Mesha');
  });

  it('finds tithi boundary times with the shared time search', () => {
    const base = Date.parse('2024-01-01T00:00:00Z');
    const currentDate = new Date('2024-01-01T18:00:00Z');
    const calculatePositions: PositionCalculator = (date) => ({
      sunLong: 0,
      moonLong: (date.getTime() - base) / (60 * 60 * 1000),
    });

    const start = calculateTithiStartTime(currentDate, 2, calculatePositions);
    const end = calculateTithiEndTime(currentDate, 2, calculatePositions);

    expect(Math.abs(start.getTime() - Date.parse('2024-01-01T12:00:00Z'))).toBeLessThan(60 * 1000);
    expect(Math.abs(end.getTime() - Date.parse('2024-01-02T00:00:00Z'))).toBeLessThan(60 * 1000);
  });
});
