package order

import (
	"context"
	"encoding/json"
	"fmt"
	pb "microservice-template/common/api"
	"microservice-template/common/broker"
	commonerrors "microservice-template/common/errors"
	commonhelpers "microservice-template/common/helpers"
	"strconv"

	"github.com/jmoiron/sqlx"
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
		fmt.Printf("Error when attemtpting to create order: %s\n", err)
		return nil, err
	}

	items := make([]*pb.Item, len(req.Items))

	for index, item := range req.Items {

		indexStr := strconv.Itoa(index)
		items[index] = &pb.Item{
			ID:       item.ID,
			Name:     "testeritem" + indexStr,
			Quantity: item.Quantity,
			PriceID:  "price_1RBggxIl3wC7xA9ZojS9Vo8v",
		}
	}

	order := &pb.Order{
		CustomerID: req.CustomerID,
		Items:      items,
	}

	// create order and order items with transaction to retain atomicitiy
	db := (s.repo).(*repository).DB

	err = commonhelpers.ExecTx(db, func(tx *sqlx.Tx) error {

		// create base order
		orderID, err := s.repo.CreateOrderTx(ctx, tx, Order{
			CustomerID: order.CustomerID,
		})

		if err != nil {
			fmt.Printf("Order had error: %s\n", err)
			return err
		}

		// create each order item under the order
		for index, item := range order.Items {

			indexStr := strconv.Itoa(index)

			newOrderItem := OrderItem{
				OrderID:  orderID,
				Name:     item.Name + indexStr,
				Quantity: int(item.Quantity),
				PriceID:  item.PriceID,
			}

			err := s.repo.CreateOrderItemTx(ctx, tx, newOrderItem)

			if err != nil {
				fmt.Printf("Error occured when creating order item: %+v for order: %s. Error was:\n%s\n", newOrderItem, orderID, err)
				return err
			}
		}

		fmt.Printf("creating order at order service: %+v\n", order)
		return nil
	})

	if err != nil {
		return nil, err
	}

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

	return nil, fmt.Errorf("Unknown error occurred when creating orders.")
}

func (s *service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) error {
	if len(req.Items) == 0 {
		return commonerrors.ErrNoItems
	}

	return nil
}
