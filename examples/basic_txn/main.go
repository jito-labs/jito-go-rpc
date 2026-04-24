package main

import (
	"context"
	"encoding/binary"
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
		jitorpc.WithRPS(1), // 1 rps
		jitorpc.WithXJitoAuth(""),
	)

	// Load wallet from local path
	walletPath := "/path/to/wallet.json"
	privateKey, err := solana.PrivateKeyFromSolanaKeygenFile(walletPath)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Set up transaction parameters
	receiver, err := solana.PublicKeyFromBase58("RECIEVER_KEY")
	if err != nil {
		log.Fatalf("Failed to parse receiver public key: %v", err)
	}

	// Get random tip account
	tipAccounts, err := jitoClient.GetTipAccounts(context.Background())
	if err != nil {
		log.Fatalf("Failed to get random tip account: %v", err)
	}
	jitoTipAccount := tipAccounts[rand.Intn(len(tipAccounts))]

	jitoTipAmount := uint64(1000)  // lamports
	transferAmount := uint64(1000) // lamports
	priorityFee := uint64(1000)    // lamports

	// Get latest blockhash
	latestBlockhash, err := solanaClient.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("Failed to get latest blockhash: %v", err)
	}

	var instructions []solana.Instruction

	instructions = append(instructions, createSetComputeUnitPriceInstruction(priorityFee))

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

	sig, err := jitoClient.SendTransaction(
		context.TODO(),
		instructions,
		latestBlockhash.Value.Blockhash,
		jitorpc.WithTransactionTip(
			privateKey,
			jitoTipAccount,
			jitoTipAmount,
		),
	)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction signature: %s\n", sig.String())
	checkTransactionStatus(solanaClient, sig)
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

func checkTransactionStatus(solanaClient *rpc.Client, sig solana.Signature) {

	for i := 0; i < 120; i++ {
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
			solscanURL := fmt.Sprintf("https://solscan.io/tx/%s", sig.String())
			fmt.Printf("View transaction on Solscan: %s\n", solscanURL)
			return
		}
	}
	log.Printf("Transaction did not reach 30 confirmations after multiple attempts")
	solscanURL := fmt.Sprintf("https://solscan.io/tx/%s", sig.String())
	fmt.Printf("View transaction on Solscan (may not have 30 confirmations yet): %s\n", solscanURL)
}
