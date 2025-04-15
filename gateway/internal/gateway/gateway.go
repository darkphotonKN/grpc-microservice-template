package gateway

import (
	"context"
	pb "microservice-template/common/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrders(context.Context) ([]*pb.Order, error)
}
