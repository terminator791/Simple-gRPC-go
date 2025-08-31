package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/terminator791/Simple-gRPC-go/pkg/user"
)

// UserServiceClient wraps the gRPC client for user service
type UserServiceClient struct {
	client user.UserServiceClient
	conn   *grpc.ClientConn
	token  string
}

// Ensure UserServiceClient implements UserValidator interface
var _ UserValidator = (*UserServiceClient)(nil)

// NewUserServiceClient creates a new user service client
func NewUserServiceClient(address, serviceToken string) (*UserServiceClient, error) {
	// Set up connection to the user service
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := user.NewUserServiceClient(conn)
	
	return &UserServiceClient{
		client: client,
		conn:   conn,
		token:  serviceToken,
	}, nil
}

// Close closes the connection to the user service
func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}

// ValidateUser checks if a user exists by calling the user service
func (c *UserServiceClient) ValidateUser(ctx context.Context, userID int64) (*user.User, error) {
	log.Printf("Validating user with ID: %d", userID)
	
	// Add authentication header
	if c.token != "" {
		md := metadata.Pairs("authorization", "Bearer "+c.token)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	
	// Set timeout for the call
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &user.GetUserRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		// Check if it's a NOT_FOUND error
		if st, ok := status.FromError(err); ok {
			log.Printf("User service returned status: %s", st.Code())
			return nil, err
		}
		return nil, fmt.Errorf("failed to validate user: %w", err)
	}

	log.Printf("User validated successfully: %s", resp.User.Email)
	return resp.User, nil
}