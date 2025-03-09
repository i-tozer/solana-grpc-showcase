package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/i-tozer/solana-grpc-exploration/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BenchmarkService implements the gRPC benchmark service
type BenchmarkService struct {
	proto.UnimplementedBenchmarkServiceServer
	solanaClient *rpc.Client
	rpcEndpoint  string
}

// NewBenchmarkService creates a new benchmark service
func NewBenchmarkService(rpcEndpoint string) *BenchmarkService {
	client := rpc.New(rpcEndpoint)
	return &BenchmarkService{
		solanaClient: client,
		rpcEndpoint:  rpcEndpoint,
	}
}

// GetAccountInfo retrieves account information and measures performance
func (s *BenchmarkService) GetAccountInfo(ctx context.Context, req *proto.AccountInfoRequest) (*proto.AccountInfoResponse, error) {
	pubkey, err := solana.PublicKeyFromBase58(req.Pubkey)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid pubkey: %v", err)
	}

	startTime := time.Now()

	// Get account info
	accountInfo, err := s.solanaClient.GetAccountInfo(ctx, pubkey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get account info: %v", err)
	}

	responseTime := time.Since(startTime).Milliseconds()

	// Convert account info to response
	response := &proto.AccountInfoResponse{
		Pubkey:         req.Pubkey,
		Data:           accountInfo.Value.Data.GetBinary(),
		Owner:          accountInfo.Value.Owner.String(),
		Lamports:       accountInfo.Value.Lamports,
		Executable:     accountInfo.Value.Executable,
		RentEpoch:      accountInfo.Value.RentEpoch,
		ResponseTimeMs: uint64(responseTime),
	}

	return response, nil
}

// GetTransaction retrieves transaction information and measures performance
func (s *BenchmarkService) GetTransaction(ctx context.Context, req *proto.TransactionRequest) (*proto.TransactionResponse, error) {
	signature, err := solana.SignatureFromBase58(req.Signature)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid signature: %v", err)
	}

	startTime := time.Now()

	// Get transaction
	tx, err := s.solanaClient.GetTransaction(ctx, signature, &rpc.GetTransactionOpts{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transaction: %v", err)
	}

	responseTime := time.Since(startTime).Milliseconds()

	// Convert transaction to response
	response := &proto.TransactionResponse{
		Signature:      req.Signature,
		Slot:           tx.Slot,
		Transaction:    []byte(fmt.Sprintf("%v", tx.Transaction)), // Simplified for demo
		Success:        tx.Meta.Err == nil,
		ResponseTimeMs: uint64(responseTime),
	}

	return response, nil
}

// GetBlock retrieves block information and measures performance
func (s *BenchmarkService) GetBlock(ctx context.Context, req *proto.BlockRequest) (*proto.BlockResponse, error) {
	startTime := time.Now()

	// Get block
	block, err := s.solanaClient.GetBlock(ctx, uint64(req.Slot))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get block: %v", err)
	}

	responseTime := time.Since(startTime).Milliseconds()

	// Extract transaction signatures
	txSigs := make([]string, 0, len(block.Transactions))
	for _, tx := range block.Transactions {
		txWithMeta, err := tx.GetTransaction()
		if err != nil {
			continue
		}
		if len(txWithMeta.Signatures) > 0 {
			txSigs = append(txSigs, txWithMeta.Signatures[0].String())
		}
	}

	// Convert block to response
	response := &proto.BlockResponse{
		Slot:              req.Slot,
		Blockhash:         block.Blockhash.String(),
		PreviousBlockhash: block.PreviousBlockhash.String(),
		ParentSlot:        block.ParentSlot,
		Transactions:      txSigs,
		ResponseTimeMs:    uint64(responseTime),
	}

	return response, nil
}

