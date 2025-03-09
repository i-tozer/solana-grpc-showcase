package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/i-tozer/solana-grpc-exploration/proto"
	"github.com/olekukonko/tablewriter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("server", "localhost:50051", "The server address in the format host:port")
	command    = flag.String("command", "benchmark", "Command to run: benchmark, account, transaction, block, stream-accounts, stream-transactions, stream-blocks")
	pubkey     = flag.String("pubkey", "", "Solana account public key")
	signature  = flag.String("signature", "", "Solana transaction signature")
	slot       = flag.Uint64("slot", 0, "Solana block slot")
	iterations = flag.Uint("iterations", 10, "Number of iterations for benchmark")
)

func main() {
	flag.Parse()

	// Set up a connection to the server
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := proto.NewBenchmarkServiceClient(conn)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Execute the requested command
	switch *command {
	case "benchmark":
		runBenchmark(ctx, client)
	case "account":
		getAccountInfo(ctx, client)
	case "transaction":
		getTransaction(ctx, client)
	case "block":
		getBlock(ctx, client)
	case "stream-accounts":
		streamAccounts(ctx, client)
	case "stream-transactions":
		streamTransactions(ctx, client)
	case "stream-blocks":
		streamBlocks(ctx, client)
	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}

func runBenchmark(ctx context.Context, client proto.BenchmarkServiceClient) {
	if *pubkey == "" && *signature == "" && *slot == 0 {
		log.Fatal("At least one of --pubkey, --signature, or --slot must be specified")
	}

	// Prepare benchmark request
	req := &proto.BenchmarkRequest{
		Iterations:      uint32(*iterations),
		RunGrpcTests:    true,
		RunJsonrpcTests: true,
	}

	// Add test accounts if provided
	if *pubkey != "" {
		req.TestAccounts = []string{*pubkey}
	}

	// Add test signatures if provided
	if *signature != "" {
		req.TestSignatures = []string{*signature}
	}

	// Add test slots if provided
	if *slot != 0 {
		req.TestSlots = []uint64{*slot}
	}

	// Run benchmark
	fmt.Println("Running benchmark...")
	startTime := time.Now()
	resp, err := client.RunBenchmark(ctx, req)
	if err != nil {
		log.Fatalf("Error running benchmark: %v", err)
	}
	duration := time.Since(startTime)

	// Print results
	fmt.Printf("\nBenchmark completed in %v\n\n", duration)

	// Print account benchmark results if available
	if len(req.TestAccounts) > 0 {
		fmt.Println("Account Benchmark Results:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Metric", "gRPC", "JSON-RPC"})
		table.Append([]string{"Avg Response Time (ms)", fmt.Sprintf("%d", resp.AccountGrpc.AvgResponseTimeMs), fmt.Sprintf("%d", resp.AccountJsonrpc.AvgResponseTimeMs)})
		table.Append([]string{"Min Response Time (ms)", fmt.Sprintf("%d", resp.AccountGrpc.MinResponseTimeMs), fmt.Sprintf("%d", resp.AccountJsonrpc.MinResponseTimeMs)})
		table.Append([]string{"Max Response Time (ms)", fmt.Sprintf("%d", resp.AccountGrpc.MaxResponseTimeMs), fmt.Sprintf("%d", resp.AccountJsonrpc.MaxResponseTimeMs)})
		table.Append([]string{"Successful Requests", fmt.Sprintf("%d", resp.AccountGrpc.SuccessfulRequests), fmt.Sprintf("%d", resp.AccountJsonrpc.SuccessfulRequests)})
		table.Append([]string{"Failed Requests", fmt.Sprintf("%d", resp.AccountGrpc.FailedRequests), fmt.Sprintf("%d", resp.AccountJsonrpc.FailedRequests)})
		table.Render()
		fmt.Println()
	}

	// Print transaction benchmark results if available
	if len(req.TestSignatures) > 0 {
		fmt.Println("Transaction Benchmark Results:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Metric", "gRPC", "JSON-RPC"})
		table.Append([]string{"Avg Response Time (ms)", fmt.Sprintf("%d", resp.TransactionGrpc.AvgResponseTimeMs), fmt.Sprintf("%d", resp.TransactionJsonrpc.AvgResponseTimeMs)})
		table.Append([]string{"Min Response Time (ms)", fmt.Sprintf("%d", resp.TransactionGrpc.MinResponseTimeMs), fmt.Sprintf("%d", resp.TransactionJsonrpc.MinResponseTimeMs)})
		table.Append([]string{"Max Response Time (ms)", fmt.Sprintf("%d", resp.TransactionGrpc.MaxResponseTimeMs), fmt.Sprintf("%d", resp.TransactionJsonrpc.MaxResponseTimeMs)})
		table.Append([]string{"Successful Requests", fmt.Sprintf("%d", resp.TransactionGrpc.SuccessfulRequests), fmt.Sprintf("%d", resp.TransactionJsonrpc.SuccessfulRequests)})
		table.Append([]string{"Failed Requests", fmt.Sprintf("%d", resp.TransactionGrpc.FailedRequests), fmt.Sprintf("%d", resp.TransactionJsonrpc.FailedRequests)})
		table.Render()
		fmt.Println()
	}

	// Print block benchmark results if available
	if len(req.TestSlots) > 0 {
		fmt.Println("Block Benchmark Results:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Metric", "gRPC", "JSON-RPC"})
		table.Append([]string{"Avg Response Time (ms)", fmt.Sprintf("%d", resp.BlockGrpc.AvgResponseTimeMs), fmt.Sprintf("%d", resp.BlockJsonrpc.AvgResponseTimeMs)})
		table.Append([]string{"Min Response Time (ms)", fmt.Sprintf("%d", resp.BlockGrpc.MinResponseTimeMs), fmt.Sprintf("%d", resp.BlockJsonrpc.MinResponseTimeMs)})
		table.Append([]string{"Max Response Time (ms)", fmt.Sprintf("%d", resp.BlockGrpc.MaxResponseTimeMs), fmt.Sprintf("%d", resp.BlockJsonrpc.MaxResponseTimeMs)})
		table.Append([]string{"Successful Requests", fmt.Sprintf("%d", resp.BlockGrpc.SuccessfulRequests), fmt.Sprintf("%d", resp.BlockJsonrpc.SuccessfulRequests)})
		table.Append([]string{"Failed Requests", fmt.Sprintf("%d", resp.BlockGrpc.FailedRequests), fmt.Sprintf("%d", resp.BlockJsonrpc.FailedRequests)})
		table.Render()
		fmt.Println()
	}

	// Print summary
	fmt.Printf("Summary: %s\n", resp.Summary.Conclusion)
	fmt.Printf("gRPC vs JSON-RPC Speedup: %.2fx\n", resp.Summary.GrpcVsJsonrpcSpeedup)
	fmt.Printf("Total Benchmark Duration: %d ms\n", resp.Summary.TotalDurationMs)
}

