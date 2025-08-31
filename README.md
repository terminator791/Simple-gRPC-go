# Simple gRPC Go Microservices

A complete, production-ready microservices architecture implementation demonstrating best practices for inter-service communication using gRPC in Go. This project showcases two distinct microservices that communicate exclusively via gRPC protocols.

## 🏗️ Architecture Overview

The system consists of two independent microservices:

1. **User Service** - Manages user data with PostgreSQL storage
2. **Order Service** - Manages orders with user validation via gRPC calls

### Core Communication Flow

```mermaid
graph LR
    Client[Client] --> OrderService[Order Service :50052]
    OrderService --> UserService[User Service :50051]
    OrderService --> OrderDB[(Order DB :5433)]
    UserService --> UserDB[(User DB :5432)]
```

**Key Flow:**
1. Client requests order creation from Order Service
2. Order Service validates user ID by calling User Service via gRPC
3. If user exists, order is created; if not, request is rejected
4. Both services maintain their own PostgreSQL databases

## 📋 Features

- ✅ **gRPC Communication**: Pure gRPC inter-service communication
- ✅ **Database Isolation**: Each service has its own PostgreSQL database
- ✅ **User Validation**: Orders validated against User Service via gRPC
- ✅ **Error Handling**: Proper gRPC status codes (NOT_FOUND, etc.)
- ✅ **Unit Testing**: Comprehensive tests with mocked gRPC clients
- ✅ **Docker Setup**: PostgreSQL databases via Docker Compose
- ✅ **Protocol Buffers**: Type-safe API definitions
- ✅ **Production Ready**: Structured logging, graceful shutdowns

## 🛠️ Technology Stack

- **Language**: Go 1.23+
- **RPC Framework**: gRPC (`google.golang.org/grpc`)
- **Interface Definition**: Protocol Buffers v3
- **Database**: PostgreSQL 15 Alpine
- **Database Client**: `jmoiron/sqlx`
- **Testing**: Go testing + `stretchr/testify`
- **Containerization**: Docker & Docker Compose

## 📁 Project Structure

```
Simple-gRPC-go/
├── proto/                      # Protocol Buffer definitions
│   ├── user/user.proto         # User service interface
│   └── order/order.proto       # Order service interface
├── pkg/                        # Generated protobuf code
│   ├── user/                   # Generated user service code
│   └── order/                  # Generated order service code
├── user-service/               # User microservice
│   ├── cmd/main.go            # User service entry point
│   └── internal/
│       ├── config/            # Configuration management
│       ├── db/                # Database layer
│       ├── handlers/          # gRPC handlers
│       └── models/            # Data models
├── order-service/             # Order microservice
│   ├── cmd/main.go           # Order service entry point
│   └── internal/
│       ├── config/           # Configuration management
│       ├── db/               # Database layer
│       ├── handlers/         # gRPC handlers
│       ├── client/           # gRPC client for user service
│       └── models/           # Data models
├── scripts/                  # Helper scripts
├── init-scripts/            # Database initialization
├── docker-compose.yml       # Database containers
└── Makefile                # Build automation
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or later
- Docker and Docker Compose
- Protocol Buffers compiler (`protoc`)
- Make

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/terminator791/Simple-gRPC-go.git
   cd Simple-gRPC-go
   ```

2. **Install dependencies:**
   ```bash
   make deps
   ```

3. **Generate Protocol Buffer code:**
   ```bash
   make proto
   ```

4. **Build services:**
   ```bash
   make build
   ```

5. **Start databases:**
   ```bash
   make db-up
   ```

### Running the Services

**Option 1: Using Make commands**
```bash
# Terminal 1: Start User Service
make user-service
./user-service/bin/user-service

# Terminal 2: Start Order Service  
make order-service
./order-service/bin/order-service
```

**Option 2: Using helper script**
```bash
./scripts/start-services.sh
```

### Testing the System

**Run unit tests:**
```bash
make test
```

**Test with grpcurl (requires services running):**
```bash
# Install grpcurl if not available
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Run automated tests
./scripts/test-services.sh
```

**Manual testing examples:**
```bash
# Create a user
grpcurl -plaintext -d '{"email": "john@example.com", "name": "John Doe"}' \
  localhost:50051 user.UserService/CreateUser

# Get user (should return user data)
grpcurl -plaintext -d '{"user_id": 1}' \
  localhost:50051 user.UserService/GetUser

# Create order for existing user (should succeed)
grpcurl -plaintext -d '{"user_id": 1, "product_name": "Laptop", "amount": 999.99, "quantity": 1}' \
  localhost:50052 order.OrderService/CreateOrder

# Try to create order for non-existent user (should fail)
grpcurl -plaintext -d '{"user_id": 999, "product_name": "Phone", "amount": 599.99, "quantity": 1}' \
  localhost:50052 order.OrderService/CreateOrder
```

