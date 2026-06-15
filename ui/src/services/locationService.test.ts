import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { popularLocations, US_UK_LOCATION_COUNT } from './locationData';
import { locationService } from './locationService';

describe('location service data boundary', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.spyOn(console, 'error').mockImplementation(() => undefined);
    vi.spyOn(console, 'info').mockImplementation(() => undefined);
    vi.spyOn(console, 'warn').mockImplementation(() => undefined);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('keeps Milpitas as the default location', () => {
    expect(popularLocations[0]).toMatchObject({
      name: 'Milpitas, California',
      latitude: 37.4323,
      longitude: -121.9066,
      timezone: 'America/Los_Angeles',
      region: 'California',
    });
  });

  it('uses the shared US/UK count when building location categories', () => {
    const categories = locationService.getLocationsByCategory();

    expect(categories.favorites).toEqual([]);
    expect(categories.usUk).toEqual(popularLocations.slice(0, US_UK_LOCATION_COUNT));
    expect(categories.popular).toEqual(popularLocations.slice(US_UK_LOCATION_COUNT));
  });

  it('returns shared popular locations without copying business logic into the service', () => {
    expect(locationService.getPopularLocations()).toBe(popularLocations);
  });

  it('uses an empty favorites list when stored favorites are invalid', () => {
    localStorage.setItem('panchangam_favorite_locations', '{bad json');

    expect(locationService.getFavoriteLocations()).toEqual([]);
    expect(console.warn).not.toHaveBeenCalled();
  });

  it('uses an empty favorites list when stored favorites are not a list', () => {
    localStorage.setItem('panchangam_favorite_locations', '{"name":"Milpitas"}');

    expect(locationService.getFavoriteLocations()).toEqual([]);
  });

  it('adds, deduplicates, and removes favorite locations by nearby coordinates', () => {
    const milpitas = popularLocations[0];
    const nearMilpitas = {
      ...milpitas,
      name: 'Near Milpitas',
      latitude: milpitas.latitude + 0.0005,
      longitude: milpitas.longitude - 0.0005,
    };

    locationService.addToFavorites(milpitas);
    locationService.addToFavorites(nearMilpitas);

    expect(locationService.getFavoriteLocations()).toEqual([milpitas]);
    expect(locationService.isFavorite(nearMilpitas)).toBe(true);

    locationService.removeFromFavorites(nearMilpitas);

    expect(locationService.getFavoriteLocations()).toEqual([]);
  });

  it('falls back to popular locations when geocoding fails without logging', async () => {
    vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce(new Error('network down'));

    const results = await locationService.searchLocations('Milpitas');

    expect(results.length).toBeGreaterThan(0);
    expect(results[0].name).toBe('Milpitas, California');
    expect(console.error).not.toHaveBeenCalled();
    expect(console.warn).not.toHaveBeenCalled();
  });
});
