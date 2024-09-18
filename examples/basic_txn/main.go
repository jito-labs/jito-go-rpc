package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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
	walletData, err := os.ReadFile(walletPath)
	if err != nil {
		log.Fatalf("Failed to read wallet file: %v", err)
	}

	var privateKeyBytes []byte
	err = json.Unmarshal(walletData, &privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to unmarshal private key: %v", err)
	}

	privateKey := solana.PrivateKey(privateKeyBytes)

	// Set up transaction parameters
	receiver, err := solana.PublicKeyFromBase58("RECEIVER_PUBKEY")
	if err != nil {
		log.Fatalf("Failed to parse receiver public key: %v", err)
	}

	// Get a random tip account
	randomTipAccount, err := jitoClient.GetRandomTipAccount()
	if err != nil {
		log.Fatalf("Failed to get random tip account: %v", err)
	}

	jitoTipAccount, err := solana.PublicKeyFromBase58(randomTipAccount.Address)
	if err != nil {
		log.Fatalf("Failed to parse Jito tip account public key: %v", err)
	}

	jitoTipAmount := uint64(1000)  // lamports
	transferAmount := uint64(1000) // lamports
	priorityFee := uint64(1000)    // lamports

	// Get latest blockhash
	latestBlockhash, err := solanaClient.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("Failed to get latest blockhash: %v", err)
	}

	// Set this to true if you want to use bundle-only mode
	bundleOnly := false

	var instructions []solana.Instruction

	if !bundleOnly {
		instructions = append(instructions, createSetComputeUnitPriceInstruction(priorityFee))
	}

	instructions = append(instructions,
		system.NewTransferInstruction(
			transferAmount,
			privateKey.PublicKey(),
			receiver,
		).Build(),
		system.NewTransferInstruction(
			jitoTipAmount,
			privateKey.PublicKey(),
			jitoTipAccount,
		).Build(),
	)

	tx, err := solana.NewTransaction(
		instructions,
		latestBlockhash.Value.Blockhash,
		solana.TransactionPayer(privateKey.PublicKey()),
	)
	if err != nil {
		log.Fatalf("Failed to create transaction: %v", err)
	}

	// Sign transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if privateKey.PublicKey().Equals(key) {
			return &privateKey
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Serialize and base58 encode the signed transaction
	serializedTx, err := tx.MarshalBinary()
	if err != nil {
		log.Fatalf("Failed to serialize transaction: %v", err)
	}
	base58EncodedTx := base58.Encode(serializedTx)

	// Prepare the transaction request
	txnRequest := []string{base58EncodedTx}

	// Send the transaction
	fmt.Printf("Sending transaction request (bundleOnly=%v): %s\n", bundleOnly, txnRequest)
	var txSignature string

	if bundleOnly {
		bundleRequest := [][]string{txnRequest}
		result, err := jitoClient.SendBundle(bundleRequest)
		if err != nil {
			log.Fatalf("Failed to send bundle: %v", err)
		}
		if err := json.Unmarshal(result, &txSignature); err != nil {
			log.Fatalf("Failed to unmarshal bundle ID: %v", err)
		}
		fmt.Printf("Bundle sent successfully. Bundle ID: %s\n", txSignature)
		checkBundleStatus(jitoClient, txSignature)
	} else {
		result, err := jitoClient.SendTxn(txnRequest, false)
		if err != nil {
			log.Fatalf("Failed to send transaction: %v", err)
		}
		txSignature = strings.Trim(string(result), "\"")
		fmt.Printf("Transaction signature: %s\n", txSignature)
		checkTransactionStatus(solanaClient, txSignature)
	}
}

func createSetComputeUnitPriceInstruction(microLamports uint64) solana.Instruction {
	data := make([]byte, 9)
	data[0] = 3 // Instruction index for SetComputeUnitPrice
	binary.LittleEndian.PutUint64(data[1:], microLamports)

	return solana.NewInstruction(
		solana.MustPublicKeyFromBase58("ComputeBudget111111111111111111111111111111"),
		solana.AccountMetaSlice{},
		data,
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

func checkTransactionStatus(solanaClient *rpc.Client, txSignature string) {
	sigBytes, err := base58.Decode(txSignature)
	if err != nil {
		log.Fatalf("Failed to decode signature: %v", err)
	}
	var sig solana.Signature
	copy(sig[:], sigBytes)

	for i := 0; i < 120; i++ { // Increased max attempts to allow more time for 30 confirmations
		time.Sleep(1 * time.Second)
		status, err := solanaClient.GetSignatureStatuses(context.Background(), true, sig)
		if err != nil {
			log.Printf("Failed to get signature status: %v", err)
			continue
		}

		if status.Value[0] == nil {
			log.Printf("Attempt %d: Transaction status not available yet", i+1)
			continue
		}

		confirmations := uint64(0)
		if status.Value[0].Confirmations != nil {
			confirmations = *status.Value[0].Confirmations
		}

		log.Printf("Attempt %d: Transaction status:", i+1)
		log.Printf("  Confirmations: %d", confirmations)
		log.Printf("  Slot: %v", status.Value[0].Slot)
		log.Printf("  Err: %v", status.Value[0].Err)

		if status.Value[0].Err != nil {
			log.Fatalf("Transaction failed: %v", status.Value[0].Err)
		}

		if confirmations >= 27 {
			fmt.Printf("Transaction confirmed with %d confirmations\n", confirmations)
			solscanURL := fmt.Sprintf("https://solscan.io/tx/%s", txSignature)
			fmt.Printf("View transaction on Solscan: %s\n", solscanURL)
			return
		}
	}
	log.Printf("Transaction did not reach 30 confirmations after multiple attempts")
	solscanURL := fmt.Sprintf("https://solscan.io/tx/%s", txSignature)
	fmt.Printf("View transaction on Solscan (may not have 30 confirmations yet): %s\n", solscanURL)
}
