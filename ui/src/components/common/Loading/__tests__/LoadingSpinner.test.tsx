import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import * as matchers from '@testing-library/jest-dom/matchers';
import { LoadingSpinner } from '../LoadingSpinner';

expect.extend(matchers);

describe('LoadingSpinner', () => {
  it('renders with default props', () => {
    render(<LoadingSpinner />);
    
    const spinner = screen.getByRole('status');
    expect(spinner).toBeInTheDocument();
    expect(spinner).toHaveAttribute('aria-label', 'Loading');
  });

  it('renders with custom message', () => {
    const message = 'Loading data...';
    render(<LoadingSpinner message={message} />);
    
    expect(screen.getByText(message)).toBeInTheDocument();
  });

  it('applies correct size classes', () => {
    const { rerender } = render(<LoadingSpinner size="sm" />);
    expect(screen.getByRole('status')).toHaveClass('h-4', 'w-4');

    rerender(<LoadingSpinner size="lg" />);
    expect(screen.getByRole('status')).toHaveClass('h-8', 'w-8');

    rerender(<LoadingSpinner size="xl" />);
    expect(screen.getByRole('status')).toHaveClass('h-12', 'w-12');
  });

  it('applies correct color classes', () => {
    const { rerender } = render(<LoadingSpinner color="blue" />);
    expect(screen.getByRole('status')).toHaveClass('border-blue-500');

    rerender(<LoadingSpinner color="green" />);
    expect(screen.getByRole('status')).toHaveClass('border-green-500');
  });

  it('renders fullscreen overlay when fullScreen is true', () => {
    const { container } = render(<LoadingSpinner fullScreen />);
    
    const overlay = container.querySelector('.fixed.inset-0.z-50');
    expect(overlay).toBeInTheDocument();
  });

  it('has spinning animation', () => {
    render(<LoadingSpinner />);
    
    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('animate-spin');
  });
});