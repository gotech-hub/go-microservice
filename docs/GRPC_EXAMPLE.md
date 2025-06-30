# gRPC Example

Hướng dẫn sử dụng gRPC trong microservice này.

## Cấu trúc

```
api/grpc/
├── handlers/          # gRPC handlers
│   ├── handler.go     # Handler cũ
│   └── user_handler.go # UserService handler
├── middlewares/       # gRPC middlewares
│   └── logging.go     # Logging middleware
├── models/           # Internal models
│   └── user.go       # User model
├── proto/            # Protocol buffer definitions
│   ├── .proto        # Proto file cũ
│   └── user.proto    # User service proto
└── server.go         # gRPC server

examples/
└── grpc_client/      # gRPC client example
    └── main.go       # Client test code
```

## Setup

### 1. Install Dependencies

```bash
# Install protobuf compiler
brew install protobuf  # macOS
# hoặc
apt-get install protobuf-compiler  # Ubuntu

# Install Go protobuf plugins
make install-deps
```

### 2. Generate Protobuf Code

```bash
make proto-gen
```

Lệnh này sẽ:
- Tạo thư mục `api/grpc/proto/gen/user/`
- Generate các file Go từ `user.proto`
- Tạo `user.pb.go` và `user_grpc.pb.go`

### 3. Run Server

```bash
make run
```

Server sẽ chạy cả HTTP (port 8080) và gRPC (port 9090).

## API Endpoints

### UserService

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `CreateUser` | `CreateUserRequest` | `CreateUserResponse` | Tạo user mới |
| `GetUser` | `GetUserRequest` | `GetUserResponse` | Lấy thông tin user |
| `UpdateUser` | `UpdateUserRequest` | `UpdateUserResponse` | Cập nhật user |
| `DeleteUser` | `DeleteUserRequest` | `DeleteUserResponse` | Xóa user |
| `ListUsers` | `ListUsersRequest` | `ListUsersResponse` | Lấy danh sách users |
| `HealthCheck` | `HealthCheckRequest` | `HealthCheckResponse` | Health check |

### Message Types

#### User
```protobuf
message User {
  string id = 1;
  string name = 2;
  string email = 3;
  string phone = 4;
  UserStatus status = 5;
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp updated_at = 7;
}
```

#### UserStatus
```protobuf
enum UserStatus {
  USER_STATUS_UNSPECIFIED = 0;
  USER_STATUS_ACTIVE = 1;
  USER_STATUS_INACTIVE = 2;
  USER_STATUS_SUSPENDED = 3;
}
```

## Testing

### 1. Using gRPC Client Example

```bash
# Build client
go build -o bin/grpc-client examples/grpc_client/main.go

# Run client
./bin/grpc-client
```

### 2. Using grpcurl

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:9090 list

# List methods
grpcurl -plaintext localhost:9090 list user.UserService

# Call HealthCheck
grpcurl -plaintext localhost:9090 user.UserService/HealthCheck

# Call CreateUser
grpcurl -plaintext -d '{"name": "John Doe", "email": "john@example.com", "phone": "+1234567890"}' localhost:9090 user.UserService/CreateUser

# Call GetUser
grpcurl -plaintext -d '{"id": "test-user-id"}' localhost:9090 user.UserService/GetUser

# Call ListUsers
grpcurl -plaintext -d '{"page": 1, "limit": 10}' localhost:9090 user.UserService/ListUsers
```

### 3. Using BloomRPC

1. Download [BloomRPC](https://github.com/bloomrpc/bloomrpc)
2. Import proto file: `api/grpc/proto/user.proto`
3. Connect to `localhost:9090`
4. Test các method

## Middleware

### Logging Middleware

Tự động log tất cả gRPC requests với:
- Method name
- Request ID
- Duration
- Status code
- Error messages

### Custom Middleware

Để thêm middleware mới:

1. Tạo file trong `api/grpc/middlewares/`
2. Implement `grpc.UnaryServerInterceptor` hoặc `grpc.StreamServerInterceptor`
3. Đăng ký trong `api/grpc/server.go`

## Error Handling

gRPC sử dụng standard error codes:

- `OK` (0): Success
- `INVALID_ARGUMENT` (3): Invalid request parameters
- `NOT_FOUND` (5): Resource not found
- `INTERNAL` (13): Internal server error
- `UNAVAILABLE` (14): Service unavailable

## Configuration

### Environment Variables

Thêm vào `local.env`:

```env
GRPC_PORT=9090
GRPC_ENABLE_REFLECTION=true
```

### Config Structure

```go
type Config struct {
    // ... existing fields
    GRPCPort int `env:"GRPC_PORT" envDefault:"9090"`
    GRPCEnableReflection bool `env:"GRPC_ENABLE_REFLECTION" envDefault:"true"`
}
```

## Best Practices

1. **Validation**: Validate requests trong handlers
2. **Error Handling**: Sử dụng appropriate gRPC status codes
3. **Logging**: Log tất cả requests và responses
4. **Metrics**: Add Prometheus metrics cho gRPC calls
5. **Tracing**: Implement distributed tracing
6. **Rate Limiting**: Add rate limiting middleware
7. **Authentication**: Implement gRPC authentication

## Next Steps

1. **Database Integration**: Connect với PostgreSQL/MongoDB
2. **Authentication**: Add JWT authentication
3. **Rate Limiting**: Implement rate limiting
4. **Metrics**: Add Prometheus metrics
5. **Tracing**: Add OpenTelemetry tracing
6. **Load Balancing**: Implement gRPC load balancing
7. **Service Discovery**: Add service discovery 