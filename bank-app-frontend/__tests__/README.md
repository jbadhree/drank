# Frontend Test Suite for Drank Banking Application

This directory contains tests for the frontend components of the Drank Banking Application.

## Running Tests

You can run the tests using the following npm scripts:

```bash
# Run all tests
npm test

# Run tests in watch mode (tests will re-run when files change)
npm run test:watch

# Run tests with coverage report
npm run test:coverage
```

## Test Structure

The tests are organized by component types:

- `__tests__/components/`: Tests for React components
- `__tests__/lib/`: Tests for utility functions and hooks

## Mock Data

Reusable mock data for tests is available in the `__tests__/mocks/` directory.

## Writing New Tests

When writing new tests:

1. Place component tests in `__tests__/components/`
2. Place utility function tests in `__tests__/lib/`
3. Use the mock data from `__tests__/mocks/` where possible
4. Follow the existing patterns for consistent testing approach

## Testing Stack

- Jest: Test runner
- React Testing Library: Component testing
- User Event: Simulating user interactions
