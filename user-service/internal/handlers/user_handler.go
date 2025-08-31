package handlers

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/terminator791/Simple-gRPC-go/pkg/user"
	"github.com/terminator791/Simple-gRPC-go/user-service/internal/db"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	userRepo *db.UserRepository
}

func NewUserServiceServer(userRepo *db.UserRepository) *UserServiceServer {
	return &UserServiceServer{
		userRepo: userRepo,
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
	
	userModel, err := s.userRepo.Create(req.Email, req.Name)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create user: %v", err))
	}
	
	userProto := &user.User{
		Id:        userModel.ID,
		Email:     userModel.Email,
		Name:      userModel.Name,
		CreatedAt: timestamppb.New(userModel.CreatedAt).String(),
	}
	
	return &user.CreateUserResponse{User: userProto}, nil
}