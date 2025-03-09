package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/i-tozer/solana-grpc-exploration/proto"
	"github.com/i-tozer/solana-grpc-exploration/server/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port        = flag.Int("port", 50051, "The server port")
	rpcEndpoint = flag.String("rpc-endpoint", "https://api.mainnet-beta.solana.com", "Solana RPC endpoint")
)

func main() {
	flag.Parse()

	// Create a listener on the specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create and register the benchmark service
	benchmarkService := services.NewBenchmarkService(*rpcEndpoint)
	proto.RegisterBenchmarkServiceServer(grpcServer, benchmarkService)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Handle graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	// Start the server
	log.Printf("Starting gRPC server on port %d...", *port)
	log.Printf("Using Solana RPC endpoint: %s", *rpcEndpoint)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
