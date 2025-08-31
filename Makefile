.PHONY: proto build test clean user-service order-service

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user/user.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/order/order.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/product/product.proto

# Build all services
build: proto
	@echo "Building user-service..."
	go build -o user-service/bin/user-service ./user-service/cmd
	@echo "Building order-service..."
	go build -o order-service/bin/order-service ./order-service/cmd
	@echo "Building product-service..."
	go build -o product-service/bin/product-service ./product-service/cmd

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Build user service
user-service: proto
	go build -o user-service/bin/user-service ./user-service/cmd

# Build order service
order-service: proto
	go build -o order-service/bin/order-service ./order-service/cmd

# Build product service
product-service: proto
	go build -o product-service/bin/product-service ./product-service/cmd

# Install dependencies
deps:
	go mod tidy

# Start databases
db-up:
	docker-compose up -d

# Stop databases
db-down:
	docker-compose down

# Clean build artifacts
clean:
	rm -rf user-service/bin order-service/bin product-service/bin
	rm -rf pkg/user/*.pb.go pkg/order/*.pb.go pkg/product/*.pb.go

# Setup development environment
setup: deps proto
	@echo "Setup complete!"