package processor

import pb "github.com/vlkhvnn/commons/api"

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
