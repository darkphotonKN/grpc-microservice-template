package order

import (
	"context"
	"encoding/json"
	"fmt"
	pb "microservice-template/common/api"
	"microservice-template/common/broker"
	commonerrors "microservice-template/common/errors"
	commonhelpers "microservice-template/common/helpers"
	commontypes "microservice-template/common/types"
	"strconv"

	"github.com/google/uuid"
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

func (s *service) GetOrderPaymentLink(ctx context.Context, req *pb.OrderId) (*pb.OrderPaymentLink, error) {
	fmt.Printf("\ngetorderpaymentlink req: %+v\n\n", req)
	order, err := s.repo.GetOrder(ctx, req)

	if err != nil {
		return nil, err
	}

	fmt.Printf("\nOrder retrieved: \n%+v\n\n", order)

	return &pb.OrderPaymentLink{
		OrderPaymentLink: order.PaymentLink,
	}, nil
}

func (s *service) GetOrderStatus(ctx context.Context, req *pb.OrderId) (*pb.OrderStatus, error) {
	order, err := s.repo.GetOrder(ctx, req)

	if err != nil {
		return nil, err
	}

	fmt.Printf("\nGetOrderStauts: Order retrieved: \n%+v\n\n", order)

	switch commontypes.OrderStatus(order.Status) {
	case commontypes.Pending:
		statusText := "pending"
		return &pb.OrderStatus{
			Status: statusText,
		}, nil
	case commontypes.Paid:
		statusText := "paid"
		return &pb.OrderStatus{
			Status: statusText}, nil
	default:
		return nil, fmt.Errorf("Unknown order status: %d", order.Status)
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

	// create order and order items with transaction to retain atomicity
	var orderID uuid.UUID
	db := (s.repo).(*repository).DB

	err = commonhelpers.ExecTx(db, func(tx *sqlx.Tx) error {

		// create base order
		orderID, err = s.repo.CreateOrderTx(ctx, tx, Order{
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
	order.ID = orderID.String() // update id based on the created id
	marshalledOrder, err := json.Marshal(order)

	if err != nil {
		return nil, err
	}

	// publish ORDER CREATED EVENT to message broker
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

func (s *service) UpdateOrderStatus(ctx context.Context, req *pb.OrderStatusUpdateRequest) (*pb.Order, error) {
	fmt.Printf("Recieved update order status: %+v\n", req)

	// parse out the id and check
	idUUID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	// parse out the status and check
	status, err := strconv.Atoi(req.Status)

	if err != nil {
		return nil, err
	}

	orderStatus := commontypes.OrderStatus(status)

	err = s.repo.UpdateOrderStatus(ctx, &UpdateOrderStatusReq{
		ID:     idUUID,
		Status: orderStatus,
	})

	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully updated order status.")

	return nil, nil
}

func (s *service) UpdateOrderPaymentLink(ctx context.Context, req *pb.OrderPaymentUpdateRequest) (*pb.Order, error) {
	fmt.Printf("Recieved update order paymentLink: %+v\n", req)

	err := s.repo.UpdateOrderPaymentLink(ctx, req)

	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully updated order payment link.")

	return nil, nil
}
