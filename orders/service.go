package main

import "context"

type service struct {
	repo OrderRepository
}

func NewService(repo OrderRepository) OrderService {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateOrder(context.Context) error {
	return nil
}
