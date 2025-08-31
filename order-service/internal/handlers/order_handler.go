package handlers

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/terminator791/Simple-gRPC-go/pkg/order"
	"github.com/terminator791/Simple-gRPC-go/order-service/internal/client"
	"github.com/terminator791/Simple-gRPC-go/order-service/internal/db"
)

type OrderServiceServer struct {
	order.UnimplementedOrderServiceServer
	orderRepo     db.OrderRepositoryInterface
	userValidator client.UserValidator
}

func NewOrderServiceServer(orderRepo db.OrderRepositoryInterface, userValidator client.UserValidator) *OrderServiceServer {
	return &OrderServiceServer{
		orderRepo:     orderRepo,
		userValidator: userValidator,
	}
}

func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	log.Printf("CreateOrder called with user_id: %d, product: %s", req.UserId, req.ProductName)
	
	// Validate input
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}
	if req.ProductName == "" {
		return nil, status.Error(codes.InvalidArgument, "product_name is required")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be positive")
	}
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
	}
	
	// THIS IS THE KEY PART: Validate user exists by calling user service
	log.Printf("Validating user %d with user service", req.UserId)
	_, err := s.userValidator.ValidateUser(ctx, req.UserId)
	if err != nil {
		// Check if it's a NOT_FOUND error from user service
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			log.Printf("User %d not found in user service", req.UserId)
			return nil, status.Error(codes.FailedPrecondition, "user not found")
		}
		log.Printf("Error validating user: %v", err)
		return nil, status.Error(codes.Internal, "failed to validate user")
	}
	
	// User exists, proceed to create order
	log.Printf("User %d validated successfully, creating order", req.UserId)
	orderModel, err := s.orderRepo.Create(req.UserId, req.ProductName, req.Amount, req.Quantity)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create order: %v", err))
	}
	
	orderProto := &order.Order{
		Id:          orderModel.ID,
		UserId:      orderModel.UserID,
		ProductName: orderModel.ProductName,
		Amount:      orderModel.Amount,
		Quantity:    orderModel.Quantity,
		Status:      orderModel.Status,
		CreatedAt:   timestamppb.New(orderModel.CreatedAt).String(),
	}
	
	log.Printf("Order created successfully with ID: %d", orderModel.ID)
	return &order.CreateOrderResponse{Order: orderProto}, nil
}

func (s *OrderServiceServer) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	log.Printf("GetOrder called with order_id: %d", req.OrderId)
	
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id must be positive")
	}
	
	orderModel, err := s.orderRepo.GetByID(req.OrderId)
	if err != nil {
		if err.Error() == "order not found" {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		log.Printf("Error getting order: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	
	orderProto := &order.Order{
		Id:          orderModel.ID,
		UserId:      orderModel.UserID,
		ProductName: orderModel.ProductName,
		Amount:      orderModel.Amount,
		Quantity:    orderModel.Quantity,
		Status:      orderModel.Status,
		CreatedAt:   timestamppb.New(orderModel.CreatedAt).String(),
	}
	
	return &order.GetOrderResponse{Order: orderProto}, nil
}