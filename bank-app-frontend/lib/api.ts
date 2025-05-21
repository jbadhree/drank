import axios from 'axios';
import { 
  LoginRequest, 
  LoginResponse, 
  User, 
  Account, 
  Transaction, 
  TransferRequest 
} from './types';

// Hardcoded default API URL that will be replaced at container startup
const API_URL = 'http://localhost:8080/api/v1';

// Log it to help with debugging
if (typeof window !== 'undefined') {
  console.log('API_URL:', API_URL);
}

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

export const getAccountByID = async (id: string): Promise<Account> => {
  const response = await api.get<Account>(`/accounts/${id}`);
  return response.data;
};

export const getAccountsByUserID = async (userId: string): Promise<Account[]> => {
  const response = await api.get<Account[]>(`/accounts/user/${userId}`);
  return response.data;
};

export const getTransactions = async (): Promise<Transaction[]> => {
  const response = await api.get<Transaction[]>('/transactions');
  return response.data;
};

export const getTransactionsByAccountID = async (accountId: string): Promise<Transaction[]> => {
  const response = await api.get<Transaction[]>(`/transactions/account/${accountId}`);
  return response.data;
};

export const transferMoney = async (transferRequest: TransferRequest): Promise<any> => {
  const response = await api.post('/transactions/transfer', transferRequest);
  return response.data;
};
