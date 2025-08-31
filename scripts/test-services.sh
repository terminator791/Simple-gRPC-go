#!/bin/bash
# Test the gRPC services using grpcurl

echo "Testing gRPC Services..."
echo "========================"

echo "1. Creating a user..."
grpcurl -plaintext -d '{"email": "test@example.com", "name": "Test User"}' localhost:50051 user.UserService/CreateUser

echo ""
echo "2. Getting user by ID (should exist)..."
grpcurl -plaintext -d '{"user_id": 1}' localhost:50051 user.UserService/GetUser

echo ""
echo "3. Creating an order for existing user..."
grpcurl -plaintext -d '{"user_id": 1, "product_name": "Test Product", "amount": 99.99, "quantity": 2}' localhost:50052 order.OrderService/CreateOrder

echo ""
echo "4. Trying to create order for non-existent user (should fail)..."
grpcurl -plaintext -d '{"user_id": 999, "product_name": "Test Product", "amount": 99.99, "quantity": 2}' localhost:50052 order.OrderService/CreateOrder

echo ""
echo "5. Getting order by ID..."
grpcurl -plaintext -d '{"order_id": 1}' localhost:50052 order.OrderService/GetOrder