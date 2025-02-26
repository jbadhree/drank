import { Transaction, TransactionType } from '@/lib/types';
import { formatCurrency, formatDate } from '@/lib/utils';

interface TransactionItemProps {
  transaction: Transaction;
}

const TransactionItem = ({ transaction }: TransactionItemProps) => {
  const getTransactionIcon = (type: TransactionType) => {
    switch (type) {
      case TransactionType.Deposit:
        return (
          <div className="flex items-center justify-center w-10 h-10 rounded-full bg-green-100 text-green-500">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 11l5-5m0 0l5 5m-5-5v12" />
            </svg>
          </div>
        );
      case TransactionType.Withdrawal:
        return (
          <div className="flex items-center justify-center w-10 h-10 rounded-full bg-red-100 text-red-500">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 13l-5 5m0 0l-5-5m5 5V6" />
            </svg>
          </div>
        );
      case TransactionType.Transfer:
        return (
          <div className="flex items-center justify-center w-10 h-10 rounded-full bg-blue-100 text-blue-500">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
            </svg>
          </div>
        );
    }
  };

  const getAmountClass = (type: TransactionType) => {
    switch (type) {
      case TransactionType.Deposit:
        return 'text-green-600';
      case TransactionType.Withdrawal:
        return 'text-red-600';
      case TransactionType.Transfer:
        // If this account is the target (money coming in)
        if (transaction.targetAccountId === transaction.accountId) {
          return 'text-green-600';
        }
        // If this account is the source (money going out)
        return 'text-red-600';
    }
  };

  const getAmountPrefix = (type: TransactionType) => {
    switch (type) {
      case TransactionType.Deposit:
        return '+';
      case TransactionType.Withdrawal:
        return '-';
      case TransactionType.Transfer:
        // If this account is the target (money coming in)
        if (transaction.targetAccountId === transaction.accountId) {
          return '+';
        }
        // If this account is the source (money going out)
        return '-';
    }
  };

  return (
    <div className="flex items-center py-4 border-b border-gray-200">
      <div className="mr-4">
        {getTransactionIcon(transaction.type)}
      </div>
      <div className="flex-1">
        <p className="text-sm font-medium text-gray-900">
          {transaction.description}
        </p>
        <p className="text-xs text-gray-500">
          {formatDate(transaction.transactionDate)}
        </p>
      </div>
      <div>
        <p className={`text-sm font-semibold ${getAmountClass(transaction.type)}`}>
          {getAmountPrefix(transaction.type)}{formatCurrency(transaction.amount)}
        </p>
        <p className="text-xs text-gray-500 text-right">
          Balance: {formatCurrency(transaction.balance)}
        </p>
      </div>
    </div>
  );
};

export default TransactionItem;
