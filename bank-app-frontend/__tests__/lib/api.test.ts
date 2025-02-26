import axios from 'axios';
import {
  login,
  getCurrentUser,
  getAccounts,
  getAccountByID,
  getAccountsByUserID,
  getTransactions,
  getTransactionsByAccountID,
  transferMoney
} from '@/lib/api';
import { 
  AccountType, 
  TransactionType, 
  LoginRequest, 
  TransferRequest
} from '@/lib/types';

// Mock localStorage
const localStorageMock = (() => {
  let store = {};
  return {
    getItem: jest.fn((key) => store[key] || null),
    setItem: jest.fn((key, value) => { store[key] = value?.toString(); }),
    removeItem: jest.fn((key) => { delete store[key]; }),
    clear: jest.fn(() => { store = {}; }),
  };
})();

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
});

// Setup simplified tests for the API
describe('API Functions', () => {
  // Mock axios implementation
  const mockGet = jest.fn();
  const mockPost = jest.fn();
  
  // Setup before tests
  beforeAll(() => {
    // Mock axios.create
    jest.spyOn(axios, 'create').mockImplementation(() => ({
      interceptors: {
        request: {
          use: jest.fn((callback) => {
            // Simple test for the interceptor
            const config = { headers: {} };
            localStorageMock.setItem('token', 'test-token');
            const result = callback(config);
            expect(result.headers.Authorization).toBe('Bearer test-token');
          })
        }
      },
      get: mockGet,
      post: mockPost
    }));
  });
  
  // Cleanup after tests
  afterAll(() => {
    jest.restoreAllMocks();
  });
  
  // Reset mocks before each test
  beforeEach(() => {
    jest.clearAllMocks();
    localStorageMock.clear();
  });
  
  // Sample test data
  const mockUser = {
    id: 1,
    firstName: 'John',
    lastName: 'Doe',
    email: 'john@example.com',
  };
  
  const mockAccount = {
    id: 1,
    userId: 1,
    accountNumber: '12345678',
    accountType: AccountType.Checking,
    balance: 1000
  };
  
  const mockAccounts = [mockAccount];
  
  const mockTransaction = {
    id: 1,
    accountId: 1,
    amount: 100,
    type: TransactionType.Deposit,
    description: 'Test transaction'
  };
  
  const mockTransactions = [mockTransaction];
  
  describe('API Endpoints', () => {
    it('login sends credentials to the correct endpoint', async () => {
      const mockResponse = { data: { token: 'token', user: mockUser } };
      mockPost.mockResolvedValueOnce(mockResponse);
      
      const credentials = { email: 'test@example.com', password: 'password' };
      const result = await login(credentials);
      
      expect(mockPost).toHaveBeenCalledWith('/auth/login', credentials);
      expect(result).toEqual(mockResponse.data);
    });
    
    it('getCurrentUser fetches from the correct endpoint', async () => {
      const mockResponse = { data: mockUser };
      mockGet.mockResolvedValueOnce(mockResponse);
      
      const result = await getCurrentUser();
      
      expect(mockGet).toHaveBeenCalledWith('/users/me');
      expect(result).toEqual(mockUser);
    });
    
    it('getAccounts fetches from the correct endpoint', async () => {
      const mockResponse = { data: mockAccounts };
      mockGet.mockResolvedValueOnce(mockResponse);
      
      const result = await getAccounts();
      
      expect(mockGet).toHaveBeenCalledWith('/accounts');
      expect(result).toEqual(mockAccounts);
    });
    
    it('getAccountByID fetches from the correct endpoint', async () => {
      const mockResponse = { data: mockAccount };
      mockGet.mockResolvedValueOnce(mockResponse);
      
      const result = await getAccountByID(1);
      
      expect(mockGet).toHaveBeenCalledWith('/accounts/1');
      expect(result).toEqual(mockAccount);
    });
    
    it('getTransactions fetches from the correct endpoint', async () => {
      const mockResponse = { data: mockTransactions };
      mockGet.mockResolvedValueOnce(mockResponse);
      
      const result = await getTransactions();
      
      expect(mockGet).toHaveBeenCalledWith('/transactions');
      expect(result).toEqual(mockTransactions);
    });
    
    it('transferMoney sends data to the correct endpoint', async () => {
      const transferData = {
        fromAccountId: 1,
        toAccountId: 2,
        amount: 100,
        description: 'Transfer'
      };
      
      const mockResponse = { data: { success: true } };
      mockPost.mockResolvedValueOnce(mockResponse);
      
      const result = await transferMoney(transferData);
      
      expect(mockPost).toHaveBeenCalledWith('/transactions/transfer', transferData);
      expect(result).toEqual({ success: true });
    });
  });
});
