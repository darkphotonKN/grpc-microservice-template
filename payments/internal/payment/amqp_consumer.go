package payment

import (
	"encoding/json"
	"fmt"
	"log"
	"microservice-template/common/broker"

	pb "microservice-template/common/api"

	amqp "github.com/rabbitmq/amqp091-go"
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
* Starts a listen for messages over rabbitmq for orders created.
**/
func (c *consumer) Listen() {
	queue, err := c.publishChan.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := c.publishChan.Consume(queue.Name, "", true, false, false, false, nil)

	var forever chan interface{}

	go func() {
		for msg := range msgs {
			fmt.Println("received message:", msg)

			var order *pb.Order

			err := json.Unmarshal(msg.Body, &order)

			if err != nil {
				fmt.Printf("Error when unmarshalling json:", err)
			}
		}
	}()

	<-forever
}
