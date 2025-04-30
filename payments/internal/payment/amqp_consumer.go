package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"microservice-template/common/broker"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "microservice-template/common/api"
)

type consumer struct {
	service     PaymentService
	publishChan *amqp.Channel
}

func NewConsumer(service PaymentService, publishChan *amqp.Channel) *consumer {

	return &consumer{
		service:     service,
		publishChan: publishChan,
	}
}

/**
* Starts a listen for messages over rabbitmq for newOrders created.
**/
func (c *consumer) Listen() {
	queueName := fmt.Sprintf("payment.%s", broker.OrderCreatedEvent)
	queue, err := c.publishChan.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	err = c.publishChan.QueueBind(
		queue.Name,
		"",
		broker.OrderCreatedEvent,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := c.publishChan.Consume(queue.Name, "", true, false, false, false, nil)

	go func() {
		for msg := range msgs {
			fmt.Println("received message:", msg)

			var newOrder *pb.Order

			err := json.Unmarshal(msg.Body, &newOrder)

			if err != nil {
				fmt.Printf("Error when unmarshalling json: %s\n", err)
				continue
			}

			paymentLink, err := c.service.CreatePayment(context.Background(), newOrder)

			if err != nil {
				fmt.Printf("Error when creating payment: %s\n", err)
				continue
			}

			// call update payment via grpc
			c.service.UpdateOrderPaymentLink(context.Background(), &pb.OrderPaymentUpdateRequest{
				ID:          newOrder.ID,
				PaymentLink: paymentLink,
			})

		}
	}()

}
