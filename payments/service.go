package main

import (
	"context"

	pb "github.com/vlkhvnn/commons/api"
	"github.com/vlkhvnn/oms-payments/gateway"
	"github.com/vlkhvnn/oms-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
	gateway   gateway.OrdersGateway
}

func NewService(processor processor.PaymentProcessor, gateway gateway.OrdersGateway) *service {
	return &service{processor, gateway}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(o)
	if err != nil {
		return "error occured", err
	}
	err = s.gateway.UpdateOrderAfterPaymentLink(ctx, o.ID, link)
	if err != nil {
		return "error occured", err
	}
	return link, nil
}
