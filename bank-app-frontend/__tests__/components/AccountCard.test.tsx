import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import AccountCard from '@/components/AccountCard';
import { Account, AccountType } from '@/lib/types';
import { formatCurrency } from '@/lib/utils';

// Mock the formatCurrency function
jest.mock('@/lib/utils', () => ({
  formatCurrency: jest.fn((amount) => `$${amount.toFixed(2)}`),
}));

describe('AccountCard Component', () => {
  const mockCheckingAccount: Account = {
    id: 1,
    userId: 1,
    accountNumber: '12345678',
    accountType: AccountType.Checking,
    balance: 1250.75,
    createdAt: '2023-01-01T00:00:00.000Z',
    updatedAt: '2023-01-01T00:00:00.000Z',
  };

  const mockSavingsAccount: Account = {
    id: 2,
    userId: 1,
    accountNumber: '87654321',
    accountType: AccountType.Savings,
    balance: 5000.50,
    createdAt: '2023-01-01T00:00:00.000Z',
    updatedAt: '2023-01-01T00:00:00.000Z',
  };

  it('renders checking account correctly', () => {
    render(<AccountCard account={mockCheckingAccount} />);
    
    // Check if the account type is displayed
    expect(screen.getByText('Checking Account')).toBeInTheDocument();
    
    // Check if the account number is displayed
    expect(screen.getByText(`#${mockCheckingAccount.accountNumber}`)).toBeInTheDocument();
    
    // Check if the "Current Balance" label is displayed
    expect(screen.getByText('Current Balance')).toBeInTheDocument();
    
    // Check if formatCurrency was called with the correct amount
    expect(formatCurrency).toHaveBeenCalledWith(mockCheckingAccount.balance);
    
    // Check if the formatted balance is displayed
    expect(screen.getByText('$1250.75')).toBeInTheDocument();
  });

  it('renders savings account correctly', () => {
    render(<AccountCard account={mockSavingsAccount} />);
    
    // Check if the account type is displayed
    expect(screen.getByText('Savings Account')).toBeInTheDocument();
    
    // Check if the account number is displayed
    expect(screen.getByText(`#${mockSavingsAccount.accountNumber}`)).toBeInTheDocument();
    
    // Check if formatCurrency was called with the correct amount
    expect(formatCurrency).toHaveBeenCalledWith(mockSavingsAccount.balance);
  });

  it('calls onClick handler when clicked', () => {
    const handleClick = jest.fn();
    render(<AccountCard account={mockCheckingAccount} onClick={handleClick} />);
    
    // Find the card element and click it
    const cardElement = screen.getByText('Checking Account').closest('div');
    fireEvent.click(cardElement!);
    
    // Check if the onClick handler was called
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('does not error when onClick is not provided', () => {
    // This test ensures that the component doesn't throw an error when onClick is not provided
    expect(() => {
      render(<AccountCard account={mockCheckingAccount} />);
      const cardElement = screen.getByText('Checking Account').closest('div');
      fireEvent.click(cardElement!);
    }).not.toThrow();
  });
});
