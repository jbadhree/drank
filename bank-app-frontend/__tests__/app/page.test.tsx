import React from 'react';
import { render, screen } from '@testing-library/react';
import HomePage from '@/app/page';
import { useRouter } from 'next/navigation';
import { isAuthenticated } from '@/lib/auth';

// Mock dependencies
jest.mock('next/navigation', () => ({
  useRouter: jest.fn(),
}));

jest.mock('@/lib/auth', () => ({
  isAuthenticated: jest.fn(),
}));

describe('HomePage Component', () => {
  const mockRouter = {
    push: jest.fn(),
  };

  beforeEach(() => {
    jest.clearAllMocks();
    (useRouter as jest.Mock).mockReturnValue(mockRouter);
  });

  it('renders the loading message correctly', () => {
    render(<HomePage />);
    
    // Check if the title is displayed
    expect(screen.getByText('Drank Banking')).toBeInTheDocument();
    
    // Check if the loading message is displayed
    expect(screen.getByText('Redirecting you to the appropriate page...')).toBeInTheDocument();
  });

  it('redirects to dashboard when user is authenticated', () => {
    // Mock user being authenticated
    (isAuthenticated as jest.Mock).mockReturnValue(true);
    
    render(<HomePage />);
    
    // Check if router.push was called with the correct path
    expect(mockRouter.push).toHaveBeenCalledWith('/dashboard');
  });

  it('redirects to login when user is not authenticated', () => {
    // Mock user not being authenticated
    (isAuthenticated as jest.Mock).mockReturnValue(false);
    
    render(<HomePage />);
    
    // Check if router.push was called with the correct path
    expect(mockRouter.push).toHaveBeenCalledWith('/login');
  });
});
