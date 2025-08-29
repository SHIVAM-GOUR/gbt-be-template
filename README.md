# GBT Backend Template

A production-ready Golang backend template following clean architecture principles. This template provides a solid foundation for building scalable web applications with modern development practices.

## 🚀 Features

- **Clean Architecture**: Modular design with clear separation of concerns
- **REST API**: Complete CRUD operations with Chi router
- **Authentication**: JWT-based authentication with login/logout
- **Database**: PostgreSQL with GORM ORM and migrations
- **Middleware**: CORS, logging, recovery, authentication, and authorization
- **Configuration**: Environment-based configuration management
- **Logging**: Structured logging with Logrus
- **Testing**: Comprehensive unit tests with mocks
- **Docker**: Multi-stage Dockerfile and docker-compose setup
- **Hot Reload**: Air integration for development
- **Graceful Shutdown**: Context-based shutdown handling

## 📁 Project Structure

```
├── cmd/app/                 # Application entrypoint
├── internal/
│   ├── config/             # Configuration management
│   ├── handlers/           # HTTP handlers (controllers)
│   ├── models/             # Data models and DTOs
│   ├── repository/         # Database access layer
│   ├── routes/             # Route definitions
│   ├── server/             # HTTP server setup
│   └── services/           # Business logic layer
├── pkg/
│   ├── logger/             # Centralized logging
│   ├── middleware/         # Reusable middleware
│   └── utils/              # Helper utilities
├── migrations/             # Database migrations
├── .air.toml              # Air configuration
├── docker-compose.yml     # Docker compose setup
├── Dockerfile             # Multi-stage Docker build
├── Makefile              # Development commands
└── README.md             # This file
```

## 🛠️ Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker and Docker Compose (optional)
- Make (optional, for using Makefile commands)

## 🚀 Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd gbt-be-template
cp .env.example .env
```

### 2. Install Dependencies

```bash
make deps
```

### 3. Install Development Tools

```bash
make tools
```

### 4. Setup Database

#### Using Docker (Recommended)
```bash
make docker-run
```

#### Manual Setup
```bash
# Create database
createdb gbt_template

# Run migrations
make migrate-up
```

### 5. Run the Application

#### Development (with hot reload)
```bash
make dev
```

#### Production build
```bash
make build
make run
```

## 🐳 Docker Usage

### Build and Run with Docker Compose
```bash
make docker-build
make docker-run
```

### Stop Services
```bash
make docker-stop
```

### View Logs
```bash
make docker-logs
```

## 📊 Database Migrations

### Create New Migration
```bash
make migrate-create name=add_new_table
```

### Run Migrations
```bash
make migrate-up
```

### Rollback Migrations
```bash
make migrate-down
```

## 🧪 Testing

### Run All Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test-coverage
```

## 📡 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout (requires auth)
- `GET /api/v1/auth/profile` - Get user profile (requires auth)

### Users
- `GET /api/v1/users` - List users (requires auth)
- `GET /api/v1/users/{id}` - Get user by ID (requires auth)
- `PUT /api/v1/users/{id}` - Update user (requires auth)
- `DELETE /api/v1/users/{id}` - Delete user (requires auth)

### Admin
- `POST /api/v1/admin/users` - Create user (admin only)

### Health Checks
- `GET /health` - Health check
- `GET /health/ready` - Readiness check
- `GET /health/live` - Liveness check

## 🔧 Configuration

Configuration is managed through environment variables. Copy `.env.example` to `.env` and modify as needed:

```bash
# Server Configuration
PORT=8080
HOST=localhost
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=gbt_template

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRY=24h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 🔐 Authentication

The API uses JWT tokens for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Example Login Request
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

## 📝 API Examples

### Register User
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

### Get User Profile
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer <your-jwt-token>"
```

## 🛠️ Development

### Available Make Commands
```bash
make help          # Show available commands
make build          # Build the application
make run            # Build and run
make dev            # Run with hot reload
make test           # Run tests
make test-coverage  # Run tests with coverage
make clean          # Clean build artifacts
make fmt            # Format code
make lint           # Run linter
```

### Code Structure Guidelines

1. **Models**: Define in `internal/models/` with GORM tags
2. **Repositories**: Implement interfaces in `internal/repository/`
3. **Services**: Business logic in `internal/services/`
4. **Handlers**: HTTP handlers in `internal/handlers/`
5. **Middleware**: Reusable middleware in `pkg/middleware/`
6. **Utils**: Helper functions in `pkg/utils/`

## 🧪 Testing Strategy

- **Unit Tests**: Test individual components in isolation
- **Mocks**: Use testify/mock for dependencies
- **Integration Tests**: Test API endpoints end-to-end
- **Coverage**: Aim for >80% test coverage

## 🚀 Deployment

### Environment Variables for Production
- Set `ENV=production`
- Use strong `JWT_SECRET`
- Configure proper database credentials
- Set appropriate `LOG_LEVEL`

### Docker Deployment
```bash
docker build -t gbt-be-template .
docker run -p 8080:8080 --env-file .env gbt-be-template
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the example API calls

---

**Happy Coding! 🎉**
