# 🚀 Go Backend API

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/chmuhammadasim/Go_Temp)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](Dockerfile)
[![SQLite](https://img.shields.io/badge/Database-SQLite### 🔄 Database Migrations
Database migrations run automatically on startup. The system creates:
- User tables with roles and permissions
- Post and comment tables with relationships
- Audit logging tables
- File management tables
- Session management tables

## 🧪 Testing

### 🔍 Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -cover ./...

# Run tests with detailed coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Using Makefile
make test
make test-coverage
```

### 📊 Test Structure
```
tests/
├── api_test.go          # API integration tests
├── unit/               # Unit tests
│   ├── services/       # Service layer tests
│   ├── handlers/       # Handler tests
│   └── utils/          # Utility tests
└── fixtures/           # Test data and fixtures
```

## 🔒 Security Features

### 🛡️ Authentication & Authorization
- **JWT Tokens**: Secure stateless authentication
- **Role-based Access Control**: Admin, Moderator, User roles
- **Password Security**: Bcrypt hashing with salt rounds
- **Account Protection**: Lockout after failed login attempts
- **Session Management**: Secure session handling with Redis

### 🔐 Security Middleware
- **CORS Protection**: Configurable allowed origins
- **Rate Limiting**: Prevent API abuse
- **Security Headers**: XSS protection, content type sniffing prevention
- **Request Validation**: Input sanitization and validation
- **Audit Logging**: Track all user actions and system events

### 🔧 Security Best Practices
- Environment-based configuration
- Secure default settings
- Regular dependency updates
- Input validation and sanitization
- Error handling without information disclosure

## 📁 Project Structure

```
go-backend/
├── 📁 cmd/                         # 🚀 Application entry points
│   └── 📁 server/
│       └── 📄 main.go              # 🎯 Main application entry point
├── 📁 internal/                    # 🔒 Private application code
│   ├── 📁 config/
│   │   └── 📄 config.go            # ⚙️ Configuration management & environment vars
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
│   │   ├── 📄 user.go              # 📋 User data models & validation
│   │   └── 📄 extended_models.go   # 📋 Additional models (posts, comments, files)
│   ├── 📁 services/
│   │   ├── 📄 user_service.go      # 👤 User business logic
│   │   ├── 📄 post_service.go      # 📝 Post management service
│   │   ├── 📄 notification_service.go # 🔔 Multi-channel notifications
│   │   ├── 📄 audit_service.go     # 📊 Audit logging service
│   │   ├── 📄 file_service.go      # 📁 File management service
│   │   ├── 📄 email_service.go     # 📧 Email service
│   │   ├── 📄 cache_service.go     # 🚀 Caching service
│   │   └── 📄 security_service.go  # 🔒 Security utilities
│   └── 📁 utils/
│       ├── 📄 jwt.go               # 🎫 JWT token utilities & validation
│       └── 📄 validation.go        # ✅ Request validation helpers
├── 📁 pkg/                         # 📦 Reusable packages
│   └── 📁 logger/
│       └── 📄 logger.go            # 📊 Structured logging utilities
├── 📁 tests/                       # 🧪 Test files
│   ├── 📄 api_test.go              # 🔍 API integration tests
│   └── 📁 fixtures/                # 📋 Test data
├── 📁 docs/                        # 📚 Documentation
│   └── 📄 swagger.go               # 📖 API documentation
├── 📄 .env                         # 🔧 Environment variables (local)
├── 📄 .env.example                 # 📝 Environment variables template
├── 📄 go.mod                       # 📋 Go module definition & dependencies
├── 📄 go.sum                       # 🔒 Dependency lock file
├── 📄 Dockerfile                   # 🐳 Docker container definition
├── 📄 Makefile                     # 🔨 Build automation & development tasks
├── 📄 .gitignore                   # 🚫 Git ignore rules
└── 📄 README.md                    # 📖 Project documentation (this file)
```greSQL-blue.svg)](https://sqlite.org)

A **production-ready**, enterprise-grade Go backend API with comprehensive authentication, authorization, file management, and audit logging. Built with modern Go technologies including Gin, GORM, Redis, and JWT for maximum performance and security.

## ✨ Features

### 🔐 Authentication & Security
- **JWT-based authentication** with secure token generation and validation
- **Role-based access control** (Admin, Moderator, User) with granular permissions
- **Two-factor authentication (2FA)** for enhanced security
- **Password hashing** with bcrypt for maximum security
- **Session management** with Redis backend
- **Account lockout protection** against brute force attacks
- **Security middleware** with rate limiting and CORS protection
- **Audit logging** for all user activities and system events

### 🗄️ Database & Storage
- **Multi-database support** (SQLite with pure Go driver, PostgreSQL)
- **GORM ORM** with automatic migrations and relationships
- **Soft deletes** for data integrity and recovery
- **Database seeding** with default admin account
- **Connection pooling** and timeout management
- **Transaction support** for data consistency
- **Cache layer** with Redis for improved performance

### 📂 File Management
- **File upload/download** with validation and security checks
- **Multiple storage backends** (local filesystem, cloud storage ready)
- **File type validation** and size limits
- **Organized file categorization** and metadata tracking
- **Secure file serving** with access control

### 🌐 API Features
- **RESTful API design** following OpenAPI/Swagger standards
- **Request validation** with comprehensive error messages
- **Structured logging** with Logrus (JSON/Text formats)
- **Global error handling** middleware with detailed responses
- **CORS support** with configurable origins
- **Rate limiting** to prevent API abuse
- **Health checks** for monitoring and deployment readiness
- **API documentation** with Swagger/OpenAPI

### 🔧 Advanced Services
- **Notification system** with multiple channels (email, SMS, push)
- **Template engine** for dynamic content generation
- **CRUD service layer** for rapid development
- **User management** with profile customization
- **Post and comment system** with moderation capabilities
- **Email service** with SMTP configuration
- **Background job processing** ready for integration

### 🏗️ Architecture & Code Quality
- **Clean architecture** with clear separation of concerns
- **Dependency injection** for testability and maintainability
- **Service layer pattern** for business logic separation
- **Repository pattern** for data access abstraction
- **Middleware pipeline** for request processing
- **Environment-based configuration** management
- **Graceful shutdown** handling for production deployments
- **Comprehensive error handling** with context propagation
- **Binary compilation** for different platforms

## 📁 Project Structure

```
go_temp/
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

## 🛠️ Technology Stack

| Technology | Purpose | Version |
|------------|---------|---------|
| **Go** | Backend Language | 1.23+ |
| **Gin** | HTTP Web Framework | v1.9.1 |
| **GORM** | ORM Library | v1.30.0 |
| **JWT** | Authentication | v4.5.0 |
| **SQLite** | Database (Pure Go) | v1.11.0 |
| **PostgreSQL** | Production Database | v1.5.4 |
| **Redis** | Caching & Sessions | v9.14.0 |
| **Logrus** | Structured Logging | v1.9.3 |
| **Bcrypt** | Password Hashing | Built-in |
| **Docker** | Containerization | Multi-stage |
| **Swagger** | API Documentation | Built-in |
| **Docker** | Containerization | Any |

## 🚀 Getting Started

### 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.23 or higher** - [Download here](https://golang.org/dl/)
- **Git** - For cloning the repository
- **PostgreSQL** (optional) - SQLite is used by default
- **Redis** (optional) - For caching and sessions
- **Docker** (optional) - For containerized deployment

### ⚡ Quick Start

1. **Clone the repository:**
```bash
git clone https://github.com/chmuhammadasim/Go_Temp.git
cd Go_Temp
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

# Edit the .env file with your configuration
# The default configuration uses SQLite and works out of the box
```

4. **Run the application:**
```bash
# Using Go directly
go run cmd/server/main.go

# Or using Makefile
make run

# Or build and run
make build
./bin/server
```

5. **Verify installation:**
```bash
# Health check
curl http://localhost:8080/health

# API documentation (if available)
curl http://localhost:8080/api/docs
```

### 🔧 Environment Configuration

Create a `.env` file based on `.env.example`:

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration
DB_TYPE=sqlite                    # or "postgres"
SQLITE_PATH=./app.db             # SQLite file path

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h

# Logging Configuration
LOG_LEVEL=info                   # debug, info, warn, error
LOG_FORMAT=json                  # json or text

# CORS Configuration
CORS_ORIGINS=http://localhost:3000,http://localhost:8080
```

### 👤 Default Admin Account

After the first run, a default admin account is automatically created:

- **Email**: admin@example.com
- **Password**: admin123
- **Role**: admin

⚠️ **Important**: Change this password immediately in production!

## 📚 API Documentation

### 🔍 Base URL
```
http://localhost:8080
```

### 🏥 Health Endpoints
| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| GET | `/health` | Basic health check | No |
| GET | `/ready` | Readiness probe | No |

### 🔐 Authentication Endpoints
| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| POST | `/api/auth/register` | User registration | No |
| POST | `/api/auth/login` | User login | No |
| POST | `/api/auth/refresh` | Refresh JWT token | JWT |
| POST | `/api/auth/logout` | User logout | JWT |

### 👤 User Management Endpoints
| Method | Endpoint | Description | Authentication | Permission |
|--------|----------|-------------|----------------|------------|
| GET | `/api/users/profile` | Get current user profile | JWT | User |
| PUT | `/api/users/profile` | Update current user profile | JWT | User |
| GET | `/api/users` | List all users | JWT | Admin |
| GET | `/api/users/:id` | Get user by ID | JWT | Admin/Owner |
| PUT | `/api/users/:id` | Update user | JWT | Admin/Owner |
| DELETE | `/api/users/:id` | Delete user | JWT | Admin |
| POST | `/api/users/:id/change-password` | Change password | JWT | Admin/Owner |

### 📝 Post Management Endpoints
| Method | Endpoint | Description | Authentication | Permission |
|--------|----------|-------------|----------------|------------|
| GET | `/api/posts` | List all posts | JWT | User |
| POST | `/api/posts` | Create new post | JWT | User |
| GET | `/api/posts/:id` | Get post by ID | JWT | User |
| PUT | `/api/posts/:id` | Update post | JWT | Admin/Owner |
| DELETE | `/api/posts/:id` | Delete post | JWT | Admin/Owner |

### 💬 Comment Management Endpoints
| Method | Endpoint | Description | Authentication | Permission |
|--------|----------|-------------|----------------|------------|
| GET | `/api/posts/:post_id/comments` | Get post comments | JWT | User |
| POST | `/api/posts/:post_id/comments` | Create comment | JWT | User |
| PUT | `/api/comments/:id` | Update comment | JWT | Admin/Owner |
| DELETE | `/api/comments/:id` | Delete comment | JWT | Admin/Owner |

### 📄 Example API Usage

#### Register a new user
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "newuser",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

#### Get user profile (with JWT token)
```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Create a new post
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first post",
    "published": true
  }'
```
```

### 🔧 Configuration Options

The application uses environment variables for configuration. Here are the key settings:

```env
# 🌐 Server Configuration
SERVER_PORT=8080                    # Port to run the server on
SERVER_HOST=localhost               # Host to bind the server to
NODE_ENV=development                # Environment (development/production)

## � Docker Deployment

### 🚀 Quick Docker Setup

1. **Build the Docker image:**
```bash
docker build -t go-backend .
```

2. **Run with Docker:**
```bash
docker run -p 8080:8080 --env-file .env go-backend
```

3. **Using Docker Compose (recommended):**
```yaml
# docker-compose.yml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_TYPE=postgres
      - DB_HOST=db
      - DB_NAME=go_backend
      - DB_USER=postgres
      - DB_PASSWORD=password
    depends_on:
      - db
      - redis

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: go_backend
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

4. **Run with Docker Compose:**
```bash
docker-compose up -d
```

### �️ Development Commands

The project includes a comprehensive Makefile for common tasks:

```bash
# Build the application
make build

# Run the application
make run

# Clean build artifacts
make clean

# Install dependencies
make deps

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Build Docker image
make docker-build

# Run with Docker
make docker-run
```

## �️ Database Configuration

### 📊 SQLite (Default)
The application uses SQLite by default with a pure Go driver (no CGO required):

```env
DB_TYPE=sqlite
SQLITE_PATH=./app.db
```

### 🐘 PostgreSQL Configuration
For production use with PostgreSQL:

```env
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=go_backend
DB_USER=your_user
DB_PASSWORD=your_password
DB_SSL_MODE=disable
```

### � Database Migrations
Database migrations run automatically on startup. The system creates:
- User tables with roles and permissions
- Post and comment tables with relationships
- Audit logging tables
- File management tables
- Session management tables
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
## 🔧 Troubleshooting

### ❌ Common Issues

#### CGO/SQLite Issues
If you encounter CGO-related errors:
```bash
# Error: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work
# Solution: The project uses pure Go SQLite driver, ensure you have latest dependencies
go mod tidy
```

#### Port Already in Use
```bash
# Error: bind: address already in use
# Solution: Change the port or kill the process using the port
lsof -ti:8080 | xargs kill -9  # Kill process on port 8080
# Or change SERVER_PORT in .env file
```

#### Database Connection Issues
```bash
# Check database configuration
cat .env | grep DB_

# For SQLite: Ensure write permissions
chmod 755 .
touch app.db
chmod 644 app.db

# For PostgreSQL: Test connection
psql -h localhost -U your_user -d go_backend
```

#### JWT Token Issues
```bash
# Ensure JWT_SECRET is set and sufficiently long
echo $JWT_SECRET
# Should be at least 32 characters for security
```

### 📊 Monitoring & Health Checks

```bash
# Health check endpoint
curl http://localhost:8080/health

# Readiness check (for Kubernetes)
curl http://localhost:8080/ready

# Check database connectivity
curl http://localhost:8080/health/db
```

### 🔍 Debugging

Enable debug logging in `.env`:
```env
LOG_LEVEL=debug
LOG_FORMAT=text
```

View detailed logs:
```bash
go run cmd/server/main.go 2>&1 | tee app.log
```

## 📈 Performance & Optimization

### 🚀 Performance Features
- **Connection Pooling**: Database connection management
- **Redis Caching**: Session and data caching
- **Middleware Optimization**: Efficient request processing
- **Graceful Shutdown**: Clean resource disposal
- **Memory Management**: Optimized Go garbage collection

### � Monitoring Recommendations
- Use **Prometheus** for metrics collection
- Implement **Health checks** for load balancers
- Monitor **Database performance** and query optimization
- Track **API response times** and error rates
- Set up **Log aggregation** with ELK stack or similar

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### 🔄 Development Workflow

1. **Fork the repository**
2. **Create a feature branch:**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Make your changes and add tests**
4. **Run tests and linting:**
   ```bash
   make test
   make lint
   make fmt
   ```
5. **Commit your changes:**
   ```bash
   git commit -m "Add amazing feature"
   ```
6. **Push to the branch:**
   ```bash
   git push origin feature/amazing-feature
   ```
7. **Open a Pull Request**

### 📋 Code Standards
- Follow **Go best practices** and **effective Go** guidelines
- Write **comprehensive tests** for new features
- Update **documentation** for API changes
- Use **conventional commits** for clear history
- Ensure **security best practices** in code

### 🐛 Bug Reports
When reporting bugs, please include:
- Go version and OS
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs and error messages
- Configuration details (without secrets)

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## � Acknowledgments

- **Gin Framework** - Fast HTTP web framework
- **GORM** - Fantastic ORM library for Go  
- **JWT-Go** - Go implementation of JSON Web Tokens
- **Logrus** - Structured logging for Go
- **Go Community** - For excellent tools and libraries

## 📞 Support & Contact

- **GitHub Issues**: [Create an issue](https://github.com/chmuhammadasim/Go_Temp/issues)
- **Email**: chmuhammadasim@gmail.com
- **Repository**: [https://github.com/chmuhammadasim/Go_Temp](https://github.com/chmuhammadasim/Go_Temp)

---

⭐ **Star this repository** if you find it helpful!

Built with ❤️ by [Muhammad Asim](https://github.com/chmuhammadasim)

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
docker build -t go_temp:latest .

# Run with Docker
docker run -p 8080:8080 --env-file .env go_temp:latest
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
   git clone https://github.com/chmuhammadasim/go_temp.git
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
- **Issues**: [GitHub Issues](https://github.com/chmuhammadasim/go_temp/issues)
- **Discussions**: [GitHub Discussions](https://github.com/chmuhammadasim/go_temp/discussions)

---

⭐ **Star this repository if you find it helpful!** ⭐

---

**Built with ❤️ using Go**