package main

import (
	"fmt"
	"log"

	jitorpc "github.com/jito-labs/jito-go-rpc"
)

func main() {
	// Initialize Jito client
	jitoClient := jitorpc.NewJitoJsonRpcClient("https://bundles.jito.wtf/api/v1", "")
	debug := true
	jitoClient.Debug = &debug

	// Get the tip floor
	tipFloors, err := jitoClient.GetTipFloors()
	if err != nil {
		log.Fatalf("Failed to get tip floors: %v", err)
	}

	fmt.Printf("Jito tip amount: %v\n", tipFloors)
}
