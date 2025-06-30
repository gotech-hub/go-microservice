# 🚀 Go Microservice - Clean Architecture Template

[![Go Version](https://img.shields.io/badge/Go-1.23.2+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/your-repo)
[![Test Coverage](https://img.shields.io/badge/Coverage-85%25-brightgreen.svg)](https://github.com/your-repo)

A production-ready Go microservice template built with clean architecture principles, featuring HTTP/gRPC APIs, comprehensive testing, monitoring, and modern development practices.

## 📋 Table of Contents

- [Overview](#-overview)
- [Architecture](#-architecture)
- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Project Structure](#-project-structure)
- [Quick Start](#-quick-start)
- [Configuration](#-configuration)
- [API Documentation](#-api-documentation)
- [Development](#-development)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)

## 🎯 Overview

This microservice template provides a robust foundation for building scalable Go services with:

- **Clean Architecture**: Separation of concerns with layers (handlers, services, repositories)
- **Dual API Support**: HTTP REST API and gRPC API
- **Database Support**: PostgreSQL, MongoDB, Redis, ClickHouse
- **Message Queue**: Kafka integration for event-driven architecture
- **Monitoring**: Prometheus metrics and structured logging
- **Testing**: Unit tests, integration tests, and mock generation
- **Security**: JWT authentication, rate limiting, and encryption utilities

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP/gRPC     │    │   Middleware    │    │   Handlers      │
│     APIs        │───▶│   (Auth, Log,   │───▶│   (Controllers) │
└─────────────────┘    │    Rate Limit)  │    └─────────────────┘
                       └─────────────────┘              │
                                                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Databases     │◀───│   Repositories  │◀───│    Services     │
│ (PostgreSQL,    │    │   (Data Access) │    │  (Business      │
│  MongoDB, etc.) │    └─────────────────┘    │    Logic)       │
└─────────────────┘                           └─────────────────┘
```

### Key Components:

- **API Layer**: HTTP REST and gRPC endpoints
- **Service Layer**: Business logic implementation
- **Repository Layer**: Data access abstraction
- **Domain Layer**: Core business entities
- **Infrastructure**: Database connections, external services

## ✨ Features

### 🔧 Core Features
- **Multi-Protocol APIs**: HTTP REST + gRPC support
- **Database Abstraction**: Support for multiple databases
- **Event-Driven**: Kafka message queue integration
- **Caching**: Redis caching layer
- **Authentication**: JWT-based authentication
- **Rate Limiting**: Request rate limiting middleware
- **Health Checks**: Built-in health check endpoints

### 📊 Monitoring & Observability
- **Structured Logging**: Zerolog with context support
- **Metrics**: Prometheus metrics collection
- **Tracing**: OpenTelemetry integration
- **Health Monitoring**: Service health status

### 🛡️ Security
- **JWT Authentication**: Secure token-based auth
- **Encryption**: AES, RSA, EdDSA utilities
- **Input Validation**: Request validation middleware
- **Rate Limiting**: DDoS protection

### 🧪 Testing
- **Unit Tests**: Comprehensive unit test coverage
- **Integration Tests**: End-to-end testing
- **Mock Generation**: Automated mock creation
- **Test Coverage**: Coverage reporting and validation

## 🛠️ Tech Stack

### Core Framework
- **Go 1.23.2+**: Latest stable Go version
- **Echo**: High-performance HTTP framework
- **gRPC**: High-performance RPC framework

### Databases
- **PostgreSQL**: Primary relational database
- **MongoDB**: Document database
- **Redis**: In-memory cache and session store
- **ClickHouse**: Analytics database

### Message Queue
- **Apache Kafka**: Event streaming platform

### Monitoring & Logging
- **Zerolog**: Structured logging
- **Prometheus**: Metrics collection
- **OpenTelemetry**: Distributed tracing

### Development Tools
- **Mockgen**: Mock generation
- **Swaggo**: API documentation
- **Gosec**: Security analysis
- **GolangCI-Lint**: Code linting

## 📁 Project Structure

```
go-microservice/
├── api/                    # API layer
│   ├── graph/             # GraphQL (future)
│   ├── grpc/              # gRPC server and handlers
│   ├── http/              # HTTP server and handlers
│   └── msg/               # Message queue handlers
├── app/                   # Application entry point
├── bootstrap/             # Dependency injection
├── config/                # Configuration management
├── internal/              # Internal application code
│   ├── domains/           # Domain entities
│   └── services/          # Business logic services
├── pkg/                   # Shared packages
│   ├── adapters/          # External service adapters
│   ├── client/            # HTTP/gRPC clients
│   ├── database/          # Database connections
│   ├── event_bus/         # Event bus implementation
│   ├── jwt/               # JWT utilities
│   ├── log/               # Logging utilities
│   ├── metric/            # Metrics collection
│   ├── middlewares/       # HTTP/gRPC middlewares
│   ├── queue/             # Message queue utilities
│   ├── utils/             # Common utilities
│   └── validate/          # Validation utilities
├── repositories/          # Data access layer
├── test/                  # Test files
│   ├── integration/       # Integration tests
│   └── mocks/             # Generated mocks
└── docs/                  # Documentation
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.23.2+**: [Download Go](https://golang.org/dl/)
- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Make**: Available on most Unix systems
- **Protobuf Compiler**: For gRPC development

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/go-microservice.git
   cd go-microservice
   ```

2. **Install dependencies**
   ```bash
   make install-deps
   ```

3. **Setup environment**
   ```bash
   cp local.env.example local.env
   # Edit local.env with your configuration
   ```

4. **Generate protobuf code**
   ```bash
   make proto-gen
   ```

5. **Run the application**
   ```bash
   make run
   ```

### Docker Setup

```bash
# Build Docker image
make docker-build

# Run with Docker
make docker-run
```

## ⚙️ Configuration

### Environment Variables

Create a `local.env` file based on `local.env.example`:

```bash
# Service Configuration
SERVICE_NAME=user-service
SERVICE_VERSION=1.0.0
ENV=development
HTTP_PORT=3000

# Database Configuration
MONGO_DB_URI=mongodb://localhost:27017
REDIS_ADDR=localhost:6379
POSTGRES_DSN=postgres://user:password@localhost:5432/dbname

# Kafka Configuration
KAFKA_BROKERS=localhost:9092

# JWT Configuration
JWT_SECRET=your-secret-key
```

### Configuration Management

The application uses environment-based configuration with automatic validation:

```go
type SystemConfig struct {
    Env            string
    HttpPort       uint64
    ServiceName    string
    ServiceVersion string
    MongoDBConfig  mongodb.MongoDBConfig
    RedisConfig    redis.RedisConfig
}
```

## 📚 API Documentation

### HTTP API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |
| POST | `/api/v1/users` | Create user |
| GET | `/api/v1/users/:id` | Get user by ID |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Delete user |
| GET | `/api/v1/users` | List users with pagination |

### gRPC API

The service exposes gRPC endpoints on port `9090`:

```protobuf
service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
}
```

### API Testing

Test gRPC endpoints:
```bash
./script/test_grpc.sh
```

## 🛠️ Development

### Available Commands

```bash
# Run application
make run

# Run in debug mode
make debug

# Build application
make build

# Clean build artifacts
make clean

# Install dependencies
make install-deps

# Generate protobuf code
make proto-gen

# Generate mocks
make mock-gen

# Generate Swagger docs
make swagger-gen

# Run security analysis
make gosec
```

### Code Generation

The project includes several code generation tools:

- **Protobuf**: Generate gRPC code from `.proto` files
- **Mocks**: Generate mock implementations for testing
- **Swagger**: Generate API documentation
- **Repository**: Generate repository interfaces

### Code Quality

```bash
# Run linter
make lint-test

# Run all tests
make all-test

# Run with coverage
make cover-test
```

## 🧪 Testing

### Test Structure

```
test/
├── integration/           # Integration tests
├── mocks/                # Generated mocks
└── utils/                # Test utilities
```

### Running Tests

```bash
# Unit tests only
make unit-test

# Integration tests
make integration-test

# All tests with coverage
make cover-test

# Run in containerized environment
./coverage.sh
```

### Test Coverage

The project maintains high test coverage with automated reporting:

```bash
# Generate coverage report
make cover-test

# View coverage in browser
go tool cover -html=reports/coverage.out
```

## 🚀 Deployment

### Docker Deployment

```bash
# Build production image
docker build -t go-microservice .

# Run with environment variables
docker run -p 3000:3000 -p 9090:9090 \
  -e SERVICE_NAME=user-service \
  -e ENV=production \
  go-microservice
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-microservice
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-microservice
  template:
    metadata:
      labels:
        app: go-microservice
    spec:
      containers:
      - name: go-microservice
        image: go-microservice:latest
        ports:
        - containerPort: 3000
        - containerPort: 9090
        env:
        - name: ENV
          value: "production"
```

### Environment-Specific Configurations

- **Development**: Local development with hot reload
- **Staging**: Pre-production testing environment
- **Production**: Production deployment with monitoring

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### Development Workflow

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Follow coding conventions**: See [CONVENTION.md](CONVENTION.md)
4. **Write tests**: Ensure all new code is tested
5. **Commit changes**: Use conventional commit messages
6. **Push to branch**: `git push origin feature/amazing-feature`
7. **Create Pull Request**: Provide detailed description

### Code Standards

- Follow Go best practices and idioms
- Maintain test coverage above 80%
- Use structured logging for all operations
- Document public APIs and interfaces
- Follow the naming conventions in [CONVENTION.md](CONVENTION.md)

### Commit Message Format

```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Contributors

- **Mai Công Trình** – [trinhbentre2013@gmail.com](mailto:trinhbentre2013@gmail.com)
- **Lương Công Văn** – [@gmail.com](mailto:@gmail.com)
- **Bùi Quốc Đạt** – [datbq.work@gmail.com](mailto:datbq.work@gmail.com)
- **Nguyễn Tiến Dũng** – [ntdung.it.2912@gmail.com](mailto:ntdung.it.2912@gmail.com)

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-repo/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-repo/discussions)
- **Documentation**: [Wiki](https://github.com/your-repo/wiki)

---

<div align="center">
Made with ❤️ by the Go Microservice Team
</div>