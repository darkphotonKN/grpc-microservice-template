package main

import (
	"log"
	"net/http"

	commonenv "microservice-template/common/env"
	"microservice-template/gateway/internal/order"

	_ "github.com/joho/godotenv/autoload" // package that loads env
)

var (
	httpAddr         = commonenv.EnvString("GATEWAY_ADDR", "2220")
	orderServiceAddr = commonenv.EnvString("GRPC_ORDER_ADDR", "2223")
)

func main() {
	// --- setup grpc connection ---

	// -- order --
	orderClient, err := order.NewClient(orderServiceAddr) // sets up client
	handler := order.NewHandler(orderClient)              // sets up handler for grpc requests

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
