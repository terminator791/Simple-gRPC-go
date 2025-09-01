package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
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
	productClient client.ProductClient
}

func NewOrderServiceServer(orderRepo db.OrderRepositoryInterface, userValidator client.UserValidator, productClient client.ProductClient) *OrderServiceServer {
	return &OrderServiceServer{
		orderRepo:     orderRepo,
		userValidator: userValidator,
		productClient: productClient,
	}
}

func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	log.Printf("CreateOrder called with user_id: %d, product_id: %d, quantity: %d", req.UserId, req.ProductId, req.Quantity)
	
	// Validate input
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}
	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id must be positive")
	}
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
	}
	
	// 1. Validate user exists by calling user service
	log.Printf("Validating user %d with user service", req.UserId)
	_, err := s.userValidator.ValidateUser(ctx, req.UserId)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			log.Printf("User %d not found in user service", req.UserId)
			return nil, status.Error(codes.FailedPrecondition, "user not found")
		}
		log.Printf("Error validating user: %v", err)
		return nil, status.Error(codes.Internal, "failed to validate user")
	}

	// 2. Get product information and validate it exists
	log.Printf("Getting product information for product %d", req.ProductId)
	product, err := s.productClient.GetProduct(ctx, req.ProductId)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			log.Printf("Product %d not found", req.ProductId)
			return nil, status.Error(codes.FailedPrecondition, "product not found")
		}
		log.Printf("Error getting product: %v", err)
		return nil, status.Error(codes.Internal, "failed to get product information")
	}

	// 3. Check inventory availability
	log.Printf("Checking inventory for product %d, quantity %d", req.ProductId, req.Quantity)
	available, currentStock, err := s.productClient.CheckInventory(ctx, req.ProductId, req.Quantity)
	if err != nil {
		log.Printf("Error checking inventory: %v", err)
		return nil, status.Error(codes.Internal, "failed to check inventory")
	}

	if !available {
		log.Printf("Insufficient inventory: requested %d, available %d", req.Quantity, currentStock)
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("insufficient inventory: requested %d, available %d", req.Quantity, currentStock))
	}

	// 4. Reserve inventory
	reservationID := uuid.New().String()
	log.Printf("Reserving inventory with reservation ID %s", reservationID)
	success, message, err := s.productClient.ReserveInventory(ctx, req.ProductId, req.Quantity, reservationID)
	if err != nil {
		log.Printf("Error reserving inventory: %v", err)
		return nil, status.Error(codes.Internal, "failed to reserve inventory")
	}

	if !success {
		log.Printf("Failed to reserve inventory: %s", message)
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("failed to reserve inventory: %s", message))
	}

	// 5. Calculate total amount
	totalAmount := product.Price * float64(req.Quantity)

	// 6. Create order with reserved inventory
	log.Printf("Creating order for user %d with reserved inventory", req.UserId)
	orderModel, err := s.orderRepo.Create(req.UserId, req.ProductId, product.Name, totalAmount, req.Quantity, reservationID)
	if err != nil {
		// If order creation fails, release the reserved inventory
		log.Printf("Order creation failed, releasing reserved inventory: %v", err)
		_, _, releaseErr := s.productClient.ReleaseInventory(ctx, reservationID)
		if releaseErr != nil {
			log.Printf("Error releasing inventory after order creation failure: %v", releaseErr)
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create order: %v", err))
	}
	
	orderProto := &order.Order{
		Id:          orderModel.ID,
		UserId:      orderModel.UserID,
		ProductId:   orderModel.ProductID,
		ProductName: orderModel.ProductName,
		Amount:      orderModel.Amount,
		Quantity:    orderModel.Quantity,
		Status:      orderModel.Status,
		CreatedAt:   timestamppb.New(orderModel.CreatedAt).String(),
		UpdatedAt:   timestamppb.New(orderModel.UpdatedAt).String(),
	}
	
	log.Printf("Order created successfully with ID: %d, reservation ID: %s", orderModel.ID, reservationID)
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
		ProductId:   orderModel.ProductID,
		ProductName: orderModel.ProductName,
		Amount:      orderModel.Amount,
		Quantity:    orderModel.Quantity,
		Status:      orderModel.Status,
		CreatedAt:   timestamppb.New(orderModel.CreatedAt).String(),
		UpdatedAt:   timestamppb.New(orderModel.UpdatedAt).String(),
	}
	
	return &order.GetOrderResponse{Order: orderProto}, nil
}