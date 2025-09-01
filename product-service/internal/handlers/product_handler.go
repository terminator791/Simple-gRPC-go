package handlers

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/terminator791/Simple-gRPC-go/pkg/product"
	"github.com/terminator791/Simple-gRPC-go/product-service/internal/db"
)

// ProductServiceServer implements the ProductService
type ProductServiceServer struct {
	product.UnimplementedProductServiceServer
	productRepo db.ProductRepositoryInterface
}

// NewProductServiceServer creates a new product service server
func NewProductServiceServer(productRepo db.ProductRepositoryInterface) *ProductServiceServer {
	return &ProductServiceServer{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product
func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	log.Printf("CreateProduct called with name: %s, price: %f", req.Name, req.Price)

	// Validate input
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "price must be positive")
	}
	if req.InitialStock < 0 {
		return nil, status.Error(codes.InvalidArgument, "initial stock cannot be negative")
	}
	if req.Category == "" {
		return nil, status.Error(codes.InvalidArgument, "category is required")
	}

	productModel, err := s.productRepo.Create(req.Name, req.Description, req.Category, req.Price, req.InitialStock)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create product: %v", err))
	}

	productProto := &product.Product{
		Id:               productModel.ID,
		Name:             productModel.Name,
		Description:      productModel.Description,
		Price:            productModel.Price,
		StockQuantity:    productModel.StockQuantity,
		ReservedQuantity: productModel.ReservedQuantity,
		Category:         productModel.Category,
		CreatedAt:        timestamppb.New(productModel.CreatedAt).String(),
		UpdatedAt:        timestamppb.New(productModel.UpdatedAt).String(),
	}

	log.Printf("Product created successfully with ID: %d", productModel.ID)
	return &product.CreateProductResponse{Product: productProto}, nil
}

// GetProduct retrieves a product by ID
func (s *ProductServiceServer) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	log.Printf("GetProduct called with product_id: %d", req.ProductId)

	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id must be positive")
	}

	productModel, err := s.productRepo.GetByID(req.ProductId)
	if err != nil {
		if err.Error() == "product not found" {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		log.Printf("Error getting product: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	productProto := &product.Product{
		Id:               productModel.ID,
		Name:             productModel.Name,
		Description:      productModel.Description,
		Price:            productModel.Price,
		StockQuantity:    productModel.StockQuantity,
		ReservedQuantity: productModel.ReservedQuantity,
		Category:         productModel.Category,
		CreatedAt:        timestamppb.New(productModel.CreatedAt).String(),
		UpdatedAt:        timestamppb.New(productModel.UpdatedAt).String(),
	}

	return &product.GetProductResponse{Product: productProto}, nil
}

// UpdateInventory updates product inventory
func (s *ProductServiceServer) UpdateInventory(ctx context.Context, req *product.UpdateInventoryRequest) (*product.UpdateInventoryResponse, error) {
	log.Printf("UpdateInventory called with product_id: %d, quantity_change: %d", req.ProductId, req.QuantityChange)

	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id must be positive")
	}
	if req.Reason == "" {
		return nil, status.Error(codes.InvalidArgument, "reason is required")
	}

	productModel, err := s.productRepo.UpdateInventory(req.ProductId, req.QuantityChange, req.Reason)
	if err != nil {
		log.Printf("Error updating inventory: %v", err)
		if err.Error() == "product not found" {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update inventory: %v", err))
	}

	productProto := &product.Product{
		Id:               productModel.ID,
		Name:             productModel.Name,
		Description:      productModel.Description,
		Price:            productModel.Price,
		StockQuantity:    productModel.StockQuantity,
		ReservedQuantity: productModel.ReservedQuantity,
		Category:         productModel.Category,
		CreatedAt:        timestamppb.New(productModel.CreatedAt).String(),
		UpdatedAt:        timestamppb.New(productModel.UpdatedAt).String(),
	}

	return &product.UpdateInventoryResponse{Product: productProto}, nil
}

// CheckInventory checks if product has sufficient inventory
func (s *ProductServiceServer) CheckInventory(ctx context.Context, req *product.CheckInventoryRequest) (*product.CheckInventoryResponse, error) {
	log.Printf("CheckInventory called with product_id: %d, required_quantity: %d", req.ProductId, req.RequiredQuantity)

	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id must be positive")
	}
	if req.RequiredQuantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "required_quantity must be positive")
	}

	available, currentStock, err := s.productRepo.CheckInventory(req.ProductId, req.RequiredQuantity)
	if err != nil {
		log.Printf("Error checking inventory: %v", err)
		if err.Error() == "product not found" {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &product.CheckInventoryResponse{
		Available:    available,
		CurrentStock: currentStock,
	}, nil
}

// ListProducts lists products with pagination
func (s *ProductServiceServer) ListProducts(ctx context.Context, req *product.ListProductsRequest) (*product.ListProductsResponse, error) {
	log.Printf("ListProducts called with page: %d, page_size: %d, category: %s", req.Page, req.PageSize, req.Category)

	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	productModels, totalCount, err := s.productRepo.List(page, pageSize, req.Category)
	if err != nil {
		log.Printf("Error listing products: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var products []*product.Product
	for _, productModel := range productModels {
		productProto := &product.Product{
			Id:               productModel.ID,
			Name:             productModel.Name,
			Description:      productModel.Description,
			Price:            productModel.Price,
			StockQuantity:    productModel.StockQuantity,
			ReservedQuantity: productModel.ReservedQuantity,
			Category:         productModel.Category,
			CreatedAt:        timestamppb.New(productModel.CreatedAt).String(),
			UpdatedAt:        timestamppb.New(productModel.UpdatedAt).String(),
		}
		products = append(products, productProto)
	}

	return &product.ListProductsResponse{
		Products:   products,
		TotalCount: totalCount,
	}, nil
}

// ReserveInventory reserves inventory for an order
func (s *ProductServiceServer) ReserveInventory(ctx context.Context, req *product.ReserveInventoryRequest) (*product.ReserveInventoryResponse, error) {
	log.Printf("ReserveInventory called with product_id: %d, quantity: %d, reservation_id: %s", 
		req.ProductId, req.Quantity, req.ReservationId)

	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id must be positive")
	}
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
	}
	if req.ReservationId == "" {
		return nil, status.Error(codes.InvalidArgument, "reservation_id is required")
	}

	err := s.productRepo.ReserveInventory(req.ProductId, req.Quantity, req.ReservationId)
	if err != nil {
		log.Printf("Error reserving inventory: %v", err)
		return &product.ReserveInventoryResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &product.ReserveInventoryResponse{
		Success: true,
		Message: "Inventory reserved successfully",
	}, nil
}

// ReleaseInventory releases reserved inventory
func (s *ProductServiceServer) ReleaseInventory(ctx context.Context, req *product.ReleaseInventoryRequest) (*product.ReleaseInventoryResponse, error) {
	log.Printf("ReleaseInventory called with reservation_id: %s", req.ReservationId)

	if req.ReservationId == "" {
		return nil, status.Error(codes.InvalidArgument, "reservation_id is required")
	}

	err := s.productRepo.ReleaseInventory(req.ReservationId)
	if err != nil {
		log.Printf("Error releasing inventory: %v", err)
		return &product.ReleaseInventoryResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &product.ReleaseInventoryResponse{
		Success: true,
		Message: "Inventory released successfully",
	}, nil
}