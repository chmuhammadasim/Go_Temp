# 🚀 Go Backend API

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](Dockerfile)

A **production-ready**, robust and well-structured Go backend API with comprehensive role-based authentication, built with modern Go technologies including Gin, GORM, and JWT.

## ✨ Features

## ✨ Features

### 🔐 Authentication & Authorization
- **JWT-based authentication** with secure token generation and validation
- **Role-based access control** (Admin, Moderator, User) with granular permissions
- **Password hashing** with bcrypt for maximum security
- **Token validation middleware** for protected routes
- **Owner-based access control** for user-specific resources

### 🗄️ Database & Storage
- **Multi-database support** (PostgreSQL and SQLite)
- **GORM ORM** with automatic migrations and relationships
- **Soft deletes** for data integrity
- **Database seeding** with default admin account
- **Connection pooling** and timeout management

### 🌐 API Features
- **RESTful API design** following industry standards
- **Request validation** with comprehensive error messages
- **Structured logging** with Logrus (JSON/Text formats)
- **Global error handling** middleware
- **CORS support** with configurable origins
- **Rate limiting** to prevent abuse
- **Security headers** for enhanced protection
- **Health checks** for monitoring and deployment

### 🏗️ Code Organization & Architecture
- **Clean architecture** with clear separation of concerns
- **Dependency injection** for testability and maintainability
- **Environment-based configuration** for different deployment stages
- **Graceful shutdown** handling for production deployments
- **Middleware pipeline** for request processing
- **Service layer** for business logic separation

### 🐳 DevOps & Deployment
- **Docker support** with multi-stage builds
- **Makefile** for common development tasks
- **GitHub Actions ready** (CI/CD pipeline compatible)
- **Environment variable management**
- **Binary compilation** for different platforms

## 📁 Project Structure

```
go-backend/
├── 📁 cmd/
│   └── 📁 server/
│       └── 📄 main.go              # 🚀 Application entry point & server setup
├── 📁 internal/                    # 🔒 Private application code
│   ├── 📁 config/
│   │   └── 📄 config.go            # ⚙️ Configuration management & env variables
│   ├── 📁 database/
│   │   └── 📄 database.go          # 🗄️ Database connection, migrations & seeding
│   ├── 📁 handlers/
│   │   ├── 📄 router.go            # 🌐 Route definitions & middleware setup
│   │   ├── 📄 user_handler.go      # 👤 User HTTP handlers & endpoints
│   │   └── 📄 health_handler.go    # ❤️ Health check endpoints
│   ├── 📁 middleware/
│   │   ├── 📄 auth.go              # 🔐 Authentication & authorization middleware
│   │   └── 📄 common.go            # 🛠️ Common middleware (CORS, logging, security)
│   ├── 📁 models/
│   │   └── 📄 user.go              # 📋 Data models, validation & relationships
│   ├── 📁 services/
│   │   └── 📄 user_service.go      # 💼 Business logic & data operations
│   └── 📁 utils/
│       ├── 📄 jwt.go               # 🎫 JWT token utilities & validation
│       └── 📄 validation.go        # ✅ Request validation helpers
├── 📁 pkg/                         # 📦 Reusable packages
│   └── 📁 logger/
│       └── 📄 logger.go            # 📊 Structured logging utilities
├── 📁 configs/                     # 📁 Configuration files directory
├── 📄 .env                         # 🔧 Environment variables (local)
├── 📄 .env.example                 # 📝 Environment variables template
├── 📄 go.mod                       # 📋 Go module definition & dependencies
├── 📄 go.sum                       # 🔒 Dependency lock file
├── 📄 Dockerfile                   # 🐳 Docker container definition
├── 📄 Makefile                     # 🔨 Build automation & development tasks
├── 📄 .gitignore                   # 🚫 Git ignore rules
└── 📄 README.md                    # 📖 Project documentation (this file)
```

## �️ Technology Stack

| Technology | Purpose | Version |
|------------|---------|---------|
| **Go** | Backend Language | 1.21+ |
| **Gin** | HTTP Web Framework | Latest |
| **GORM** | ORM Library | Latest |
| **JWT** | Authentication | v4 |
| **PostgreSQL** | Primary Database | Any |
| **SQLite** | Development Database | Built-in |
| **Logrus** | Structured Logging | Latest |
| **Bcrypt** | Password Hashing | Built-in |
| **Docker** | Containerization | Any |

## 🚀 Getting Started

### 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.21 or higher** - [Download here](https://golang.org/dl/)
- **Git** - For cloning the repository
- **PostgreSQL** (optional) - SQLite is used by default
- **Docker** (optional) - For containerized deployment

### ⚡ Quick Start

1. **Clone the repository:**
```bash
git clone https://github.com/your-username/go-backend.git
cd go-backend
```

2. **Install dependencies:**
```bash
go mod download
go mod tidy
```

3. **Set up environment:**
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your preferred settings
# Default configuration works out of the box with SQLite
```

4. **Run the application:**
```bash
# Using Go directly
go run cmd/server/main.go

