package order

import (
	"context"
	"fmt"
	pb "microservice-template/common/api"

	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrderService
}

func NewGrpcHandler(grpcServer *grpc.Server, service OrderService) {
	newGrpcHandler := grpcHandler{
		service: service,
	}

	// create server
	pb.RegisterOrderServiceServer(grpcServer, &newGrpcHandler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	fmt.Println("Order received!")

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
