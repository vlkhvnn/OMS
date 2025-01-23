package main

import (
	"context"

	pb "github.com/vlkhvnn/commons/api"
)

type PaymentService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
