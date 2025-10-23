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

**Prerequisites:**

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed and running
- Git

**Environment Setup:**

**Create root `.env` file** in the project root directory:

```env
# Database Configuration
POSTGRES_DB=toasted_coffee
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_PORT=5432

# Service Ports
API_PORT=8080
FRONTEND_PORT=5173
ADMIN_PORT=5174

# Application Secrets
JWT_SECRET=your-production-jwt-secret-here
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:5174
```

\*\*

# Start all services (PostgreSQL, Backend, Frontend, Admin)

docker-compose -f docker-compose.dev.yml up --build

# Access Applications:

Customer Frontend: http://localhost:5173
Admin Dashboard: http://localhost:5174
Backend API: http://localhost:8080
API Health Check: http://localhost:8080/health

# Stop all services

docker-compose -f docker-compose.dev.yml down

# Common Docker commands

'''

# Stop all services

docker-compose -f docker-compose.dev.yml down

# Restart specific service (after code changes)

docker-compose -f docker-compose.dev.yml restart backend
docker-compose -f docker-compose.dev.yml restart admin-frontend

# View logs

docker-compose -f docker-compose.dev.yml logs
docker-compose -f docker-compose.dev.yml logs backend
docker-compose -f docker-compose.dev.yml logs -f # Follow logs

# Fresh start (clears database)

docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up --build

# Check service status

docker-compose -f docker-compose.dev.yml ps
'''

# Database Operations:

'''

# Connect to PostgreSQL

docker-compose -f docker-compose.dev.yml exec postgres psql -U postgres -d toasted_coffee

# Check tables

docker-compose -f docker-compose.dev.yml exec postgres psql -U postgres -d toasted_coffee -c "\dt"

# View booking data

docker-compose -f docker-compose.dev.yml exec postgres psql -U postgres -d toasted_coffee -c "SELECT \* FROM bookings;"

# Exit PostgreSQL

\q
'''

# Troubleshooting:

'''

# Check container status

docker-compose -f docker-compose.dev.yml ps

# View specific service logs

docker-compose -f docker-compose.dev.yml logs backend | grep ERROR

# Rebuild specific service

docker-compose -f docker-compose.dev.yml build backend
docker-compose -f docker-compose.dev.yml up backend

# Check environment variables

docker-compose -f docker-compose.dev.yml exec backend env | grep DATABASE

# Test database connection

docker-compose -f docker-compose.dev.yml exec postgres pg_isready -U postgres
'''
