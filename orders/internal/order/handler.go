package order

import (
	"context"
	"fmt"
	pb "microservice-template/common/api"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrderService
}

func NewGrpcHandler(service OrderService) *grpcHandler {
	return &grpcHandler{
		service: service,
	}
}

func (h *grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	fmt.Println("Order received!")

	return h.service.CreateOrder(ctx, req)
}

func (h *grpcHandler) GetOrders(ctx context.Context) (*pb.Orders, error) {
	return h.service.GetOrders(ctx)
}
