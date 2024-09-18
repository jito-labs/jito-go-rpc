# jito-go-rpc

[![Discord](https://img.shields.io/discord/938287290806042626?label=Discord&logo=discord&style=flat&color=7289DA)](https://discord.gg/jTSmEzaR)
![Go](https://img.shields.io/badge/Go-1.22.2-blue?logo=go&logoColor=white)

The Jito JSON-RPC Go SDK provides an interface for interacting with Jito's enhanced Solana infrastructure. This SDK supports methods for managing bundles and transactions, offering improved performance and additional features while interacting with the Block Engine.

## Features

### Bundles
- `GetInflightBundleStatuses`: Retrieve the status of in-flight bundles.
- `GetBundleStatuses`: Fetch the statuses of submitted bundles.
- `GetTipAccounts`: Get accounts eligible for tips.
- `SendBundle`: Submit bundles to the Jito Block Engine.

### Transactions
- `SendTransaction`: Submit transactions with enhanced priority and speed.

## Installation

### Prerequisites

This project requires Go 1.22.2 or higher. If you haven't installed Go yet, follow these steps:

1. **Install Go**:
   Download and install Go from [golang.org](https://golang.org/dl/)

2. Verify the installation:
   ```bash
   go version
   ```

### Installing jito-go-rpc

Install the SDK using go get:

```bash
go get github.com/jito-labs/jito-go-rpc
```

## Usage Examples

### Basic Transaction Example

To run the basic transaction example:

1. Ensure your environment is set up in `basic_txn.go`:

   ```go
   // Load the sender's keypair
   walletPath := "/path/to/wallet.json"
   privateKey, err := solana.PrivateKeyFromSolanaKeygenFile(walletPath)
   if err != nil {
       log.Fatalf("Failed to load private key: %v", err)
   }

   // Set up receiver pubkey
   receiver, err := solana.PublicKeyFromBase58("YOUR_RECEIVER_KEY")
   if err != nil {
       log.Fatalf("Failed to parse receiver public key: %v", err)
   }
   ```

2. Run the example:
   ```bash
   go run basic_txn.go
   ```

### Basic Bundle Example

To run the basic bundle example:

1. Ensure your environment is set up in `basic_bundle.go`:

   ```go
   // Load the sender's keypair
   walletPath := "/path/to/wallet.json"
   privateKey, err := solana.PrivateKeyFromSolanaKeygenFile(walletPath)
   if err != nil {
       log.Fatalf("Failed to load private key: %v", err)
   }

   // Set up receiver pubkey
   receiver, err := solana.PublicKeyFromBase58("YOUR_RECEIVER_KEY")
   if err != nil {
       log.Fatalf("Failed to parse receiver public key: %v", err)
   }
   ```

2. Run the example:
   ```bash
   go run basic_bundle.go
   ```

These examples demonstrate how to set up and run basic transactions and bundles using the Jito Go RPC SDK. Make sure to replace the wallet path and receiver key with your actual values.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For support, please join our [Discord community](https://discord.gg/jTSmEzaR).
