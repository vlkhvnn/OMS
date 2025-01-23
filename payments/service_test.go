package main

import (
	"context"
	"testing"

	"github.com/vlkhvnn/commons/api"
	inmemRegistry "github.com/vlkhvnn/commons/discovery/inmem"
	"github.com/vlkhvnn/oms-payments/gateway"
	"github.com/vlkhvnn/oms-payments/processor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	registry := inmemRegistry.NewRegistry()
	gateway := gateway.NewGateway(registry)
	svc := NewService(processor, gateway)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &api.Order{})
		if err != nil {
			t.Errorf("CreatePayment() error = %v, want nil", err)
		}
		if link == "" {
			t.Errorf("CreatePayment() link is empty")
		}

	})
}
