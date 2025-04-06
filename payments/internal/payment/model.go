package payment

import (
	"context"
	pb "microservice-template/common/api"
)

type PaymentService interface {
	CreatePayment(context.Context, *pb.Order) (*pb.Payment, error)
}
