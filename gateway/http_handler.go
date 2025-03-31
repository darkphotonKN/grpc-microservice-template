package main

import (
	"encoding/json"
	"fmt"
	pb "microservice-template/common/api"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	order, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerId,
		Items:      items,
	})

	// handle errors from GRPC, using grpc's status convert helper
	errStatus := status.Convert(err)

	if errStatus != nil {
		// matching for invalid argument with grpc's codes helper
		if errStatus.Code() != codes.InvalidArgument {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error when attempting to create an order: " + err.Error()))
			return
		}

		// other error codes
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error when attempting to create an order: " + err.Error()))
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error when attempting to create an order: " + err.Error()))
		return
	}

	fmt.Println("Order created:", order)
}
