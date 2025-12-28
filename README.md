# API Stockflow

A RESTful API for stock flow management system built with Go, Fiber, and PostgreSQL using Clean Architecture principles.

## 🚀 Features

### ✅ Authentication & Authorization (Implemented)
- JWT-based authentication with access and refresh tokens
- Role-based access control (Admin/Manager/Staff)
- Password hashing with bcrypt
- Secure token management

**User Roles:**
- **Admin/Manager**: Full access to all features and can approve any amount
- **Staff**: Limited access, cannot approve purchases over 10,000,000

### 📝 Coming Soon
- Transaction management
- Approval workflows  
- Master data management
- Reporting

## 🏗️ Architecture

This project follows **Clean Architecture** principles:

```
├── cmd/
│   ├── api/              # Application entry point
│   ├── migrate/          # Database migration tool
│   └── seed/             # Database seeder
├── internal/
│   ├── domain/           # Business entities & rules (Enterprise Business Rules)
│   ├── repository/       # Data access interfaces & implementations (Interface Adapters)
│   ├── usecase/          # Business logic (Application Business Rules)
│   ├── handler/          # HTTP handlers (Interface Adapters)
│   ├── middleware/       # HTTP middlewares (Interface Adapters)
│   ├── security/         # Security utilities (Frameworks & Drivers)
│   ├── database/         # Database connection (Frameworks & Drivers)
│   └── server/           # Server setup (Frameworks & Drivers)
```

## 🛠️ Tech Stack

- **Go** 1.24.3
- **Fiber** v2 - Web framework
- **PostgreSQL** - Database
- **JWT** - Authentication
- **Bcrypt** - Password hashing
- **Testcontainers** - Integration testing

## 📋 Prerequisites

- Go 1.24.3 or higher
- PostgreSQL 12 or higher
- Docker (for running tests)
- Make (optional, for using Makefile commands)

## 🚀 Quick Start

### 1. Clone the repository
```bash
cd d:\Project\GoWorkspace\api-stockflow
```

### 2. Install dependencies
```bash
go mod download
# or
make deps
```

### 3. Configure environment
Copy `.env.example` to `.env` and update the values:
```bash
cp .env.example .env
```

Key environment variables:
```env
PORT=8080
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=stockflow
BLUEPRINT_DB_USERNAME=postgres
BLUEPRINT_DB_PASSWORD=postgres
BLUEPRINT_DB_SCHEMA=public
ACCESS_TOKEN_KEY=your-secret-access-token-key
REFRESH_TOKEN_KEY=your-secret-refresh-token-key
ACCESS_TOKEN_AGE=3600
```

### 4. Setup database
```bash
# Create database
createdb stockflow

# Run migrations
make migrate-up

# Seed initial data (optional)
make seed
```

Default seeded users:
- **Admin**: `admin` / `admin123`
- **Manager**: `manager` / `manager123`
- **Staff**: `staff` / `staff123`

### 5. Run the application
```bash
make run
# or with live reload
make dev
```

Server will start at `http://localhost:8080`

## 📚 API Documentation

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/auth/register` | Register new user | No |
| POST | `/api/auth/login` | Login user | No |
| PUT | `/api/auth/refresh` | Refresh access token | No |
| DELETE | `/api/auth/logout` | Logout user | No |
| GET | `/api/auth/profile` | Get user profile | Yes |

### Example Requests

See [api-requests.http](api-requests.http) for complete API examples.

**Register:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "password123",
    "fullname": "John Doe",
    "role": "staff"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "password123"
  }'
```

**Get Profile:**
```bash
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 📦 Makefile Commands

```bash
make help         # Show all available commands
make build        # Build the application
make run          # Run the application
make dev          # Run with live reload
make migrate-up   # Run database migrations
make migrate-down # Rollback database migrations
make seed         # Seed database with initial data
make db-setup     # Setup database (migrate + seed)
make docker-run   # Start Docker containers
make docker-down  # Stop Docker containers
make clean        # Clean build artifacts
make deps         # Install dependencies
make fmt          # Format code
```

## 📁 Project Structure

```
api-stockflow/
├── cmd/
│   ├── api/
│   │   └── main.go                      # Application entry point
│   ├── migrate/
│   │   └── migrate.go                   # Database migration tool
│   └── seed/
│       └── seed.go                      # Database seeder
├── internal/
│   ├── domain/
│   │   ├── user.go                      # User entity
│   │   ├── authentication.go            # Authentication entity
│   │   └── errors.go                    # Domain errors
│   ├── repository/
│   │   ├── user_repository.go           # User repository interface
│   │   ├── user_repository_postgres.go  # User repository implementation
│   │   ├── authentication_repository.go # Auth repository interface
│   │   └── authentication_repository_postgres.go # Auth repository implementation
│   ├── usecase/
│   │   └── auth_usecase.go             # Authentication use case
│   ├── handler/
│   │   └── auth_handler.go             # Authentication HTTP handler
│   ├── middleware/
│   │   └── auth_middleware.go          # Authentication & authorization middleware
│   ├── security/
│   │   ├── token_manager.go            # JWT token manager
│   │   └── password_hash.go            # Password hashing
│   ├── database/
│   │   ├── database.go                 # Database connection
│   │   ├── database_test.go            # Database tests
│   │   └── migrations/                 # SQL migration files
│   │       ├── 000001_create_users_table.up.sql
│   │       ├── 000001_create_users_table.down.sql
│   │       ├── 000002_create_authentications_table.up.sql
│   │       └── 000002_create_authentications_table.down.sql
│   └── server/
│       ├── server.go                   # Server setup
│       ├── routes.go                   # Route definitions
│       └── routes_test.go              # Route tests
├── .env                                # Environment variables
├── .env.example                        # Environment variables example
├── api-requests.http                   # API request examples
├── docker-compose.yml                  # Docker compose configuration
├── go.mod                              # Go module file
├── go.sum                              # Go dependencies checksum
├── Makefile                            # Build automation
├── README.md                           # This file
├── README_AUTH.md                      # Authentication documentation
├── TESTING.md                          # Testing documentation
└── QUICKSTART.md                       # Quick start guide
```

## 🔒 Security

- Passwords are hashed using bcrypt with salt
- JWT tokens for stateless authentication
- Access tokens expire after configured time (default: 1 hour)
- Refresh tokens for long-lived sessions (30 days)
- Role-based access control for authorization
- SQL injection prevention through parameterized queries

## 🎯 Middleware

### AuthMiddleware
Validates JWT access token and sets user in context.

```go
authMiddleware := middleware.AuthMiddleware(tokenManager, userRepo)
app.Get("/protected", authMiddleware, handler)
```

### RoleMiddleware
Restricts access based on user roles.

```go
roleMiddleware := middleware.RoleMiddleware(domain.RoleAdmin, domain.RoleManager)
app.Get("/admin", authMiddleware, roleMiddleware, handler)
```

### ApprovalMiddleware
Checks if user can approve based on purchase amount.

```go
approvalMiddleware := middleware.ApprovalMiddleware("amount")
app.Post("/approve", authMiddleware, approvalMiddleware, handler)
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License.

## 📞 Support

For more information:
- [Quick Start Guide](QUICKSTART.md)
- [Authentication Documentation](README_AUTH.md)

## 🗺️ Roadmap

- [x] Authentication & Authorization
- [x] JWT Token Management
- [x] Role-based Access Control
- [ ] Transaction Management
- [ ] Approval Workflows
- [ ] Master Data Management
- [ ] API Documentation (Swagger)
- [ ] Logging & Monitoring
- [ ] Rate Limiting
- [ ] API Versioning
