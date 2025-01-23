package main

import pb "github.com/vlkhvnn/commons/api"

type CreateOrderRequest struct {
	Order         *pb.Order `"json": order`
	RedirectToURL string    `"json": redirectToURL`
}