# Or using Make
make run
```

5. **Verify it's working:**
```bash
# Health check
curl http://localhost:8080/health

# Expected response: {"status":"healthy","message":"API is running","version":"1.0.0"}
```

### 🔧 Configuration Options

The application uses environment variables for configuration. Here are the key settings:

```env
# 🌐 Server Configuration
SERVER_PORT=8080                    # Port to run the server on
SERVER_HOST=localhost               # Host to bind the server to
NODE_ENV=development                # Environment (development/production)

# 🗄️ Database Configuration
DB_TYPE=sqlite                      # Database type (sqlite/postgres)
DB_HOST=localhost                   # Database host (for PostgreSQL)
DB_PORT=5432                       # Database port (for PostgreSQL)
DB_NAME=go_backend                 # Database name
DB_USER=your_user                  # Database username
DB_PASSWORD=your_password          # Database password
DB_SSL_MODE=disable               # SSL mode for PostgreSQL
SQLITE_PATH=./app.db              # SQLite database file path

# 🔐 JWT Configuration
JWT_SECRET=your-super-secret-key   # JWT signing secret (CHANGE IN PRODUCTION!)
JWT_EXPIRY=24h                    # Token expiration time

# 📊 Logging Configuration
LOG_LEVEL=info                    # Log level (debug/info/warn/error)
LOG_FORMAT=json                   # Log format (json/text)

# 🌍 CORS Configuration
CORS_ORIGINS=http://localhost:3000,http://localhost:8080  # Allowed origins
```

### 🐘 Using PostgreSQL

To use PostgreSQL instead of SQLite:

1. **Install PostgreSQL** and create a database
2. **Update your `.env` file:**
```env
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=go_backend
DB_USER=your_username
DB_PASSWORD=your_password
DB_SSL_MODE=disable
```

3. **Run the application** - it will automatically create tables

## 📚 API Documentation

### 🌐 Base URL
```
http://localhost:8080/api/v1
```

### 🔓 Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/ready` | Readiness check |
| `POST` | `/api/v1/auth/register` | Register new user |
| `POST` | `/api/v1/auth/login` | User login |

### 🔒 Protected Endpoints (Require Authentication)

| Method | Endpoint | Description | Role Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/v1/user/profile` | Get current user profile | Any |
| `PUT` | `/api/v1/user/profile` | Update current user profile | Any |
| `POST` | `/api/v1/user/change-password` | Change password | Any |

### 👑 Admin Endpoints (Admin Role Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/admin/users` | Get all users (paginated) |
| `GET` | `/api/v1/admin/users/:id` | Get user by ID |
| `PUT` | `/api/v1/admin/users/:id` | Update user |
| `DELETE` | `/api/v1/admin/users/:id` | Delete user |

### 🛡️ Moderator Endpoints (Admin or Moderator Role Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/mod/users` | Get all users (view only) |

## 🎯 API Examples

### 📝 User Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "username": "johndoe",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

**Response:**
```json
{
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "john.doe@example.com",
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user",
      "is_active": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  }
}
```

### 🔑 User Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepassword123"
  }'
```

### 👤 Get User Profile

```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 👥 Get All Users (Admin Only)

```bash
curl -X GET "http://localhost:8080/api/v1/admin/users?page=1&limit=10" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

## 👥 User Roles & Permissions

### 🔐 Role Hierarchy

```
👑 Admin
├── Full system access
├── User management (CRUD)
├── System configuration
└── All moderator permissions

🛡️ Moderator  
├── View all users
├── Content moderation
└── All user permissions

👤 User
├── Profile management
├── Password change
└── Basic API access
```

### 🚪 Default Admin Account

The application automatically creates a default admin account:

- **Email**: `admin@example.com`
- **Password**: `admin123`
- **Role**: `admin`

⚠️ **Security Warning**: Change the default admin password immediately in production!

## 🔧 Development

### 🛠️ Development Commands

```bash
# Install dependencies
make deps

# Run in development mode
### 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 📦 Building for Production

```bash
# Build optimized binary
go build -ldflags="-w -s" -o bin/server cmd/server/main.go

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go
GOOS=windows GOARCH=amd64 go build -o bin/server.exe cmd/server/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/server-mac cmd/server/main.go
```

## 🐳 Docker Deployment

### 🏗️ Building Docker Image

```bash
# Build the Docker image
docker build -t go-backend:latest .

# Run with Docker
docker run -p 8080:8080 --env-file .env go-backend:latest
```

### 🐙 Docker Compose (Optional)

Create a `docker-compose.yml` file:

