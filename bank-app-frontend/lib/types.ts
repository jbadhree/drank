export interface User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  createdAt: string;
  updatedAt: string;
}

export enum AccountType {
  Checking = "CHECKING",
  Savings = "SAVINGS"
}

export interface Account {
  id: string;
  userId: string;
  accountNumber: string;
  accountType: AccountType;
  balance: number;
  createdAt: string;
  updatedAt: string;
}

export enum TransactionType {
  Deposit = "DEPOSIT",
  Withdrawal = "WITHDRAWAL",
  Transfer = "TRANSFER"
}

export interface Transaction {
  id: string;
  accountId: string;
  sourceAccountId?: string;
  targetAccountId?: string;
  amount: number;
  balance: number;
  type: TransactionType;
  description: string;
  transactionDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface TransferRequest {
  fromAccountId: string;
  toAccountId: string;
  amount: number;
  description: string;
}
