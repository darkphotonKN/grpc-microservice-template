package payment

import (
	"context"
	"encoding/json"
	"fmt"
	pb "microservice-template/common/api"
	"microservice-template/common/broker"
	commontypes "microservice-template/common/types"
	"microservice-template/payments/internal/order"
	"microservice-template/payments/processor"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	// stripe service injection
	paymentProcessor processor.PaymentProcessor

	// config
	stripeWebhookSecret string

	// channel for communicating on message broker
	ch *amqp.Channel

	// for communicating with order service via grpc
	orderClient *order.Client
}

func NewService(paymentProcessor processor.PaymentProcessor, stripeWebhookSecret string, orderClient *order.Client) *service {
	return &service{
		paymentProcessor:    paymentProcessor,
		stripeWebhookSecret: stripeWebhookSecret,
		orderClient:         orderClient,
	}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	link, err := s.paymentProcessor.CreatePaymentLink(order)

	if err != nil {
		return "", err
	}

	return link, nil
}

func (s *service) GetWebhookSecret() string {
	return s.stripeWebhookSecret
}

func (s *service) UpdateOrderStatus(ctx context.Context, update UpdateOrderStatus) (*pb.Order, error) {

	payload := &pb.OrderStatusUpdateRequest{
		ID:     update.OrderId,
		Status: strconv.Itoa(int(commontypes.Paid)),
	}

	// update status directly via grpc - oder service takes precedence
	fmt.Println("firing update order status to order service")
	order, err := s.orderClient.UpdateOrderStatus(ctx, payload)

	marshalledPayload, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("Error on marshalling payment status for message broker.")
		return nil, fmt.Errorf("Error on marshalling payment status for message broker.")
	}

	fmt.Println("Publishing order paid event to message broker")

	// publish PAYMENT PAID event to message broker
	err = s.ch.PublishWithContext(
		ctx,
		broker.OrderPaidEvent,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         marshalledPayload,
			DeliveryMode: amqp.Persistent,
		})

	if err != nil {
		return nil, err
	}

	return order, nil

}

func (s *service) UpdateOrderPaymentLink(ctx context.Context, req *pb.OrderPaymentUpdateRequest) (*pb.Order, error) {
	fmt.Println("firing update order payment link")

	return s.orderClient.UpdateOrderPaymentLink(ctx, req)
}
