package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/terminator791/Simple-gRPC-go/pkg/auth"
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

	// Create JWT manager for inter-service communication
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)
	
	// Generate service token for inter-service communication
	serviceToken, err := jwtManager.GenerateServiceToken("order-service")
	if err != nil {
		log.Fatalf("Failed to generate service token: %v", err)
	}

	// Create user service client
	userClient, err := client.NewUserServiceClient(cfg.UserServiceAddr, serviceToken)
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}
	defer userClient.Close()
	log.Printf("Connected to user service at %s", cfg.UserServiceAddr)

	// Create product service client
	productClient, err := client.NewProductServiceClient(cfg.ProductServiceAddr, serviceToken)
	if err != nil {
		log.Fatalf("Failed to create product service client: %v", err)
	}
	defer productClient.Close()
	log.Printf("Connected to product service at %s", cfg.ProductServiceAddr)

	// Create repository
	orderRepo := db.NewOrderRepository(database)

	// Create gRPC server with authentication
	s := grpc.NewServer(
		grpc.UnaryInterceptor(jwtManager.UnaryAuthInterceptor()),
	)
	
	// Register service
	orderHandler := handlers.NewOrderServiceServer(orderRepo, userClient, productClient)
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