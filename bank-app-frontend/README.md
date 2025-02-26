# Drank Banking App Frontend

This is the frontend for the Drank Banking application, built with Next.js and TypeScript.

## Features

- User authentication with JWT
- Dashboard to view account summaries
- View detailed account information and transaction history
- Transfer money between accounts

## Technologies Used

- Next.js 14
- React 18
- TypeScript
- Tailwind CSS for styling
- SWR for data fetching
- React Hook Form for form handling
- Zod for validation

## Getting Started

1. Install dependencies:
```bash
npm install
```

2. Create a `.env.local` file with:
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

3. Run the development server:
```bash
npm run dev
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser.

## Demo Credentials

Use these credentials to log in:

```
Email: john.doe@example.com
Password: password123
```

or

```
Email: jane.smith@example.com
Password: password123
```

## API Endpoints Used

The frontend interacts with the following API endpoints:

- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/users/me` - Get current user
- `GET /api/v1/accounts/user/:userId` - Get accounts by user ID
- `GET /api/v1/transactions/account/:accountId` - Get transactions by account ID
- `POST /api/v1/transactions/transfer` - Transfer money between accounts
