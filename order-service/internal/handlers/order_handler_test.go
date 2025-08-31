package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/terminator791/Simple-gRPC-go/pkg/order"
	"github.com/terminator791/Simple-gRPC-go/pkg/user"
	"github.com/terminator791/Simple-gRPC-go/order-service/internal/models"
)

// MockUserValidator is a mock implementation of UserValidator interface
type MockUserValidator struct {
	mock.Mock
}

func (m *MockUserValidator) ValidateUser(ctx context.Context, userID int64) (*user.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

// MockOrderRepository is a mock implementation of OrderRepositoryInterface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(userID int64, productName string, amount float64, quantity int32) (*models.Order, error) {
	args := m.Called(userID, productName, amount, quantity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByID(orderID int64) (*models.Order, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func TestCreateOrder_Success(t *testing.T) {
	// Arrange
	mockUserValidator := new(MockUserValidator)
	mockOrderRepo := new(MockOrderRepository)
	
	server := NewOrderServiceServer(mockOrderRepo, mockUserValidator)
	
	ctx := context.Background()
	req := &order.CreateOrderRequest{
		UserId:      1,
		ProductName: "Test Product",
		Amount:      99.99,
		Quantity:    2,
	}
	
	expectedUser := &user.User{
		Id:    1,
		Email: "test@example.com",
		Name:  "Test User",
	}
	
	expectedOrder := &models.Order{
		ID:          1,
		UserID:      1,
		ProductName: "Test Product",
		Amount:      99.99,
		Quantity:    2,
		Status:      "pending",
	}
	
	// Mock expectations
	mockUserValidator.On("ValidateUser", ctx, int64(1)).Return(expectedUser, nil)
	mockOrderRepo.On("Create", int64(1), "Test Product", 99.99, int32(2)).Return(expectedOrder, nil)
	
	// Act
	response, err := server.CreateOrder(ctx, req)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.Order)
	assert.Equal(t, int64(1), response.Order.Id)
	assert.Equal(t, int64(1), response.Order.UserId)
	assert.Equal(t, "Test Product", response.Order.ProductName)
	assert.Equal(t, 99.99, response.Order.Amount)
	assert.Equal(t, int32(2), response.Order.Quantity)
	assert.Equal(t, "pending", response.Order.Status)
	
	// Verify mock expectations
	mockUserValidator.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func TestCreateOrder_UserNotFound(t *testing.T) {
	// Arrange
	mockUserValidator := new(MockUserValidator)
	mockOrderRepo := new(MockOrderRepository)
	
	server := NewOrderServiceServer(mockOrderRepo, mockUserValidator)
	
	ctx := context.Background()
	req := &order.CreateOrderRequest{
		UserId:      999,
		ProductName: "Test Product",
		Amount:      99.99,
		Quantity:    2,
	}
	
	// Mock user not found error
	userNotFoundErr := status.Error(codes.NotFound, "user not found")
	mockUserValidator.On("ValidateUser", ctx, int64(999)).Return(nil, userNotFoundErr)
	
	// Act
	response, err := server.CreateOrder(ctx, req)
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	
	// Check that the error is FailedPrecondition (user not found)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
	assert.Contains(t, st.Message(), "user not found")
	
	// Verify that Create was NOT called since user validation failed
	mockUserValidator.AssertExpectations(t)
	mockOrderRepo.AssertNotCalled(t, "Create")
}

func TestCreateOrder_InvalidInput(t *testing.T) {
	// Arrange
	mockUserValidator := new(MockUserValidator)
	mockOrderRepo := new(MockOrderRepository)
	
	server := NewOrderServiceServer(mockOrderRepo, mockUserValidator)
	
	ctx := context.Background()
	
	testCases := []struct {
		name        string
		request     *order.CreateOrderRequest
		expectedErr string
	}{
		{
			name: "Invalid user ID",
			request: &order.CreateOrderRequest{
				UserId:      0,
				ProductName: "Test Product",
				Amount:      99.99,
				Quantity:    2,
			},
			expectedErr: "user_id must be positive",
		},
		{
			name: "Empty product name",
			request: &order.CreateOrderRequest{
				UserId:      1,
				ProductName: "",
				Amount:      99.99,
				Quantity:    2,
			},
			expectedErr: "product_name is required",
		},
		{
			name: "Invalid amount",
			request: &order.CreateOrderRequest{
				UserId:      1,
				ProductName: "Test Product",
				Amount:      0,
				Quantity:    2,
			},
			expectedErr: "amount must be positive",
		},
		{
			name: "Invalid quantity",
			request: &order.CreateOrderRequest{
				UserId:      1,
				ProductName: "Test Product",
				Amount:      99.99,
				Quantity:    0,
			},
			expectedErr: "quantity must be positive",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			response, err := server.CreateOrder(ctx, tc.request)
			
			// Assert
			assert.Error(t, err)
			assert.Nil(t, response)
			
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, codes.InvalidArgument, st.Code())
			assert.Contains(t, st.Message(), tc.expectedErr)
		})
	}
	
	// Verify no mock methods were called for invalid input
	mockUserValidator.AssertNotCalled(t, "ValidateUser")
	mockOrderRepo.AssertNotCalled(t, "Create")
}