func getAccountInfo(ctx context.Context, client proto.BenchmarkServiceClient) {
	if *pubkey == "" {
		log.Fatal("--pubkey is required")
	}

	// Get account info
	fmt.Printf("Getting account info for %s...\n", *pubkey)
	resp, err := client.GetAccountInfo(ctx, &proto.AccountInfoRequest{
		Pubkey:         *pubkey,
		Commitment:     "finalized",
		EncodingBinary: true,
	})
	if err != nil {
		log.Fatalf("Error getting account info: %v", err)
	}

	// Print results
	fmt.Printf("\nAccount Info:\n")
	fmt.Printf("Pubkey: %s\n", resp.Pubkey)
	fmt.Printf("Owner: %s\n", resp.Owner)
	fmt.Printf("Lamports: %d\n", resp.Lamports)
	fmt.Printf("Executable: %t\n", resp.Executable)
	fmt.Printf("Rent Epoch: %d\n", resp.RentEpoch)
	fmt.Printf("Data Length: %d bytes\n", len(resp.Data))
	fmt.Printf("Response Time: %d ms\n", resp.ResponseTimeMs)
}

func getTransaction(ctx context.Context, client proto.BenchmarkServiceClient) {
	if *signature == "" {
		log.Fatal("--signature is required")
	}

	// Get transaction
	fmt.Printf("Getting transaction %s...\n", *signature)
	resp, err := client.GetTransaction(ctx, &proto.TransactionRequest{
		Signature:  *signature,
		Commitment: "finalized",
	})
	if err != nil {
		log.Fatalf("Error getting transaction: %v", err)
	}

	// Print results
	fmt.Printf("\nTransaction Info:\n")
	fmt.Printf("Signature: %s\n", resp.Signature)
	fmt.Printf("Slot: %d\n", resp.Slot)
	fmt.Printf("Success: %t\n", resp.Success)
	fmt.Printf("Transaction Data Length: %d bytes\n", len(resp.Transaction))
	fmt.Printf("Response Time: %d ms\n", resp.ResponseTimeMs)
}

