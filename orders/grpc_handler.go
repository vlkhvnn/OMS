package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
	pb "github.com/vlkhvnn/commons/api"
	"github.com/vlkhvnn/commons/broker"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrderService
	channel *amqp091.Channel
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrderService, channel *amqp091.Channel) {
	handler := &grpcHandler{service: service, channel: channel}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	log.Printf("Getting an Order %v", p)

	return h.service.GetOrder(ctx, p)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received!, Order %v", p)

	items, err := h.service.ValidateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	o, err := h.service.CreateOrder(ctx, p, items)

	// Marshal the full Order to JSON
	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	// Publish to RabbitMQ
	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp091.Persistent,
	})
	if err != nil {
		log.Printf("Failed to publish order to RabbitMQ: %v", err)
		return nil, err
	}

	log.Printf("Order published to RabbitMQ: %v", o)
	return o, nil
}
