package main

import (
	"context"
	"log"
	commonenv "microservice-template/common/env"
	"microservice-template/orders/internal/order"
	"net"

	_ "github.com/joho/godotenv/autoload" // package that loads env
	"google.golang.org/grpc"
)

var (
	grpcAddr = commonenv.EnvString("GRPC_ORDER_ADDR", "2221")
)

func main() {
	// create the grpc server instance
	grpcServer := grpc.NewServer()

	// create a network listener to this service
	l, err := net.Listen("tcp", "localhost:"+grpcAddr)

	if err != nil {
		log.Fatalf(
			"Failed to listen at port: %s\nError: %s\n", grpcAddr, err,
		)
	}

	defer l.Close()

	// service setup
	repo := order.NewRepository()
	service := order.NewService(repo)

	// start grpc server
	order.NewGrpcHandler(grpcServer, service)

	service.CreateOrder(context.Background())

	log.Printf("grpc Order Server started on PORT: %s\n", grpcAddr)
	// start serving requests
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal("Can't connect to grpc server. Error:", err.Error())
	}
}
