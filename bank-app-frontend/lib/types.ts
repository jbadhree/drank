export interface User {
  id: number;
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
  id: number;
  userId: number;
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
  id: number;
  accountId: number;
  sourceAccountId?: number;
  targetAccountId?: number;
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
  fromAccountId: number;
  toAccountId: number;
  amount: number;
  description: string;
}
