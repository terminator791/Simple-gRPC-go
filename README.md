# Advanced gRPC Go Microservices with Security & Observability

A production-ready microservices architecture implementation demonstrating enterprise-grade patterns including authentication, authorization, inventory management, circuit breakers, and comprehensive observability using gRPC in Go.

## 🏗️ Architecture Overview

The system consists of three independent microservices with sophisticated inter-service communication:

1. **User Service** - User management with JWT authentication
2. **Order Service** - Order processing with inventory validation  
3. **Product Service** - Product catalog and inventory management

### Enhanced Communication Flow

```mermaid
graph TB
    Client[Client Application] --> LB[Load Balancer]
    LB --> Gateway[API Gateway]
    Gateway --> Auth[JWT Auth Layer]
    
    Auth --> UserService[User Service :50051]
    Auth --> OrderService[Order Service :50052]
    Auth --> ProductService[Product Service :50053]
    
    OrderService --> |gRPC + Service Token| UserService
    OrderService --> |gRPC + Circuit Breaker| ProductService
    
    UserService --> UserDB[(User DB :5432)]
    OrderService --> OrderDB[(Order DB :5433)]
    ProductService --> ProductDB[(Product DB :5434)]
    
    ProductService --> |Inventory Mgmt| Cache[Redis Cache]
    OrderService --> |Events| Queue[Message Queue]
    
    Gateway --> Metrics[Prometheus :9090]
    Gateway --> Logs[ELK Stack]
    Gateway --> Traces[Jaeger]
```

## 🚀 Advanced Features

### Security Features
- ✅ **JWT Authentication**: Role-based access control with secure token management
- ✅ **Inter-Service Authentication**: Service-to-service authentication tokens
- ✅ **Input Validation**: Comprehensive request validation and sanitization
- ✅ **Public/Private Endpoints**: Granular access control per endpoint
- ⚡ **Rate Limiting**: Request throttling and abuse prevention
- ⚡ **mTLS Support**: Mutual TLS for encrypted inter-service communication

### Business Logic Complexity
- ✅ **Product Catalog**: Complete product management with categories
- ✅ **Inventory Management**: Real-time stock tracking with reservations
- ✅ **Order Processing**: Multi-step workflow with validation
- ✅ **Reservation System**: Temporary inventory holds with expiration
- ⚡ **Payment Integration**: Payment service simulation
- ⚡ **Order Lifecycle**: Sophisticated status management (pending → processing → shipped → delivered)

### Observability & Resilience
- ✅ **Circuit Breaker**: Fault tolerance for service-to-service calls
- ✅ **Health Checks**: Comprehensive health monitoring endpoints
- ⚡ **Prometheus Metrics**: Request latency, error rates, and business metrics
- ⚡ **Distributed Tracing**: Request tracing across services
- ⚡ **Structured Logging**: Correlation IDs and contextual logging
- ⚡ **Graceful Shutdowns**: Clean service termination

### Data Management
- ✅ **Database Per Service**: Complete data isolation
- ✅ **ACID Transactions**: Consistent inventory operations
- ✅ **Connection Pooling**: Optimized database performance
- ✅ **Migration System**: Versioned schema changes
- ⚡ **Event Sourcing**: Audit trail for all operations
- ⚡ **CQRS Pattern**: Command/Query responsibility segregation

## 🛠️ Technology Stack

- **Language**: Go 1.23+
- **RPC Framework**: gRPC (`google.golang.org/grpc`)
- **Authentication**: JWT (`github.com/golang-jwt/jwt/v5`)
- **Circuit Breaker**: Sony GoBreaker (`github.com/sony/gobreaker`)
- **Metrics**: Prometheus (`github.com/prometheus/client_golang`)
- **Database**: PostgreSQL 15 Alpine (3 instances)
- **Database Client**: `jmoiron/sqlx`
- **Testing**: Go testing + `stretchr/testify`
- **Containerization**: Docker & Docker Compose

## 📁 Enhanced Project Structure

