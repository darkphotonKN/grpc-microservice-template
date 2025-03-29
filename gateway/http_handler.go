package main

import (
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

func (h *HttpHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating Order")

	var items []*pb.ItemsWithQuantity

	items = append(items, &pb.ItemsWithQuantity{
		ID:       "1",
		Quantity: 10,
	})

	items = append(items, &pb.ItemsWithQuantity{
		ID:       "2",
		Quantity: 7,
	})

	h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: "user1",
		Items:      items,
	})

}
