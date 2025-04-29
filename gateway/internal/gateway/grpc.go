package gateway

import (
	"context"
	"errors"
	"fmt"
	"log"
	pb "microservice-template/common/api"
	"microservice-template/common/discovery"
	commonerrors "microservice-template/common/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (g *gateway) GetOrderStatus(ctx context.Context, req *pb.OrderId) (*pb.OrderStatus, error) {
	// discovery
	conn, err := discovery.ServiceConnection(ctx, serviceName, g.registry)

	if err != nil {
		log.Fatalf("Failed to dial to server. Error: %s\n", err)
	}

	// create client to interface with through service discovery connection
	client := pb.NewOrderServiceClient(conn)
	orderStatus, err := client.GetOrderStatus(ctx, req)

	// custom error mapping
	if err != nil {
		if errors.Is(err, commonerrors.ErrNoItemFound) {
			return orderStatus, status.Error(codes.NotFound, "Order not found.")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get order status: %v", err)
	}

	fmt.Printf("getting order status:  %+v through gateway after service discovery\n", orderStatus)

	return orderStatus, nil
}

func (g *gateway) GetOrderPaymentLink(ctx context.Context, req *pb.OrderId) (*pb.OrderPaymentLink, error) {
	// discovery
	conn, err := discovery.ServiceConnection(ctx, serviceName, g.registry)

	if err != nil {
		log.Fatalf("Failed to dial to server. Error: %s\n", err)
	}

	// create client to interface with through service discovery connection
	client := pb.NewOrderServiceClient(conn)
	orderPaymentLink, err := client.GetOrderPaymentLink(ctx, req)

	// custom error mapping
	if err != nil {
		if errors.Is(err, commonerrors.ErrNoItemFound) {
			return orderPaymentLink, status.Error(codes.NotFound, "Order not found.")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get order payment link: %v", err)
	}

	fmt.Printf("getting order paymentLink:  %+v through gateway after service discovery\n", orderPaymentLink)

	return orderPaymentLink, nil
}
