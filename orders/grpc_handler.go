package main

import (
	"context"
	"fmt"
	pb "microservice-template/common/api"

	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
}

func NewGrpcHandler(grpcServer *grpc.Server) {
	newGrpcHandler := grpcHandler{}

	// create server
	pb.RegisterOrderServiceServer(grpcServer, &newGrpcHandler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	order := &pb.Order{
		ID: "111",
	}

	fmt.Println("Order received:", order)

	return order, nil
}
