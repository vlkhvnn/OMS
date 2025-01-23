package main

import (
	"context"
	"log"

	common "github.com/vlkhvnn/commons"
	pb "github.com/vlkhvnn/commons/api"
)

type service struct {
	store OrderStore
}

func NewService(store OrderStore) *service {
	return &service{store: store}
}

func (s *service) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	return s.store.Get(ctx, p.OrderID, p.CustomerID)
}

func (s *service) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	err := s.store.Update(ctx, o.ID, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	id, err := s.store.Create(ctx, p, items)
	if err != nil {
		return nil, err
	}

	o := &pb.Order{
		ID:          id,
		CustomerID:  p.CustomerID,
		Status:      "pending",
		Items:       items,
		PaymentLink: "",
	}
	return o, nil
}
func (s *service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(p.Items)
	log.Print(mergedItems)
	//validate with the stock service

	// Temporary:
	var itemsWithPrice []*pb.Item
	for _, i := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			PriceID:  "price_1Qj5t1P8Nkuey8SGulWOi8Oq",
			ID:       i.ID,
			Quantity: i.Quantity,
		})
	}

	return itemsWithPrice, nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)
	for _, item := range items {
		found := false
		for _, finalItem := range merged {
			if item.ID == finalItem.ID {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, item)
		}
	}
	return merged
}
