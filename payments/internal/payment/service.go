package payment

import (
	"context"
	"fmt"
	pb "microservice-template/common/api"
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

	// TODO: remove after testing, should fire from stripe payment callback
	fmt.Println("firing update order status")
	_, err = s.UpdateOrderStatus(ctx, &pb.OrderStatusUpdateRequest{
		ID:          "b9a9dbce-a8b4-402a-8b68-23cad60d3b27",
		Status:      strconv.Itoa(int(commontypes.Paid)),
		PaymentLink: link,
	})
	if err != nil {
		return "", err
	}

	return link, nil
}

func (s *service) GetWebhookSecret() string {
	return s.stripeWebhookSecret
}

func (s *service) UpdateOrderStatus(ctx context.Context, req *pb.OrderStatusUpdateRequest) (*pb.Order, error) {
	return s.orderClient.UpdateOrderStatus(ctx, req)
}
