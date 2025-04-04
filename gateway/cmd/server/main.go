package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"microservice-template/common/discovery"
	"microservice-template/common/discovery/consul"
	commonenv "microservice-template/common/env"
	"microservice-template/gateway/internal/gateway"
	"microservice-template/gateway/internal/order"

	_ "github.com/joho/godotenv/autoload" // package that loads env
)

var (
	httpAddr         = commonenv.EnvString("GATEWAY_ADDR", "2220")
	orderServiceAddr = commonenv.EnvString("GRPC_ORDER_ADDR", "2223")
	consulAddr       = commonenv.EnvString("CONSUL_ADDR", "localhost:8500")
	serviceName      = "gateway"
)

func main() {
	// --- service discovery setup ---

	// -- consul client --
	registry, err := consul.NewRegistry(consulAddr, serviceName)

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	// -- discovery --
	if err := registry.Register(ctx, instanceID, serviceName, "localhost:"+httpAddr); err != nil {
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

	// --- setup grpc connection ---
	// sets up grpc connection with registry from service discovery injected

	// TODO: REMOVE AFTER DEBUG
	fmt.Printf("\ncurrent registry: \n\n%+v\n\n", *registry)

	orderGateway := gateway.NewGRPCGateway(registry)
	handler := order.NewHandler(orderGateway)

	if err != nil {
		log.Printf("Error occured when attempting to establish grpc connection to order service through the gateway service: %s", err)
	}

	// --- routes ---
	mux := http.NewServeMux()

	// -- interface requests to FE client via HTTP --
	mux.HandleFunc("POST /api/customers/{customerID}/orders", handler.HandleCreateOrder)

	// --- server initialization ---
	log.Printf("Server started on port %s", httpAddr)

	// -- start server and capture errors --
	if err := http.ListenAndServe(":"+httpAddr, mux); err != nil {
		log.Fatal("Failed to start server")
	}
}
