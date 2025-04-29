package payment

import (
	"context"
	pb "microservice-template/common/api"

	"github.com/gin-gonic/gin"
)

type PaymentHandler interface {
	HandleStripeWebhook(c *gin.Context)
}

type PaymentService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
	GetWebhookSecret() string
	UpdateOrderStatus(ctx context.Context, update UpdateOrderStatus) (*pb.Order, error)
	UpdateOrderPaymentLink(ctx context.Context, req *pb.OrderPaymentUpdateRequest) (*pb.Order, error)
}

// Request / Response

type UpdateOrderStatus struct {
	OrderId     string
	PaymentLink string
}
