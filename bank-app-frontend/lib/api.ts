import axios from 'axios';
import { 
  LoginRequest, 
  LoginResponse, 
  User, 
  Account, 
  Transaction, 
  TransferRequest 
} from './types';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_URL,
});

// Intercept requests to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export const login = async (credentials: LoginRequest): Promise<LoginResponse> => {
  const response = await api.post<LoginResponse>('/auth/login', credentials);
  return response.data;
};

export const getCurrentUser = async (): Promise<User> => {
  const response = await api.get<User>('/users/me');
  return response.data;
};

export const getAccounts = async (): Promise<Account[]> => {
  const response = await api.get<Account[]>('/accounts');
  return response.data;
};

export const getAccountByID = async (id: number): Promise<Account> => {
  const response = await api.get<Account>(`/accounts/${id}`);
  return response.data;
};

export const getAccountsByUserID = async (userId: number): Promise<Account[]> => {
  const response = await api.get<Account[]>(`/accounts/user/${userId}`);
  return response.data;
};

export const getTransactions = async (): Promise<Transaction[]> => {
  const response = await api.get<Transaction[]>('/transactions');
  return response.data;
};

export const getTransactionsByAccountID = async (accountId: number): Promise<Transaction[]> => {
  const response = await api.get<Transaction[]>(`/transactions/account/${accountId}`);
  return response.data;
};

export const transferMoney = async (transferRequest: TransferRequest): Promise<any> => {
  const response = await api.post('/transactions/transfer', transferRequest);
  return response.data;
};
