package gateway

import (
	"context"

	pb "github.com/vlkhvnn/commons/api"
)

type KitchenGateway interface {
	UpdateOrder(context.Context, *pb.Order) error
}
