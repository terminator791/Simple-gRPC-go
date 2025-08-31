package db

import (
	"database/sql"
	"fmt"

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
	query := `SELECT id, user_id, product_name, amount, quantity, status, created_at FROM orders WHERE id = $1`
	
	err := r.db.Get(&order, query, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	return &order, nil
}

func (r *OrderRepository) Create(userID int64, productName string, amount float64, quantity int32) (*models.Order, error) {
	var order models.Order
	query := `
		INSERT INTO orders (user_id, product_name, amount, quantity, status) 
		VALUES ($1, $2, $3, $4, 'pending') 
		RETURNING id, user_id, product_name, amount, quantity, status, created_at`
	
	err := r.db.Get(&order, query, userID, productName, amount, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	return &order, nil
}