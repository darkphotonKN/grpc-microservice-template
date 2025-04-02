package order

import (
	"context"
	"fmt"
	pb "microservice-template/common/api"
	commonerrors "microservice-template/common/errors"
)

type service struct {
	repo OrderRepository
}

func NewService(repo OrderRepository) OrderService {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	// validation
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
			PriceID:  "rando",
		}
	}

	order := &pb.Order{
		ID:         "1",
		CustomerID: req.CustomerID,
		Status:     "initiated",
		Items:      items,
	}

	fmt.Printf("Outgoing order: %+v\n", order)

	return order, nil
}

func (s *service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) error {
	if len(req.Items) == 0 {
		return commonerrors.ErrNoItems
	}

	return nil
}
