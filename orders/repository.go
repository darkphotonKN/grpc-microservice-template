package main

import "context"

type repository struct {
}

func NewRepository() OrderRepository {
	return &repository{}
}

func (s *repository) Create(context.Context) error {
	return nil
}
