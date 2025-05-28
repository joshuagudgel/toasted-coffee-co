# Toasted Coffee Co
Cold brew coffee bar booking platform
https://toasted-coffee-frontend.onrender.com/

## Project Structure

- `frontend/`: React application built with TypeScript and Vite
- `backend/`: Go with PostgreSQL database

## Development Setup

### Prerequisites

- Node.js 18+
- Go 1.22+
- PostgreSQL 14+

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Backend

Create a `.env` file in the backend directory with the following variables:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=toasted_coffee
DB_SSLMODE=disable
ALLOW_ORIGINS=http://localhost:5173
```

1. Create the Database

```bash
psql -U postgres -c "CREATE DATABASE toasted_coffee;"
```

2. Run the migration

```bash
psql -U postgres -d toasted_coffee -f internal/database/migrations/01_create_bookings_table.sql
```

3. Start the Server

```bash
cd backend
go run cmd/api/main.go
```
