package order

import (
	"context"
	pb "microservice-template/common/api"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrders(ctx context.Context, empty *emptypb.Empty) (*pb.Orders, error)
	ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) error
}

type OrderRepository interface {
	Create(ctx context.Context) error
	Get(ctx context.Context, id, customerId string) (*pb.Order, error)
}

// Entity

type Order struct {
	ID uuid.UUID `json:"id" db:"id"`
}
