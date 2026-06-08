import { describe, expect, it } from 'vitest';

import {
  countCalendarDaysInclusive,
  getCalendarDayDifference,
  getCalendarDayOfYear,
  parseApiDate,
} from './dateHelpers';

describe('dateHelpers', () => {
  it('parses API dates as local calendar dates', () => {
    const parsed = parseApiDate('2024-01-15');

    expect(parsed.getFullYear()).toBe(2024);
    expect(parsed.getMonth()).toBe(0);
    expect(parsed.getDate()).toBe(15);
  });

  it('counts calendar days across daylight-saving changes', () => {
    expect(countCalendarDaysInclusive(new Date(2024, 10, 1), new Date(2024, 10, 30))).toBe(30);
  });

  it('measures signed calendar day differences across daylight-saving changes', () => {
    const beforeFallBack = new Date(2024, 10, 1);
    const afterFallBack = new Date(2024, 10, 6, 12);

    expect(getCalendarDayDifference(beforeFallBack, afterFallBack)).toBe(5);
    expect(getCalendarDayDifference(afterFallBack, beforeFallBack)).toBe(-5);
  });

  it('gets the calendar day of year across daylight-saving changes', () => {
    expect(getCalendarDayOfYear(new Date(2024, 2, 11))).toBe(71);
  });
});
