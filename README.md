[![Discord](https://img.shields.io/discord/938287290806042626?label=Discord&logo=discord&style=flat&color=7289DA)](https://discord.gg/jTSmEzaR)
![Go](https://img.shields.io/badge/Go-1.22.2-blue?logo=go&logoColor=white)

# jito-go-rpc
Jito JSON-RPC Go SDK

## Example

```
go mod init example
touch main.go
```

## main.go

```
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
```

Ensure that your `go.mod` has the recent package:
```
require github.com/jito-labs/jito-go-rpc v0.1.2
```

Pull in the go package:
```
go get github.com/jito-labs/jito-go-rpc
```

Run the code:
```
go run main.go
```

