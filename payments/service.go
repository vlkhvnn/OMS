package main

import (
	"context"

	pb "github.com/vlkhvnn/commons/api"
	"github.com/vlkhvnn/oms-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
}

func NewService(processor processor.PaymentProcessor) *service {
	return &service{processor}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(o)
	if err != nil {
		return "error occured", err
	}
	return link, nil
}
