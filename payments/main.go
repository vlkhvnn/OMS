package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stripe/stripe-go/v81"
	common "github.com/vlkhvnn/commons"
	"github.com/vlkhvnn/commons/broker"
	"github.com/vlkhvnn/commons/discovery"
	"github.com/vlkhvnn/commons/discovery/consul"
	"github.com/vlkhvnn/oms-payments/gateway"
	st "github.com/vlkhvnn/oms-payments/processor/stripe"
	"google.golang.org/grpc"
)

var (
	grpcAddr             = common.GetString("GRPC_ADDR", "localhost:2001")
	consulAddr           = common.GetString("CONSUL_ADDR", "localhost:8500")
	serviceName          = "payment"
	amqpUser             = common.GetString("RABBITMQ_USER", "guest")
	amqpPass             = common.GetString("RABBITMQ_PASS", "guest")
	amqpHost             = common.GetString("RABBITMQ_HOST", "localhost")
	amqpPort             = common.GetString("RABBITMQ_PORT", "5672")
	stripeKey            = common.GetString("STRIPE_KEY", "")
	httpAddr             = common.GetString("HTTP_ADDR", "localhost:8081")
	endpointStripeSecret = common.GetString("STRIPE_ENDPOINT_KEY", "")
	jaegerAddr           = common.GetString("JAEGER_ADDR", "localhost:4318")
)

func main() {
	if err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr); err != nil {
		log.Fatal("could set global tracer", err)
	}
	// Register consul
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

	// stripe setup
	stripe.Key = stripeKey

	// Broker connection
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	stripeProcessor := st.NewProcessor()
	gateway := gateway.NewGateway(registry)
	svc := NewService(stripeProcessor, gateway)

	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(ch)

	// http server
	mux := http.NewServeMux()

	httpServer := NewPaymentHTTPHandler(ch)
	httpServer.registerRoutes(mux)

	go func() {
		log.Printf("Starting HTTP server at %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Fatal("failed to start http server")
		}
	}()

	// gRPC server
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen %v", grpcAddr)
	}
	defer l.Close()

	log.Println("gprc server started at:", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