```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_TYPE=postgres
      - DB_HOST=db
      - DB_PORT=5432
      - DB_NAME=go_backend
      - DB_USER=postgres
      - DB_PASSWORD=password
    depends_on:
      - db
    
  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=go_backend
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

```bash
# Run with Docker Compose
docker-compose up -d
```

## 🔐 Security Features

### 🛡️ Implemented Security Measures

- **🔒 Password Security**
  - Bcrypt hashing with salt
  - Minimum password requirements
  - Password change functionality

- **🎫 JWT Security**
  - Secure token generation
  - Token expiration handling
  - Token validation middleware

- **🌐 HTTP Security**
  - CORS protection with configurable origins
  - Security headers (XSS, CSRF, etc.)
  - Rate limiting to prevent abuse
  - Request size limiting

- **🗄️ Database Security**
  - SQL injection prevention (GORM)
  - Parameterized queries
  - Connection security

- **📝 Input Validation**
  - Request payload validation
  - Data sanitization
  - Error message standardization

### 🔧 Security Best Practices

1. **Change default credentials** in production
2. **Use strong JWT secrets** (minimum 32 characters)
3. **Enable HTTPS** in production
4. **Configure proper CORS** origins
5. **Monitor and log** security events
6. **Regular dependency updates**
7. **Environment variable protection**

## 📊 Monitoring & Observability

### 📈 Health Checks

```bash
# Basic health check
curl http://localhost:8080/health

# Readiness check (includes database)
curl http://localhost:8080/ready
```

### 📋 Logging

The application provides structured logging with different levels:

```json
{
  "level": "info",
  "msg": "HTTP Request",
  "time": "2024-01-01T12:00:00Z",
  "method": "GET",
  "path": "/api/v1/user/profile",
  "status": 200,
  "latency": "2.5ms",
  "ip": "127.0.0.1",
  "user_agent": "curl/7.68.0"
}
```

### 🎯 Metrics (Future Enhancement)

Consider adding:
- Prometheus metrics
- Performance monitoring
- Error rate tracking
- Database connection monitoring

## 🚨 Error Handling

### 📋 Standard Error Responses

All API endpoints return consistent error responses:

```json
{
  "error": "Validation failed",
  "errors": {
    "email": "Must be a valid email address",
    "password": "Must be at least 6 characters long"
  }
}
```

### 🔢 HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| `200` | Success |
| `201` | Created |
| `400` | Bad Request (validation errors) |
| `401` | Unauthorized (authentication required) |
| `403` | Forbidden (insufficient permissions) |
| `404` | Not Found |
| `429` | Too Many Requests (rate limited) |
| `500` | Internal Server Error |

## 🎛️ Environment Variables Reference

### 📋 Complete Configuration

```env
# ========================================
# 🌐 SERVER CONFIGURATION
# ========================================
SERVER_PORT=8080                          # HTTP server port
SERVER_HOST=localhost                     # Server host binding
NODE_ENV=development                      # Environment mode

# ========================================
# 🗄️ DATABASE CONFIGURATION  
# ========================================
DB_TYPE=sqlite                           # Database type (sqlite/postgres)
DB_HOST=localhost                        # Database host
DB_PORT=5432                            # Database port
DB_NAME=go_backend                      # Database name
DB_USER=your_username                   # Database username
DB_PASSWORD=your_password               # Database password
DB_SSL_MODE=disable                     # SSL mode (disable/require/verify-full)
SQLITE_PATH=./app.db                    # SQLite file path

# ========================================
# 🔐 JWT CONFIGURATION
# ========================================
JWT_SECRET=your-super-secret-jwt-key-min-32-chars    # JWT signing secret
JWT_EXPIRY=24h                                       # Token expiration

# ========================================
# 📊 LOGGING CONFIGURATION
# ========================================
LOG_LEVEL=info                          # Log level (debug/info/warn/error)
LOG_FORMAT=json                         # Log format (json/text)

# ========================================
# 🌍 CORS CONFIGURATION
# ========================================
CORS_ORIGINS=http://localhost:3000,http://localhost:8080    # Allowed origins
```

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### 🚀 Getting Started

1. **Fork** the repository
2. **Clone** your fork:
   ```bash
   git clone https://github.com/your-username/go-backend.git
   ```
3. **Create** a feature branch:
   ```bash
   git checkout -b feature/amazing-feature
   ```

### 📝 Development Guidelines

- **Code Style**: Follow Go conventions and use `gofmt`
- **Testing**: Write tests for new features
- **Documentation**: Update README and code comments
- **Commits**: Use conventional commit messages

### 🔍 Pull Request Process

1. **Update** documentation if needed
2. **Add tests** for new functionality
3. **Ensure** all tests pass
4. **Create** a detailed pull request description

### 🐛 Bug Reports

Please include:
- Go version
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Error logs (if any)

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2024 Go Backend API

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## 🙏 Acknowledgments

- **Gin Web Framework** - Fast HTTP web framework
- **GORM** - Fantastic ORM library for Go
- **JWT-Go** - Go implementation of JSON Web Tokens
- **Logrus** - Structured logger for Go
- **Go Community** - For excellent documentation and support

## 📞 Support

- **Documentation**: This README and inline code comments
- **Issues**: [GitHub Issues](https://github.com/your-username/go-backend/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/go-backend/discussions)

---

⭐ **Star this repository if you find it helpful!** ⭐

---

**Built with ❤️ using Go**