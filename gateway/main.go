package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/vlkhvnn/commons"
	"github.com/vlkhvnn/commons/discovery"
	"github.com/vlkhvnn/commons/discovery/consul"
	"github.com/vlkhvnn/oms-gateway/gateway"
)

var (
	httpAddr    = common.GetString("HTTP_ADDR", ":8080")
	consulAddr  = common.GetString("CONSUL_ADDR", "localhost:8500")
	serviceName = "gateway"
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatalf("failed to health check. error: %s", err.Error())
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	mux := http.NewServeMux()

	ordersGateway := gateway.NewGRPCGateway(registry)
	handler := NewHandler(ordersGateway)
	handler.registerRoutes(mux)

	log.Printf("Starting http server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server")
	}
}
