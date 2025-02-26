# Drank Banking Application Tests

This directory contains both unit tests and functional tests for the Drank Banking Application.

## Test Structure

- `unit/`: Contains unit tests for individual components (models, services, etc.)
- `functional/`: Contains functional API tests that test the entire API endpoints
- `integration/`: Contains integration tests (reserved for future use)

## Running Tests

### Prerequisites

Before running the tests, ensure you have:

1. Go installed (version 1.18 or higher recommended)
2. PostgreSQL installed and running
3. A test database named `drank_test` created

### Environment Setup

The tests use the same configuration as your main application but will append `_test` to the database name. Make sure your `.env` file has the correct database credentials.

### Running Unit Tests

```bash
cd bank-app-backend
go test ./tests/unit/... -v
```

### Running Functional Tests

For functional tests, you need a test database running on port 5435:

1. Start the test database container:
```bash
cd databse_setup
docker-compose -f postgres_docker_compose_test.yml up -d
```

2. Run the tests:
```bash
cd bank-app-backend
go test ./tests/functional/... -v
```

3. When you're done, you can stop the test database:
```bash
cd databse_setup
docker-compose -f postgres_docker_compose_test.yml down
```

### Running All Tests

To run all tests:

```bash
cd bank-app-backend
go test ./tests/... -v
```

## Test Coverage

To run tests with coverage:

```bash
cd bank-app-backend
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Data

The functional tests will automatically:
1. Set up a clean test database environment
2. Create test users, accounts, and transactions as needed
3. Clean up after tests are complete

## Writing New Tests

### Unit Tests

Unit tests should focus on testing individual components in isolation. Use mocks for dependencies.

Example:
```go
func TestUserService_CreateUser(t *testing.T) {
    // Setup mock repositories
    // Test the service function
    // Assert expected results
}
```

### Functional Tests

Functional tests should test the API endpoints from an external perspective. They should verify that:
1. Proper HTTP status codes are returned
2. Response bodies contain expected data
3. Authentication and authorization are enforced
4. Database changes are actually persisted

Example:
```go
func TestTransferAPI(t *testing.T) {
    // Set up test data
    // Make API request
    // Verify response
    // Verify database state
}
```

## Best Practices

1. Always clean up test data to avoid interference between tests
2. Use meaningful test names that describe what's being tested
3. Follow the Arrange-Act-Assert pattern
4. Test both success and failure cases
5. Test edge cases and boundary conditions
