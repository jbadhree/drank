import { User } from '@/lib/types';

export const mockUser: User = {
  id: 1,
  firstName: 'John',
  lastName: 'Doe',
  email: 'john.doe@example.com',
  createdAt: '2023-01-01T00:00:00.000Z',
  updatedAt: '2023-01-01T00:00:00.000Z'
};

export const mockAuthContextValue = {
  user: mockUser,
  login: jest.fn(),
  logout: jest.fn(),
  isLoading: false,
  error: null
};

export const mockAuthContextValueLoggedOut = {
  user: null,
  login: jest.fn(),
  logout: jest.fn(),
  isLoading: false,
  error: null
};

export const mockAuthContextValueLoading = {
  user: null,
  login: jest.fn(),
  logout: jest.fn(),
  isLoading: true,
  error: null
};

export const mockAuthContextValueError = {
  user: null,
  login: jest.fn(),
  logout: jest.fn(),
  isLoading: false,
  error: 'Invalid credentials'
};
