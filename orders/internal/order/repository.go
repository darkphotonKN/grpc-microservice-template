package order

import (
	"context"
	pb "microservice-template/common/api"
)

type repository struct {
}

func NewRepository() OrderRepository {
	return &repository{}
}

func (s *repository) Create(context.Context) error {
	return nil
}

func (s *repository) Get(ctx context.Context, id, customerId string) (*pb.Order, error) {
	return nil, nil
}
