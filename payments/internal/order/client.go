package order

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
)

const (
	serviceName = "orders"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) *Client {
	return &Client{
		registry: registry,
	}
}

func (c *Client) UpdateOrderStatus(ctx context.Context, req *pb.OrderStatusUpdateRequest) (*pb.Order, error) {
	// discovery
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		log.Fatalf("Failed to dial to server. Error: %s\n", err)
	}

	// create client to interface with through service discovery connection
	client := pb.NewOrderServiceClient(conn)
	orderStatus, err := client.UpdateOrderStatus(ctx, req)

	// custom error mapping
	if err != nil {
		if errors.Is(err, commonerrors.ErrNoItemFound) {
			return orderStatus, status.Error(codes.NotFound, "Order not found.")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get order status: %v", err)
	}

	fmt.Printf("Updating order status %+v through payment service after service discovery\n", orderStatus)

	return orderStatus, nil
}
