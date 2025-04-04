package main

import (
	"context"
	"fmt"
	"log"
	pb "microservice-template/common/api"
	"microservice-template/common/discovery"
	"microservice-template/common/discovery/consul"
	commonenv "microservice-template/common/env"
	"microservice-template/orders/internal/order"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload" // package that loads env
	"google.golang.org/grpc"
)

var (
	grpcAddr     = commonenv.EnvString("GRPC_ORDER_ADDR", "2221")
	amqpUser     = commonenv.EnvString("RABBITMQ_USER", "guest")
	amqpPassword = commonenv.EnvString("RABBITMQ_USER", "guest")
	amqpHost     = commonenv.EnvString("RABBITMQ_USER", "localhost")
	amqpPort     = commonenv.EnvString("RABBITMQ_USER", "5672")
	consulAddr   = commonenv.EnvString("CONSUL_ADDR", "localhost:8500")
	serviceName  = "orders"
)

func main() {
	// --- service discovery setup ---
	registry, err := consul.NewRegistry(consulAddr, serviceName)

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	// -- discovery --
	if err := registry.Register(ctx, instanceID, serviceName, "localhost:"+grpcAddr); err != nil {
		// TODO: REMOVE AFTER DEBUG
		fmt.Printf("\nError when registering service:\n\n%s\n\n", err)
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("Health check failed.")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	// --- server initialization ---
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
	handler := order.NewGrpcHandler(service)

	// create server
	pb.RegisterOrderServiceServer(grpcServer, handler)

	log.Printf("grpc Order Server started on PORT: %s\n", grpcAddr)
	// start serving requests
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal("Can't connect to grpc server. Error:", err.Error())
	}

}
