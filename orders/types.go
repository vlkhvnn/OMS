package main

import (
	"context"

	pb "github.com/vlkhvnn/commons/api"
)

type OrderService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	UpdateOrder(context.Context, *pb.Order) (*pb.Order, error)
}

type OrderStore interface {
	Create(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (string, error)
	Get(ctx context.Context, id, customerID string) (*pb.Order, error)
	Update(ctx context.Context, id string, order *pb.Order) error
}
