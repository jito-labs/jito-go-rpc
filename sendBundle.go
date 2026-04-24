package jitorpc

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/jito-labs/jito-go-rpc/rpc"
)

// curl https://mainnet.block-engine.jito.wtf:443/api/v1/bundles -X POST -H "Content-Type: application/json" -d '
//
//	{
//	  "jsonrpc": "2.0",
//	  "id": 1,
//	  "method": "sendBundle",
//	  "params": [
//	    [
//	      "AT2AqtlokikUWgGNnSX5xrmdvBjSaiIPxvFz6zc5Abn5Z0CPFW5GO+Y3rXceLnqLgQFnGw0yTk3NtJdFNsbrwwQBAAIEsXPDJ9cMVbpFQYClVM7PGLh8JOfCD6E2vz5VNmBCF+p4Uhyxec67hYm1VqLV7JTSSYaC/fm7KvWtZOSRzEFT2gAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABUpTWpkpIQZNJOhxYNo4fHw1td28kruB5B+oQEEFRI1i3Wzl2VfewCI8oYXParnP78725sKFzYheTEn8v865YQIDABhqaXRvIGJ1bmRsZSAwOiBqaXRvIHRlc3QCAgABDAIAAACghgEAAAAAAA==",
//	      "AS6fOZuGDsmyYdd+RC0fiFUgNe1BYTOYT+1hkRXHAeroC8R60h3g34EPF5Ys8sGzVBMP9MDSTVgy1/SSTqpCtA4BAAIEsXPDJ9cMVbpFQYClVM7PGLh8JOfCD6E2vz5VNmBCF+p4Uhyxec67hYm1VqLV7JTSSYaC/fm7KvWtZOSRzEFT2gAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABUpTWpkpIQZNJOhxYNo4fHw1td28kruB5B+oQEEFRI1i3Wzl2VfewCI8oYXParnP78725sKFzYheTEn8v865YQIDABhqaXRvIGJ1bmRsZSAxOiBqaXRvIHRlc3QCAgABDAIAAACghgEAAAAAAA==",
//	      "AZBed/TbhHt8WTUbfdT9Ibi8vr9ZpPCPtav8HIUrlDaqqKOLuJldEPpb899WvnfrCm4Qtq7i/MNLDTde1sHVYwEBAAIEsXPDJ9cMVbpFQYClVM7PGLh8JOfCD6E2vz5VNmBCF+p4Uhyxec67hYm1VqLV7JTSSYaC/fm7KvWtZOSRzEFT2gAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABUpTWpkpIQZNJOhxYNo4fHw1td28kruB5B+oQEEFRI1i3Wzl2VfewCI8oYXParnP78725sKFzYheTEn8v865YQIDABhqaXRvIGJ1bmRsZSAyOiBqaXRvIHRlc3QCAgABDAIAAACghgEAAAAAAA==",
//	      "AbS0zOISjAjVgUqs+bBh8x93oNOBzVcBnwicjyXhKLJ6E0/EPi+HbLcBN1u/87Kf4nNS+PoK3lLmhHtLjTyyDw0BAAIEsXPDJ9cMVbpFQYClVM7PGLh8JOfCD6E2vz5VNmBCF+p4Uhyxec67hYm1VqLV7JTSSYaC/fm7KvWtZOSRzEFT2gAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABUpTWpkpIQZNJOhxYNo4fHw1td28kruB5B+oQEEFRI1i3Wzl2VfewCI8oYXParnP78725sKFzYheTEn8v865YQIDABhqaXRvIGJ1bmRsZSAzOiBqaXRvIHRlc3QCAgABDAIAAACghgEAAAAAAA==",
//	      "AVfm3snz2ZU2xJlj4tBYLNKfuNUC57NgV+ci/2eL4vUJt+NfstiUY+m1kriNMtDXhiBf28wR6WKVuoLkuiegNQABAAIEsXPDJ9cMVbpFQYClVM7PGLh8JOfCD6E2vz5VNmBCF+p4Uhyxec67hYm1VqLV7JTSSYaC/fm7KvWtZOSRzEFT2gAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABUpTWpkpIQZNJOhxYNo4fHw1td28kruB5B+oQEEFRI1i3Wzl2VfewCI8oYXParnP78725sKFzYheTEn8v865YQIDABhqaXRvIGJ1bmRsZSA0OiBqaXRvIHRlc3QCAgABDAIAAACghgEAAAAAAA=="
//	    ],
//	    {
//	      "encoding": "base64"
//	    }
//	  ]
//	}'
func (jito *Jito) SendBundle(ctx context.Context, transactions []*solana.Transaction, recentBlockHash solana.Hash, opts ...TransactionOption) (bundleID string, err error) {

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
	if tx, err = solana.NewTransaction([]solana.Instruction{options.tipIx}, recentBlockHash, solana.TransactionPayer(payer.PublicKey())); err != nil {
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

	transactions = append(transactions, tx)

	txs := make([]string, len(transactions))

	for k, v := range transactions {
		txs[k] = v.MustToBase64()
	}

	err = (*rpc.Client)(jito).RPCCallForInto(ctx, &bundleID, "/api/v1/bundles", "sendBundle", append([]any{}, txs, map[string]string{"encoding": "base64"}))
	return
}
