import React from 'react';
import { render, screen } from '@testing-library/react';
import Spinner from '@/components/ui/Spinner';

describe('Spinner Component', () => {
  it('renders with default props', () => {
    render(<Spinner />);
    
    // Check if the loading span is in the document
    expect(screen.getByText('Loading...')).toBeInTheDocument();
    
    // Get the spinner element
    const spinnerElement = screen.getByRole('status');
    
    // Check if it has the default classes
    expect(spinnerElement).toHaveClass('w-6 h-6'); // medium size
    expect(spinnerElement).toHaveClass('text-primary-600'); // primary color
    expect(spinnerElement).toHaveClass('animate-spin');
  });

  it('renders with small size', () => {
    render(<Spinner size="small" />);
    
    const spinnerElement = screen.getByRole('status');
    expect(spinnerElement).toHaveClass('w-4 h-4');
    expect(spinnerElement).not.toHaveClass('w-6 h-6');
    expect(spinnerElement).not.toHaveClass('w-8 h-8');
  });

  it('renders with medium size', () => {
    render(<Spinner size="medium" />);
    
    const spinnerElement = screen.getByRole('status');
    expect(spinnerElement).toHaveClass('w-6 h-6');
    expect(spinnerElement).not.toHaveClass('w-4 h-4');
    expect(spinnerElement).not.toHaveClass('w-8 h-8');
  });

  it('renders with large size', () => {
    render(<Spinner size="large" />);
    
    const spinnerElement = screen.getByRole('status');
    expect(spinnerElement).toHaveClass('w-8 h-8');
    expect(spinnerElement).not.toHaveClass('w-4 h-4');
    expect(spinnerElement).not.toHaveClass('w-6 h-6');
  });

  it('renders with primary color', () => {
    render(<Spinner color="primary" />);
    
    const spinnerElement = screen.getByRole('status');
    expect(spinnerElement).toHaveClass('text-primary-600');
    expect(spinnerElement).not.toHaveClass('text-white');
  });

  it('renders with white color', () => {
    render(<Spinner color="white" />);
    
    const spinnerElement = screen.getByRole('status');
    expect(spinnerElement).toHaveClass('text-white');
    expect(spinnerElement).not.toHaveClass('text-primary-600');
  });

  it('combines size and color correctly', () => {
    render(<Spinner size="small" color="white" />);
    
    const spinnerElement = screen.getByRole('status');
    expect(spinnerElement).toHaveClass('w-4 h-4');
    expect(spinnerElement).toHaveClass('text-white');
  });
});
