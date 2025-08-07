package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mlucas4330/orderflow-pro/pkg/model"
	pb "github.com/mlucas4330/orderflow-pro/pkg/productpb"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order *model.Order, items []model.OrderItem) error {
	args := m.Called(ctx, order, items)
	return args.Error(0)
}

func (m *MockOrderRepository) FindOrderById(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) FindOrders(ctx context.Context) ([]model.Order, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateOrder(ctx context.Context, id uuid.UUID, status model.Status) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

type MockIdempotencyRepository struct {
	mock.Mock
}

func (m *MockIdempotencyRepository) GetResponse(ctx context.Context, key, userID uuid.UUID) (*model.IdempotencyResponse, error) {
	args := m.Called(ctx, key, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.IdempotencyResponse), args.Error(1)
}
func (m *MockIdempotencyRepository) SaveResponse(ctx context.Context, key, userID uuid.UUID, response *model.IdempotencyResponse) error {
	args := m.Called(ctx, key, userID, response)
	return args.Error(0)
}

type MockProductServiceClient struct {
	mock.Mock
}

func (m *MockProductServiceClient) GetProductDetails(ctx context.Context, req *pb.GetProductDetailsRequest, opts ...grpc.CallOption) (*pb.GetProductDetailsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetProductDetailsResponse), args.Error(1)
}
