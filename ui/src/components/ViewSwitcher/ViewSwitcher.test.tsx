import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ViewSwitcher, ViewMode } from './ViewSwitcher';

describe('ViewSwitcher', () => {
  it('renders all view options', () => {
    const onViewChange = vi.fn();
    render(
      <ViewSwitcher
        currentView="calendar"
        onViewChange={onViewChange}
      />
    );

    expect(screen.getByText('Calendar')).toBeInTheDocument();
    expect(screen.getByText('Table')).toBeInTheDocument();
    expect(screen.getByText('Analytics')).toBeInTheDocument();
  });

  it('highlights the current view', () => {
    const onViewChange = vi.fn();
    render(
      <ViewSwitcher
        currentView="calendar"
        onViewChange={onViewChange}
      />
    );

    const calendarButton = screen.getByText('Calendar').closest('button');
    expect(calendarButton).toHaveClass('from-orange-400');
  });

  it('calls onViewChange when clicking a different view', () => {
    const onViewChange = vi.fn();
    render(
      <ViewSwitcher
        currentView="calendar"
        onViewChange={onViewChange}
      />
    );

    const tableButton = screen.getByText('Table');
    fireEvent.click(tableButton);

    expect(onViewChange).toHaveBeenCalledWith('table');
  });

  it('switches to graph view when analytics is clicked', () => {
    const onViewChange = vi.fn();
    render(
      <ViewSwitcher
        currentView="calendar"
        onViewChange={onViewChange}
      />
    );

    const graphButton = screen.getByText('Analytics');
    fireEvent.click(graphButton);

    expect(onViewChange).toHaveBeenCalledWith('graph');
  });

  it('applies custom className', () => {
    const onViewChange = vi.fn();
    const { container } = render(
      <ViewSwitcher
        currentView="calendar"
        onViewChange={onViewChange}
        className="custom-class"
      />
    );

    expect(container.firstChild).toHaveClass('custom-class');
  });

  it('sets correct aria-pressed attribute', () => {
    const onViewChange = vi.fn();
    render(
      <ViewSwitcher
        currentView="table"
        onViewChange={onViewChange}
      />
    );

    const calendarButton = screen.getByLabelText('Switch to Calendar view');
    const tableButton = screen.getByLabelText('Switch to Table view');
    const analyticsButton = screen.getByLabelText('Switch to Analytics view');

    expect(calendarButton).toHaveAttribute('aria-pressed', 'false');
    expect(tableButton).toHaveAttribute('aria-pressed', 'true');
    expect(analyticsButton).toHaveAttribute('aria-pressed', 'false');
  });

  it('has proper accessibility labels', () => {
    const onViewChange = vi.fn();
    render(
      <ViewSwitcher
        currentView="calendar"
        onViewChange={onViewChange}
      />
    );

    expect(screen.getByLabelText('Switch to Calendar view')).toBeInTheDocument();
    expect(screen.getByLabelText('Switch to Table view')).toBeInTheDocument();
    expect(screen.getByLabelText('Switch to Analytics view')).toBeInTheDocument();
  });

  it('renders icons for each view mode', () => {
    const onViewChange = vi.fn();
    const { container } = render(
      <ViewSwitcher
        currentView="calendar"
        onViewChange={onViewChange}
      />
    );

    // Check that SVG icons are rendered
    const svgs = container.querySelectorAll('svg');
    expect(svgs.length).toBeGreaterThanOrEqual(3); // At least 3 icons
  });

  it('maintains view state across multiple clicks', () => {
    const onViewChange = vi.fn();
    const { rerender } = render(
      <ViewSwitcher
        currentView="calendar"
        onViewChange={onViewChange}
      />
    );

    fireEvent.click(screen.getByText('Table'));
    expect(onViewChange).toHaveBeenCalledWith('table');

    rerender(
      <ViewSwitcher
        currentView="table"
        onViewChange={onViewChange}
      />
    );

    fireEvent.click(screen.getByText('Analytics'));
    expect(onViewChange).toHaveBeenCalledWith('graph');
  });
});
