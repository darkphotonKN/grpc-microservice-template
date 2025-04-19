package order

import (
	"context"
	pb "microservice-template/common/api"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrders(ctx context.Context, empty *emptypb.Empty) (*pb.Orders, error)
	GetOrderStatus(ctx context.Context, req *pb.OrderId) (OrderStatus, error)
	ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) error
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order Order) (uuid.UUID, error)
	CreateOrderTx(ctx context.Context, tx *sqlx.Tx, order Order) (uuid.UUID, error)
	CreateOrderItem(ctx context.Context, item OrderItem) error
	CreateOrderItemTx(ctx context.Context, tx *sqlx.Tx, item OrderItem) error
	GetAll(ctx context.Context) ([]*pb.Order, error)
}

// Entity

type Order struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CustomerID string    `json:"customer_id" db:"customer_id"`
	Status     int       `json:"status" db:"status"`
}

type OrderItem struct {
	ID       uuid.UUID `json:"id" db:"id"`
	OrderID  uuid.UUID `json:"order_id" db:"order_id"`
	Name     string    `json:"name" db:"name"`
	Quantity int       `json:"quantity" db:"quantity"`
	PriceID  string    `json:"price_id" db:"price_id"`
}

// Shared Types
type OrderStatus string

const (
	pending OrderStatus = "pending"
	paid    OrderStatus = "paid"
)
