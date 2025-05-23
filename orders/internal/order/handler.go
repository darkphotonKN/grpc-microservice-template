package order

import (
	"context"
	pb "microservice-template/common/api"

	"google.golang.org/protobuf/types/known/emptypb"
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
	return h.service.CreateOrder(ctx, req)
}

func (h *grpcHandler) GetOrders(ctx context.Context, empty *emptypb.Empty) (*pb.Orders, error) {

	return h.service.GetOrders(ctx, empty)
}

func (h *grpcHandler) GetOrderStatus(ctx context.Context, req *pb.OrderId) (*pb.OrderStatus, error) {

	return h.service.GetOrderStatus(ctx, req)
}

func (h *grpcHandler) UpdateOrderStatus(ctx context.Context, req *pb.OrderStatusUpdateRequest) (*pb.Order, error) {

	return h.service.UpdateOrderStatus(ctx, req)
}

func (h *grpcHandler) UpdateOrderPaymentLink(ctx context.Context, req *pb.OrderPaymentUpdateRequest) (*pb.Order, error) {
	return h.service.UpdateOrderPaymentLink(ctx, req)
}

func (h *grpcHandler) GetOrderPaymentLink(ctx context.Context, req *pb.OrderId) (*pb.OrderPaymentLink, error) {
	return h.service.GetOrderPaymentLink(ctx, req)
}
