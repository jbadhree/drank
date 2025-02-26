import { mockUser, mockAuthContextValue } from './authContext';

describe('Mocks', () => {
  it('exports mock data correctly', () => {
    expect(mockUser).toBeDefined();
    expect(mockAuthContextValue).toBeDefined();
  });
});
