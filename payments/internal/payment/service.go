package payment

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	pb "microservice-template/common/api"
)

type service struct {
	// stripe service injection

	// channel for communicating on message broker
	ch *amqp.Channel
}

func NewService() *service {
	return &service{}
}

func (s *service) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.Payment, error) {

	// TODO: connect to payment service

	return nil, nil
}
