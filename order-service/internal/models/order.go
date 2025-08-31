package models

import (
	"time"
)

type Order struct {
	ID          int64     `db:"id" json:"id"`
	UserID      int64     `db:"user_id" json:"user_id"`
	ProductName string    `db:"product_name" json:"product_name"`
	Amount      float64   `db:"amount" json:"amount"`
	Quantity    int32     `db:"quantity" json:"quantity"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}