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
	fmt.Println("Webhook called.")

	// Log all headers for debugging
	fmt.Println("All request headers:")
	for key, values := range c.Request.Header {
		fmt.Printf("  %s: %v\n", key, values)
	}

	// read webhook response data
	body, err := c.GetRawData()

	if err != nil {
		fmt.Printf("\nErr when parsing from body: %+v\n\n", err)

		c.JSON(http.StatusBadRequest, gin.H{"error": "Error when attempting fetch body of webhook." + err.Error()})
		return
	}

	// Log all headers for debugging
	fmt.Println("All request headers:")
	for key, values := range c.Request.Header {
		fmt.Printf("  %s: %v\n", key, values)
	}

	stripeSignature := c.GetHeader("Stripe-Signature")

	fmt.Printf("stripeSignature being used: %s\n", stripeSignature)

	// verification
	stripeWebhookSecret := h.service.GetWebhookSecret()

	fmt.Printf("Webhook secret being used: %s\n", stripeWebhookSecret)

	event, err := webhook.ConstructEventWithOptions(body, stripeSignature, stripeWebhookSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})

	if err != nil {
		fmt.Printf("\nErr when acquiring event from construct event method: %+v\n\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook verification failed." + err.Error()})
		return
	}

	fmt.Printf("\nEvent from response: %+v\n\n", event)

	switch event.Type {
	// succesfully completed payment
	case "checkout.session.completed":
		fmt.Println("a client successfully completed a stripe payment.")

		// --- extract data ---
		var data map[string]interface{}

		// -- raw event data --
		if err := json.Unmarshal(event.Data.Raw, &data); err != nil {
			fmt.Printf("Error when trying to unmarshal data from payment success: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error when trying to unmarshal data from payment success: %v", err)})
			return
		}

		fmt.Printf("\nUnmarshalled Data: %+v\n\n", data)

		// -- metadata --
		metadata, ok := data["metadata"].(map[string]interface{})

		if !ok {
			fmt.Println("No metadata found in session or wrong format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metadata in session"})
			return
		}

		fmt.Printf("\nUnmarshalled Metadata: %+v\n\n", metadata)

		// -- order id --
		orderID, ok := metadata["orderID"].(string)

		if !ok {
			fmt.Println("No orderID found found in the metadata.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metadata in session"})
			return
		}

		fmt.Printf("orderID from response %s\n", orderID)

		// update order status to order service
		_, err := h.service.UpdateOrderStatus(c.Request.Context(), UpdateOrderStatus{OrderId: orderID})
		if err != nil {
			fmt.Println("Error on order status update from payment handler.")
			c.JSON(http.StatusBadRequest, gin.H{"result": "Error on order status update from payment handler."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "success"})
	}

	fmt.Println("passed check with no matching event type.")
}
