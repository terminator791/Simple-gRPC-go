package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/terminator791/Simple-gRPC-go/pkg/auth"
	"github.com/terminator791/Simple-gRPC-go/pkg/product"
	"github.com/terminator791/Simple-gRPC-go/product-service/internal/config"
	"github.com/terminator791/Simple-gRPC-go/product-service/internal/db"
	"github.com/terminator791/Simple-gRPC-go/product-service/internal/handlers"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting product service on port %s", cfg.Port)

	// Connect to database
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("Connected to product database successfully")

	// Create repository
	productRepo := db.NewProductRepository(database)

	// Create JWT manager for authentication
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// Create gRPC server with authentication interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(jwtManager.UnaryAuthInterceptor()),
	)

	// Register service
	productHandler := handlers.NewProductServiceServer(productRepo)
	product.RegisterProductServiceServer(s, productHandler)

	// Enable reflection for grpcurl and other tools
	reflection.Register(s)

	// Start listening
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Product service listening on port %s", cfg.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}