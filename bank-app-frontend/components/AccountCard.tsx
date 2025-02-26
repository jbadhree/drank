import { Account, AccountType } from '@/lib/types';
import { formatCurrency } from '@/lib/utils';

interface AccountCardProps {
  account: Account;
  onClick?: () => void;
}

const AccountCard = ({ account, onClick }: AccountCardProps) => {
  const accountTypeLabels = {
    [AccountType.Checking]: 'Checking Account',
    [AccountType.Savings]: 'Savings Account',
  };

  return (
    <div 
      className="bg-white rounded-lg shadow-md p-6 cursor-pointer transition-shadow hover:shadow-lg"
      onClick={onClick}
    >
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold text-gray-800">
          {accountTypeLabels[account.accountType]}
        </h3>
        <span className="text-xs text-gray-500">
          #{account.accountNumber}
        </span>
      </div>
      <div className="mt-2">
        <p className="text-sm text-gray-600">Current Balance</p>
        <p className="text-2xl font-bold text-primary-600">
          {formatCurrency(account.balance)}
        </p>
      </div>
    </div>
  );
};

export default AccountCard;
