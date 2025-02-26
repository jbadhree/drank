# Banking Application Demo

This is a full-stack banking application demo with a Next.js frontend and Go backend. The application allows users to view their accounts, check transaction history, and transfer money between accounts.

## Features

- User authentication with JWT
- Dashboard to view account summaries
- View detailed account information and transaction history
- Transfer money between accounts
- RESTful API with Swagger documentation

## Project Structure

The project is divided into two main parts:

1. Backend (Go)
2. Frontend (Next.js)

### Backend

The backend is built using Go with the following components:

- Gin for HTTP routing
- GORM for database ORM
- PostgreSQL for the database (running on port 5434)
- Swagger for API documentation

### Frontend

The frontend is built using Next.js with the following technologies:

- React 18
- TypeScript
- Tailwind CSS for styling
- SWR for data fetching
- React Hook Form for form handling
- Zod for validation

## Setup and Installation

### Prerequisites

- Go (1.17 or later)
- Node.js (16 or later)
- PostgreSQL

### Backend Setup

1. Create a PostgreSQL database:

```sql
-- Note: Using the existing database 'drank' on port 5434
-- The database is already set up via docker-compose
```

2. Navigate to the backend directory:

```bash
cd bank-app-backend
```

3. Initialize the Swagger documentation:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

4. Run database migrations and seed data:

```bash
go run main.go --seed
```

5. Start the backend server:

```bash
go run main.go
```

The API will be available at http://localhost:8080, and Swagger documentation at http://localhost:8080/swagger/index.html.

Note: The application is configured to connect to a PostgreSQL database running on port 5434 with the following credentials:
- Database: drank
- Username: postgres
- Password: Demo123!

Make sure this database is running before starting the application.

### Frontend Setup

To set up the frontend, you would:

1. Navigate to the frontend directory
2. Run `npm install` to install dependencies
3. Create a `.env.local` file with the API URL
4. Run `npm run dev` to start the development server

The frontend code is not yet implemented in this repository, but the structure and API services are provided in this README.

## Demo Credentials

Use the following credentials to log in:

```
Email: john.doe@example.com
Password: password123
```

or

```
Email: jane.smith@example.com
Password: password123
```

## API Endpoints

Here are the main API endpoints:

### Authentication

- `POST /api/v1/auth/login` - Login user

### Users

- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `GET /api/v1/users/me` - Get current user

### Accounts

- `GET /api/v1/accounts` - Get all accounts
- `GET /api/v1/accounts/:id` - Get account by ID
- `GET /api/v1/accounts/user/:userId` - Get accounts by user ID

### Transactions

- `GET /api/v1/transactions` - Get all transactions
- `GET /api/v1/transactions/:id` - Get transaction by ID
- `GET /api/v1/transactions/account/:accountId` - Get transactions by account ID
- `POST /api/v1/transactions/transfer` - Transfer money between accounts
