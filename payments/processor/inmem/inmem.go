package inmem

import (
	common "microservice-template/common/api"
	"microservice-template/payments/processor"
)

type InMemProcessor struct {
}

func NewInMemProcessor() processor.PaymentProcessor {
	return &InMemProcessor{}
}

func (p *InMemProcessor) CreatePaymentLink(*common.Order) (string, error) {
	return "test link", nil
}
