import { describe, expect, it } from 'vitest';

import { calculateChartPanchangamForDate } from './chartCalculations';

describe('chart calculations', () => {
  it('derives chart panchangam positions from the selected date', () => {
    const panchangam = calculateChartPanchangamForDate(new Date(2024, 0, 1, 12));

    expect(panchangam.sunPosition.longitude).toBe(281);
    expect(panchangam.moonPosition.longitude).toBeCloseTo(19.8);
    expect(panchangam.tithi.number).toBe(9);
  });
});
