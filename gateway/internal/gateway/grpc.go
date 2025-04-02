package gateway

import (
	"context"
	"log"
	pb "microservice-template/common/api"
	"microservice-template/common/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func (g *gateway) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial to server.")
	}
	client := pb.NewOrderServiceClient(conn)

	return client.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})
}
