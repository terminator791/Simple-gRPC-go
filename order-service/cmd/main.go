package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/terminator791/Simple-gRPC-go/pkg/order"
	"github.com/terminator791/Simple-gRPC-go/order-service/internal/client"
	"github.com/terminator791/Simple-gRPC-go/order-service/internal/config"
	"github.com/terminator791/Simple-gRPC-go/order-service/internal/db"
	"github.com/terminator791/Simple-gRPC-go/order-service/internal/handlers"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting order service on port %s", cfg.Port)

	// Connect to database
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("Connected to order database successfully")

	// Create user service client
	userClient, err := client.NewUserServiceClient(cfg.UserServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}
	defer userClient.Close()
	log.Printf("Connected to user service at %s", cfg.UserServiceAddr)

	// Create repository
	orderRepo := db.NewOrderRepository(database)

	// Create gRPC server
	s := grpc.NewServer()
	
	// Register service
	orderHandler := handlers.NewOrderServiceServer(orderRepo, userClient)
	order.RegisterOrderServiceServer(s, orderHandler)
	
	// Enable reflection for grpcurl and other tools
	reflection.Register(s)

	// Start listening
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Order service listening on port %s", cfg.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}