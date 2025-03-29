package main

import "context"

type OrderService interface {
	CreateOrder(context.Context) error
}

type OrderRepository interface {
	Create(context.Context) error
}
