package jitorpc

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/jito-labs/jito-go-rpc/rpc"
)

// curl https://mainnet.block-engine.jito.wtf/api/v1/transactions -X POST -H "Content-Type: application/json" -d '
//
//	{
//	  "id": 1,
//	  "jsonrpc": "2.0",
//	  "method": "sendTransaction",
//	  "params": [
//	    "AVXo5X7UNzpuOmYzkZ+fqHDGiRLTSMlWlUCcZKzEV5CIKlrdvZa3/2GrJJfPrXgZqJbYDaGiOnP99tI/sRJfiwwBAAEDRQ/n5E5CLbMbHanUG3+iVvBAWZu0WFM6NoB5xfybQ7kNwwgfIhv6odn2qTUu/gOisDtaeCW1qlwW/gx3ccr/4wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAvsInicc+E3IZzLqeA+iM5cn9kSaeFzOuClz1Z2kZQy0BAgIAAQwCAAAAAPIFKgEAAAA=",
//	    {
//	      "encoding": "base64"
//	    }
//	  ]
//	}'
func (jito *Jito) SendTransaction(ctx context.Context, instructions []solana.Instruction, recentBlockHash solana.Hash, opts ...TransactionOption) (signature solana.Signature, err error) {

	options := &transactionOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if options.tipIx == nil {
		err = fmt.Errorf("transaction must write lock at least one tip account")
		return
	}
	payer := options.payer

	var tx *solana.Transaction
	if tx, err = solana.NewTransaction(append(instructions, options.tipIx), recentBlockHash, solana.TransactionPayer(payer.PublicKey())); err != nil {
		return
	}

	if _, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		switch {
		case key.Equals(payer.PublicKey()):
			return &payer
		default:
			return nil
		}
	}); err != nil {
		return
	}

	err = (*rpc.Client)(jito).RPCCallForInto(ctx, &signature, "/api/v1/transactions", "sendTransaction", append([]any{}, tx.MustToBase64(), map[string]string{"encoding": "base64"}))
	return
}
