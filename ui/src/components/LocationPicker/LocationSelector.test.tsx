import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';

import { LocationSelector } from './LocationSelector';
import { locationService } from '../../services/locationService';
import type { Location } from '../../types/panchangam';

vi.mock('../../services/locationService', () => ({
  locationService: {
    addToFavorites: vi.fn(),
    getCurrentLocation: vi.fn(),
    getLocationsByCategory: vi.fn(),
    isFavorite: vi.fn(),
    removeFromFavorites: vi.fn(),
    searchLocations: vi.fn(),
  },
}));

const currentLocation: Location = {
  name: 'Milpitas, California',
  latitude: 37.4323,
  longitude: -121.9066,
  timezone: 'America/Los_Angeles',
  region: 'California',
};

const mockedLocationService = vi.mocked(locationService);

describe('LocationSelector', () => {
  const onLocationSelect = vi.fn();
  const onClose = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    mockedLocationService.getLocationsByCategory.mockReturnValue({
      favorites: [],
      usUk: [currentLocation],
      popular: [],
    });
    mockedLocationService.isFavorite.mockReturnValue(false);
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  it('shows empty search results without logging failed searches', async () => {
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    mockedLocationService.searchLocations.mockRejectedValueOnce(new Error('search failed'));

    render(
      <LocationSelector
        currentLocation={currentLocation}
        onLocationSelect={onLocationSelect}
        onClose={onClose}
      />
    );

    fireEvent.change(screen.getByPlaceholderText('Search cities worldwide...'), {
      target: { value: 'Paris' },
    });

    await waitFor(() => {
      expect(mockedLocationService.searchLocations).toHaveBeenCalledWith('Paris');
    });

    expect(screen.getByText('No locations found')).toBeInTheDocument();
    expect(consoleError).not.toHaveBeenCalled();
  });

  it('keeps the manual-selection alert without logging failed GPS lookup', async () => {
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    const alert = vi.fn();
    vi.stubGlobal('alert', alert);
    mockedLocationService.getCurrentLocation.mockRejectedValueOnce(new Error('gps failed'));

    render(
      <LocationSelector
        currentLocation={currentLocation}
        onLocationSelect={onLocationSelect}
        onClose={onClose}
      />
    );

    fireEvent.click(screen.getByRole('button', { name: /use current location/i }));

    await waitFor(() => {
      expect(alert).toHaveBeenCalledWith('Unable to get your location. Please select manually.');
    });

    expect(consoleError).not.toHaveBeenCalled();
    expect(onLocationSelect).not.toHaveBeenCalled();
    expect(onClose).not.toHaveBeenCalled();
  });
});
