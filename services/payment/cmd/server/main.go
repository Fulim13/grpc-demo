package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Fulim13/grpc-demo/services/payment/internal/handler"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Another Goroutine listen on the terminate signal to graceful shutdowwn
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		port = "1323"
	}
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen on port %s, error %v", port, err)
	}

	// service := service.NewService() - In Real Scenario, will have service

	grpcServer := grpc.NewServer()
	// handler.NewGRPCHandler(grpcServer, service) - In Real Scenario, will have service
	handler.NewGRPCHandler(grpcServer)

	log.Printf("Staring gRPC server Payment service on port %v", listen.Addr().String())

	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Printf("Shutting down the server...")
	grpcServer.GracefulStop()
}
