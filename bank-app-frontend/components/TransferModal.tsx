import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Account } from '@/lib/types';
import { transferMoney } from '@/lib/api';

const transferSchema = z.object({
  fromAccountId: z.number().min(1, 'Source account is required'),
  toAccountId: z.number().min(1, 'Destination account is required'),
  amount: z.number().min(0.01, 'Amount must be greater than 0'),
  description: z.string().min(3, 'Description must be at least 3 characters')
}).refine(data => data.fromAccountId !== data.toAccountId, {
  message: "Source and destination accounts must be different",
  path: ["toAccountId"]
});

type TransferFormData = z.infer<typeof transferSchema>;

interface TransferModalProps {
  accounts: Account[];
  onClose: () => void;
  onSuccess: () => void;
}

const TransferModal = ({ accounts, onClose, onSuccess }: TransferModalProps) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const { register, handleSubmit, formState: { errors }, watch } = useForm<TransferFormData>({
    resolver: zodResolver(transferSchema),
    defaultValues: {
      fromAccountId: accounts.length > 0 ? accounts[0].id : 0,
      toAccountId: accounts.length > 1 ? accounts[1].id : 0,
      amount: 0,
      description: ''
    }
  });

  const fromAccountId = watch('fromAccountId');
  const selectedFromAccount = accounts.find(a => a.id === Number(fromAccountId));
  
  const onSubmit = async (data: TransferFormData) => {
    setIsSubmitting(true);
    setError(null);
    
    try {
      await transferMoney({
        fromAccountId: Number(data.fromAccountId),
        toAccountId: Number(data.toAccountId),
        amount: Number(data.amount),
        description: data.description
      });
      onSuccess();
    } catch (error: any) {
      setError(error.response?.data?.message || 'Transfer failed. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-lg p-6 w-full max-w-md">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold text-gray-800">Transfer Money</h2>
          <button 
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        {error && (
          <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-md">
            {error}
          </div>
        )}
        
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="mb-4">
            <label className="block text-gray-700 text-sm font-medium mb-2">
              From Account
            </label>
            <select 
              {...register('fromAccountId', { valueAsNumber: true })}
              className="input"
            >
              {accounts.map(account => (
                <option key={account.id} value={account.id}>
                  {account.accountType} - {account.accountNumber} (Balance: ${account.balance.toFixed(2)})
                </option>
              ))}
            </select>
            {errors.fromAccountId && (
              <p className="mt-1 text-sm text-red-600">{errors.fromAccountId.message}</p>
            )}
          </div>
          
          <div className="mb-4">
            <label className="block text-gray-700 text-sm font-medium mb-2">
              To Account
            </label>
            <select 
              {...register('toAccountId', { valueAsNumber: true })}
              className="input"
            >
              {accounts.map(account => (
                <option key={account.id} value={account.id}>
                  {account.accountType} - {account.accountNumber}
                </option>
              ))}
            </select>
            {errors.toAccountId && (
              <p className="mt-1 text-sm text-red-600">{errors.toAccountId.message}</p>
            )}
          </div>
          
          <div className="mb-4">
            <label className="block text-gray-700 text-sm font-medium mb-2">
              Amount
            </label>
            <input 
              type="number" 
              step="0.01"
              {...register('amount', { valueAsNumber: true })}
              className="input"
              max={selectedFromAccount?.balance || 0}
            />
            {errors.amount && (
              <p className="mt-1 text-sm text-red-600">{errors.amount.message}</p>
            )}
            {selectedFromAccount && (
              <p className="mt-1 text-xs text-gray-600">
                Available balance: ${selectedFromAccount.balance.toFixed(2)}
              </p>
            )}
          </div>
          
          <div className="mb-6">
            <label className="block text-gray-700 text-sm font-medium mb-2">
              Description
            </label>
            <input 
              type="text" 
              {...register('description')}
              className="input"
              placeholder="What's this transfer for?"
            />
            {errors.description && (
              <p className="mt-1 text-sm text-red-600">{errors.description.message}</p>
            )}
          </div>
          
          <div className="flex justify-end">
            <button
              type="button"
              onClick={onClose}
              className="btn btn-secondary mr-2"
              disabled={isSubmitting}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={isSubmitting}
            >
              {isSubmitting ? 'Processing...' : 'Transfer Money'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default TransferModal;
