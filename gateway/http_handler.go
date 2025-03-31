package main

import (
	"encoding/json"
	"fmt"
	pb "microservice-template/common/api"
	"net/http"
)

type HttpHandler struct {
	client pb.OrderServiceClient
}

func NewHttpHandler(client pb.OrderServiceClient) *HttpHandler {
	return &HttpHandler{
		client: client,
	}
}

// makes an order to the order service
func (h *HttpHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating Order")

	var items []*pb.ItemsWithQuantity

	err := json.NewDecoder(r.Body).Decode(&items)

	if err != nil {
		fmt.Println("Error when creating order:", err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	customerId := r.PathValue("customerID")

	// makes order to order service through GRPC
	order, _ := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerId,
		Items:      items,
	})

	fmt.Println("Order created:", order)
}
