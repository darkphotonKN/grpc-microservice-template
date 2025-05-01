package stripe

import (
	"fmt"
	pb "microservice-template/common/api"
	commonenv "microservice-template/common/env"
	"microservice-template/payments/processor"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

var (
	httpGatewayAddr = commonenv.EnvString("GATEWAY_HTTP_ADDRESS", "http://localhost:3009")
)

type stripeProcessor struct {
}

func NewStripeProcessor() processor.PaymentProcessor {
	return &stripeProcessor{}
}

func (s *stripeProcessor) CreatePaymentLink(o *pb.Order) (string, error) {

	gatewaySuccessUrl := fmt.Sprintf("%s/success?orderId=%s&customerId=%s", httpGatewayAddr, o.ID, o.CustomerID)

	// create stripe line items from our order
	var items []*stripe.CheckoutSessionLineItemParams

	for _, item := range o.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			// prod_S5sRbXrdUHcRfd
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gatewaySuccessUrl),
		Metadata:   map[string]string{"orderID": o.ID},
	}

	result, err := session.New(params)

	if err != nil {
		return "", err
	}

	fmt.Printf("\nResult of payment link %+v\n\n", result)

	return result.URL, nil
}
