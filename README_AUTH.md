# API Stockflow - Authentication Feature

This project implements JWT-based authentication with role-based access control.

## Features

- JWT Authentication (Access Token & Refresh Token)
- Role-based Authorization (Admin/Manager & Staff)
- Clean Architecture implementation
- Comprehensive test coverage (Unit, Integration, Functional tests)

## User Roles

### Admin/Manager
- Full access to all menus (master, transactions, approvals)
- Can approve purchases of any amount

### Staff
- Limited access (can only input transactions)
- Cannot approve purchases over 10,000,000

## Tech Stack

- Go 1.24.3
- Fiber Web Framework
- PostgreSQL
- JWT for authentication
- Bcrypt for password hashing
- Testcontainers for integration testing

## Project Structure

```
api-stockflow/
├── cmd/
│   ├── api/              # Main application entry point
│   └── migrate/          # Database migration tool
├── internal/
│   ├── domain/           # Domain entities and business rules
│   ├── repository/       # Data access layer
│   ├── usecase/          # Business logic layer
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # HTTP middlewares
│   ├── security/         # Security utilities (JWT, password hash)
│   ├── database/         # Database connection
│   └── server/           # Server setup and routes
└── .env                  # Environment configuration
```

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Copy `.env.example` to `.env` and configure your settings:
```bash
cp .env.example .env
```

3. Run database migrations:
```bash
go run cmd/migrate/migrate.go up
```

4. Run the application:
```bash
go run cmd/api/main.go
```

## Database Migrations

Run migrations up:
```bash
go run cmd/migrate/migrate.go up
```

Run migrations down (rollback):
```bash
go run cmd/migrate/migrate.go down
```

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `PUT /api/auth/refresh` - Refresh access token
- `DELETE /api/auth/logout` - Logout user
- `GET /api/auth/profile` - Get user profile (requires authentication)

### Example Requests

#### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123",
    "fullname": "Admin User",
    "role": "admin"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

#### Get Profile
```bash
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Testing

### Run all tests:
```bash
go test ./... -v
```

### Run specific test types:

Unit tests:
```bash
go test ./internal/domain/... -v
go test ./internal/security/... -v
go test ./internal/usecase/... -v
```

Integration tests:
```bash
go test ./internal/repository/... -v
```

Functional tests:
```bash
go test ./internal/handler/... -v
```

## Database Migrations

Migrations are located in `internal/database/migrations/` directory:

**Run migrations:**
```bash
go run cmd/migrate/migrate.go up
```

**Rollback migrations:**
```bash
go run cmd/migrate/migrate.go down
```

The migration tool automatically reads migration files and executes them in order.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| BLUEPRINT_DB_HOST | Database host | localhost |
| BLUEPRINT_DB_PORT | Database port | 5432 |
| BLUEPRINT_DB_DATABASE | Database name | stockflow |
| BLUEPRINT_DB_USERNAME | Database username | postgres |
| BLUEPRINT_DB_PASSWORD | Database password | postgres |
| BLUEPRINT_DB_SCHEMA | Database schema | public |
| ACCESS_TOKEN_KEY | JWT access token secret key | (required) |
| REFRESH_TOKEN_KEY | JWT refresh token secret key | (required) |
| ACCESS_TOKEN_AGE | Access token expiration in seconds | 3600 |

## Security

- Passwords are hashed using bcrypt
- JWT tokens for authentication
- Access tokens expire after configured time (default: 1 hour)
- Refresh tokens for long-lived sessions (30 days)
- Role-based access control for authorization

## Middleware

### AuthMiddleware
Validates JWT access token and sets user in context.

Usage:
```go
authMiddleware := middleware.AuthMiddleware(tokenManager, userRepo)
app.Get("/protected", authMiddleware, handler)
```

### RoleMiddleware
Restricts access based on user roles.

Usage:
```go
roleMiddleware := middleware.RoleMiddleware(domain.RoleAdmin, domain.RoleManager)
app.Get("/admin", authMiddleware, roleMiddleware, handler)
```

### ApprovalMiddleware
Checks if user can approve based on purchase amount.

Usage:
```go
approvalMiddleware := middleware.ApprovalMiddleware("amount")
app.Post("/approve", authMiddleware, approvalMiddleware, handler)
```

## Understanding Clean Architecture

Project ini menggunakan Clean Architecture dengan layer:

1. **Domain Layer** (`internal/domain/`)
   - Entity dan business rules
   - Tidak bergantung pada layer lain
   - Contoh: User entity dengan method CanApprove()

2. **Repository Layer** (`internal/repository/`)
   - Interface untuk akses data
   - Implementation untuk PostgreSQL
   - Memisahkan business logic dari database

3. **Use Case Layer** (`internal/usecase/`)
   - Business logic aplikasi
   - Orchestrasi antara repository dan domain
   - Contoh: Login, Register, Refresh Token

4. **Handler Layer** (`internal/handler/`)
   - HTTP request/response handling
   - Validation input
   - Memanggil use case

5. **Infrastructure** (`internal/security/`, `internal/database/`)
   - Tool dan utility
   - Database connection
   - JWT, password hashing

**Alur Request:**
```
HTTP Request → Handler → Use Case → Repository → Database
                  ↓         ↓          ↓
              Validation  Business   Data Access
                         Logic
```

## License

MIT
