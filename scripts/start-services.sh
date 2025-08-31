#!/bin/bash
# Start both services in separate terminals

echo "Starting User Service and Order Service..."

# Start user service in background
echo "Starting User Service on port 50051..."
./user-service/bin/user-service &
USER_PID=$!

# Wait a bit for user service to start
sleep 2

# Start order service in background  
echo "Starting Order Service on port 50052..."
./order-service/bin/order-service &
ORDER_PID=$!

echo "Services started!"
echo "User Service PID: $USER_PID"
echo "Order Service PID: $ORDER_PID"
echo ""
echo "To stop services, run: kill $USER_PID $ORDER_PID"
echo ""
echo "Services are running on:"
echo "- User Service:  localhost:50051"
echo "- Order Service: localhost:50052"

# Keep script running
wait