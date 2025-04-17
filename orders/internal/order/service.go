package order

import (
	"context"
	"encoding/json"
	"fmt"
	pb "microservice-template/common/api"
	"microservice-template/common/broker"
	commonerrors "microservice-template/common/errors"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	repo      OrderRepository
	publishCh *amqp.Channel
}

func NewService(repo OrderRepository, publishCh *amqp.Channel) OrderService {
	return &service{
		repo:      repo,
		publishCh: publishCh,
	}
}

func (s *service) GetOrders(ctx context.Context, pb *emptypb.Empty) (*pb.Orders, error) {
	fmt.Println("Getting orders!")
	return nil, nil
}

func (s *service) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	err := s.ValidateOrder(ctx, req)

	if err != nil {
		return nil, err
	}

	items := make([]*pb.Item, len(req.Items))

	for index, item := range req.Items {
		items[index] = &pb.Item{
			ID:       item.ID,
			Name:     "testeritem",
			Quantity: item.Quantity,
			PriceID:  "price_1RBggxIl3wC7xA9ZojS9Vo8v",
		}
	}

	order := &pb.Order{
		ID:         "1",
		CustomerID: req.CustomerID,
		Status:     "initiated",
		Items:      items,
	}

	fmt.Printf("creating order at order service: %+v\n", order)

	// publish created order via rabbitmq
	marshalledOrder, err := json.Marshal(order)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	s.publishCh.PublishWithContext(
		ctx,
		broker.OrderCreatedEvent,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        marshalledOrder,
			// persist message
			DeliveryMode: amqp.Persistent,
		})

	return order, nil
}

func (s *service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) error {
	if len(req.Items) == 0 {
		return commonerrors.ErrNoItems
	}

	return nil
}
