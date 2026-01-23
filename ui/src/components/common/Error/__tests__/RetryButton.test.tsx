import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { RetryButton } from '../RetryButton';

describe('RetryButton', () => {
  it('renders with default props', () => {
    const mockRetry = vi.fn();
    render(<RetryButton onRetry={mockRetry} />);
    
    const button = screen.getByRole('button');
    expect(button).toBeInTheDocument();
    expect(button).toHaveTextContent('Retry');
    expect(button).not.toBeDisabled();
  });

  it('calls onRetry when clicked', () => {
    const mockRetry = vi.fn();
    render(<RetryButton onRetry={mockRetry} />);
    
    const button = screen.getByRole('button');
    fireEvent.click(button);
    
    expect(mockRetry).toHaveBeenCalledTimes(1);
  });

  it('shows retrying state correctly', () => {
    const mockRetry = vi.fn();
    render(<RetryButton onRetry={mockRetry} isRetrying={true} />);
    
    const button = screen.getByRole('button');
    expect(button).toHaveTextContent('Retrying...');
    expect(button).toBeDisabled();
    expect(button).toHaveAttribute('aria-label', 'Retrying...');
  });

  it('shows countdown correctly', async () => {
    const mockRetry = vi.fn();
    render(<RetryButton onRetry={mockRetry} countdown={3} />);
    
    const button = screen.getByRole('button');
    expect(button).toHaveTextContent('Retry in 3s');
    expect(button).toBeDisabled();
    
    // Wait for countdown to tick down
    await waitFor(
      () => {
        expect(screen.getByText('Retry in 2s')).toBeInTheDocument();
      },
      { timeout: 1200 }
    );
  });

  it('becomes enabled after countdown finishes', async () => {
    const mockRetry = vi.fn();
    render(<RetryButton onRetry={mockRetry} countdown={1} />);
    
    const button = screen.getByRole('button');
    expect(button).toBeDisabled();
    
    // Wait for countdown to finish
    await waitFor(
      () => {
        expect(button).not.toBeDisabled();
        expect(button).toHaveTextContent('Retry');
      },
      { timeout: 1200 }
    );
  });

  it('applies variant classes correctly', () => {
    const mockRetry = vi.fn();
    const { rerender } = render(<RetryButton onRetry={mockRetry} variant="primary" />);
    let button = screen.getByRole('button');
    expect(button).toHaveClass('bg-orange-500', 'text-white');
    
    rerender(<RetryButton onRetry={mockRetry} variant="secondary" />);
    button = screen.getByRole('button');
    expect(button).toHaveClass('bg-white', 'text-orange-600', 'border-orange-500');
    
    rerender(<RetryButton onRetry={mockRetry} variant="ghost" />);
    button = screen.getByRole('button');
    expect(button).toHaveClass('bg-transparent', 'text-orange-600');
  });

  it('renders custom children', () => {
    const mockRetry = vi.fn();
    render(<RetryButton onRetry={mockRetry}>Try Again</RetryButton>);
    
    const button = screen.getByRole('button');
    expect(button).toHaveTextContent('Try Again');
  });
});