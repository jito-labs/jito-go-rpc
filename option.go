package jitorpc

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

type transactionOptions struct {
	payer solana.PrivateKey
	tipIx *system.Instruction
}

type TransactionOption func(opts *transactionOptions)

func WithTransactionTip(from solana.PrivateKey, to TipAccount, lamports uint64) TransactionOption {
	return func(opts *transactionOptions) {
		opts.payer = from
		opts.tipIx = system.NewTransferInstruction(
			lamports,
			from.PublicKey(),
			solana.PublicKey(to),
		).Build()
	}
}
