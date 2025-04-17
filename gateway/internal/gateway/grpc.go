package gateway

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	pb "microservice-template/common/api"
	"microservice-template/common/discovery"
)

const (
	serviceName = "orders"
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
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, g.registry)

	if err != nil {
		log.Fatalf("Failed to dial to server. Error: %s\n", err)
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

func (g *gateway) GetOrders(ctx context.Context) (*pb.Orders, error) {
	// discovery
	conn, err := discovery.ServiceConnection(ctx, serviceName, g.registry)

	if err != nil {
		log.Fatalf("Failed to dial to server. Error: %s\n", err)
	}

	// create client to interface with through service discovery connection
	client := pb.NewOrderServiceClient(conn)
	order, err := client.GetOrders(ctx, &emptypb.Empty{})

	fmt.Printf("Creating order %+v through gateway after service discovery\n", order)

	return order, nil
}
