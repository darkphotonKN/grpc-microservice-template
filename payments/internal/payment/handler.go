package payment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78/webhook"
)

type handler struct {
	service PaymentService
}

func NewHandler(service PaymentService) PaymentHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) HandleStripeWebhook(c *gin.Context) {
	// read webhook response data
	body, err := c.GetRawData()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error when attempting fetch body of webhook." + err.Error()})
		return
	}

	stripeSignature := c.GetHeader("Stripe-Signature")

	// verification
	stripeWebhookSecret := h.service.GetWebhookSecret()

	event, err := webhook.ConstructEvent(body, stripeSignature, stripeWebhookSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook verification failed." + err.Error()})
		return
	}

	switch event.Type {
	// succesfully completed payment
	case "checkout.session.completed":
		fmt.Println("a client successfully completed a stripe payment.")
		var session map[string]interface{}
		json.Unmarshal(event.Data.Raw, &session)

		// Extract order ID from metadata
		metadata := session["metadata"].(map[string]interface{})
		orderID := metadata["order_id"].(string)
		fmt.Println("orderID:", orderID)

		// Update order status
		// You'll handle this via gRPC to Order Service

		// Publish payment completed event
		// This is what you're already familiar with
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
