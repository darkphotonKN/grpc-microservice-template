package order

import (
	"context"
	"fmt"
	pb "microservice-template/common/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

/*
- Establishing the gRPC connection to the order service
- Creating and managing the gRPC client stub
- Providing a clean interface for making remote procedure calls
- Handling connection lifecycle (initialization, closing)

This of this file like an adapter that transforms your local functions
into remote service calls.

Our handler file will then call this like "normal" functions, and the grpc
stuff is abstracted away here.
*/

type Client struct {
	conn   *grpc.ClientConn
	client pb.OrderServiceClient
}

func NewClient(addr string) (*Client, error) {

	fullAddr := fmt.Sprintf("localhost:%s", addr)

	conn, err := grpc.Dial(fullAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	newGrpcClient := pb.NewOrderServiceClient(conn)

	return &Client{
		conn:   conn,
		client: newGrpcClient,
	}, nil
}

// CreateOrder forwards a create order request to the order service
func (c *Client) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	return c.client.CreateOrder(ctx, req)
}

func (c *Client) GetOrders(ctx context.Context) (*pb.Orders, error) {
	return c.client.GetOrders(ctx, &emptypb.Empty{})
}

func (c *Client) GetOrderStatus(ctx context.Context, orderId *pb.OrderId) (*pb.OrderStatus, error) {
	return c.client.GetOrderStatus(ctx, orderId)
}
