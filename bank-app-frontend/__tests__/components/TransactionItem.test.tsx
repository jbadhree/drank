import React from 'react';
import { render, screen } from '@testing-library/react';
import TransactionItem from '@/components/TransactionItem';
import { Transaction, TransactionType } from '@/lib/types';
import { formatCurrency, formatDate } from '@/lib/utils';

// Mock the utility functions
jest.mock('@/lib/utils', () => ({
  formatCurrency: jest.fn((amount) => `$${amount.toFixed(2)}`),
  formatDate: jest.fn((date) => '01/01/2023, 10:00 AM'),
}));

describe('TransactionItem Component', () => {
  // Common transaction props
  const baseTransaction: Partial<Transaction> = {
    id: 1,
    accountId: 1,
    amount: 100.50,
    balance: 1250.75,
    description: 'Test Transaction',
    transactionDate: '2023-01-01T10:00:00.000Z',
    createdAt: '2023-01-01T10:00:00.000Z',
    updatedAt: '2023-01-01T10:00:00.000Z',
  };

  it('renders deposit transaction correctly', () => {
    const depositTransaction: Transaction = {
      ...baseTransaction as Transaction,
      type: TransactionType.Deposit,
    };

    render(<TransactionItem transaction={depositTransaction} />);
    
    // Check if description is displayed
    expect(screen.getByText('Test Transaction')).toBeInTheDocument();
    
    // Check if formatted date is displayed
    expect(screen.getByText('01/01/2023, 10:00 AM')).toBeInTheDocument();
    
    // Check if amount has the correct prefix and class
    const amountElement = screen.getByText('+$100.50');
    expect(amountElement).toBeInTheDocument();
    expect(amountElement).toHaveClass('text-green-600');
    
    // Check if the balance is displayed
    expect(screen.getByText('Balance: $1250.75')).toBeInTheDocument();
  });

  it('renders withdrawal transaction correctly', () => {
    const withdrawalTransaction: Transaction = {
      ...baseTransaction as Transaction,
      type: TransactionType.Withdrawal,
    };

    render(<TransactionItem transaction={withdrawalTransaction} />);
    
    // Check if amount has the correct prefix and class
    const amountElement = screen.getByText('-$100.50');
    expect(amountElement).toBeInTheDocument();
    expect(amountElement).toHaveClass('text-red-600');
  });

  it('renders transfer transaction (outgoing) correctly', () => {
    const transferOutTransaction: Transaction = {
      ...baseTransaction as Transaction,
      type: TransactionType.Transfer,
      sourceAccountId: 1, // Same as accountId (source of the transfer)
      targetAccountId: 2,
    };

    render(<TransactionItem transaction={transferOutTransaction} />);
    
    // Check if amount has the correct prefix and class
    const amountElement = screen.getByText('-$100.50');
    expect(amountElement).toBeInTheDocument();
    expect(amountElement).toHaveClass('text-red-600');
  });

  it('renders transfer transaction (incoming) correctly', () => {
    const transferInTransaction: Transaction = {
      ...baseTransaction as Transaction,
      type: TransactionType.Transfer,
      sourceAccountId: 2,
      targetAccountId: 1, // Same as accountId (target of the transfer)
    };

    render(<TransactionItem transaction={transferInTransaction} />);
    
    // Check if amount has the correct prefix and class
    const amountElement = screen.getByText('+$100.50');
    expect(amountElement).toBeInTheDocument();
    expect(amountElement).toHaveClass('text-green-600');
  });

  it('renders the correct transaction icon for each type', () => {
    // Test deposit icon
    const depositTransaction: Transaction = {
      ...baseTransaction as Transaction,
      type: TransactionType.Deposit,
    };
    const { rerender } = render(<TransactionItem transaction={depositTransaction} />);
    
    // For deposit, we expect an up arrow (path with d attribute containing "M7 11l5-5m0 0l5 5m-5-5v12")
    const depositSVG = document.querySelector('path[d="M7 11l5-5m0 0l5 5m-5-5v12"]');
    expect(depositSVG).toBeInTheDocument();
    
    // Test withdrawal icon
    const withdrawalTransaction: Transaction = {
      ...baseTransaction as Transaction,
      type: TransactionType.Withdrawal,
    };
    rerender(<TransactionItem transaction={withdrawalTransaction} />);
    
    // For withdrawal, we expect a down arrow (path with d attribute containing "M17 13l-5 5m0 0l-5-5m5 5V6")
    const withdrawalSVG = document.querySelector('path[d="M17 13l-5 5m0 0l-5-5m5 5V6"]');
    expect(withdrawalSVG).toBeInTheDocument();
    
    // Test transfer icon
    const transferTransaction: Transaction = {
      ...baseTransaction as Transaction,
      type: TransactionType.Transfer,
      sourceAccountId: 1,
      targetAccountId: 2,
    };
    rerender(<TransactionItem transaction={transferTransaction} />);
    
    // For transfer, we expect a double arrow (path with d attribute containing "M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4")
    const transferSVG = document.querySelector('path[d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"]');
    expect(transferSVG).toBeInTheDocument();
  });
});
