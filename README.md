# Go Backend API

A robust and well-structured Go backend API with role-based authentication, built with Gin, GORM, and JWT.

## Features

- **Authentication & Authorization**
  - JWT-based authentication
  - Role-based access control (Admin, Moderator, User)
  - Password hashing with bcrypt
  - Token validation middleware

- **Database**
  - Support for PostgreSQL and SQLite
  - GORM ORM with migrations
  - Soft deletes
  - Database seeding

- **API Features**
  - RESTful API design
  - Request validation
  - Structured logging with Logrus
  - Error handling middleware
  - CORS support
  - Rate limiting
  - Security headers

- **Code Organization**
  - Clean architecture
  - Dependency injection
  - Environment-based configuration
  - Graceful shutdown

## Project Structure

```
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database connection and migrations
│   ├── handlers/
│   │   ├── router.go            # Route definitions
│   │   ├── user_handler.go      # User HTTP handlers
│   │   └── health_handler.go    # Health check handlers
│   ├── middleware/
│   │   ├── auth.go              # Authentication middleware
│   │   └── common.go            # Common middleware (CORS, logging, etc.)
│   ├── models/
│   │   └── user.go              # Data models
│   ├── services/
│   │   └── user_service.go      # Business logic
│   └── utils/
│       ├── jwt.go               # JWT utilities
│       └── validation.go        # Request validation
├── pkg/
│   └── logger/
│       └── logger.go            # Logging utilities
├── .env                         # Environment variables
├── .env.example                 # Environment variables template
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL (optional, SQLite is default)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-backend
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Edit `.env` file with your configuration:
```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration (SQLite by default)
DB_TYPE=sqlite
SQLITE_PATH=./app.db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h
```

4. Install dependencies:
```bash
go mod tidy
```

5. Run the application:
```bash
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`.

### Using PostgreSQL

To use PostgreSQL instead of SQLite, update your `.env` file:

```env
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=your_database
DB_USER=your_username
DB_PASSWORD=your_password
DB_SSL_MODE=disable
```

## API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `GET /ready` - Readiness check
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### Protected Endpoints (Require Authentication)

- `GET /api/v1/user/profile` - Get current user profile
- `PUT /api/v1/user/profile` - Update current user profile
- `POST /api/v1/user/change-password` - Change password

### Admin Endpoints (Admin Role Required)

- `GET /api/v1/admin/users` - Get all users (paginated)
- `GET /api/v1/admin/users/:id` - Get user by ID
- `PUT /api/v1/admin/users/:id` - Update user
- `DELETE /api/v1/admin/users/:id` - Delete user

### Moderator Endpoints (Admin or Moderator Role Required)

- `GET /api/v1/mod/users` - Get all users (view only)

## Authentication

### Register a new user:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "testuser",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Use the returned JWT token in subsequent requests:
```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Default Admin Account

The application creates a default admin account on first run:
- **Email**: `admin@example.com`
- **Password**: `admin123`
- **Role**: `admin`

**Important**: Change the default admin password in production!

## User Roles

- **Admin**: Full access to all endpoints
- **Moderator**: Can view users and moderate content
- **User**: Basic access to their own profile

## Configuration

The application uses environment variables for configuration. See `.env.example` for all available options.

### Key Configuration Options

- `JWT_SECRET`: Secret key for JWT token signing (change in production!)
- `DB_TYPE`: Database type (`sqlite` or `postgres`)
- `LOG_LEVEL`: Logging level (`debug`, `info`, `warn`, `error`)
- `CORS_ORIGINS`: Allowed CORS origins (comma-separated)

## Development

### Running in Development Mode

```bash
# Set environment to development
export NODE_ENV=development

# Run with live reload (using air)
air
```

### Building for Production

```bash
# Build binary
go build -o bin/server cmd/server/main.go

# Run binary
./bin/server
```

## Security Features

- Password hashing with bcrypt
- JWT token authentication
- Role-based authorization
- CORS protection
- Rate limiting
- Security headers
- Input validation
- SQL injection prevention (GORM)

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Error message",
  "errors": {
    "field": "Validation error message"
  }
}
```

## Logging

Structured logging with Logrus:
- JSON format in production
- Configurable log levels
- Request/response logging
- Error tracking with context

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.