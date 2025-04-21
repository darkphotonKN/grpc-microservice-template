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
}
