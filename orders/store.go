package main

import (
	"context"
	"errors"

	pb "github.com/vlkhvnn/commons/api"
)

var orders = make([]*pb.Order, 0)

type store struct {
	// add here our MongoDB
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (string, error) {
	id := "42"
	orders = append(orders, &pb.Order{
		ID:          id,
		CustomerID:  p.CustomerID,
		Status:      "pending",
		Items:       items,
		PaymentLink: "",
	})
	return id, nil
}

func (s *store) Get(ctx context.Context, id, customerID string) (*pb.Order, error) {
	for _, o := range orders {
		if o.ID == id && o.CustomerID == customerID {
			return o, nil
		}
	}
	return nil, errors.New("order not found")
}

func (s *store) Update(ctx context.Context, id string, order *pb.Order) error {
	for i, o := range orders {
		if o.ID == id {
			orders[i].Status = order.Status
			orders[i].PaymentLink = order.PaymentLink
			return nil
		}
	}
	return nil
}
