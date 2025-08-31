package db

import "github.com/terminator791/Simple-gRPC-go/order-service/internal/models"

// OrderRepository interface defines methods for order data access
type OrderRepositoryInterface interface {
	GetByID(orderID int64) (*models.Order, error)
	Create(userID int64, productName string, amount float64, quantity int32) (*models.Order, error)
}