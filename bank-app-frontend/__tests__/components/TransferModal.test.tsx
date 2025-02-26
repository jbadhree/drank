import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import TransferModal from '@/components/TransferModal';
import { Account, AccountType } from '@/lib/types';
import { transferMoney } from '@/lib/api';

// Mock the API function
jest.mock('@/lib/api', () => ({
  transferMoney: jest.fn(),
}));

describe('TransferModal Component', () => {
  // Mock accounts for testing
  const mockAccounts: Account[] = [
    {
      id: 1,
      userId: 1,
      accountNumber: '12345678',
      accountType: AccountType.Checking,
      balance: 1000,
      createdAt: '2023-01-01T00:00:00.000Z',
      updatedAt: '2023-01-01T00:00:00.000Z',
    },
    {
      id: 2,
      userId: 1,
      accountNumber: '87654321',
      accountType: AccountType.Savings,
      balance: 5000,
      createdAt: '2023-01-01T00:00:00.000Z',
      updatedAt: '2023-01-01T00:00:00.000Z',
    },
  ];

  // Mock callbacks
  const mockOnClose = jest.fn();
  const mockOnSuccess = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders the transfer form with account options', () => {
    render(
      <TransferModal
        accounts={mockAccounts}
        onClose={mockOnClose}
        onSuccess={mockOnSuccess}
      />
    );
    
    // Check if title is displayed
    expect(screen.getByText('Transfer Money')).toBeInTheDocument();
    
    // Check if form fields are rendered
    expect(screen.getByText('From Account')).toBeInTheDocument();
    expect(screen.getByText('To Account')).toBeInTheDocument();
    expect(screen.getByText('Amount')).toBeInTheDocument();
    expect(screen.getByText('Description')).toBeInTheDocument();
    
    // Check if buttons are rendered
    expect(screen.getByText('Cancel')).toBeInTheDocument();
    expect(screen.getByText('Transfer Money')).toBeInTheDocument();
    
    // Check if account options are available
    const fromAccountSelect = screen.getByText('From Account').closest('label')?.nextElementSibling;
    expect(fromAccountSelect).toBeInTheDocument();
    
    // Check for available balance display
    expect(screen.getByText(/Available balance/)).toBeInTheDocument();
  });

  it('calls onClose when close button is clicked', async () => {
    render(
      <TransferModal
        accounts={mockAccounts}
        onClose={mockOnClose}
        onSuccess={mockOnSuccess}
      />
    );
    
    // Click the X button
    const closeButton = screen.getByRole('button', { name: '' }); // The X button doesn't have text
    fireEvent.click(closeButton);
    
    // Verify onClose was called
    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('calls onClose when Cancel button is clicked', async () => {
    render(
      <TransferModal
        accounts={mockAccounts}
        onClose={mockOnClose}
        onSuccess={mockOnSuccess}
      />
    );
    
    // Click the Cancel button
    const cancelButton = screen.getByText('Cancel');
    fireEvent.click(cancelButton);
    
    // Verify onClose was called
    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('validates form inputs and prevents submission with invalid data', async () => {
    render(
      <TransferModal
        accounts={mockAccounts}
        onClose={mockOnClose}
        onSuccess={mockOnSuccess}
      />
    );
    
    // Submit form without entering required data (description is empty by default)
    const submitButton = screen.getByText('Transfer Money');
    fireEvent.click(submitButton);
    
    // Wait for validation errors
    await waitFor(() => {
      expect(screen.getByText('Description must be at least 3 characters')).toBeInTheDocument();
    });
    
    // Verify API was not called
    expect(transferMoney).not.toHaveBeenCalled();
  });

  it('prevents transfer if source and destination accounts are the same', async () => {
    render(
      <TransferModal
        accounts={mockAccounts}
        onClose={mockOnClose}
        onSuccess={mockOnSuccess}
      />
    );
    
    // Set both selects to the same account
    const toAccountSelect = screen.getAllByRole('combobox')[1]; // Second select is To Account
    const amountInput = screen.getByText('Amount').closest('label')?.nextElementSibling as HTMLInputElement;
    fireEvent.change(amountInput, { target: { value: '100' } });
    
    // Submit the form
    const submitButton = screen.getByText('Transfer Money');
    fireEvent.click(submitButton);
    
    // Wait for validation error
    await waitFor(() => {
      expect(screen.getByText('Source and destination accounts must be different')).toBeInTheDocument();
    });
    
    // Verify API was not called
    expect(transferMoney).not.toHaveBeenCalled();
  });

  it('successfully submits the form with valid data', async () => {
    // Set up the mock to resolve successfully
    (transferMoney as jest.Mock).mockResolvedValue({ success: true });
    
    render(
      <TransferModal
        accounts={mockAccounts}
        onClose={mockOnClose}
        onSuccess={mockOnSuccess}
      />
    );
    
    // Fill in the form
    const fromAccountSelect = screen.getAllByRole('combobox')[0];
    const toAccountSelect = screen.getAllByRole('combobox')[1];
    userEvent.selectOptions(fromAccountSelect, ['1']);
    userEvent.selectOptions(toAccountSelect, ['2']);
    
    const amountInput = screen.getByLabelText('Amount');
    fireEvent.change(amountInput, { target: { value: '100' } });
    
    const descriptionInput = screen.getByPlaceholderText("What's this transfer for?");
    fireEvent.change(descriptionInput, { target: { value: 'Test transfer' } });
    
    // Submit the form
    const submitButton = screen.getByText('Transfer Money');
    fireEvent.click(submitButton);
    
    // Wait for the API call to be made
    await waitFor(() => {
      expect(transferMoney).toHaveBeenCalledWith({
        fromAccountId: 1,
        toAccountId: 2,
        amount: 100,
        description: 'Test transfer'
      });
    });
    
    // Verify onSuccess was called
    expect(mockOnSuccess).toHaveBeenCalledTimes(1);
  });

  it('displays an error message when the API call fails', async () => {
    // Set up the mock to reject with an error
    const errorMessage = 'Insufficient funds';
    (transferMoney as jest.Mock).mockRejectedValue({
      response: {
        data: {
          message: errorMessage
        }
      }
    });
    
    render(
      <TransferModal
        accounts={mockAccounts}
        onClose={mockOnClose}
        onSuccess={mockOnSuccess}
      />
    );
    
    // Fill in the form
    const fromAccountSelect = screen.getAllByRole('combobox')[0];
    const toAccountSelect = screen.getAllByRole('combobox')[1];
    userEvent.selectOptions(fromAccountSelect, ['1']);
    userEvent.selectOptions(toAccountSelect, ['2']);
    
    const amountInput = screen.getByLabelText('Amount');
    fireEvent.change(amountInput, { target: { value: '100' } });
    
    const descriptionInput = screen.getByPlaceholderText("What's this transfer for?");
    fireEvent.change(descriptionInput, { target: { value: 'Test transfer' } });
    
    // Submit the form
    const submitButton = screen.getByText('Transfer Money');
    fireEvent.click(submitButton);
    
    // Wait for the error message to appear
    await waitFor(() => {
      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });
    
    // Verify onSuccess was not called
    expect(mockOnSuccess).not.toHaveBeenCalled();
  });

  it('disables buttons while submitting', async () => {
    // Set up the mock to resolve after a delay to test loading state
    (transferMoney as jest.Mock).mockImplementation(() => {
      return new Promise(resolve => {
        setTimeout(() => resolve({ success: true }), 100);
      });
    });
    
    render(
      <TransferModal
        accounts={mockAccounts}
        onClose={mockOnClose}
        onSuccess={mockOnSuccess}
      />
    );
    
    // Fill in the form
    const fromAccountSelect = screen.getAllByRole('combobox')[0];
    const toAccountSelect = screen.getAllByRole('combobox')[1];
    userEvent.selectOptions(fromAccountSelect, ['1']);
    userEvent.selectOptions(toAccountSelect, ['2']);
    
    const amountInput = screen.getByLabelText('Amount');
    fireEvent.change(amountInput, { target: { value: '100' } });
    
    const descriptionInput = screen.getByPlaceholderText("What's this transfer for?");
    fireEvent.change(descriptionInput, { target: { value: 'Test transfer' } });
    
    // Submit the form
    const submitButton = screen.getByText('Transfer Money');
    fireEvent.click(submitButton);
    
    // Check if the button text changed to "Processing..."
    expect(screen.getByText('Processing...')).toBeInTheDocument();
    
    // Check if buttons are disabled
    expect(screen.getByText('Cancel')).toBeDisabled();
    expect(screen.getByText('Processing...')).toBeDisabled();
    
    // Wait for the form submission to complete
    await waitFor(() => {
      expect(mockOnSuccess).toHaveBeenCalled();
    });
  });

  it('handles API error with no response data gracefully', async () => {
    // Mock a network error with no response data
    (transferMoney as jest.Mock).mockRejectedValue(new Error('Network error'));
    
    render(
      <TransferModal
        accounts={mockAccounts}
        onClose={mockOnClose}
        onSuccess={mockOnSuccess}
      />
    );
    
    // Fill in the form
    const amountInput = screen.getByLabelText('Amount');
    fireEvent.change(amountInput, { target: { value: '100' } });
    
    const descriptionInput = screen.getByPlaceholderText("What's this transfer for?");
    fireEvent.change(descriptionInput, { target: { value: 'Test transfer' } });
    
    // Submit the form
    const submitButton = screen.getByText('Transfer Money');
    fireEvent.click(submitButton);
    
    // Wait for the default error message to appear
    await waitFor(() => {
      expect(screen.getByText('Transfer failed. Please try again.')).toBeInTheDocument();
    });
  });

  it('restricts amount to be less than or equal to available balance', () => {
    render(
      <TransferModal
        accounts={mockAccounts}
        onClose={mockOnClose}
        onSuccess={mockOnSuccess}
      />
    );
    
    // Get the amount input
    const amountInput = screen.getByLabelText('Amount');
    
    // Check if max attribute is set to the balance of the selected account
    expect(amountInput).toHaveAttribute('max', '1000');
    
    // Change the from account to the second account
    const fromAccountSelect = screen.getAllByRole('combobox')[0];
    userEvent.selectOptions(fromAccountSelect, ['2']);
    
    // Check if max attribute is updated to the new account's balance
    expect(amountInput).toHaveAttribute('max', '5000');
  });
});
