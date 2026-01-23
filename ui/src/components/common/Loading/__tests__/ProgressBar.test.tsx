import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { ProgressBar } from '../ProgressBar';

describe('ProgressBar', () => {
  it('renders with default props', () => {
    render(<ProgressBar progress={50} />);
    
    const progressbar = screen.getByRole('progressbar');
    expect(progressbar).toBeInTheDocument();
    expect(progressbar).toHaveAttribute('aria-valuenow', '50');
    expect(progressbar).toHaveAttribute('aria-valuemin', '0');
    expect(progressbar).toHaveAttribute('aria-valuemax', '100');
  });

  it('displays label and percentage when provided', () => {
    render(
      <ProgressBar 
        progress={75} 
        label="Loading data..." 
        showPercentage={true} 
      />
    );
    
    expect(screen.getByText('Loading data...')).toBeInTheDocument();
    expect(screen.getByText('75%')).toBeInTheDocument();
  });

  it('clamps progress values correctly', () => {
    const { rerender } = render(<ProgressBar progress={150} />);
    let progressbar = screen.getByRole('progressbar');
    expect(progressbar).toHaveAttribute('aria-valuenow', '100');
    
    rerender(<ProgressBar progress={-10} />);
    progressbar = screen.getByRole('progressbar');
    expect(progressbar).toHaveAttribute('aria-valuenow', '0');
  });

  it('applies correct size classes', () => {
    const { rerender } = render(<ProgressBar progress={50} size="sm" />);
    let container = screen.getByRole('progressbar').parentElement;
    expect(container).toHaveClass('h-2');
    
    rerender(<ProgressBar progress={50} size="lg" />);
    container = screen.getByRole('progressbar').parentElement;
    expect(container).toHaveClass('h-4');
  });

  it('applies correct color classes', () => {
    const { rerender } = render(<ProgressBar progress={50} color="blue" />);
    let container = screen.getByRole('progressbar').parentElement;
    expect(container).toHaveClass('bg-blue-100');
    
    let progressElement = container?.querySelector('div[style*="width"]');
    expect(progressElement).toHaveClass('bg-blue-500');
    
    rerender(<ProgressBar progress={50} color="green" />);
    container = screen.getByRole('progressbar').parentElement;
    expect(container).toHaveClass('bg-green-100');
    
    progressElement = container?.querySelector('div[style*="width"]');
    expect(progressElement).toHaveClass('bg-green-500');
  });
});