package jitorpc

import (
	"context"

	"github.com/jito-labs/jito-go-rpc/rpc"
)

type BundleStatus string

var (
	ConfirmationStatusProcessed BundleStatus = "processed"
	ConfirmationStatusConfirmed BundleStatus = "confirmed"
	ConfirmationStatusFinalized BundleStatus = "finalized"
)

type GetBundleStatusesResp struct {
	Context struct {
		Slot int64 `json:"slot"`
	} `json:"context"`
	Value []*BundleNode `json:"value"`
}

type BundleNode struct {
	BundleID     string       `json:"bundle_id"`
	Transactions []string     `json:"transactions"`
	Slot         int64        `json:"slot"`
	Status       BundleStatus `json:"confirmation_status"`
	Err          struct {
		Ok any `json:"Ok"`
	} `json:"err"`
}

// curl https://mainnet.block-engine.jito.wtf/api/v1/getBundleStatuses -X POST -H "Content-Type: application/json" -d '
//
//	{
//	    "jsonrpc": "2.0",
//	    "id": 1,
//	    "method": "getBundleStatuses",
//	    "params": [
//	      [
//	        "892b79ed49138bfb3aa5441f0df6e06ef34f9ee8f3976c15b323605bae0cf51d"
//	      ]
//	    ]
//	}
//
// '
func (jito *Jito) GetBundleStatuses(ctx context.Context, bundleIDs []string) (resp GetBundleStatusesResp, err error) {
	err = (*rpc.Client)(jito).RPCCallForInto(ctx, &resp, "/api/v1/getBundleStatuses", "getBundleStatuses", append([]any{}, bundleIDs))
	return
}
