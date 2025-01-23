package main

import (
	"context"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/vlkhvnn/commons"
	"github.com/vlkhvnn/commons/broker"
	"github.com/vlkhvnn/commons/discovery"
	"github.com/vlkhvnn/commons/discovery/consul"
	"google.golang.org/grpc"
)

var (
	grpcAddr    = common.GetString("GRPC_ADDR", "localhost:2000")
	consulAddr  = common.GetString("CONSUL_ADDR", "localhost:8500")
	serviceName = "orders"
	amqpUser    = common.GetString("RABBITMQ_USER", "guest")
	amqpPass    = common.GetString("RABBITMQ_PASS", "guest")
	amqpHost    = common.GetString("RABBITMQ_HOST", "localhost")
	amqpPort    = common.GetString("RABBITMQ_PORT", "5672")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
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

	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen %v", grpcAddr)
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)
	NewGRPCHandler(grpcServer, svc, ch)

	log.Println("gprc server started at:", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
