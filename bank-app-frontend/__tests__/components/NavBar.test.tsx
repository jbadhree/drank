import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import NavBar from '@/components/NavBar';
import { useAuth } from '@/lib/auth';

// Mock the useAuth hook
jest.mock('@/lib/auth', () => ({
  useAuth: jest.fn(),
}));

// Mock Next.js Link component
jest.mock('next/link', () => {
  return ({ children, href }: { children: React.ReactNode; href: string }) => {
    return <a href={href}>{children}</a>;
  };
});

describe('NavBar Component', () => {
  it('renders the bank name and link to dashboard', () => {
    // Mock a user being logged in
    (useAuth as jest.Mock).mockReturnValue({
      user: { id: 1, firstName: 'John', lastName: 'Doe', email: 'john@example.com' },
      logout: jest.fn(),
    });

    render(<NavBar />);
    
    // Check if the bank name is displayed
    const bankNameLink = screen.getByText('Drank Banking');
    expect(bankNameLink).toBeInTheDocument();
    expect(bankNameLink.closest('a')).toHaveAttribute('href', '/dashboard');
  });

  it('renders user greeting when logged in', () => {
    // Mock a user being logged in
    (useAuth as jest.Mock).mockReturnValue({
      user: { id: 1, firstName: 'John', lastName: 'Doe', email: 'john@example.com' },
      logout: jest.fn(),
    });

    render(<NavBar />);
    
    // Check if the user greeting is displayed
    expect(screen.getByText('Hello, John')).toBeInTheDocument();
    
    // Check if the logout button is displayed
    expect(screen.getByText('Logout')).toBeInTheDocument();
  });

  it('does not render user greeting or logout button when not logged in', () => {
    // Mock no user logged in
    (useAuth as jest.Mock).mockReturnValue({
      user: null,
      logout: jest.fn(),
    });

    render(<NavBar />);
    
    // Verify no greeting or logout button
    expect(screen.queryByText(/Hello,/)).not.toBeInTheDocument();
    expect(screen.queryByText('Logout')).not.toBeInTheDocument();
  });

  it('calls logout function when logout button is clicked', () => {
    // Create a mock logout function
    const mockLogout = jest.fn();
    
    // Mock a user being logged in
    (useAuth as jest.Mock).mockReturnValue({
      user: { id: 1, firstName: 'John', lastName: 'Doe', email: 'john@example.com' },
      logout: mockLogout,
    });

    render(<NavBar />);
    
    // Get the logout button and click it
    const logoutButton = screen.getByText('Logout');
    fireEvent.click(logoutButton);
    
    // Verify logout function was called
    expect(mockLogout).toHaveBeenCalledTimes(1);
  });
});
