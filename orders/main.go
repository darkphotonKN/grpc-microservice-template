package main

import "context"

func main() {

	repo := NewRepository()
	service := NewService(repo)

	service.CreateOrder(context.Background())

}
