package gateway

import (
	"context"
	"fmt"
	"log"
	pb "microservice-template/common/api"
	"microservice-template/common/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) OrdersGateway {
	return &gateway{
		registry: registry,
	}
}

func (g *gateway) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {

	// connection instance created through service discovery first
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)

	if err != nil {
		log.Fatalf("Failed to dial to server.")
	}

	// create client to interface with through service discovery connection
	client := pb.NewOrderServiceClient(conn)
	order, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})

	fmt.Printf("Creating order %+v through gateway after service discovery\n", order)

	return order, nil
}
