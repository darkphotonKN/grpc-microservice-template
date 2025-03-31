package order

import (
	"encoding/json"
	"fmt"
	pb "microservice-template/common/api"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler handles HTTP requests for orders
type Handler struct {
	client *Client
}

// NewHandler creates a new order handler
func NewHandler(client *Client) *Handler {
	return &Handler{
		client: client,
	}
}

// HandleCreateOrder creates a new order
func (h *Handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating Order")

	var items []*pb.ItemsWithQuantity
	err := json.NewDecoder(r.Body).Decode(&items)

	if err != nil {
		fmt.Println("Error when creating order:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	customerID := r.PathValue("customerID")

	// call order service
	order, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
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

	fmt.Printf("Successfully created order %+v", order)

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
