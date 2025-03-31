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

	// validation
	err := h.service.ValidateOrder(ctx, req)

	if err != nil {
		return nil, err
	}

	items := make([]*pb.Item, len(req.Items))

	for index, item := range req.Items {
		items[index] = &pb.Item{
			ID:       item.ID,
			Name:     "testeritem",
			Quantity: item.Quantity,
			PriceID:  "rando",
		}
	}

	order := &pb.Order{
		ID:         "1",
		CustomerID: req.CustomerID,
		Status:     "initiated",
		Items:      items,
	}

	fmt.Printf("Outgoing order: %+v\n", order)

	return order, nil
}
