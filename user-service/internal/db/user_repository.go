package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/terminator791/Simple-gRPC-go/user-service/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
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

func (r *UserRepository) GetByID(userID int64) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, name, created_at FROM users WHERE id = $1`
	
	err := r.db.Get(&user, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return &user, nil
}

func (r *UserRepository) Create(email, name string) (*models.User, error) {
	var user models.User
	query := `
		INSERT INTO users (email, name) 
		VALUES ($1, $2) 
		RETURNING id, email, name, created_at`
	
	err := r.db.Get(&user, query, email, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	return &user, nil
}