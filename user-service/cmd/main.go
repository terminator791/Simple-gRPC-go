package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/terminator791/Simple-gRPC-go/pkg/user"
	"github.com/terminator791/Simple-gRPC-go/user-service/internal/config"
	"github.com/terminator791/Simple-gRPC-go/user-service/internal/db"
	"github.com/terminator791/Simple-gRPC-go/user-service/internal/handlers"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting user service on port %s", cfg.Port)

	// Connect to database
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("Connected to database successfully")

	// Create repository
	userRepo := db.NewUserRepository(database)

	// Create gRPC server
	s := grpc.NewServer()
	
	// Register service
	userHandler := handlers.NewUserServiceServer(userRepo)
	user.RegisterUserServiceServer(s, userHandler)
	
	// Enable reflection for grpcurl and other tools
	reflection.Register(s)

	// Start listening
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("User service listening on port %s", cfg.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}