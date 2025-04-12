package main

import (
	"context"
	"log"
	"microservice-template/common/broker"
	"microservice-template/common/discovery"
	"microservice-template/common/discovery/consul"
	commonenv "microservice-template/common/env"
	"microservice-template/context/internal/notification"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload" // package that loads env
	"google.golang.org/grpc"
)

var (
	serviceName  = "notification"
	grpcAddr     = commonenv.EnvString("GRPC_ORDER_ADDR", "2229")
	amqpUser     = commonenv.EnvString("RABBITMQ_USER", "guest")
	amqpPassword = commonenv.EnvString("RABBITMQ_PASS", "guest")
	amqpHost     = commonenv.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort     = commonenv.EnvString("RABBITMQ_PORT", "5672")
	consulAddr   = commonenv.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	// --- service discovery ---

	// -- setup --
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	// -- register --
	if err := registry.Register(ctx, instanceID, serviceName, "localhost:"+grpcAddr); err != nil {
		panic(err)
	}

	// -- health check --
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("Health check failed.")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	// --- message broker ---
	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)

	defer func() {
		close()
		ch.Close()
	}()

	// --- server initialization ---
	grpcServer := grpc.NewServer()
	paymentService := notification.NewService()
	paymentConsumer := notification.NewConsumer(paymentService, ch)
	paymentConsumer.Listen()

	// create a network listener to this service
	l, err := net.Listen("tcp", "localhost:"+grpcAddr)

	if err != nil {
		log.Fatalf(
			"Failed to listen at port: %s\nError: %s\n", grpcAddr, err,
		)
	}

	defer l.Close()

	log.Printf("grpc Order Server started on PORT: %s\n", grpcAddr)
	// start serving requests
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal("Can't connect to grpc server. Error:", err.Error())
	}
}
