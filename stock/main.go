package main

import (
	"context"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/vlkhvnn/commons"
	"github.com/vlkhvnn/commons/broker"
	"github.com/vlkhvnn/commons/discovery"
	"github.com/vlkhvnn/commons/discovery/consul"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName = "stock"
	grpcAddr    = common.GetString("GRPC_ADDR", "localhost:2002")
	consulAddr  = common.GetString("CONSUL_ADDR", "localhost:8500")
	amqpUser    = common.GetString("RABBITMQ_USER", "guest")
	amqpPass    = common.GetString("RABBITMQ_PASS", "guest")
	amqpHost    = common.GetString("RABBITMQ_HOST", "localhost")
	amqpPort    = common.GetString("RABBITMQ_PORT", "5672")
	jaegerAddr  = common.GetString("JAEGER_ADDR", "localhost:4318")
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	if err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr); err != nil {
		logger.Fatal("could set global tracer", zap.Error(err))
	}

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
				logger.Error("Failed to health check", zap.Error(err))
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
		logger.Fatal("failed to listen", zap.Error(err))
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)
	svcWithTelemetry := NewTelemetryMiddleware(svc)

	NewGRPCHandler(grpcServer, ch, svcWithTelemetry)

	consumer := NewConsumer()
	go consumer.Listen(ch)

	logger.Info("Starting gRPC server", zap.String("port", grpcAddr))

	if err := grpcServer.Serve(l); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
