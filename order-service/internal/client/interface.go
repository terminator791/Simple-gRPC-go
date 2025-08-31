package client

import (
	"context"
	"github.com/terminator791/Simple-gRPC-go/pkg/user"
)

// UserValidator defines the interface for user validation
// This allows us to mock the user service client in tests
type UserValidator interface {
	ValidateUser(ctx context.Context, userID int64) (*user.User, error)
}