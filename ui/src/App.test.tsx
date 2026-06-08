import { describe, expect, it, beforeEach, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';

import App from './App';
import type { Settings } from './types/panchangam';

const hookMocks = vi.hoisted(() => ({
  useProgressivePanchangam: vi.fn(),
  useDayDetail: vi.fn(),
}));

vi.mock('./hooks/useProgressivePanchangam', () => ({
  useProgressivePanchangam: hookMocks.useProgressivePanchangam,
}));

vi.mock('./hooks/useDayDetail', () => ({
  useDayDetail: hookMocks.useDayDetail,
}));

vi.mock('./components/Calendar/CalendarDisplayManager', () => ({
  CalendarDisplayManager: () => <div data-testid="calendar-display" />,
}));

vi.mock('./components/Settings/SettingsPanel', () => ({
  SettingsPanel: ({
    settings,
    onSettingsChange,
  }: {
    settings: Settings;
    onSettingsChange: (settings: Settings) => void;
  }) => (
    <div role="dialog" aria-label="Settings panel">
      <span>Calendar system: {settings.calendar_system ?? 'purnimanta'}</span>
      <button
        type="button"
        onClick={() => onSettingsChange({ ...settings, calendar_system: 'solar' })}
      >
        Use solar calendar
      </button>
    </div>
  ),
}));

const getLatestSettings = (): Settings => {
  const calls = hookMocks.useProgressivePanchangam.mock.calls;
  return calls[calls.length - 1][2] as Settings;
};

describe('App settings', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    hookMocks.useProgressivePanchangam.mockReturnValue({
      data: {},
      loading: false,
      isProgressiveLoading: false,
      progress: 0,
      todayLoaded: false,
      loadedCount: 0,
      totalCount: 0,
      error: null,
      errorState: {
        hasError: false,
        message: null,
        isNetworkError: false,
      },
      loadingPhase: {
        phase: 'complete',
        description: 'All data loaded',
      },
      retry: vi.fn(),
    });
    hookMocks.useDayDetail.mockReturnValue({
      data: null,
      isLoading: false,
      error: null,
      retry: vi.fn(),
    });
  });

  it('keeps the same settings object when unrelated App state changes', async () => {
    render(<App />);

    const firstSettings = getLatestSettings();

    fireEvent.click(screen.getByRole('button', { name: 'Settings' }));

    await waitFor(() => {
      expect(hookMocks.useProgressivePanchangam.mock.calls.length).toBeGreaterThan(1);
    });
    expect(getLatestSettings()).toBe(firstSettings);
  });

  it('keeps calendar system changes in the settings passed to data loaders', async () => {
    render(<App />);

    fireEvent.click(screen.getByRole('button', { name: 'Settings' }));
    fireEvent.click(await screen.findByRole('button', { name: 'Use solar calendar' }));

    await waitFor(() => {
      expect(getLatestSettings().calendar_system).toBe('solar');
    });
  });
});
