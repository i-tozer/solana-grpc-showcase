# Solana gRPC Exploration

A project exploring the implementation of gRPC with Solana blockchain.

## Overview

This project explores how to build a gRPC service layer on top of Solana blockchain, enabling efficient, type-safe, and bidirectional streaming communication.

### What is gRPC?

gRPC is an open-source RPC framework that uses HTTP/2 for transport and Protocol Buffers as the interface description language.

### What is Solana?

Solana is a blockchain platform designed for decentralized applications with high transaction throughput.

## Features

- Bidirectional streaming for real-time blockchain data
- Type-safe interfaces using Protocol Buffers
- Direct interaction with Solana blockchain
- Performance benchmarking comparing gRPC vs JSON-RPC

## Project Structure

```
solana-grpc-exploration/
├── proto/                  # Protocol Buffer definitions
├── server/                 # gRPC server implementation
│   ├── solana/             # Solana blockchain integration
│   └── services/           # gRPC service implementations
├── client/                 # Sample client implementations
│   ├── go/                 # Go client example
│   ├── js/                 # JavaScript client example (coming soon)
│   └── python/             # Python client example (coming soon)
├── tests/                  # Integration and unit tests (coming soon)
└── docs/                   # Documentation (coming soon)
```

## Getting Started

### Prerequisites

- Go 1.16+
- Protocol Buffers compiler
- Solana CLI tools (optional)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/i-tozer/solana-grpc-exploration.git
   cd solana-grpc-exploration
   ```

2. Install Protocol Buffers compiler:
   ```bash
   # macOS
   brew install protobuf

   # Ubuntu
   sudo apt-get install protobuf-compiler

   # Windows (using Chocolatey)
   choco install protoc
   ```

3. Install Go dependencies:
   ```bash
   go mod tidy
   ```

4. Generate Go code from Protocol Buffers:
   ```bash
   make proto
   ```

5. Build the server and client:
   ```bash
   make all
   ```

### Running the Server

Start the gRPC server:

```bash
make run-server
```

You can specify a different Solana RPC endpoint:

```bash
./bin/server --port=50051 --rpc-endpoint=https://api.devnet.solana.com
```

### Running the Client

The client provides several commands to interact with the gRPC server:

#### Benchmark

Run a performance benchmark comparing gRPC vs JSON-RPC:

```bash
make run-benchmark
```

Or with custom parameters:

```bash
./bin/client --command=benchmark --pubkey=SRMuApVNdxXokk5GT7XD5cUUgXMBCoAz2LHeuAoKWRt4 --iterations=10
```

#### Get Account Info

Retrieve information about a Solana account:

```bash
make run-account
```

Or with custom parameters:

```bash
./bin/client --command=account --pubkey=SRMuApVNdxXokk5GT7XD5cUUgXMBCoAz2LHeuAoKWRt4
```

#### Get Transaction Info

Retrieve information about a Solana transaction:

```bash
./bin/client --command=transaction --signature=YOUR_TRANSACTION_SIGNATURE
```

#### Get Block Info

Retrieve information about a Solana block:

```bash
make run-block
```

Or with custom parameters:

```bash
./bin/client --command=block --slot=150000000
```

#### Stream Account Updates

Stream real-time updates for a Solana account:

```bash
make run-stream-accounts
```

Or with custom parameters:

```bash
./bin/client --command=stream-accounts --pubkey=SRMuApVNdxXokk5GT7XD5cUUgXMBCoAz2LHeuAoKWRt4
```

#### Stream Transaction Updates

Stream real-time transaction updates:

```bash
make run-stream-transactions
```

#### Stream Block Updates

Stream real-time block updates:

```bash
make run-stream-blocks
```

## Performance Benchmarking

This project includes a benchmarking tool to compare the performance of gRPC vs JSON-RPC for Solana operations. The benchmark measures:

- Account information retrieval
- Transaction information retrieval
- Block information retrieval

The benchmark results include:
- Average response time
- Minimum response time
- Maximum response time
- Success/failure counts
- Overall speedup factor

## Technical Details

### Protocol Buffers

The project defines several Protocol Buffer message types for Solana entities:

- `AccountInfoRequest/Response`: For account information retrieval
- `TransactionRequest/Response`: For transaction information retrieval
- `BlockRequest/Response`: For block information retrieval
- `AccountUpdate`: For real-time account updates
- `TransactionUpdate`: For real-time transaction updates
- `BlockUpdate`: For real-time block updates
- `BenchmarkRequest/Results`: For performance benchmarking

### gRPC Services

The following gRPC services are implemented:

- `BenchmarkService`: Provides methods for retrieving Solana data and benchmarking performance

## Contributing

Contributions are welcome. Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.