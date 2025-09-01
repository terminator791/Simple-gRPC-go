package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/terminator791/Simple-gRPC-go/product-service/internal/models"
)

// ProductRepositoryInterface defines the methods for product repository
type ProductRepositoryInterface interface {
	Create(name, description, category string, price float64, initialStock int32) (*models.Product, error)
	GetByID(productID int64) (*models.Product, error)
	UpdateInventory(productID int64, quantityChange int32, reason string) (*models.Product, error)
	CheckInventory(productID int64, requiredQuantity int32) (bool, int32, error)
	List(page, pageSize int32, category string) ([]*models.Product, int32, error)
	ReserveInventory(productID int64, quantity int32, reservationID string) error
	ReleaseInventory(reservationID string) error
	GetReservation(reservationID string) (*models.Reservation, error)
}

// ProductRepository implements ProductRepositoryInterface
type ProductRepository struct {
	db *sqlx.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product
func (r *ProductRepository) Create(name, description, category string, price float64, initialStock int32) (*models.Product, error) {
	query := `
		INSERT INTO products (name, description, price, stock_quantity, reserved_quantity, category, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, description, price, stock_quantity, reserved_quantity, category, created_at, updated_at`
	
	now := time.Now()
	product := &models.Product{}
	
	err := r.db.QueryRowx(query, name, description, price, initialStock, 0, category, now, now).StructScan(product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Log the initial inventory
	if initialStock > 0 {
		r.logInventoryChange(product.ID, initialStock, "Initial stock")
	}
	
	return product, nil
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(productID int64) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock_quantity, reserved_quantity, category, created_at, updated_at
		FROM products WHERE id = $1`
	
	product := &models.Product{}
	err := r.db.Get(product, query, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	
	return product, nil
}

// UpdateInventory updates product inventory
func (r *ProductRepository) UpdateInventory(productID int64, quantityChange int32, reason string) (*models.Product, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current product with lock
	var currentStock int32
	err = tx.Get(&currentStock, "SELECT stock_quantity FROM products WHERE id = $1 FOR UPDATE", productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current stock: %w", err)
	}

	newStock := currentStock + quantityChange
	if newStock < 0 {
		return nil, fmt.Errorf("insufficient inventory: current=%d, requested_change=%d", currentStock, quantityChange)
	}

	// Update the stock
	updateQuery := `
		UPDATE products 
		SET stock_quantity = $1, updated_at = $2 
		WHERE id = $3`
	
	_, err = tx.Exec(updateQuery, newStock, time.Now(), productID)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	// Log the inventory change
	err = r.logInventoryChangeInTx(tx, productID, quantityChange, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to log inventory change: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.GetByID(productID)
}

// CheckInventory checks if product has sufficient inventory
func (r *ProductRepository) CheckInventory(productID int64, requiredQuantity int32) (bool, int32, error) {
	query := `SELECT stock_quantity FROM products WHERE id = $1`
	
	var currentStock int32
	err := r.db.Get(&currentStock, query, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, fmt.Errorf("product not found")
		}
		return false, 0, fmt.Errorf("failed to check inventory: %w", err)
	}
	
	available := currentStock >= requiredQuantity
	return available, currentStock, nil
}

// List retrieves products with pagination
func (r *ProductRepository) List(page, pageSize int32, category string) ([]*models.Product, int32, error) {
	offset := (page - 1) * pageSize
	
	var products []*models.Product
	var totalCount int32
	
	// Build query with optional category filter
	whereClause := ""
	args := []interface{}{}
	if category != "" {
		whereClause = "WHERE category = $1"
		args = append(args, category)
	}
	
	// Count total products
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	err := r.db.Get(&totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}
	
	// Get products
	listQuery := fmt.Sprintf(`
		SELECT id, name, description, price, stock_quantity, reserved_quantity, category, created_at, updated_at
		FROM products %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, len(args)+1, len(args)+2)
	
	args = append(args, pageSize, offset)
	err = r.db.Select(&products, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	
	return products, totalCount, nil
}

// ReserveInventory reserves inventory for an order
func (r *ProductRepository) ReserveInventory(productID int64, quantity int32, reservationID string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check current stock
	var currentStock, reserved int32
	err = tx.Get(&struct {
		Stock    int32 `db:"stock_quantity"`
		Reserved int32 `db:"reserved_quantity"`
	}{currentStock, reserved}, 
		"SELECT stock_quantity, reserved_quantity FROM products WHERE id = $1 FOR UPDATE", productID)
	if err != nil {
		return fmt.Errorf("failed to get product stock: %w", err)
	}

	availableStock := currentStock - reserved
	if availableStock < quantity {
		return fmt.Errorf("insufficient available inventory: available=%d, requested=%d", availableStock, quantity)
	}

	// Update reserved quantity
	_, err = tx.Exec("UPDATE products SET reserved_quantity = reserved_quantity + $1 WHERE id = $2", quantity, productID)
	if err != nil {
		return fmt.Errorf("failed to update reserved quantity: %w", err)
	}

	// Create reservation record
	expiresAt := time.Now().Add(15 * time.Minute) // Reservations expire in 15 minutes
	_, err = tx.Exec(`
		INSERT INTO reservations (product_id, quantity, reservation_id, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)`,
		productID, quantity, reservationID, time.Now(), expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	return tx.Commit()
}

// ReleaseInventory releases reserved inventory
func (r *ProductRepository) ReleaseInventory(reservationID string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get reservation details
	reservation := &models.Reservation{}
	err = tx.Get(reservation, "SELECT product_id, quantity FROM reservations WHERE reservation_id = $1", reservationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("reservation not found")
		}
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	// Update reserved quantity
	_, err = tx.Exec("UPDATE products SET reserved_quantity = reserved_quantity - $1 WHERE id = $2", 
		reservation.Quantity, reservation.ProductID)
	if err != nil {
		return fmt.Errorf("failed to update reserved quantity: %w", err)
	}

	// Delete reservation
	_, err = tx.Exec("DELETE FROM reservations WHERE reservation_id = $1", reservationID)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}

	return tx.Commit()
}

// GetReservation retrieves a reservation by ID
func (r *ProductRepository) GetReservation(reservationID string) (*models.Reservation, error) {
	reservation := &models.Reservation{}
	err := r.db.Get(reservation, 
		"SELECT product_id, quantity, reservation_id, created_at, expires_at FROM reservations WHERE reservation_id = $1", 
		reservationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reservation not found")
		}
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}
	return reservation, nil
}

// logInventoryChange logs an inventory change
func (r *ProductRepository) logInventoryChange(productID int64, change int32, reason string) error {
	_, err := r.db.Exec(`
		INSERT INTO inventory_logs (product_id, change_quantity, reason, created_at)
		VALUES ($1, $2, $3, $4)`,
		productID, change, reason, time.Now())
	return err
}

// logInventoryChangeInTx logs an inventory change within a transaction
func (r *ProductRepository) logInventoryChangeInTx(tx *sqlx.Tx, productID int64, change int32, reason string) error {
	_, err := tx.Exec(`
		INSERT INTO inventory_logs (product_id, change_quantity, reason, created_at)
		VALUES ($1, $2, $3, $4)`,
		productID, change, reason, time.Now())
	return err
}