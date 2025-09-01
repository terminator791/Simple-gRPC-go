package client

import (
	"context"
	"time"

	"github.com/sony/gobreaker"
	"github.com/terminator791/Simple-gRPC-go/pkg/product"
)

// ProductClientWithCircuitBreaker wraps ProductClient with circuit breaker functionality
type ProductClientWithCircuitBreaker struct {
	client ProductClient
	cb     *gobreaker.CircuitBreaker
}

// NewProductClientWithCircuitBreaker creates a new product client with circuit breaker
func NewProductClientWithCircuitBreaker(address, serviceToken string) (*ProductClientWithCircuitBreaker, error) {
	client, err := NewProductServiceClient(address, serviceToken)
	if err != nil {
		return nil, err
	}

	// Configure circuit breaker
	settings := gobreaker.Settings{
		Name:        "ProductService",
		MaxRequests: 3,
		Interval:    30 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// Log state changes for monitoring
			// In production, this could send metrics to monitoring system
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)

	return &ProductClientWithCircuitBreaker{
		client: client,
		cb:     cb,
	}, nil
}

// CheckInventory with circuit breaker
func (c *ProductClientWithCircuitBreaker) CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, int32, error) {
	result, err := c.cb.Execute(func() (interface{}, error) {
		available, stock, err := c.client.CheckInventory(ctx, productID, quantity)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"available": available,
			"stock":     stock,
		}, nil
	})
	
	if err != nil {
		return false, 0, err
	}
	
	data := result.(map[string]interface{})
	return data["available"].(bool), data["stock"].(int32), nil
}

// ReserveInventory with circuit breaker
func (c *ProductClientWithCircuitBreaker) ReserveInventory(ctx context.Context, productID int64, quantity int32, reservationID string) (bool, string, error) {
	result, err := c.cb.Execute(func() (interface{}, error) {
		success, message, err := c.client.ReserveInventory(ctx, productID, quantity, reservationID)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": success,
			"message": message,
		}, nil
	})
	
	if err != nil {
		return false, "", err
	}
	
	data := result.(map[string]interface{})
	return data["success"].(bool), data["message"].(string), nil
}

// ReleaseInventory with circuit breaker
func (c *ProductClientWithCircuitBreaker) ReleaseInventory(ctx context.Context, reservationID string) (bool, string, error) {
	result, err := c.cb.Execute(func() (interface{}, error) {
		success, message, err := c.client.ReleaseInventory(ctx, reservationID)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": success,
			"message": message,
		}, nil
	})
	
	if err != nil {
		return false, "", err
	}
	
	data := result.(map[string]interface{})
	return data["success"].(bool), data["message"].(string), nil
}

// GetProduct with circuit breaker
func (c *ProductClientWithCircuitBreaker) GetProduct(ctx context.Context, productID int64) (*product.Product, error) {
	result, err := c.cb.Execute(func() (interface{}, error) {
		return c.client.GetProduct(ctx, productID)
	})
	
	if err != nil {
		return nil, err
	}
	
	return result.(*product.Product), nil
}

// Close closes the underlying client
func (c *ProductClientWithCircuitBreaker) Close() error {
	return c.client.Close()
}

// GetCircuitBreakerState returns the current state of the circuit breaker
func (c *ProductClientWithCircuitBreaker) GetCircuitBreakerState() gobreaker.State {
	return c.cb.State()
}