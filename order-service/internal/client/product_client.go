package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/terminator791/Simple-gRPC-go/pkg/product"
)

// ProductClient interface for product service operations
type ProductClient interface {
	CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, int32, error)
	ReserveInventory(ctx context.Context, productID int64, quantity int32, reservationID string) (bool, string, error)
	ReleaseInventory(ctx context.Context, reservationID string) (bool, string, error)
	GetProduct(ctx context.Context, productID int64) (*product.Product, error)
	Close() error
}

// ProductServiceClient implements ProductClient
type ProductServiceClient struct {
	conn   *grpc.ClientConn
	client product.ProductServiceClient
	token  string
}

// NewProductServiceClient creates a new product service client
func NewProductServiceClient(address, serviceToken string) (*ProductServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := product.NewProductServiceClient(conn)

	return &ProductServiceClient{
		conn:   conn,
		client: client,
		token:  serviceToken,
	}, nil
}

// CheckInventory checks if a product has sufficient inventory
func (c *ProductServiceClient) CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, int32, error) {
	ctx = c.addAuthHeader(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.CheckInventory(ctx, &product.CheckInventoryRequest{
		ProductId:        productID,
		RequiredQuantity: quantity,
	})
	if err != nil {
		return false, 0, err
	}

	return resp.Available, resp.CurrentStock, nil
}

// ReserveInventory reserves inventory for an order
func (c *ProductServiceClient) ReserveInventory(ctx context.Context, productID int64, quantity int32, reservationID string) (bool, string, error) {
	ctx = c.addAuthHeader(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.ReserveInventory(ctx, &product.ReserveInventoryRequest{
		ProductId:     productID,
		Quantity:      quantity,
		ReservationId: reservationID,
	})
	if err != nil {
		return false, "", err
	}

	return resp.Success, resp.Message, nil
}

// ReleaseInventory releases reserved inventory
func (c *ProductServiceClient) ReleaseInventory(ctx context.Context, reservationID string) (bool, string, error) {
	ctx = c.addAuthHeader(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.ReleaseInventory(ctx, &product.ReleaseInventoryRequest{
		ReservationId: reservationID,
	})
	if err != nil {
		return false, "", err
	}

	return resp.Success, resp.Message, nil
}

// GetProduct retrieves a product by ID
func (c *ProductServiceClient) GetProduct(ctx context.Context, productID int64) (*product.Product, error) {
	ctx = c.addAuthHeader(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.GetProduct(ctx, &product.GetProductRequest{
		ProductId: productID,
	})
	if err != nil {
		return nil, err
	}

	return resp.Product, nil
}

// Close closes the gRPC connection
func (c *ProductServiceClient) Close() error {
	return c.conn.Close()
}

// addAuthHeader adds authentication header to the context
func (c *ProductServiceClient) addAuthHeader(ctx context.Context) context.Context {
	if c.token != "" {
		md := metadata.Pairs("authorization", "Bearer "+c.token)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}