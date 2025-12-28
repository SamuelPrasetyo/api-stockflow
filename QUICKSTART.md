# Quick Start Guide - API Stockflow Authentication

## Prerequisites
- Go 1.24.3 or higher
- PostgreSQL 12 or higher
- Docker (for running tests)

## Setup Steps

### 1. Clone and Navigate to Project
```bash
cd d:\Project\GoWorkspace\api-stockflow
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Configure Environment
Create `.env` file (or use the existing one):
```env
PORT=8080

# Database Configuration
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=stockflow
BLUEPRINT_DB_USERNAME=postgres
BLUEPRINT_DB_PASSWORD=postgres
BLUEPRINT_DB_SCHEMA=public

# JWT Configuration
ACCESS_TOKEN_KEY=your-secret-access-token-key-change-this-in-production
REFRESH_TOKEN_KEY=your-secret-refresh-token-key-change-this-in-production
ACCESS_TOKEN_AGE=3600
```

### 4. Create Database
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE stockflow;

# Exit
\q
```

### 5. Run Migrations
```bash
go run cmd/migrate/migrate.go up
```

### 6. Seed Initial Data (Optional)
```bash
go run cmd/seed/seed.go
```

This will create 3 default users:
- **Admin**: username=`admin`, password=`admin123`
- **Manager**: username=`manager`, password=`manager123`
- **Staff**: username=`staff`, password=`staff123`

### 7. Run the Application
```bash
go run cmd/api/main.go
```

Server will start at `http://localhost:8080`

## Testing the API

### Option 1: Using cURL

**Register a new user:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "fullname": "Test User",
    "role": "staff"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

Save the `access_token` from the response.

**Get Profile:**
```bash
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Option 2: Using VS Code REST Client
Open `api-requests.http` and use the REST Client extension to send requests.

### Option 3: Using Postman
Import the requests from `api-requests.http` into Postman.

## Understanding the Flow

### Registration Flow:
1. User kirim data ke `/api/auth/register`
2. Handler validate input
3. Use Case hash password dengan bcrypt
4. Repository simpan ke database
5. Return success response

### Login Flow:
1. User kirim username & password
2. Handler validate input
3. Use Case:
   - Cari user by username
   - Verify password dengan bcrypt
   - Generate JWT access token (1 hour)
   - Generate JWT refresh token (30 days)
   - Simpan refresh token ke database
4. Return tokens

### Protected Route Flow:
1. Client kirim request dengan Authorization header
2. Auth Middleware:
   - Extract token dari header
   - Verify JWT signature
   - Decode user info dari token
   - Set user ke context
3. Handler bisa akses user info

### Role-Based Authorization:
1. Auth Middleware validate token
2. Role Middleware check user role
3. Reject jika role tidak sesuai

## Project Structure Overview

```
api-stockflow/
├── cmd/
│   ├── api/              # Main application
│   ├── migrate/          # Database migration tool
│   └── seed/             # Database seeder
├── internal/
│   ├── domain/           # Business entities & rules
│   ├── repository/       # Database layer
│   ├── usecase/          # Business logic
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # HTTP middlewares
│   ├── security/         # JWT & password utils
│   ├── database/         # DB connection
│   └── server/           # Server setup
├── .env                  # Environment config
├── api-requests.http     # API examples
└── README_AUTH.md        # Full documentation
```

## Common Issues

### Database Connection Failed
- Make sure PostgreSQL is running
- Check database credentials in `.env`
- Ensure database exists

### Migrations Failed
- Check if database user has CREATE TABLE permissions
- Verify database connection settings

## Next Steps

1. ✅ Authentication is ready
2. 📝 Build transaction features
3. 📝 Add approval workflows
4. 📝 Implement master data management

## Useful Commands

```bash
# Run migrations up
go run cmd/migrate/migrate.go up

# Run migrations down (rollback)
go run cmd/migrate/migrate.go down

# Seed database
go run cmd/seed/seed.go

# Run application
go run cmd/api/main.go

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy
```

## API Endpoints Summary

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/auth/register` | Register new user | No |
| POST | `/api/auth/login` | Login user | No |
| PUT | `/api/auth/refresh` | Refresh access token | No |
| DELETE | `/api/auth/logout` | Logout user | No |
| GET | `/api/auth/profile` | Get user profile | Yes |

## Support

For more details, see:
- [README_AUTH.md](README_AUTH.md) - Complete authentication documentation
- [README.md](README.md) - Full project documentation
