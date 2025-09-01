package models

import (
	"time"
)

// Product represents the product data model
type Product struct {
	ID               int64     `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`
	Description      string    `db:"description" json:"description"`
	Price            float64   `db:"price" json:"price"`
	StockQuantity    int32     `db:"stock_quantity" json:"stock_quantity"`
	ReservedQuantity int32     `db:"reserved_quantity" json:"reserved_quantity"`
	Category         string    `db:"category" json:"category"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// Reservation represents an inventory reservation
type Reservation struct {
	ID            string    `db:"id" json:"id"`
	ProductID     int64     `db:"product_id" json:"product_id"`
	Quantity      int32     `db:"quantity" json:"quantity"`
	ReservationID string    `db:"reservation_id" json:"reservation_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	ExpiresAt     time.Time `db:"expires_at" json:"expires_at"`
}

// InventoryLog represents inventory change log
type InventoryLog struct {
	ID        int64     `db:"id" json:"id"`
	ProductID int64     `db:"product_id" json:"product_id"`
	Change    int32     `db:"change_quantity" json:"change"`
	Reason    string    `db:"reason" json:"reason"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}