package jitorpc

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/jito-labs/jito-go-rpc/rpc"
)

type TipAccount = solana.PublicKey

var (
	accounts = []TipAccount{
		solana.MustPublicKeyFromBase58("96gYZGLnJYVFmbjzopPSU6QiEV5fGqZNyN9nmNhvrZU5"),
		solana.MustPublicKeyFromBase58("DttWaMuVvTiduZRnguLF7jNxTgiMBZ1hyAumKUiL2KRL"),
		solana.MustPublicKeyFromBase58("Cw8CFyM9FkoMi7K7Crf6HNQqf4uEMzpKw6QNghXLvLkY"),
		solana.MustPublicKeyFromBase58("DfXygSm4jCyNCybVYYK6DwvWqjKee8pbDmJGcLWNDXjh"),
		solana.MustPublicKeyFromBase58("3AVi9Tg9Uo68tJfuvoKvqKNWKkC5wPdSSdeBnizKZ6jT"),
		solana.MustPublicKeyFromBase58("HFqU5x63VTqvQss8hp11i4wVV8bD44PvwucfZ2bU7gRe"),
		solana.MustPublicKeyFromBase58("ADaUMid9yfUytqMBgopwjb2DTLSokTSzL1zt6iGPaS49"),
		solana.MustPublicKeyFromBase58("ADuUkR4vqLUMWXxW9gh6D6L8pMSawimctcNZ5pGwDcEt"),
	}
)

func GetTipAccounts() []TipAccount {
	return accounts
}

// curl https://mainnet.block-engine.jito.wtf/api/v1/getTipAccounts -X POST -H "Content-Type: application/json" -d '

// 	{
// 	    "jsonrpc": "2.0",
// 	    "id": 1,
// 	    "method": "getTipAccounts",
// 	    "params": []
// 	}

// '

func (jito *Jito) GetTipAccounts(ctx context.Context) (tipAccounts []TipAccount, err error) {
	err = (*rpc.Client)(jito).RPCCallForInto(ctx, &tipAccounts, "/api/v1/getTipAccounts", "getTipAccounts", []any{})
	return
}