// StreamAccountUpdates streams account updates in real-time
func (s *BenchmarkService) StreamAccountUpdates(req *proto.AccountStreamRequest, stream proto.BenchmarkService_StreamAccountUpdatesServer) error {
	// Convert pubkeys to solana.PublicKey
	pubkeys := make([]solana.PublicKey, 0, len(req.Pubkeys))
	for _, pubkeyStr := range req.Pubkeys {
		pubkey, err := solana.PublicKeyFromBase58(pubkeyStr)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid pubkey: %v", err)
		}
		pubkeys = append(pubkeys, pubkey)
	}

	// For demo purposes, we'll simulate account updates
	// In a real implementation, you would use WebSocket subscriptions
	for i := 0; i < 10; i++ {
		for _, pubkey := range pubkeys {
			// Get account info
			accountInfo, err := s.solanaClient.GetAccountInfo(context.Background(), pubkey)
			if err != nil {
				log.Printf("Error getting account info: %v", err)
				continue
			}

			// Send account update
			err = stream.Send(&proto.AccountUpdate{
				Pubkey:    pubkey.String(),
				Data:      accountInfo.Value.Data.GetBinary(),
				Owner:     accountInfo.Value.Owner.String(),
				Lamports:  accountInfo.Value.Lamports,
				Slot:      accountInfo.Context.Slot,
				Timestamp: uint64(time.Now().Unix()),
			})
			if err != nil {
				return status.Errorf(codes.Internal, "failed to send account update: %v", err)
			}
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

// StreamTransactions streams transactions in real-time
func (s *BenchmarkService) StreamTransactions(req *proto.TransactionStreamRequest, stream proto.BenchmarkService_StreamTransactionsServer) error {
	// For demo purposes, we'll simulate transaction updates
	// In a real implementation, you would use WebSocket subscriptions
	for i := 0; i < 10; i++ {
		// Send transaction update
		err := stream.Send(&proto.TransactionUpdate{
			Signature:   "simulated_signature_" + fmt.Sprint(i),
			Slot:        uint64(100000 + i),
			Transaction: []byte("simulated_transaction_data"),
			Success:     true,
			Timestamp:   uint64(time.Now().Unix()),
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to send transaction update: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

// StreamBlocks streams blocks in real-time
func (s *BenchmarkService) StreamBlocks(req *proto.BlockStreamRequest, stream proto.BenchmarkService_StreamBlocksServer) error {
	// For demo purposes, we'll simulate block updates
	// In a real implementation, you would use WebSocket subscriptions
	for i := 0; i < 10; i++ {
		// Send block update
		err := stream.Send(&proto.BlockUpdate{
			Slot:              uint64(100000 + i),
			Blockhash:         fmt.Sprintf("simulated_blockhash_%d", i),
			PreviousBlockhash: fmt.Sprintf("simulated_previous_blockhash_%d", i-1),
			ParentSlot:        uint64(100000 + i - 1),
			Timestamp:         uint64(time.Now().Unix()),
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to send block update: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

// RunBenchmark runs a comprehensive benchmark suite and returns results
func (s *BenchmarkService) RunBenchmark(ctx context.Context, req *proto.BenchmarkRequest) (*proto.BenchmarkResults, error) {
	startTime := time.Now()

	results := &proto.BenchmarkResults{
		AccountGrpc:        &proto.AccountBenchmark{},
		AccountJsonrpc:     &proto.AccountBenchmark{},
		TransactionGrpc:    &proto.TransactionBenchmark{},
		TransactionJsonrpc: &proto.TransactionBenchmark{},
		BlockGrpc:          &proto.BlockBenchmark{},
		BlockJsonrpc:       &proto.BlockBenchmark{},
		Summary:            &proto.BenchmarkSummary{},
	}

	var wg sync.WaitGroup

	// Run account benchmarks
	if len(req.TestAccounts) > 0 && req.RunGrpcTests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.runAccountGrpcBenchmark(ctx, req, results.AccountGrpc)
		}()
	}

	if len(req.TestAccounts) > 0 && req.RunJsonrpcTests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.runAccountJsonRpcBenchmark(ctx, req, results.AccountJsonrpc)
		}()
	}

	// Run transaction benchmarks
	if len(req.TestSignatures) > 0 && req.RunGrpcTests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.runTransactionGrpcBenchmark(ctx, req, results.TransactionGrpc)
		}()
	}

	if len(req.TestSignatures) > 0 && req.RunJsonrpcTests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.runTransactionJsonRpcBenchmark(ctx, req, results.TransactionJsonrpc)
		}()
	}

	// Run block benchmarks
	if len(req.TestSlots) > 0 && req.RunGrpcTests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.runBlockGrpcBenchmark(ctx, req, results.BlockGrpc)
		}()
	}

	if len(req.TestSlots) > 0 && req.RunJsonrpcTests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.runBlockJsonRpcBenchmark(ctx, req, results.BlockJsonrpc)
		}()
	}

	wg.Wait()

	// Calculate summary
	totalDuration := time.Since(startTime).Milliseconds()
	results.Summary.TotalDurationMs = uint64(totalDuration)

	// Calculate speedup
	if results.AccountJsonrpc.AvgResponseTimeMs > 0 && results.AccountGrpc.AvgResponseTimeMs > 0 {
		speedup := float64(results.AccountJsonrpc.AvgResponseTimeMs) / float64(results.AccountGrpc.AvgResponseTimeMs)
		results.Summary.GrpcVsJsonrpcSpeedup = speedup

		if speedup > 1.0 {
			results.Summary.Conclusion = fmt.Sprintf("gRPC is %.2fx faster than JSON-RPC for Solana operations", speedup)
		} else {
			results.Summary.Conclusion = fmt.Sprintf("JSON-RPC is %.2fx faster than gRPC for Solana operations", 1/speedup)
		}
	}

	return results, nil
}

