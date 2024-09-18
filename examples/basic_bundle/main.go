package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	jitorpc "github.com/jito-labs/jito-go-rpc"
	"github.com/mr-tron/base58"
)

func main() {
	// Initialize Solana client
	solanaClient := rpc.New("https://api.mainnet-beta.solana.com")

	// Initialize Jito client
	jitoClient := jitorpc.NewJitoJsonRpcClient("https://mainnet.block-engine.jito.wtf/api/v1", "")

	// Load wallet from local path
	walletPath := "/path/to/wallet.json"
	privateKey, err := solana.PrivateKeyFromSolanaKeygenFile(walletPath)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Get latest blockhash
	latestBlockhash, err := solanaClient.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("Failed to get latest blockhash: %v", err)
	}

	// Get random tip account
	tipAccount, err := jitoClient.GetRandomTipAccount()
	if err != nil {
		log.Fatalf("Failed to get random tip account: %v", err)
	}

	// Create tip transaction
	tipAmount := uint64(1000) // lamports
	tipTx, err := createTipTransaction(privateKey, tipAmount, latestBlockhash.Value.Blockhash, tipAccount.Address)
	if err != nil {
		log.Fatalf("Failed to create tip transaction: %v", err)
	}

	// Create main transaction
	mainTx, err := createMainTransaction(privateKey, latestBlockhash.Value.Blockhash)
	if err != nil {
		log.Fatalf("Failed to create main transaction: %v", err)
	}

	// Prepare the bundle request
	bundleRequest := [][]string{{
		encodeTransaction(tipTx),
		encodeTransaction(mainTx),
	}}

	// Send the bundle
	fmt.Printf("Sending bundle request: %v\n", bundleRequest)
	bundleIdRaw, err := jitoClient.SendBundle(bundleRequest)
	if err != nil {
		log.Fatalf("Failed to send bundle: %v", err)
	}

	var bundleId string
	if err := json.Unmarshal(bundleIdRaw, &bundleId); err != nil {
		log.Fatalf("Failed to unmarshal bundle ID: %v", err)
	}

	fmt.Printf("Bundle sent successfully. Bundle ID: %s\n", bundleId)

	// Check the bundle status
	checkBundleStatus(jitoClient, bundleId)
}

func createTipTransaction(privateKey solana.PrivateKey, amount uint64, recentBlockhash solana.Hash, tipAddress string) (*solana.Transaction, error) {
	tipAccount, err := solana.PublicKeyFromBase58(tipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tip account: %v", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amount,
				privateKey.PublicKey(),
				tipAccount,
			).Build(),
		},
		recentBlockhash,
		solana.TransactionPayer(privateKey.PublicKey()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tip transaction: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if privateKey.PublicKey().Equals(key) {
			return &privateKey
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sign tip transaction: %v", err)
	}

	return tx, nil
}

func createMainTransaction(privateKey solana.PrivateKey, recentBlockhash solana.Hash) (*solana.Transaction, error) {
	receiver, err := solana.PublicKeyFromBase58("RECIEVER_PUBKEY")
	if err != nil {
		return nil, fmt.Errorf("failed to parse receiver public key: %v", err)
	}

	transferAmount := uint64(1000) // lamports

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				transferAmount,
				privateKey.PublicKey(),
				receiver,
			).Build(),
			createMemoInstruction("Hello, Jito!"),
		},
		recentBlockhash,
		solana.TransactionPayer(privateKey.PublicKey()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create main transaction: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if privateKey.PublicKey().Equals(key) {
			return &privateKey
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sign main transaction: %v", err)
	}

	return tx, nil
}

func encodeTransaction(tx *solana.Transaction) string {
	serializedTx, err := tx.MarshalBinary()
	if err != nil {
		log.Fatalf("Failed to serialize transaction: %v", err)
	}
	return base58.Encode(serializedTx)
}

func createMemoInstruction(message string) solana.Instruction {
	memoProgramID, _ := solana.PublicKeyFromBase58("MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr")
	return solana.NewInstruction(
		memoProgramID,
		solana.AccountMetaSlice{},
		[]byte(message),
	)
}

func checkBundleStatus(jitoClient *jitorpc.JitoJsonRpcClient, bundleId string) {
	maxAttempts := 60
	pollInterval := 5 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		time.Sleep(pollInterval)

		statusResponse, err := jitoClient.GetBundleStatuses([]string{bundleId})
		if err != nil {
			log.Printf("Attempt %d: Failed to get bundle status: %v", attempt, err)
			continue
		}

		if len(statusResponse.Value) == 0 {
			log.Printf("Attempt %d: No bundle status available", attempt)
			continue
		}

		bundleStatus := statusResponse.Value[0]
		log.Printf("Attempt %d: Bundle status: %s", attempt, bundleStatus.ConfirmationStatus)

		switch bundleStatus.ConfirmationStatus {
		case "processed":
			fmt.Println("Bundle has been processed by the cluster. Continuing to poll...")
		case "confirmed":
			fmt.Println("Bundle has been confirmed by the cluster. Continuing to poll...")
		case "finalized":
			fmt.Printf("Bundle has been finalized by the cluster in slot %d.\n", bundleStatus.Slot)
			if bundleStatus.Err.Ok == nil {
				fmt.Println("Bundle executed successfully.")
				fmt.Println("Transaction URLs:")
				for _, txID := range bundleStatus.Transactions {
					solscanURL := fmt.Sprintf("https://solscan.io/tx/%s", txID)
					fmt.Printf("- %s\n", solscanURL)
				}
			} else {
				fmt.Printf("Bundle execution failed with error: %v\n", bundleStatus.Err.Ok)
			}
			return
		default:
			fmt.Printf("Unexpected status: %s. Please check the bundle manually.\n", bundleStatus.ConfirmationStatus)
			return
		}
	}

	log.Printf("Maximum polling attempts reached. Final status unknown.")
}
