package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/terminator791/Simple-gRPC-go/order-service/internal/models"
)

type OrderRepository struct {
	db *sqlx.DB
}

// Ensure OrderRepository implements the interface
var _ OrderRepositoryInterface = (*OrderRepository)(nil)

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func Connect(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return db, nil
}

func (r *OrderRepository) GetByID(orderID int64) (*models.Order, error) {
	var order models.Order
	query := `SELECT id, user_id, product_id, product_name, amount, quantity, status, reservation_id, created_at, updated_at FROM orders WHERE id = $1`
	
	err := r.db.Get(&order, query, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	return &order, nil
}

func (r *OrderRepository) Create(userID, productID int64, productName string, amount float64, quantity int32, reservationID string) (*models.Order, error) {
	var order models.Order
	now := time.Now()
	query := `
		INSERT INTO orders (user_id, product_id, product_name, amount, quantity, status, reservation_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, 'pending', $6, $7, $8) 
		RETURNING id, user_id, product_id, product_name, amount, quantity, status, reservation_id, created_at, updated_at`
	
	err := r.db.Get(&order, query, userID, productID, productName, amount, quantity, reservationID, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	return &order, nil
}

// UpdateStatus updates the order status
func (r *OrderRepository) UpdateStatus(orderID int64, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}