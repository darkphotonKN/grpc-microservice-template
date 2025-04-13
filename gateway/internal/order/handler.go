package order

import (
	"fmt"
	pb "microservice-template/common/api"
	"microservice-template/gateway/internal/gateway"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gateway gateway.OrdersGateway
}

func NewHandler(gateway gateway.OrdersGateway) *Handler {
	return &Handler{
		// inject gatway here
		gateway: gateway,
	}
}

func (h *Handler) HandleCreateOrder(c *gin.Context) {
	fmt.Println("Creating Order")

	var items []*pb.ItemsWithQuantity

	if err := c.ShouldBindJSON(&items); err != nil {

		fmt.Println("Error when creating order:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	customerID := c.Param("customerID")

	// call order service
	order, err := h.gateway.CreateOrder(c.Request.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

	// handle errors from GRPC, using grpc's status convert helper
	errStatus := status.Convert(err)

	if errStatus != nil {
		// matching for invalid argument with grpc's codes helper
		if errStatus.Code() != codes.InvalidArgument {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error when attempting to create an order:" + err.Error()})
			return
		}

		// other error codes
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when attempting to create an order: " + err.Error()})
	}

	fmt.Printf("Successfully created order %+v", order)
	c.JSON(http.StatusOK, gin.H{"result": order})
}
