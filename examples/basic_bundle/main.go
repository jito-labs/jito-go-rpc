package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	jitorpc "github.com/jito-labs/jito-go-rpc"
)

func main() {
	// Initialize Solana client
	solanaClient := rpc.New("https://api.mainnet-beta.solana.com")

	// Initialize Jito client
	jitoClient := jitorpc.NewJito(
		jitorpc.WithRegion(jitorpc.GetRegions()[0]),
		jitorpc.WithRPS(1),
		jitorpc.WithXJitoAuth(""),
	)

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
	tipAccounts, err := jitoClient.GetTipAccounts(context.Background())
	if err != nil {
		log.Fatalf("Failed to get random tip account: %v", err)
	}
	tipAccount := tipAccounts[rand.Intn(len(tipAccounts))]

	tipAmount := uint64(1000) // lamports

	// Create main transaction
	mainTx, err := createMainTransaction(privateKey, latestBlockhash.Value.Blockhash)
	if err != nil {
		log.Fatalf("Failed to create main transaction: %v", err)
	}

	// Send the bundle
	bundleID, err := jitoClient.SendBundle(
		context.TODO(),
		[]*solana.Transaction{mainTx},
		latestBlockhash.Value.Blockhash,
		jitorpc.WithTransactionTip(
			privateKey,
			tipAccount,
			tipAmount,
		),
	)
	if err != nil {
		log.Fatalf("Failed to send bundle: %v", err)
	}

	fmt.Printf("Bundle sent successfully. Bundle ID: %s\n", bundleID)

	// Check the bundle status
	checkBundleStatus(jitoClient, bundleID)
}

func createMainTransaction(privateKey solana.PrivateKey, recentBlockhash solana.Hash) (*solana.Transaction, error) {
	receiver, err := solana.PublicKeyFromBase58("RECIEVE_KEY")
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

func createMemoInstruction(message string) solana.Instruction {
	memoProgramID, _ := solana.PublicKeyFromBase58("MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr")
	return solana.NewInstruction(
		memoProgramID,
		solana.AccountMetaSlice{},
		[]byte(message),
	)
}

func checkBundleStatus(jitoClient *jitorpc.Jito, bundleID string) {
	maxAttempts := 60
	pollInterval := 5 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		time.Sleep(pollInterval)

		statusResponse, err := jitoClient.GetBundleStatuses(context.Background(), []string{bundleID})
		if err != nil {
			log.Printf("Attempt %d: Failed to get bundle status: %v", attempt, err)
			continue
		}

		if len(statusResponse.Value) == 0 {
			log.Printf("Attempt %d: No bundle status available", attempt)
			continue
		}

		bundleStatus := statusResponse.Value[0]
		log.Printf("Attempt %d: Bundle status: %s", attempt, bundleStatus.Status)

		switch bundleStatus.Status {
		case jitorpc.ConfirmationStatusProcessed:
			fmt.Println("Bundle has been processed by the cluster. Continuing to poll...")
		case jitorpc.ConfirmationStatusConfirmed:
			fmt.Println("Bundle has been confirmed by the cluster. Continuing to poll...")
		case jitorpc.ConfirmationStatusFinalized:
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
			fmt.Printf("Unexpected status: %s. Please check the bundle manually.\n", bundleStatus.Status)
			return
		}
	}

	log.Printf("Maximum polling attempts reached. Final status unknown.")
}
