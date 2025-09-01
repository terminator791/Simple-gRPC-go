package healthservice

import (
	"context"
	"database/sql"
	"sync"

	"github.com/terminator791/Simple-gRPC-go/pkg/health"
)

// HealthChecker interface for checking service health
type HealthChecker interface {
	Check(ctx context.Context, service string) (*health.HealthCheckResponse, error)
}

// DatabaseHealthChecker checks database connectivity
type DatabaseHealthChecker struct {
	db *sql.DB
}

// NewDatabaseHealthChecker creates a new database health checker
func NewDatabaseHealthChecker(db *sql.DB) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{db: db}
}

// Check checks database connectivity
func (d *DatabaseHealthChecker) Check(ctx context.Context, service string) (*health.HealthCheckResponse, error) {
	if err := d.db.PingContext(ctx); err != nil {
		return &health.HealthCheckResponse{
			Status:  health.HealthCheckResponse_NOT_SERVING,
			Message: "Database connection failed: " + err.Error(),
		}, nil
	}

	return &health.HealthCheckResponse{
		Status:  health.HealthCheckResponse_SERVING,
		Message: "Database is healthy",
	}, nil
}

// HealthServiceServer implements the health service
type HealthServiceServer struct {
	health.UnimplementedHealthServiceServer
	checkers map[string]HealthChecker
	mu       sync.RWMutex
}

// NewHealthServiceServer creates a new health service server
func NewHealthServiceServer() *HealthServiceServer {
	return &HealthServiceServer{
		checkers: make(map[string]HealthChecker),
	}
}

// RegisterChecker registers a health checker for a service
func (h *HealthServiceServer) RegisterChecker(serviceName string, checker HealthChecker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers[serviceName] = checker
}

// Check performs a health check
func (h *HealthServiceServer) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// If no service is specified, check overall health
	if req.Service == "" {
		return h.checkOverallHealth(ctx)
	}

	// Check specific service
	checker, exists := h.checkers[req.Service]
	if !exists {
		return &health.HealthCheckResponse{
			Status:  health.HealthCheckResponse_SERVICE_UNKNOWN,
			Message: "Service not found: " + req.Service,
		}, nil
	}

	return checker.Check(ctx, req.Service)
}

// Watch streams health status changes (simplified implementation)
func (h *HealthServiceServer) Watch(req *health.HealthCheckRequest, stream health.HealthService_WatchServer) error {
	// For simplicity, we'll just send the current status once
	// In a real implementation, you would watch for changes and stream updates
	resp, err := h.Check(stream.Context(), req)
	if err != nil {
		return err
	}

	return stream.Send(resp)
}

// checkOverallHealth checks the health of all registered services
func (h *HealthServiceServer) checkOverallHealth(ctx context.Context) (*health.HealthCheckResponse, error) {
	allHealthy := true
	messages := []string{}

	for serviceName, checker := range h.checkers {
		resp, err := checker.Check(ctx, serviceName)
		if err != nil {
			allHealthy = false
			messages = append(messages, serviceName+": error - "+err.Error())
			continue
		}

		if resp.Status != health.HealthCheckResponse_SERVING {
			allHealthy = false
			messages = append(messages, serviceName+": "+resp.Message)
		}
	}

	if allHealthy {
		return &health.HealthCheckResponse{
			Status:  health.HealthCheckResponse_SERVING,
			Message: "All services are healthy",
		}, nil
	}

	return &health.HealthCheckResponse{
		Status:  health.HealthCheckResponse_NOT_SERVING,
		Message: "Some services are unhealthy",
	}, nil
}