package handlers

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/terminator791/Simple-gRPC-go/pkg/auth"
	"github.com/terminator791/Simple-gRPC-go/pkg/user"
	"github.com/terminator791/Simple-gRPC-go/user-service/internal/db"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	userRepo   *db.UserRepository
	jwtManager *auth.JWTManager
}

func NewUserServiceServer(userRepo *db.UserRepository, jwtManager *auth.JWTManager) *UserServiceServer {
	return &UserServiceServer{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	log.Printf("GetUser called with user_id: %d", req.UserId)
	
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}
	
	userModel, err := s.userRepo.GetByID(req.UserId)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		log.Printf("Error getting user: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	
	userProto := &user.User{
		Id:        userModel.ID,
		Email:     userModel.Email,
		Name:      userModel.Name,
		Role:      userModel.Role,
		CreatedAt: timestamppb.New(userModel.CreatedAt).String(),
	}
	
	return &user.GetUserResponse{User: userProto}, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	log.Printf("CreateUser called with email: %s, name: %s", req.Email, req.Name)
	
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	
	userModel, err := s.userRepo.Create(req.Email, req.Name, req.Password)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create user: %v", err))
	}
	
	userProto := &user.User{
		Id:        userModel.ID,
		Email:     userModel.Email,
		Name:      userModel.Name,
		Role:      userModel.Role,
		CreatedAt: timestamppb.New(userModel.CreatedAt).String(),
	}
	
	return &user.CreateUserResponse{User: userProto}, nil
}

// LoginUser authenticates a user and returns a JWT token
func (s *UserServiceServer) LoginUser(ctx context.Context, req *user.LoginUserRequest) (*user.LoginUserResponse, error) {
	log.Printf("LoginUser called with email: %s", req.Email)
	
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	
	// Get user by email
	userModel, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, status.Error(codes.NotFound, "invalid email or password")
		}
		log.Printf("Error getting user by email: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	
	// Verify password
	if !s.userRepo.VerifyPassword(userModel.PasswordHash, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}
	
	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(userModel.ID, userModel.Email, userModel.Role)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	
	userProto := &user.User{
		Id:        userModel.ID,
		Email:     userModel.Email,
		Name:      userModel.Name,
		Role:      userModel.Role,
		CreatedAt: timestamppb.New(userModel.CreatedAt).String(),
	}
	
	log.Printf("User %s logged in successfully", req.Email)
	return &user.LoginUserResponse{
		User:  userProto,
		Token: token,
	}, nil
}