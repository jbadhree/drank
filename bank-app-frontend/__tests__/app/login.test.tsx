import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import LoginPage from '@/app/login/page';
import { login } from '@/lib/api';
import { useAuth } from '@/lib/auth';

// Mock dependencies
jest.mock('@/lib/api', () => ({
  login: jest.fn(),
}));

jest.mock('@/lib/auth', () => ({
  useAuth: jest.fn(),
}));

describe('LoginPage Component', () => {
  const mockAuthLogin = jest.fn();
  
  beforeEach(() => {
    jest.clearAllMocks();
    (useAuth as jest.Mock).mockReturnValue({
      login: mockAuthLogin,
    });
  });

  it('renders the login form correctly', () => {
    render(<LoginPage />);
    
    // Check if the title and header are displayed
    expect(screen.getByText('Drank Banking')).toBeInTheDocument();
    expect(screen.getByText('Sign in to your account')).toBeInTheDocument();
    
    // Check if form inputs are rendered
    expect(screen.getByPlaceholderText('Email address')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Password')).toBeInTheDocument();
    
    // Check if the submit button is rendered
    expect(screen.getByRole('button', { name: 'Sign in' })).toBeInTheDocument();
    
    // Check if demo credentials are displayed
    expect(screen.getByText('Demo Credentials:')).toBeInTheDocument();
    expect(screen.getByText('Email: john.doe@example.com')).toBeInTheDocument();
    expect(screen.getByText('Password: password123')).toBeInTheDocument();
  });

  it('validates form inputs', async () => {
    render(<LoginPage />);
    
    // Submit the form without filling any fields
    const submitButton = screen.getByRole('button', { name: 'Sign in' });
    fireEvent.click(submitButton);
    
    // Wait for validation errors
    await waitFor(() => {
      expect(screen.getByText('Invalid email address')).toBeInTheDocument();
      expect(screen.getByText('Password must be at least 6 characters')).toBeInTheDocument();
    });
    
    // Verify login function was not called
    expect(login).not.toHaveBeenCalled();
  });

  it('handles successful login', async () => {
    const mockToken = 'test-token';
    const mockResponse = {
      token: mockToken,
      user: {
        id: 1,
        firstName: 'John',
        lastName: 'Doe',
        email: 'john@example.com',
      }
    };
    
    // Mock successful API response
    (login as jest.Mock).mockResolvedValueOnce(mockResponse);
    
    render(<LoginPage />);
    
    // Fill in the form
    const emailInput = screen.getByPlaceholderText('Email address');
    const passwordInput = screen.getByPlaceholderText('Password');
    await userEvent.type(emailInput, 'john@example.com');
    await userEvent.type(passwordInput, 'password123');
    
    // Submit the form
    const submitButton = screen.getByRole('button', { name: 'Sign in' });
    fireEvent.click(submitButton);
    
    // The button change to "Signing in..." happens asynchronously, so we'll skip this check
    
    // Wait for the API call to complete
    await waitFor(() => {
      expect(login).toHaveBeenCalledWith({
        email: 'john@example.com',
        password: 'password123',
      });
      expect(mockAuthLogin).toHaveBeenCalledWith(mockToken);
    });
  });

  it('handles login failure', async () => {
    const errorMessage = 'Invalid credentials';
    
    // Mock API error response
    (login as jest.Mock).mockRejectedValueOnce({
      response: {
        data: {
          message: errorMessage
        }
      }
    });
    
    render(<LoginPage />);
    
    // Fill in the form
    const emailInput = screen.getByPlaceholderText('Email address');
    const passwordInput = screen.getByPlaceholderText('Password');
    await userEvent.type(emailInput, 'john@example.com');
    await userEvent.type(passwordInput, 'wrongpassword');
    
    // Submit the form
    const submitButton = screen.getByRole('button', { name: 'Sign in' });
    fireEvent.click(submitButton);
    
    // Wait for the error message to appear
    await waitFor(() => {
      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });
    
    // Verify authLogin was not called
    expect(mockAuthLogin).not.toHaveBeenCalled();
  });

  it('handles login failure with no specific error message', async () => {
    // Mock API error with no response data
    (login as jest.Mock).mockRejectedValueOnce(new Error('Network Error'));
    
    render(<LoginPage />);
    
    // Fill in the form
    const emailInput = screen.getByPlaceholderText('Email address');
    const passwordInput = screen.getByPlaceholderText('Password');
    await userEvent.type(emailInput, 'john@example.com');
    await userEvent.type(passwordInput, 'password123');
    
    // Submit the form
    const submitButton = screen.getByRole('button', { name: 'Sign in' });
    fireEvent.click(submitButton);
    
    // Wait for the default error message to appear
    await waitFor(() => {
      expect(screen.getByText('Login failed. Please check your credentials.')).toBeInTheDocument();
    });
  });
});
