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
	fmt.Println("Listening for Order Create messages.")

	// --- Order Created Event ---

	// -- declare unique queue  --
	// declare queue with a unique name, different than payment service's consumer
	queueName := fmt.Sprintf("notification.%s", broker.OrderCreatedEvent)
	queue, err := c.publishChan.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	// -- bind queue to exchange --
	// here we bind our newly declared queue to the order created event queue to react on orders being created
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

	// -- consumption --
	msgs, err := c.publishChan.Consume(queue.Name, "", true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	var forever chan interface{}

	go func() {
		fmt.Printf("\nStarted go routine for receiving messages:\n %+v\n\n", msgs)

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
