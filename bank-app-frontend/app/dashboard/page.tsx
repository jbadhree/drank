'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/auth';
import { getAccountsByUserID, getTransactionsByAccountID } from '@/lib/api';
import { Account, Transaction } from '@/lib/types';
import NavBar from '@/components/NavBar';
import ApiInfoWrapper from '@/components/ApiInfoWrapper';
import AccountCard from '@/components/AccountCard';
import TransactionItem from '@/components/TransactionItem';
import TransferModal from '@/components/TransferModal';

export default function DashboardPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [selectedAccount, setSelectedAccount] = useState<Account | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [isTransferModalOpen, setIsTransferModalOpen] = useState<boolean>(false);
  
  const { user, loading: authLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
    }
  }, [user, authLoading, router]);

  useEffect(() => {
    const fetchAccounts = async () => {
      if (!user) return;
      
      try {
        setLoading(true);
        const fetchedAccounts = await getAccountsByUserID(user.id);
        setAccounts(fetchedAccounts);
        
        if (fetchedAccounts.length > 0) {
          setSelectedAccount(fetchedAccounts[0]);
        }
      } catch (error) {
        console.error('Error fetching accounts:', error);
        setError('Failed to load account information. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    if (user) {
      fetchAccounts();
    }
  }, [user]);

  useEffect(() => {
    const fetchTransactions = async () => {
      if (!selectedAccount) return;
      
      try {
        const fetchedTransactions = await getTransactionsByAccountID(selectedAccount.id);
        setTransactions(fetchedTransactions);
      } catch (error) {
        console.error('Error fetching transactions:', error);
      }
    };

    if (selectedAccount) {
      fetchTransactions();
    }
  }, [selectedAccount]);

  const handleAccountSelect = (account: Account) => {
    setSelectedAccount(account);
  };

  const handleTransferComplete = async () => {
    setIsTransferModalOpen(false);
    
    if (!user) return;
    
    try {
      // Refetch accounts to update balances
      const fetchedAccounts = await getAccountsByUserID(user.id);
      setAccounts(fetchedAccounts);
      
      // Update selected account if necessary
      if (selectedAccount) {
        const updatedSelectedAccount = fetchedAccounts.find(a => a.id === selectedAccount.id);
        if (updatedSelectedAccount) {
          setSelectedAccount(updatedSelectedAccount);
          
          // Refetch transactions for this account
          const fetchedTransactions = await getTransactionsByAccountID(updatedSelectedAccount.id);
          setTransactions(fetchedTransactions);
        }
      }
    } catch (error) {
      console.error('Error refreshing data after transfer:', error);
    }
  };

  if (authLoading || (loading && !user)) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-gray-700">Loading...</h2>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <NavBar />
      <ApiInfoWrapper />
      
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="md:flex md:justify-between md:items-center mb-6">
          <h1 className="text-2xl font-bold text-gray-900">Your Dashboard</h1>
          <button
            onClick={() => setIsTransferModalOpen(true)}
            className="mt-3 md:mt-0 w-full md:w-auto px-4 py-2 bg-primary-600 text-white font-medium rounded hover:bg-primary-700 transition-colors"
          >
            Transfer Money
          </button>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-100 text-red-700 rounded-md">
            {error}
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
          {accounts.map(account => (
            <AccountCard
              key={account.id}
              account={account}
              onClick={() => handleAccountSelect(account)}
            />
          ))}
        </div>

        {selectedAccount && (
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">
              Recent Transactions
              {selectedAccount && (
                <span className="text-base font-normal text-gray-600 ml-2">
                  ({selectedAccount.accountType} - {selectedAccount.accountNumber})
                </span>
              )}
            </h2>
            
            {transactions.length === 0 ? (
              <p className="text-gray-500 py-4">No transactions found for this account.</p>
            ) : (
              <div className="divide-y divide-gray-200">
                {transactions.map(transaction => (
                  <TransactionItem
                    key={transaction.id}
                    transaction={transaction}
                  />
                ))}
              </div>
            )}
          </div>
        )}
      </main>

      {isTransferModalOpen && accounts.length >= 2 && (
        <TransferModal
          accounts={accounts}
          onClose={() => setIsTransferModalOpen(false)}
          onSuccess={handleTransferComplete}
        />
      )}
    </div>
  );
}