// Helper methods for benchmarking

func (s *BenchmarkService) runAccountGrpcBenchmark(ctx context.Context, req *proto.BenchmarkRequest, result *proto.AccountBenchmark) {
	var totalTime uint64
	var minTime uint64 = ^uint64(0) // Max uint64 value
	var maxTime uint64
	var successCount uint32

	for i := 0; i < int(req.Iterations); i++ {
		for _, account := range req.TestAccounts {
			resp, err := s.GetAccountInfo(ctx, &proto.AccountInfoRequest{
				Pubkey:         account,
				Commitment:     "finalized",
				EncodingBinary: true,
			})

			if err == nil {
				totalTime += resp.ResponseTimeMs
				successCount++

				if resp.ResponseTimeMs < minTime {
					minTime = resp.ResponseTimeMs
				}

				if resp.ResponseTimeMs > maxTime {
					maxTime = resp.ResponseTimeMs
				}
			}
		}
	}

	if successCount > 0 {
		result.AvgResponseTimeMs = totalTime / uint64(successCount)
		result.MinResponseTimeMs = minTime
		result.MaxResponseTimeMs = maxTime
		result.SuccessfulRequests = successCount
		result.FailedRequests = uint32(len(req.TestAccounts)*int(req.Iterations)) - successCount
	}
}

func (s *BenchmarkService) runAccountJsonRpcBenchmark(ctx context.Context, req *proto.BenchmarkRequest, result *proto.AccountBenchmark) {
	var totalTime uint64
	var minTime uint64 = ^uint64(0) // Max uint64 value
	var maxTime uint64
	var successCount uint32

	for i := 0; i < int(req.Iterations); i++ {
		for _, accountStr := range req.TestAccounts {
			account, err := solana.PublicKeyFromBase58(accountStr)
			if err != nil {
				continue
			}

			startTime := time.Now()

			_, err = s.solanaClient.GetAccountInfo(ctx, account)

			responseTime := time.Since(startTime).Milliseconds()

			if err == nil {
				totalTime += uint64(responseTime)
				successCount++

				if uint64(responseTime) < minTime {
					minTime = uint64(responseTime)
				}

				if uint64(responseTime) > maxTime {
					maxTime = uint64(responseTime)
				}
			}
		}
	}

	if successCount > 0 {
		result.AvgResponseTimeMs = totalTime / uint64(successCount)
		result.MinResponseTimeMs = minTime
		result.MaxResponseTimeMs = maxTime
		result.SuccessfulRequests = successCount
		result.FailedRequests = uint32(len(req.TestAccounts)*int(req.Iterations)) - successCount
	}
}

func (s *BenchmarkService) runTransactionGrpcBenchmark(ctx context.Context, req *proto.BenchmarkRequest, result *proto.TransactionBenchmark) {
	var totalTime uint64
	var minTime uint64 = ^uint64(0) // Max uint64 value
	var maxTime uint64
	var successCount uint32

	for i := 0; i < int(req.Iterations); i++ {
		for _, signature := range req.TestSignatures {
			resp, err := s.GetTransaction(ctx, &proto.TransactionRequest{
				Signature:  signature,
				Commitment: "finalized",
			})

			if err == nil {
				totalTime += resp.ResponseTimeMs
				successCount++

				if resp.ResponseTimeMs < minTime {
					minTime = resp.ResponseTimeMs
				}

				if resp.ResponseTimeMs > maxTime {
					maxTime = resp.ResponseTimeMs
				}
			}
		}
	}

	if successCount > 0 {
		result.AvgResponseTimeMs = totalTime / uint64(successCount)
		result.MinResponseTimeMs = minTime
		result.MaxResponseTimeMs = maxTime
		result.SuccessfulRequests = successCount
		result.FailedRequests = uint32(len(req.TestSignatures)*int(req.Iterations)) - successCount
	}
}

