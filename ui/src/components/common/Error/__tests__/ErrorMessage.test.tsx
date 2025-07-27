import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ErrorMessage } from '../ErrorMessage';

describe('ErrorMessage', () => {
  it('renders error message with default props', () => {
    const message = 'Something went wrong';
    render(<ErrorMessage message={message} />);
    
    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText(message)).toBeInTheDocument();
  });

  it('renders with custom title', () => {
    const title = 'Custom Error';
    const message = 'Error details';
    render(<ErrorMessage title={title} message={message} />);
    
    expect(screen.getByText(title)).toBeInTheDocument();
    expect(screen.getByText(message)).toBeInTheDocument();
  });

  it('applies correct styles for different types', () => {
    const message = 'Test message';
    const { rerender } = render(<ErrorMessage message={message} type="error" />);
    
    let container = screen.getByRole('alert');
    expect(container).toHaveClass('bg-red-50', 'border-red-200');
    
    rerender(<ErrorMessage message={message} type="warning" />);
    container = screen.getByRole('alert');
    expect(container).toHaveClass('bg-yellow-50', 'border-yellow-200');
    
    rerender(<ErrorMessage message={message} type="info" />);
    container = screen.getByRole('alert');
    expect(container).toHaveClass('bg-blue-50', 'border-blue-200');
  });

  it('calls onRetry when retry button is clicked', () => {
    const onRetry = vi.fn();
    const message = 'Network error';
    
    render(<ErrorMessage message={message} onRetry={onRetry} />);
    
    const retryButton = screen.getByText('Try Again');
    fireEvent.click(retryButton);
    
    expect(onRetry).toHaveBeenCalledTimes(1);
  });

  it('calls onDismiss when dismiss button is clicked', () => {
    const onDismiss = vi.fn();
    const message = 'Dismissible error';
    
    render(<ErrorMessage message={message} onDismiss={onDismiss} />);
    
    const dismissButton = screen.getByLabelText('Dismiss');
    fireEvent.click(dismissButton);
    
    expect(onDismiss).toHaveBeenCalledTimes(1);
  });

  it('toggles detail text when showDetails is true', () => {
    const message = 'Error message';
    const details = 'Detailed error information';
    
    render(
      <ErrorMessage 
        message={message} 
        showDetails={true} 
        details={details} 
      />
    );
    
    // Details should be hidden initially
    expect(screen.queryByText(details)).not.toBeInTheDocument();
    
    // Click show details button
    const showDetailsButton = screen.getByText('Show details');
    fireEvent.click(showDetailsButton);
    
    // Details should now be visible
    expect(screen.getByText(details)).toBeInTheDocument();
    
    // Button text should change
    expect(screen.getByText('Hide details')).toBeInTheDocument();
    
    // Click hide details button
    const hideDetailsButton = screen.getByText('Hide details');
    fireEvent.click(hideDetailsButton);
    
    // Details should be hidden again
    expect(screen.queryByText(details)).not.toBeInTheDocument();
  });

  it('does not render details toggle when showDetails is false', () => {
    const message = 'Error message';
    const details = 'Hidden details';
    
    render(
      <ErrorMessage 
        message={message} 
        showDetails={false} 
        details={details} 
      />
    );
    
    expect(screen.queryByText('Show details')).not.toBeInTheDocument();
    expect(screen.queryByText(details)).not.toBeInTheDocument();
  });

  it('renders appropriate icons for different types', () => {
    const message = 'Test message';
    const { rerender } = render(<ErrorMessage message={message} type="error" />);
    
    expect(screen.getByText('❌')).toBeInTheDocument();
    
    rerender(<ErrorMessage message={message} type="warning" />);
    expect(screen.getByText('⚠️')).toBeInTheDocument();
    
    rerender(<ErrorMessage message={message} type="info" />);
    expect(screen.getByText('ℹ️')).toBeInTheDocument();
  });
});