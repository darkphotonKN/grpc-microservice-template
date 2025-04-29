package gateway

import (
	"context"
	pb "microservice-template/common/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrders(context.Context) (*pb.Orders, error)
	GetOrderStatus(ctx context.Context, req *pb.OrderId) (*pb.OrderStatus, error)
	GetOrderPaymentLink(ctx context.Context, req *pb.OrderId) (*pb.OrderPaymentLink, error)
}