func (s *BenchmarkService) runTransactionJsonRpcBenchmark(ctx context.Context, req *proto.BenchmarkRequest, result *proto.TransactionBenchmark) {
	var totalTime uint64
	var minTime uint64 = ^uint64(0) // Max uint64 value
	var maxTime uint64
	var successCount uint32

	for i := 0; i < int(req.Iterations); i++ {
		for _, signatureStr := range req.TestSignatures {
			signature, err := solana.SignatureFromBase58(signatureStr)
			if err != nil {
				continue
			}

			startTime := time.Now()

			_, err = s.solanaClient.GetTransaction(ctx, signature, &rpc.GetTransactionOpts{})

			responseTime := time.Since(startTime).Milliseconds()

			if err == nil {
				totalTime += uint64(responseTime)
				successCount++

				if uint64(responseTime) < minTime {
					minTime = uint64(responseTime)
				}

				if uint64(responseTime) > maxTime {
					maxTime = uint64(responseTime)
				}
			}
		}
	}

	if successCount > 0 {
		result.AvgResponseTimeMs = totalTime / uint64(successCount)
		result.MinResponseTimeMs = minTime
		result.MaxResponseTimeMs = maxTime
		result.SuccessfulRequests = successCount
		result.FailedRequests = uint32(len(req.TestSignatures)*int(req.Iterations)) - successCount
	}
}

func (s *BenchmarkService) runBlockGrpcBenchmark(ctx context.Context, req *proto.BenchmarkRequest, result *proto.BlockBenchmark) {
	var totalTime uint64
	var minTime uint64 = ^uint64(0) // Max uint64 value
	var maxTime uint64
	var successCount uint32

	for i := 0; i < int(req.Iterations); i++ {
		for _, slot := range req.TestSlots {
			resp, err := s.GetBlock(ctx, &proto.BlockRequest{
				Slot:       slot,
				Commitment: "finalized",
			})

			if err == nil {
				totalTime += resp.ResponseTimeMs
				successCount++

				if resp.ResponseTimeMs < minTime {
					minTime = resp.ResponseTimeMs
				}

				if resp.ResponseTimeMs > maxTime {
					maxTime = resp.ResponseTimeMs
				}
			}
		}
	}

	if successCount > 0 {
		result.AvgResponseTimeMs = totalTime / uint64(successCount)
		result.MinResponseTimeMs = minTime
		result.MaxResponseTimeMs = maxTime
		result.SuccessfulRequests = successCount
		result.FailedRequests = uint32(len(req.TestSlots)*int(req.Iterations)) - successCount
	}
}

func (s *BenchmarkService) runBlockJsonRpcBenchmark(ctx context.Context, req *proto.BenchmarkRequest, result *proto.BlockBenchmark) {
	var totalTime uint64
	var minTime uint64 = ^uint64(0) // Max uint64 value
	var maxTime uint64
	var successCount uint32

	for i := 0; i < int(req.Iterations); i++ {
		for _, slot := range req.TestSlots {
			startTime := time.Now()

			_, err := s.solanaClient.GetBlock(ctx, slot)

			responseTime := time.Since(startTime).Milliseconds()

			if err == nil {
				totalTime += uint64(responseTime)
				successCount++

				if uint64(responseTime) < minTime {
					minTime = uint64(responseTime)
				}

				if uint64(responseTime) > maxTime {
					maxTime = uint64(responseTime)
				}
			}
		}
	}

	if successCount > 0 {
		result.AvgResponseTimeMs = totalTime / uint64(successCount)
		result.MinResponseTimeMs = minTime
		result.MaxResponseTimeMs = maxTime
		result.SuccessfulRequests = successCount
		result.FailedRequests = uint32(len(req.TestSlots)*int(req.Iterations)) - successCount
	}
}