```
Simple-gRPC-go/
├── proto/                          # Protocol Buffer definitions
│   ├── user/user.proto             # User service interface + auth
│   ├── order/order.proto           # Order service interface + workflow
│   ├── product/product.proto       # Product service + inventory
│   └── health/health.proto         # Health check interface
├── pkg/                            # Generated protobuf + shared code
│   ├── auth/                       # JWT authentication middleware
│   ├── health/                     # Health check service
│   ├── user/                       # Generated user service code
│   ├── order/                      # Generated order service code
│   └── product/                    # Generated product service code
├── user-service/                   # User microservice
│   ├── cmd/main.go                # Service entry point + auth
│   └── internal/
│       ├── config/                # Configuration management
│       ├── db/                    # Database layer + migrations
│       ├── handlers/              # gRPC handlers + JWT
│       └── models/                # Data models + validation
├── order-service/                  # Order microservice
│   ├── cmd/main.go               # Service entry point + circuit breaker
│   └── internal/
│       ├── config/               # Configuration management
│       ├── db/                   # Database layer + transactions
│       ├── handlers/             # gRPC handlers + workflow
│       ├── client/               # Service clients + circuit breaker
│       └── models/               # Data models + business logic
├── product-service/               # Product microservice
│   ├── cmd/main.go              # Service entry point + inventory
│   └── internal/
│       ├── config/              # Configuration management
│       ├── db/                  # Database layer + reservations
│       ├── handlers/            # gRPC handlers + inventory logic
│       └── models/              # Product + inventory models
├── init-scripts/                  # Database initialization
│   ├── user-db.sql             # User schema + roles + security
│   ├── order-db.sql            # Order schema + indexes + triggers
│   └── product-db.sql          # Product schema + inventory + audit
├── scripts/                      # Development + testing scripts
└── docker-compose.yml           # Multi-database setup + networking
```

## 🚀 Quick Start

### Prerequisites
- Go 1.23+
- Docker & Docker Compose
- Protocol Buffer Compiler (`protoc`)
- gRPC tools for Go

### Installation & Setup

```bash
# Clone repository
git clone https://github.com/terminator791/Simple-gRPC-go.git
cd Simple-gRPC-go

# Install dependencies
make deps

# Start databases
make db-up

# Build all services
make build

# Generate protobuf code (if needed)
make proto
```

### Running the Enhanced System

**Terminal 1: User Service (with JWT)**
```bash
./user-service/bin/user-service
# Starts on :50051 with authentication middleware
```

**Terminal 2: Product Service (with Inventory)**
```bash
./product-service/bin/product-service  
# Starts on :50053 with inventory management
```

**Terminal 3: Order Service (with Circuit Breaker)**
```bash
./order-service/bin/order-service
# Starts on :50052 with service integration
```

## 🧪 Testing the Enhanced System

### 1. User Authentication Flow

**Create User Account:**
```bash
grpcurl -plaintext -d '{
  "email": "user@example.com", 
  "name": "Test User",
  "password": "securepass123"
}' localhost:50051 user.UserService/CreateUser
```

**Login and Get JWT Token:**
```bash
grpcurl -plaintext -d '{
  "email": "user@example.com",
  "password": "securepass123"
}' localhost:50051 user.UserService/LoginUser
```

### 2. Product Management

**Create Product (requires auth):**
```bash
grpcurl -plaintext -H "authorization: Bearer YOUR_JWT_TOKEN" -d '{
  "name": "iPhone 15 Pro",
  "description": "Latest Apple smartphone", 
  "price": 999.99,
  "initial_stock": 50,
  "category": "Electronics"
}' localhost:50053 product.ProductService/CreateProduct
```

**Check Inventory:**
```bash
grpcurl -plaintext -H "authorization: Bearer YOUR_JWT_TOKEN" -d '{
  "product_id": 1,
  "required_quantity": 2
}' localhost:50053 product.ProductService/CheckInventory
```

### 3. Advanced Order Processing

**Create Order (with inventory check & reservation):**
```bash
grpcurl -plaintext -H "authorization: Bearer YOUR_JWT_TOKEN" -d '{
  "user_id": 1,
  "product_id": 1,
  "quantity": 2
}' localhost:50052 order.OrderService/CreateOrder
```

### 4. Health Monitoring

