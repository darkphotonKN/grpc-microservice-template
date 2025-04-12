package notification

import (
	"context"
	"fmt"
	"microservice-template/context/internal/model"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	// channel for communicating on message broker
	ch *amqp.Channel
}

func NewService() NotificationService {
	return &service{}
}

func (s *service) SendMessage(ctx context.Context, message string) error {
	fmtedMsg := model.Message{
		ID:      uuid.New(),
		Content: message,
	}

	fmt.Printf("\nSend Message:\n%+v\n\n", fmtedMsg)

	return nil
}
