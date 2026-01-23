import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { CalendarDisplayManager } from '../CalendarDisplayManager';
import '@testing-library/jest-dom';

// Mock child components
vi.mock('../CalendarGrid', () => ({
  CalendarGrid: ({ year, month }: { year: number; month: number }) => (
    <div data-testid="calendar-grid">Calendar {year}-{month}</div>
  ),
}));

vi.mock('../../common/Loading/SkeletonCalendar', () => ({
  SkeletonCalendar: () => <div data-testid="skeleton-calendar">Loading...</div>,
}));

vi.mock('../../common/Error', () => ({
  NetworkError: ({ customMessage }: { customMessage: string }) => (
    <div data-testid="network-error">{customMessage}</div>
  ),
  ApiError: ({ error }: { error: string }) => (
    <div data-testid="api-error">{error}</div>
  ),
}));

const mockCalendarProps = {
  year: 2024,
  month: 0,
  panchangamData: {},
  settings: {
    calculation_method: 'Drik',
    locale: 'en',
    region: 'California',
    time_format: '12',
    location: {
      name: 'Test Location',
      latitude: 37.4323,
      longitude: -121.9066,
      timezone: 'America/Los_Angeles',
      region: 'California'
    }
  },
  onDateClick: vi.fn()
};

const mockErrorState = {
  hasError: false,
  message: null,
  isNetworkError: false
};

describe('CalendarDisplayManager', () => {
  it('renders skeleton calendar when loading and no data', () => {
    render(
      <CalendarDisplayManager
        loading={true}
        hasData={false}
        error={null}
        errorState={mockErrorState}
        isProgressiveLoading={false}
        progress={0}
        loadedCount={0}
        totalCount={35}
        retry={vi.fn()}
        calendarProps={mockCalendarProps}
      />
    );

    expect(screen.getByTestId('skeleton-calendar')).toBeInTheDocument();
    expect(screen.queryByTestId('calendar-grid')).not.toBeInTheDocument();
  });

  it('renders calendar when has data', () => {
    render(
      <CalendarDisplayManager
        loading={false}
        hasData={true}
        error={null}
        errorState={mockErrorState}
        isProgressiveLoading={false}
        progress={100}
        loadedCount={35}
        totalCount={35}
        retry={vi.fn()}
        calendarProps={mockCalendarProps}
      />
    );

    expect(screen.getByTestId('calendar-grid')).toBeInTheDocument();
    expect(screen.queryByTestId('skeleton-calendar')).not.toBeInTheDocument();
  });

  it('renders network error when network error occurs', () => {
    const errorState = {
      hasError: true,
      message: 'Network Error',
      isNetworkError: true
    };

    render(
      <CalendarDisplayManager
        loading={false}
        hasData={false}
        error="Network Error"
        errorState={errorState}
        isProgressiveLoading={false}
        progress={0}
        loadedCount={0}
        totalCount={35}
        retry={vi.fn()}
        calendarProps={mockCalendarProps}
      />
    );

    expect(screen.getByTestId('network-error')).toBeInTheDocument();
    expect(screen.queryByTestId('calendar-grid')).not.toBeInTheDocument();
    expect(screen.queryByTestId('skeleton-calendar')).not.toBeInTheDocument();
  });

  it('renders API error when API error occurs', () => {
    const errorState = {
      hasError: true,
      message: 'API Error',
      isNetworkError: false,
      statusCode: 500
    };

    render(
      <CalendarDisplayManager
        loading={false}
        hasData={false}
        error="API Error"
        errorState={errorState}
        isProgressiveLoading={false}
        progress={0}
        loadedCount={0}
        totalCount={35}
        retry={vi.fn()}
        calendarProps={mockCalendarProps}
      />
    );

    expect(screen.getByTestId('api-error')).toBeInTheDocument();
    expect(screen.queryByTestId('calendar-grid')).not.toBeInTheDocument();
    expect(screen.queryByTestId('skeleton-calendar')).not.toBeInTheDocument();
  });

  it('renders calendar with progress indicator when progressively loading', () => {
    render(
      <CalendarDisplayManager
        loading={false}
        hasData={true}
        error={null}
        errorState={mockErrorState}
        isProgressiveLoading={true}
        progress={60}
        loadedCount={21}
        totalCount={35}
        retry={vi.fn()}
        calendarProps={mockCalendarProps}
      />
    );

    expect(screen.getByTestId('calendar-grid')).toBeInTheDocument();
    expect(screen.getByText('Loading calendar data...')).toBeInTheDocument();
    expect(screen.getByText('21/35 days loaded (60%)')).toBeInTheDocument();
    expect(screen.queryByTestId('skeleton-calendar')).not.toBeInTheDocument();
  });

  it('ensures only one calendar type renders at a time', () => {
    const { rerender } = render(
      <CalendarDisplayManager
        loading={true}
        hasData={false}
        error={null}
        errorState={mockErrorState}
        isProgressiveLoading={false}
        progress={0}
        loadedCount={0}
        totalCount={35}
        retry={vi.fn()}
        calendarProps={mockCalendarProps}
      />
    );

    // Initially loading - should show skeleton
    expect(screen.getByTestId('skeleton-calendar')).toBeInTheDocument();
    expect(screen.queryByTestId('calendar-grid')).not.toBeInTheDocument();

    // With data - should show calendar, not skeleton
    rerender(
      <CalendarDisplayManager
        loading={false}
        hasData={true}
        error={null}
        errorState={mockErrorState}
        isProgressiveLoading={false}
        progress={100}
        loadedCount={35}
        totalCount={35}
        retry={vi.fn()}
        calendarProps={mockCalendarProps}
      />
    );

    expect(screen.getByTestId('calendar-grid')).toBeInTheDocument();
    expect(screen.queryByTestId('skeleton-calendar')).not.toBeInTheDocument();
  });

  it('includes proper accessibility attributes', () => {
    render(
      <CalendarDisplayManager
        loading={false}
        hasData={true}
        error={null}
        errorState={mockErrorState}
        isProgressiveLoading={true}
        progress={75}
        loadedCount={26}
        totalCount={35}
        retry={vi.fn()}
        calendarProps={mockCalendarProps}
      />
    );

    // Check for accessibility attributes
    const progressBar = screen.getByRole('progressbar');
    expect(progressBar).toHaveAttribute('aria-valuenow', '75');
    expect(progressBar).toHaveAttribute('aria-valuemin', '0');
    expect(progressBar).toHaveAttribute('aria-valuemax', '100');
    expect(progressBar).toHaveAttribute('aria-label', 'Loading calendar data: 75% complete');

    const calendar = screen.getByRole('main');
    expect(calendar).toHaveAttribute('aria-label', 'Panchangam calendar');
  });
});