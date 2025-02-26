import { mockUser, mockAuthContextValue, mockAuthContextValueLoggedOut, mockAuthContextValueLoading, mockAuthContextValueError } from './authContext';

describe('Auth Context Mocks', () => {
  it('exports mock user data', () => {
    expect(mockUser).toBeDefined();
    expect(mockUser.firstName).toBe('John');
    expect(mockUser.lastName).toBe('Doe');
  });

  it('exports mock auth context values', () => {
    expect(mockAuthContextValue).toBeDefined();
    expect(mockAuthContextValue.user).toEqual(mockUser);
    expect(mockAuthContextValue.login).toBeDefined();
    expect(mockAuthContextValue.logout).toBeDefined();
    expect(mockAuthContextValue.isLoading).toBe(false);
    expect(mockAuthContextValue.error).toBeNull();
  });

  it('exports logged out auth context value', () => {
    expect(mockAuthContextValueLoggedOut).toBeDefined();
    expect(mockAuthContextValueLoggedOut.user).toBeNull();
  });

  it('exports loading auth context value', () => {
    expect(mockAuthContextValueLoading).toBeDefined();
    expect(mockAuthContextValueLoading.isLoading).toBe(true);
  });

  it('exports error auth context value', () => {
    expect(mockAuthContextValueError).toBeDefined();
    expect(mockAuthContextValueError.error).toBe('Invalid credentials');
  });
});
