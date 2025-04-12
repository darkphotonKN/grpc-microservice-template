package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	pb "microservice-template/common/api"
	"microservice-template/common/broker"

	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service     NotificationService
	publishChan *amqp.Channel
}

func NewConsumer(service NotificationService, publishChan *amqp.Channel) *consumer {
	return &consumer{
		service:     service,
		publishChan: publishChan,
	}
}

/**
* Starts a listen for events that require a notification messge to be sent.
**/
func (c *consumer) Listen() {

	// --- Order Created Event ---
	// -- queue --
	queue, err := c.publishChan.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	// -- consumption --
	msgs, err := c.publishChan.Consume(queue.Name, "", true, false, false, false, nil)

	var forever chan interface{}

	go func() {
		for msg := range msgs {
			fmt.Println("received message:", msg)

			var order *pb.Order

			err := json.Unmarshal(msg.Body, &order)

			if err != nil {
				fmt.Printf("Error when unmarshalling json: %s\n", err)
				continue
			}

			err = c.service.SendMessage(context.Background(), fmt.Sprintf("Order with CustomerID %s was sent.", order.CustomerID))

			if err != nil {
				fmt.Printf("Error when sending notification.: %s\n", err)
				continue
			}

		}
	}()

	<-forever
}