## 🧪 Testing Strategy

### Unit Tests

The project includes comprehensive unit tests that demonstrate the core requirement: **order creation with mocked user service validation**.

**Key Test: `TestCreateOrder_UserNotFound`**
```go
// This test verifies that when the User Service returns NOT_FOUND,
// the Order Service properly rejects the order creation
func TestCreateOrder_UserNotFound(t *testing.T) {
    // Mock user service to return NOT_FOUND
    userNotFoundErr := status.Error(codes.NotFound, "user not found")
    mockUserValidator.On("ValidateUser", ctx, int64(999)).Return(nil, userNotFoundErr)
    
    // Attempt to create order
    response, err := server.CreateOrder(ctx, req)
    
    // Verify order creation was rejected
    assert.Error(t, err)
    assert.Equal(t, codes.FailedPrecondition, st.Code())
}
```

**Run specific tests:**
```bash
# Run order service tests with verbose output
go test -v ./order-service/internal/handlers

# Run all tests
go test ./...
```

## 🔧 Configuration

### Environment Variables

**User Service:**
- `DATABASE_URL`: PostgreSQL connection string (default: `postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable`)
- `PORT`: Service port (default: `50051`)

**Order Service:**
- `DATABASE_URL`: PostgreSQL connection string (default: `postgres://postgres:postgres@localhost:5433/order_service?sslmode=disable`)
- `PORT`: Service port (default: `50052`)
- `USER_SERVICE_ADDR`: User service address (default: `localhost:50051`)

### Database Setup

The system uses two separate PostgreSQL databases:

- **User DB**: Port 5432, Database: `user_service`
- **Order DB**: Port 5433, Database: `order_service`

Both are automatically initialized with sample data when using Docker Compose.

## 📊 API Reference

### User Service (Port 50051)

**CreateUser**
```protobuf
rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);

message CreateUserRequest {
  string email = 1;
  string name = 2;
}
```

**GetUser**
```protobuf
rpc GetUser(GetUserRequest) returns (GetUserResponse);

message GetUserRequest {
  int64 user_id = 1;
}
```

### Order Service (Port 50052)

**CreateOrder**
```protobuf
rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);

message CreateOrderRequest {
  int64 user_id = 1;
  string product_name = 2;
  double amount = 3;
  int32 quantity = 4;
}
```

**GetOrder**
```protobuf
rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);

message GetOrderRequest {
  int64 order_id = 1;
}
```

## 🏭 Production Considerations

### Monitoring & Observability
- Structured logging with service identification
- gRPC reflection enabled for debugging
- Health checks via database ping

### Error Handling
- Proper gRPC status codes
- Timeout handling (5-second timeout for user validation)
- Database connection pooling

### Security
- Input validation on all endpoints
- SQL injection prevention via parameterized queries
- Connection timeouts and retries

### Scalability
- Stateless service design
- Database per service pattern
- Interface-based design for easy mocking/testing

## 🛠️ Development

### Make Commands

```bash
make proto          # Generate protobuf code
make build          # Build both services
make test           # Run all tests
make user-service   # Build user service only
make order-service  # Build order service only
make deps           # Install dependencies
make db-up          # Start databases
make db-down        # Stop databases
make clean          # Clean build artifacts
make setup          # Complete development setup
```

### Adding New Features

1. **Update Protocol Buffers**: Modify `.proto` files in `proto/` directory
2. **Regenerate Code**: Run `make proto`
3. **Implement Handlers**: Update service handlers in respective `internal/handlers/`
4. **Add Tests**: Create unit tests with mocked dependencies
5. **Update Documentation**: Update this README

### Database Migrations

To add new database schema changes:

1. Update SQL files in `init-scripts/`
2. For existing databases, create migration scripts
3. Test with fresh Docker containers: `make db-down && make db-up`

## 📚 Learning Resources

This project demonstrates several key concepts:

- **Microservices Architecture**: Service isolation and communication
- **gRPC**: High-performance RPC framework usage
- **Protocol Buffers**: Type-safe API definitions
- **Database per Service**: Data isolation in microservices
- **Unit Testing**: Mocking external dependencies
- **Go Best Practices**: Project structure and error handling

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make changes and add tests
4. Run tests: `make test`
5. Submit a pull request

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

## 👨‍💻 Author

Built with ❤️ as a demonstration of microservices best practices using Go and gRPC.

---

**Note**: This project is designed for educational and demonstration purposes, showcasing production-ready patterns for microservices architecture with gRPC communication in Go.