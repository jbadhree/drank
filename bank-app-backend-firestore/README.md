# Drank Banking App - Firebase Version

This is a migration of the original Drank Banking App from PostgreSQL to Firebase Firestore.

## Overview

This project demonstrates migrating a Go backend application from a relational database (PostgreSQL) to a NoSQL document database (Firebase Firestore). The application maintains the same API and business logic but changes the data persistence layer.

## Key Features

- **Firestore Database**: Uses Google Cloud Firestore for data storage
- **Firebase Authentication**: Integration with Firebase Auth
- **Domain-Driven Design**: Clear separation between models, repositories, services, and handlers
- **Repository Pattern**: Interface-based repository design for better testability
- **Atomic Transactions**: Implements Firestore transactions for ensuring data integrity
- **JWT Authentication**: Session management using JWT tokens
- **Unit & Integration Tests**: Comprehensive test suite using Firebase emulators

## Prerequisites

- Go 1.17+
- Node.js 14+ (for Firebase Emulators)
- Firebase CLI

## Migration Highlights

### Key Changes from PostgreSQL Version
1. **Model Modifications**:
   - Changed ID fields from `uint` to `string` to use Firestore document IDs
   - Added Firestore-specific struct tags (e.g., `firestore:"fieldName"`)
   - Updated DTOs to reflect the new ID type

2. **Repository Layer**:
   - Implemented Firestore-specific CRUD operations
   - Added interface definitions for repositories to improve testability
   - Utilized Firestore transactions for operations requiring atomicity

3. **Authentication**:
   - Added Firebase Authentication integration
   - Maintained JWT token generation for API authentication

## Setup Instructions

### 1. Install Firebase CLI

If you haven't already installed the Firebase CLI, you can do so with npm:

```bash
npm install -g firebase-tools
```

### 2. Login to Firebase

```bash
firebase login
```

### 3. Start Firebase Emulators

From the root directory (where firebase.json is located), run:

```bash
firebase emulators:start
```

Or use the provided script:

```bash
./start_emulators.sh
```

This will start:
- Firestore Emulator (port 8091)
- Authentication Emulator (port 9099)
- Emulator UI (usually on port 4000)

### 4. Build and Run the Backend

```bash
cd bank-app-backend-firestore
go build
./bank-app-backend-firestore
```

Or with the startup script:

```bash
./startup.sh
```

The server will start on port 8080 by default.

### 5. Seed the Database (optional)

```bash
cd bank-app-backend-firestore
go run main.go --seed
```

This will populate the database with test users, accounts, and transactions.

## Running Tests

### Unit Tests

```bash
go test ./tests/unit/...
```

### Integration Tests (requires emulators running)

```bash
go test ./tests/integration/...
```

### Functional Tests (requires emulators running)

```bash
go test ./tests/functional/...
```

Functional tests perform end-to-end testing of the API endpoints using the Firebase emulators.

## API Endpoints

### Authentication

- `POST /api/v1/auth/login` - Login with email and password
- `POST /api/v1/auth/register` - Register a new user

### Users

- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `GET /api/v1/users/me` - Get current user

### Accounts

- `GET /api/v1/accounts` - Get all accounts
- `GET /api/v1/accounts/:id` - Get account by ID
- `GET /api/v1/accounts/user/:userId` - Get accounts by user ID
- `POST /api/v1/accounts` - Create a new account

### Transactions

- `GET /api/v1/transactions` - Get all transactions
- `GET /api/v1/transactions/:id` - Get transaction by ID
- `GET /api/v1/transactions/account/:accountId` - Get transactions by account ID
- `POST /api/v1/transactions/transfer` - Transfer funds between accounts
- `POST /api/v1/transactions/deposit` - Create a deposit transaction
- `POST /api/v1/transactions/withdrawal` - Create a withdrawal transaction

## Environment Variables

Create a `.env` file in the `bank-app-backend-firestore` directory with the following variables:

```
FIREBASE_PROJECT_ID=drank-firebase
FIREBASE_AUTH_EMULATOR_HOST=localhost:9099
FIRESTORE_EMULATOR_HOST=localhost:8091
PORT=8080
JWT_SECRET=your-very-secret-jwt-key-change-in-production
```

## Architecture

### Repository Pattern with Interfaces

The application uses a repository pattern with interfaces to facilitate testing and maintain separation of concerns:

```go
// Interface definition
type UserRepository interface {
    Create(user models.User) (models.User, error)
    FindByID(id string) (models.User, error)
    // other methods...
}

// Implementation
type UserRepositoryImpl struct {
    client *firestore.Client
    ctx    context.Context
}
```

