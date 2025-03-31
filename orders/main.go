package main

import (
	"context"
	"log"
	"microservice-template/common"
	"net"

	"google.golang.org/grpc"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "2221")
)

func main() {
	// create the grpc server instance
	grpcServer := grpc.NewServer()

	// create a network listener
	l, err := net.Listen("tcp", "localhost:"+grpcAddr)

	if err != nil {
		log.Fatalf(
			"Failed to listen at port: %s\nError: %s\n", grpcAddr, err,
		)
	}

	defer l.Close()

	// service setup
	repo := NewRepository()
	service := NewService(repo)

	// start grpc server
	NewGrpcHandler(grpcServer, service)

	service.CreateOrder(context.Background())

	log.Printf("grpc Order Server started on PORT: %s\n", grpcAddr)
	// start serving requests
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal("Can't connect to grpc server. Error:", err.Error())
	}
}
