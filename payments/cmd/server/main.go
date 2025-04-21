package main

import (
	"context"
	"fmt"
	"log"
	"microservice-template/common/broker"
	"microservice-template/common/discovery"
	"microservice-template/common/discovery/consul"
	commonenv "microservice-template/common/env"
	"microservice-template/payments/internal/payment"
	stripeProcessor "microservice-template/payments/processor/stripe"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload" // package that loads env
	"github.com/stripe/stripe-go/v78"
	"google.golang.org/grpc"
)

var (
	serviceName  = "payment"
	grpcAddr     = commonenv.EnvString("GRPC_ORDER_ADDR", "2222")
	httpAddr     = commonenv.EnvString("PAYMENT_ADDR", "8070")
	amqpUser     = commonenv.EnvString("RABBITMQ_USER", "guest")
	amqpPassword = commonenv.EnvString("RABBITMQ_PASS", "guest")
	amqpHost     = commonenv.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort     = commonenv.EnvString("RABBITMQ_PORT", "5672")
	consulAddr   = commonenv.EnvString("CONSUL_ADDR", "localhost:8500")
	stripeKey    = commonenv.EnvString("STRIPE_KEY", "testkey")
)

func main() {
	// --- service discovery ---

	// -- setup --
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	// -- register --
	if err := registry.Register(ctx, instanceID, serviceName, "localhost:"+grpcAddr); err != nil {
		panic(err)
	}

	// -- health check --
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("Health check failed.")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	// --- third party services setup ---

	// -- stripe --
	stripe.Key = stripeKey

	fmt.Println("stripekey:", stripeKey)
	fmt.Println("stripe.Key:", stripe.Key)

	processor := stripeProcessor.NewStripeProcessor()

	// --- message broker ---
	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)

	defer func() {
		close()
		ch.Close()
	}()

	// --- server initialization ---
	stripeWebhookSecret := commonenv.EnvString("STRIPE_WEBHOOK_SECRET", "")

	grpcServer := grpc.NewServer()
	paymentService := payment.NewService(processor, stripeWebhookSecret)
	paymentHandler := payment.NewHandler(paymentService)
	paymentConsumer := payment.NewConsumer(paymentService, ch) // listen through channel from message broker
	paymentConsumer.Listen()                                   // listen to the channel for messages

	// create a local network listener to this service
	l, err := net.Listen("tcp", "localhost:"+grpcAddr)

	if err != nil {
		log.Fatalf(
			"Failed to listen at port: %s\nError: %s\n", grpcAddr, err,
		)
	}

	defer l.Close()

	// start a http server for exposing webhook endpoint to stripe
	router := gin.Default()

	router.GET("/api/payment/webhook", paymentHandler.HandleStripeWebhook)

	// -- start server and capture errors --
	if err := router.Run(":" + httpAddr); err != nil {
		log.Fatal("Failed to start server")
	}

	log.Printf("http payment server started on PORT: %s\n", httpAddr)
	log.Printf("grpc Order Server started on PORT: %s\n", grpcAddr)

	// start serving requests
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal("Can't connect to grpc server. Error:", err.Error())
	}

}
