package client

import (
	"log"
	"os"

	"github.com/Fulim13/grpc-demo/gen/golang/v1/payment"
	"github.com/Fulim13/grpc-demo/services/order/internal/middleware"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type paymentServiceClient struct {
	Client payment.PaymentClient
	conn   *grpc.ClientConn
}

func NewPaymentServiceClient() (*paymentServiceClient, error) {
	paymentServiceUrl := os.Getenv("PAYMENT_SERVICE_URL")
	if paymentServiceUrl == "" {
		paymentServiceUrl = "localhost:1324"
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// Retry Pattern
	// opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
	// 	grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
	// 	grpc_retry.WithMax(5),
	// 	grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)))))

	// Circuit Breaker Pattern
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "cbPayment", // Unique name of circuit breaker
		MaxRequests: 0,           // Allowed number of requests for half open circult
		Timeout:     4,           // Timeout for an open to half-open transition
		ReadyToTrip: func(counts gobreaker.Counts) bool { // Decide if the circuit will be open
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.6
		},
		OnStateChange: func(name string, from, to gobreaker.State) { // Execute on each state changes
			log.Printf("Circuit Breaker: %s, changed from %v, to %v", name, from, to)
		},
	})
	opts = append(opts, grpc.WithUnaryInterceptor(middleware.CircuitBreakerClientInterceptor(cb)))
	conn, err := grpc.NewClient(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

	client := payment.NewPaymentClient(conn)
	return &paymentServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}
