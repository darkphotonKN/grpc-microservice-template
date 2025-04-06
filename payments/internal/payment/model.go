package payment

import (
	"context"
	pb "microservice-template/common/api"
)

type PaymentService interface {
	CreatePayment(context.Context, *pb.CreatePaymentRequest) (*pb.Payment, error)
}
