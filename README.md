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

## Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Go 1.17 or later**: Required for the backend
  - [Download Go](https://golang.org/dl/)
  - Verify installation: `go version`

- **Node.js 16 or later**: Required for the frontend
  - [Download Node.js](https://nodejs.org/)
  - Verify installation: `node -v` and `npm -v`

- **Docker and Docker Compose**: Required for the database
  - [Install Docker](https://docs.docker.com/get-docker/)
  - [Install Docker Compose](https://docs.docker.com/compose/install/)
  - Verify installation: `docker --version` and `docker-compose --version`

## Setup and Installation

### 1. Database Setup

First, start the PostgreSQL database using Docker Compose:

```bash
# Navigate to the database setup directory
cd databse_setup

# Start the PostgreSQL container
docker-compose -f postgres_docker_compose.yml up -d
```

This will start a PostgreSQL instance on port 5434 with the following credentials:
- Database: drank
- Username: postgres
- Password: Demo123!

You can access the Adminer database management tool at http://localhost:8070.

### 2. Backend Setup

1. Navigate to the backend directory:

```bash
cd bank-app-backend
```

2. Install the Swagger CLI tool:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

3. Initialize the Swagger documentation:

```bash
~/go/bin/swag init  # or just 'swag init' if it's in your PATH
```

4. Seed the database with initial data:

```bash
go run main.go --seed
```

5. Start the backend server:

```bash
go run main.go
```

The API will be available at http://localhost:8080, and Swagger documentation at http://localhost:8080/swagger/index.html.

### 3. Frontend Setup

1. Navigate to the frontend directory:

```bash
cd bank-app-frontend
```

2. Install dependencies:

```bash
npm install
```

3. Ensure the `.env.local` file exists with the API URL:

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

4. Start the development server:

```bash
npm run dev
```

The frontend will be available at http://localhost:3000.

## Testing the Application

1. Open your browser and navigate to http://localhost:3000
2. You will be redirected to the login page
3. Use one of the demo credentials to log in

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

## Troubleshooting

- **Database connection issues**: Ensure the Docker container is running with `docker ps`
- **Swagger initialization errors**: Make sure you have the latest version of swag installed
- **CORS errors**: Check that the frontend is connecting to the correct API URL

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
