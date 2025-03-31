package main

import (
	"log"
	"net/http"

	pb "microservice-template/common/api"
	commonenv "microservice-template/common/env"

	_ "github.com/joho/godotenv/autoload" // package that loads env
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr         = commonenv.EnvString("PORT", "2220")
	orderServiceAddr = commonenv.EnvString("GRPC_ADDR", "2221")
)

func main() {
	// --- setup grpc connection ---

	// -- order --
	conn, err := grpc.Dial("localhost:"+orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	defer conn.Close()

	if err != nil {
		log.Fatalf("Could not connect to Order Service on port: %s", orderServiceAddr)
	}

	// new order service client - sets up a client that knows how to talk to the other service
	c := pb.NewOrderServiceClient(conn)

	// setup order client handler
	handler := NewHttpHandler(c)

	// routes
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/customers/{customerID}/orders", handler.handleCreateOrder)

	// start server
	log.Printf("Server started on port %s", httpAddr)

	if err := http.ListenAndServe(":"+httpAddr, mux); err != nil {
		log.Fatal("Failed to start server")
	}
}
