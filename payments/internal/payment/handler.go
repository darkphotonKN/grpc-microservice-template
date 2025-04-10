package payment

import (
	"context"
	"fmt"
	pb "microservice-template/common/api"
)

type grpcHandler struct {
	// TODO: update to payment
	pb.UnimplementedPaymentServiceServer
	service PaymentService
}

func NewGrpcHandler(service PaymentService) *grpcHandler {
	return &grpcHandler{
		service: service,
	}
}

func (h *grpcHandler) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	fmt.Println("Order received!")
	return h.service.CreatePayment(ctx, order)
}