**Check Service Health:**
```bash
grpcurl -plaintext localhost:50051 health.HealthService/Check
grpcurl -plaintext localhost:50052 health.HealthService/Check  
grpcurl -plaintext localhost:50053 health.HealthService/Check
```

## 📊 Enhanced API Reference

### User Service (Port 50051) - Authentication & Users

**Public Endpoints:**
- `CreateUser` - User registration
- `LoginUser` - Authentication with JWT response

**Protected Endpoints (require JWT):**
- `GetUser` - Retrieve user information

### Product Service (Port 50053) - Catalog & Inventory

**All endpoints require authentication**

- `CreateProduct` - Add products to catalog
- `GetProduct` - Retrieve product details
- `ListProducts` - Paginated product listing with filters
- `UpdateInventory` - Modify stock levels with audit
- `CheckInventory` - Real-time availability check
- `ReserveInventory` - Temporary inventory holds
- `ReleaseInventory` - Release reserved stock

### Order Service (Port 50052) - Order Processing

**All endpoints require authentication**

- `CreateOrder` - Multi-step order creation with:
  - User validation via User Service
  - Product validation via Product Service  
  - Inventory availability check
  - Automatic inventory reservation
  - Rollback on failure
- `GetOrder` - Retrieve order details

## 🏭 Production Considerations

### Security
- JWT tokens with configurable expiration
- Service-to-service authentication
- Input validation and sanitization
- SQL injection prevention
- Rate limiting per endpoint
- CORS and security headers

### Observability
- Health check endpoints for all services
- Circuit breaker state monitoring
- Database connection health
- Structured logging with levels
- Request/response tracing
- Business metrics collection

### Scalability & Performance
- Stateless service design
- Database per service pattern
- Connection pooling with limits
- Circuit breaker for fault tolerance
- Graceful service degradation
- Horizontal scaling ready

### Data Management
- ACID transactions for critical operations
- Inventory reservation expiration
- Audit logging for all changes
- Database migration versioning
- Backup and disaster recovery
- Data retention policies

## 🔧 Configuration

All services support environment-based configuration:

```bash
# User Service
export PORT=50051
export DATABASE_URL="postgres://..."
export JWT_SECRET="your-secret-key"

# Order Service  
export PORT=50052
export USER_SERVICE_ADDR="localhost:50051"
export PRODUCT_SERVICE_ADDR="localhost:50053"

# Product Service
export PRODUCT_SERVICE_PORT=50053
export PRODUCT_DATABASE_URL="postgres://..."
```

## 🛠️ Development

### Make Commands

```bash
make proto          # Generate protobuf code
make build          # Build all services  
make test           # Run all tests
make user-service   # Build user service only
make order-service  # Build order service only
make product-service # Build product service only
make deps           # Install dependencies
make db-up          # Start databases
make db-down        # Stop databases
make clean          # Clean build artifacts
make setup          # Complete development setup
```

### Adding New Features

1. **Update Protocol Buffers**: Modify `.proto` files
2. **Regenerate Code**: Run `make proto`
3. **Implement Handlers**: Update service handlers
4. **Add Database Changes**: Update SQL migration scripts
5. **Add Tests**: Create unit/integration tests
6. **Update Documentation**: Update this README

### Testing Strategy

- Unit tests with mocked dependencies
- Integration tests with test databases
- gRPC client testing with testify
- Circuit breaker testing scenarios
- Load testing for performance validation

## 📚 Key Learning Concepts

This repository demonstrates:

1. **Microservices Architecture**: Service decomposition and communication
2. **gRPC Communication**: Type-safe, efficient inter-service calls
3. **Authentication & Authorization**: JWT-based security
4. **Data Consistency**: Transaction management across services
5. **Fault Tolerance**: Circuit breaker and graceful degradation
6. **Observability**: Health checks, metrics, and monitoring
7. **Domain-Driven Design**: Business logic organization
8. **Infrastructure as Code**: Docker-based development environment

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Update documentation
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License. See LICENSE file for details.

## 👨‍💻 Author

Built with ❤️ to demonstrate production-ready microservices patterns in Go.

---

**Note**: This is an educational project showcasing advanced microservices patterns. For production use, consider additional security hardening, comprehensive monitoring, and deployment automation.

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