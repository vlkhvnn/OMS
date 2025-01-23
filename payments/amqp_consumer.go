package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
	pb "github.com/vlkhvnn/commons/api"
	"github.com/vlkhvnn/commons/broker"
)

type consumer struct {
	service PaymentService
}

func NewConsumer(service PaymentService) *consumer {
	return &consumer{service}
}

func (c *consumer) Listen(ch *amqp091.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}
	go func() {
		for d := range msgs {
			log.Printf("Received message: %v", string(d.Body))

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				log.Printf("Failed to unmarshal order: %v", err)
				continue
			}
			log.Printf("ORDER IS %v", o)

			paymentLink, err := c.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("Failed to create payment: %v", err)
				continue
			}
			log.Printf("Payment link created %s", paymentLink)
		}
	}()

	<-forever
}
