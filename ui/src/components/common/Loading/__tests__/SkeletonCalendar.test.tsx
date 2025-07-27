import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import * as matchers from '@testing-library/jest-dom/matchers';
import { SkeletonCalendar } from '../SkeletonCalendar';

expect.extend(matchers);

describe('SkeletonCalendar', () => {
  it('renders calendar skeleton structure', () => {
    const { container } = render(<SkeletonCalendar />);
    
    // Check for main container
    const mainContainer = container.querySelector('.bg-white.rounded-lg.shadow-lg');
    expect(mainContainer).toBeInTheDocument();
  });

  it('renders 7 weekday headers by default', () => {
    const { container } = render(<SkeletonCalendar />);
    
    // Check weekday headers grid
    const weekdayGrid = container.querySelector('.grid-cols-7');
    expect(weekdayGrid).toBeInTheDocument();
    
    // Count weekday skeleton items (should be 7)
    const weekdaySkeletons = container.querySelectorAll('.grid-cols-7:first-child > div');
    expect(weekdaySkeletons).toHaveLength(7);
  });

  it('renders correct number of calendar cells based on rows prop', () => {
    const rows = 5;
    const { container } = render(<SkeletonCalendar rows={rows} />);
    
    // Should render rows * 7 cells (5 * 7 = 35)
    const calendarCells = container.querySelectorAll('.grid-cols-7:last-child > div');
    expect(calendarCells).toHaveLength(rows * 7);
  });

  it('renders default 6 rows when no rows prop provided', () => {
    const { container } = render(<SkeletonCalendar />);
    
    // Should render 6 * 7 = 42 cells by default
    const calendarCells = container.querySelectorAll('.grid-cols-7:last-child > div');
    expect(calendarCells).toHaveLength(42);
  });

  it('has pulse animation', () => {
    const { container } = render(<SkeletonCalendar />);
    
    const mainContainer = container.firstChild;
    expect(mainContainer).toHaveClass('animate-pulse');
  });

  it('renders skeleton elements within each cell', () => {
    const { container } = render(<SkeletonCalendar rows={1} />);
    
    // Check first calendar cell for skeleton elements
    const firstCell = container.querySelector('.grid-cols-7:last-child > div');
    expect(firstCell).toBeInTheDocument();
    
    // Should have date number skeleton
    const dateNumberSkeleton = firstCell?.querySelector('.h-4.bg-gray-200.w-6');
    expect(dateNumberSkeleton).toBeInTheDocument();
    
    // Should have tithi skeleton
    const tithiSkeleton = firstCell?.querySelector('.h-3.bg-gray-100.w-12');
    expect(tithiSkeleton).toBeInTheDocument();
    
    // Should have event indicators
    const eventIndicators = firstCell?.querySelectorAll('.h-2.w-2.bg-gray-100.rounded-full');
    expect(eventIndicators).toHaveLength(2);
  });
});