### Service Layer

Services encapsulate business logic and use repository interfaces:

```go
type UserService struct {
    repo interfaces.UserRepository
}
```

### Firestore Transactions

For operations requiring atomicity (like fund transfers), Firestore transactions are used:

```go
// Use a Firestore transaction for atomic operation
return r.client.RunTransaction(r.ctx, func(ctx context.Context, tx *firestore.Transaction) error {
    // Transaction code here
})
```

## Data Model

### Users
- ID (string) - Firestore document ID
- Email (string)
- Password (string) - Hashed
- FirstName (string)
- LastName (string)
- CreatedAt (timestamp)
- UpdatedAt (timestamp)

### Accounts
- ID (string) - Firestore document ID
- UserID (string) - References Users collection
- AccountNumber (string)
- AccountType (CHECKING or SAVINGS)
- Balance (float64)
- CreatedAt (timestamp)
- UpdatedAt (timestamp)

### Transactions
- ID (string) - Firestore document ID
- AccountID (string) - References Accounts collection
- SourceAccountID (string, optional) - For transfers
- TargetAccountID (string, optional) - For transfers
- Amount (float64)
- Balance (float64) - Account balance after transaction
- Type (DEPOSIT, WITHDRAWAL, or TRANSFER)
- Description (string)
- TransactionDate (timestamp)
- CreatedAt (timestamp)
- UpdatedAt (timestamp)

## Firebase Security Rules

The application uses Firestore security rules to ensure proper access control:

```
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    // Users collection
    match /users/{userId} {
      allow read: if request.auth != null;
      allow create: if request.auth != null;
      allow update, delete: if request.auth != null && request.auth.uid == userId;
    }
    
    // Accounts collection
    match /accounts/{accountId} {
      allow read: if request.auth != null;
      allow create: if request.auth != null;
      allow update, delete: if request.auth != null && 
        request.resource.data.userId == request.auth.uid;
    }
    
    // Transactions collection
    match /transactions/{transactionId} {
      allow read: if request.auth != null;
      allow create: if request.auth != null;
      // Only allow updates and deletes if the user owns the account
      allow update, delete: if request.auth != null && exists(/databases/$(database)/documents/accounts/$(resource.data.accountId)) &&
        get(/databases/$(database)/documents/accounts/$(resource.data.accountId)).data.userId == request.auth.uid;
    }
  }
}
```

## Production Deployment

To deploy the application to production:

1. **Create a Firebase project**:
   - Go to the [Firebase Console](https://console.firebase.google.com/)
   - Create a new project
   - Enable Firestore and Authentication

2. **Update environment variables**:
   ```
   FIREBASE_PROJECT_ID=your-production-project-id
   # Remove emulator hosts for production
   # FIREBASE_AUTH_EMULATOR_HOST=
   # FIRESTORE_EMULATOR_HOST=
   PORT=8080
   JWT_SECRET=your-production-secret-key
   ```

3. **Deploy the application**:
   - Deploy the application to your preferred hosting platform
   - Ensure the application has network access to Firebase services

4. **Configure Firebase Authentication**:
   - Set up authentication providers in Firebase Console
   - Update security rules for production

## Testing Strategy

### Unit Tests
- Test individual components in isolation
- Mock repository interfaces for service tests
- Focus on business logic validation

### Integration Tests
- Test interaction with Firestore emulator
- Verify data persistence and retrieval
- Test transaction integrity

### End-to-End Tests
- Test complete API flows
- Verify authentication and authorization
- Test real-world scenarios

## Common Migration Challenges

1. **Document Design**: NoSQL databases require different thinking about data structure
   - Solution: Designed documents around query patterns
   - Duplicated some data for performance

2. **Transactions**: Ensuring atomic operations
   - Solution: Implemented Firestore transactions
   - Limited transaction scope to related documents

3. **Query Capabilities**: Firestore has different querying capabilities
   - Solution: Created proper indexes
   - Redesigned complex queries to work with Firestore

4. **Testing**: Different testing approach needed
   - Solution: Created test helpers for emulator integration
   - Implemented clean-up routines for test isolation

## Future Improvements

1. **Scalability**:
   - Add caching layer (Redis)
   - Implement pagination for large collections

2. **Security**:
   - Enhance security rules
   - Add rate limiting

3. **Monitoring**:
   - Implement logging and metrics
   - Set up alerts for abnormal patterns

4. **CI/CD**:
   - Add GitHub Actions workflow
   - Automate deployment
