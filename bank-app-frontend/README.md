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

## Prerequisites

Before running the frontend, make sure you have the following installed:

- **Node.js 16 or later**
  - [Download Node.js](https://nodejs.org/)
  - Verify installation: `node -v` and `npm -v`

- **Backend API Running**
  - The Go backend should be running at http://localhost:8080
  - See the main project README for backend setup instructions

- **PostgreSQL Database**
  - Ensure the database is running via Docker Compose
  - Required for the backend to function properly

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

## Build for Production

To create a production build:

```bash
npm run build
```

To run the production build:

```bash
npm start
```

## Project Structure

- `/app` - Next.js App Router pages and layouts
- `/components` - Reusable UI components
- `/lib` - Utility functions, types, and API services
- `/public` - Static assets

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

## Troubleshooting

- **API Connection Issues**: Make sure the backend server is running and accessible
- **Authentication Problems**: Try clearing localStorage and logging in again
- **UI Rendering Issues**: Ensure you're using a supported browser (Chrome, Firefox, Safari, Edge)
