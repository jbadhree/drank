# Banking App Firebase Migration: Developer Notes

## Current Status

We've completed the following steps in migrating the banking app from PostgreSQL to Firebase Firestore:

1. Set up Firebase emulators for local development
2. Migrated the models to work with Firestore (using string IDs instead of uint)
3. Implemented repository interfaces and concrete implementations
4. Updated service layer to work with the repository interfaces
5. Implemented mocks for testing
6. Added integration test infrastructure for Firestore emulator

## Next Steps

### 1. Make Scripts Executable

Run the following commands:

```bash
chmod +x start_emulators.sh
chmod +x startup.sh
```

### 2. Start Emulators for Testing

```bash
# In one terminal
./start_emulators.sh

# Verify emulators are running by checking:
# - Firestore at http://localhost:8091/
# - Auth at http://localhost:9099/
# - Emulator UI at http://localhost:4000/
```

### 3. Run Unit Tests

```bash
# Run all unit tests
go test ./tests/unit/...

# Or run specific tests
go test ./tests/unit/user_service_test.go
```

### 4. Run Integration Tests

```bash
# Make sure emulators are running first
go test ./tests/integration/...
```

### 5. Start the Application

```bash
# With the emulators running in another terminal
./startup.sh

# Or directly
go run main.go
```

### 6. Test API Endpoints

You can use tools like curl, Postman, or the frontend application to test the APIs:

#### Register a new user:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"password123","firstName":"Test","lastName":"User"}'
```

#### Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"password123"}'
```

#### Create an account (with auth token):
```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN_HERE' \
  -d '{"userId":"YOUR_USER_ID","accountType":"CHECKING"}'
```

### 7. Complete Remaining Tests

- Finish the account repository integration tests
- Finish the transaction repository integration tests
- Add handler tests

### 8. Firestore Security Rules Testing

Test the Firestore security rules using the Firebase emulator:

```bash
# In the parent directory where firebase.json is located
firebase emulators:exec --only firestore "npm test"
```

You'll need to create test scripts in JavaScript to test the security rules.

### 9. Production Deployment Preparation

1. Update the configuration to connect to a real Firebase project:
   - Create a Firebase project in the Firebase Console
   - Update the environment variables to point to the real project
   - Generate and download service account credentials if needed

2. Set up continuous integration:
   - Add GitHub Actions or other CI/CD workflow
   - Automate testing with emulators
   - Configure deployment process

3. Set up monitoring and logging:
   - Implement structured logging
   - Configure Firebase performance monitoring
   - Set up alerts for critical events

### 10. Documentation Updates

- Update API documentation
- Create deployment guide
- Document Firestore schema design decisions
- Create runbook for common operations

## Known Issues and Limitations

1. **Limited Query Capabilities**: Firestore has different query capabilities than PostgreSQL. Complex queries might need to be redesigned.

2. **No SQL-like Joins**: Data may need to be denormalized for efficient querying.

3. **Transaction Limitations**: Firestore transactions have limitations on the number of operations and document paths.

4. **Costs**: Be aware of Firestore billing - reads, writes, and deletes are all billable operations.

## Performance Considerations

1. **Minimize Document Size**: Keep documents small and focused.

2. **Batch Operations**: Use batched writes for multiple operations.

3. **Index Management**: Create only necessary indexes to reduce costs.

4. **Collection Group Queries**: Use collection group queries for querying across multiple collections.

## Security Considerations

1. **Authentication**: Ensure Firebase Authentication is properly integrated.

2. **Security Rules**: Test and verify security rules thoroughly.

3. **Data Validation**: Implement validation both server-side and in security rules.

4. **API Security**: Protect your API endpoints with proper authentication and authorization.
