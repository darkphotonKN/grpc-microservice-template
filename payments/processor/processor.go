package processor

import (
	pb "microservice-template/common/api"
)

type PaymentProcessor interface {
	CreatePaymentLink(o *pb.Order) (string, error)
}
