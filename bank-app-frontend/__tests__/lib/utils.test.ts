import { formatCurrency, formatDate, formatAccountNumber, paginateArray } from '@/lib/utils';

describe('Utils Functions', () => {
  describe('formatCurrency', () => {
    it('formats positive numbers with $ sign and 2 decimal places', () => {
      expect(formatCurrency(1234.56)).toBe('$1,234.56');
      expect(formatCurrency(0)).toBe('$0.00');
      expect(formatCurrency(0.1)).toBe('$0.10');
      expect(formatCurrency(1000000)).toBe('$1,000,000.00');
    });

    it('formats negative numbers correctly', () => {
      expect(formatCurrency(-1234.56)).toBe('-$1,234.56');
      expect(formatCurrency(-0.5)).toBe('-$0.50');
    });

    it('handles floating point precision issues', () => {
      expect(formatCurrency(0.1 + 0.2)).toBe('$0.30'); // Not $0.30000000000000004
    });
  });

  describe('formatDate', () => {
    it('formats date strings in a user-friendly format', () => {
      // Mock the DateTimeFormat to ensure consistent output across environments
      const originalDateTimeFormat = Intl.DateTimeFormat;
      
      // Mock implementation
      global.Intl.DateTimeFormat = jest.fn().mockImplementation(() => ({
        format: () => 'Jan 1, 2023, 10:00 AM'
      }));

      expect(formatDate('2023-01-01T10:00:00Z')).toBe('Jan 1, 2023, 10:00 AM');
      
      // Restore original implementation
      global.Intl.DateTimeFormat = originalDateTimeFormat;
    });

    it('handles invalid date strings gracefully', () => {
      // Providing an invalid date should still return something and not throw
      expect(() => formatDate('invalid-date')).not.toThrow();
    });
  });

  describe('formatAccountNumber', () => {
    it('masks account numbers, showing only the last 4 digits', () => {
      expect(formatAccountNumber('1234567890123456')).toBe('xxxx-xxxx-3456');
      expect(formatAccountNumber('9876543210')).toBe('xxxx-xxxx-3210');
    });

    it('returns original string if less than 4 characters', () => {
      expect(formatAccountNumber('123')).toBe('123');
      expect(formatAccountNumber('')).toBe('');
    });

    it('handles null or undefined gracefully', () => {
      // @ts-ignore - Testing runtime behavior with invalid input
      expect(formatAccountNumber(null)).toBe(null);
      // @ts-ignore - Testing runtime behavior with invalid input
      expect(formatAccountNumber(undefined)).toBe(undefined);
    });
  });

  describe('paginateArray', () => {
    const testArray = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];

    it('returns correct slice of array for given page and pageSize', () => {
      expect(paginateArray(testArray, 1, 3)).toEqual([1, 2, 3]);
      expect(paginateArray(testArray, 2, 3)).toEqual([4, 5, 6]);
      expect(paginateArray(testArray, 3, 3)).toEqual([7, 8, 9]);
      expect(paginateArray(testArray, 4, 3)).toEqual([10]);
    });

    it('returns empty array for out-of-bounds pages', () => {
      expect(paginateArray(testArray, 5, 3)).toEqual([]);
      expect(paginateArray(testArray, 0, 3)).toEqual([]);
    });

    it('handles empty arrays gracefully', () => {
      expect(paginateArray([], 1, 10)).toEqual([]);
    });

    it('returns all items if pageSize is larger than array length', () => {
      expect(paginateArray(testArray, 1, 20)).toEqual(testArray);
    });
  });
});
