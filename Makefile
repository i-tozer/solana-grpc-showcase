.PHONY: all proto server client clean

# Default target
all: proto server client

# Generate Go code from Protocol Buffers
proto:
	@echo "Generating Go code from Protocol Buffers..."
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/solana_benchmark.proto

# Build the server
server:
	@echo "Building server..."
	@go build -o bin/server server/main.go

# Build the client
client:
	@echo "Building client..."
	@go build -o bin/client client/go/main.go

# Run the server
run-server:
	@echo "Running server..."
	@./bin/server --port=50051 --rpc-endpoint=https://api.mainnet-beta.solana.com

# Run the client with benchmark command
run-benchmark:
	@echo "Running benchmark..."
	@./bin/client --command=benchmark --pubkey=CKJCVxuM99Rn3v6SBxCQ5osdwuKkWBWbdKG38pYXdfrj --iterations=5

# Run the client with account command
run-account:
	@echo "Getting account info..."
	@./bin/client --command=account --pubkey=CKJCVxuM99Rn3v6SBxCQ5osdwuKkWBWbdKG38pYXdfrj

# Run the client with transaction command
run-transaction:
	@echo "Getting transaction info..."
	@./bin/client --command=transaction --signature=YOUR_TRANSACTION_SIGNATURE

# Run the client with block command
run-block:
	@echo "Getting block info..."
	@./bin/client --command=block --slot=150000000

# Run the client with stream-accounts command
run-stream-accounts:
	@echo "Streaming account updates..."
	@./bin/client --command=stream-accounts --pubkey=CKJCVxuM99Rn3v6SBxCQ5osdwuKkWBWbdKG38pYXdfrj

# Run the client with stream-transactions command
run-stream-transactions:
	@echo "Streaming transaction updates..."
	@./bin/client --command=stream-transactions

# Run the client with stream-blocks command
run-stream-blocks:
	@echo "Streaming block updates..."
	@./bin/client --command=stream-blocks

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f proto/*.pb.go 