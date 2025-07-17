# Toasted Coffee Co

Cold brew coffee bar catering service
Customer View:
https://toasted-coffee-frontend.onrender.com/
Administrator View:
https://toasted-coffee-admin.onrender.com/

## Project Structure

- `frontend/`: React application built with TypeScript and Express

  - **Frontend Stack**: React, TypeScript, Express
  - **Purpose**: Customer-facing website for booking coffee bar services
  - **Dev Server**: Runs on http://localhost:5173

- `admin-frontend/`: React application built with TypeScript and Express

  - **Admin Stack**: React, TypeScript, Express
  - **Purpose**: Administration portal for managing bookings, inventory, and services
  - **Dev Server**: Runs on http://localhost:5174

- `backend/`: Go with PostgreSQL database
  - **Backend Stack**: Go, PostgreSQL
  - **Purpose**: RESTful API service handling data persistence, authentication, and business logic
  - **Features**: JWT authentication, email notifications, database operations
  - **Server**: Runs on http://localhost:8080

## Development Setup

### Prerequisites

- Node.js 18+
- Go 1.22+
- PostgreSQL 14+

### Backend

Create a `.env` file in the backend directory with the following variables:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=toasted_coffee
DB_SSLMODE=disable
ALLOW_ORIGINS=http://localhost:5173,http://localhost:5174
JWT_SECRET=your-token-secret
JWT_REFRESH_SECRET=your-secure-refresh-token-secret
TOKEN_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=7d
ENVIRONMENT=development
```

If you'd like to use the email service then include these variables as well.

```env
SMTP_HOST=smtp.gmail.com
SMTP_USER=youremail@gmail.com
SMTP_PASSWORD=your-password
NOTIFICATION_EMAIL=youremail@gmail.com
```

1. Create the Database

```bash
psql -U postgres -c "CREATE DATABASE toasted_coffee;"
```

2. Start the Server

```bash
cd backend
go run cmd/api/main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Admin Frontend

```bash
cd frontend
npm install
npm run dev
```
