package main

import (
	"fmt"
	"log"

	jitorpc "github.com/jito-labs/jito-go-rpc"
)

func main() {
	//Implementation with no UUID(default)
	client := jitorpc.NewJitoJsonRpcClient("https://mainnet.block-engine.jito.wtf/api/v1", "")

	//Implementation with UUID
	//client := jitorpc.NewJitoJsonRpcClient("https://mainnet.block-engine.jito.wtf/api/v1", "YOUR_UUID")

	// Make an RPC call to get tip accounts
	tipAccounts, err := client.GetTipAccounts()
	if err != nil {
		log.Fatalf("Error getting tip accounts: %v", err)
	}

	// Print the prettified JSON response
	fmt.Println("Tip Accounts:")
	fmt.Println(jitorpc.PrettifyJSON(tipAccounts))

}
