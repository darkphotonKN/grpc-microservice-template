package main

import (
	"context"
	"microservice-template/common"
	pb "microservice-template/common/api"
)

type service struct {
	repo OrderRepository
}

func NewService(repo OrderRepository) OrderService {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateOrder(ctx context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) error {
	if len(req.Items) == 0 {
		return common.ErrNoItems
	}

	return nil
}
