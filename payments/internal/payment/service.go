package payment

import (
	"context"
	pb "microservice-template/common/api"
	"microservice-template/payments/processor"

	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	// stripe service injection
	paymentProcessor processor.PaymentProcessor

	// channel for communicating on message broker
	ch *amqp.Channel
}

func NewService(paymentProcessor processor.PaymentProcessor) *service {
	return &service{
		paymentProcessor: paymentProcessor,
	}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	link, err := s.paymentProcessor.CreatePaymentLink(order)

	if err != nil {
		return "", err
	}

	return link, nil
}