func getBlock(ctx context.Context, client proto.BenchmarkServiceClient) {
	if *slot == 0 {
		log.Fatal("--slot is required")
	}

	// Get block
	fmt.Printf("Getting block at slot %d...\n", *slot)
	resp, err := client.GetBlock(ctx, &proto.BlockRequest{
		Slot:       *slot,
		Commitment: "finalized",
	})
	if err != nil {
		log.Fatalf("Error getting block: %v", err)
	}

	// Print results
	fmt.Printf("\nBlock Info:\n")
	fmt.Printf("Slot: %d\n", resp.Slot)
	fmt.Printf("Blockhash: %s\n", resp.Blockhash)
	fmt.Printf("Previous Blockhash: %s\n", resp.PreviousBlockhash)
	fmt.Printf("Parent Slot: %d\n", resp.ParentSlot)
	fmt.Printf("Transactions: %d\n", len(resp.Transactions))
	fmt.Printf("Response Time: %d ms\n", resp.ResponseTimeMs)
}

func streamAccounts(ctx context.Context, client proto.BenchmarkServiceClient) {
	if *pubkey == "" {
		log.Fatal("--pubkey is required")
	}

	// Stream account updates
	fmt.Printf("Streaming account updates for %s...\n", *pubkey)
	stream, err := client.StreamAccountUpdates(ctx, &proto.AccountStreamRequest{
		Pubkeys:    []string{*pubkey},
		Commitment: "finalized",
	})
	if err != nil {
		log.Fatalf("Error streaming account updates: %v", err)
	}

	// Receive updates
	for {
		update, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving account update: %v", err)
		}

		// Print update
		fmt.Printf("\nAccount Update at %s:\n", time.Unix(int64(update.Timestamp), 0).Format(time.RFC3339))
		fmt.Printf("Pubkey: %s\n", update.Pubkey)
		fmt.Printf("Owner: %s\n", update.Owner)
		fmt.Printf("Lamports: %d\n", update.Lamports)
		fmt.Printf("Slot: %d\n", update.Slot)
		fmt.Printf("Data Length: %d bytes\n", len(update.Data))
	}
}

func streamTransactions(ctx context.Context, client proto.BenchmarkServiceClient) {
	// Stream transaction updates
	fmt.Println("Streaming transaction updates...")
	stream, err := client.StreamTransactions(ctx, &proto.TransactionStreamRequest{
		Accounts:      []string{},
		IncludeFailed: false,
		Commitment:    "finalized",
	})
	if err != nil {
		log.Fatalf("Error streaming transaction updates: %v", err)
	}

	// Receive updates
	for {
		update, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving transaction update: %v", err)
		}

		// Print update
		fmt.Printf("\nTransaction Update at %s:\n", time.Unix(int64(update.Timestamp), 0).Format(time.RFC3339))
		fmt.Printf("Signature: %s\n", update.Signature)
		fmt.Printf("Slot: %d\n", update.Slot)
		fmt.Printf("Success: %t\n", update.Success)
		fmt.Printf("Transaction Data Length: %d bytes\n", len(update.Transaction))
	}
}

func streamBlocks(ctx context.Context, client proto.BenchmarkServiceClient) {
	// Stream block updates
	fmt.Println("Streaming block updates...")
	stream, err := client.StreamBlocks(ctx, &proto.BlockStreamRequest{
		Commitment: "finalized",
	})
	if err != nil {
		log.Fatalf("Error streaming block updates: %v", err)
	}

	// Receive updates
	for {
		update, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving block update: %v", err)
		}

		// Print update
		fmt.Printf("\nBlock Update at %s:\n", time.Unix(int64(update.Timestamp), 0).Format(time.RFC3339))
		fmt.Printf("Slot: %d\n", update.Slot)
		fmt.Printf("Blockhash: %s\n", update.Blockhash)
		fmt.Printf("Previous Blockhash: %s\n", update.PreviousBlockhash)
		fmt.Printf("Parent Slot: %d\n", update.ParentSlot)
	}
}
