package processor

import (
	pb "microservice-template/common/api"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

type PaymentProcessor interface {
	CreatePaymentLink(o *pb.Order) (string, error)
}

type stripeProcessor struct {
}

func NewStripeProcessor() PaymentProcessor {
	return &stripeProcessor{}
}

func (s *stripeProcessor) CreatePaymentLink(o *pb.Order) (string, error) {

	// create stripe line items from our order
	var items []*stripe.CheckoutSessionLineItemParams

	for _, item := range o.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("https://example.com/success"),
	}

	result, err := session.New(params)

	if err != nil {
		return "", err
	}

	return result.URL, nil

}
