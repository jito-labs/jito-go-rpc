package jitorpc

import (
	"context"

	"github.com/jito-labs/jito-go-rpc/rpc"
)

type InflightBundleStatus string

var (
	InflightBundleStatusInvalid InflightBundleStatus = "Invalid"
	InflightBundleStatusPending InflightBundleStatus = "Pending"
	InflightBundleStatusFailed  InflightBundleStatus = "Failed"
	InflightBundleStatusLanded  InflightBundleStatus = "Landed"
)

type GetInflightBundleStatusesResp struct {
	Context struct {
		Slot int64 `json:"slot"`
	} `json:"context"`
	Value []*InflightBundleNode `json:"value"`
}

type InflightBundleNode struct {
	BundleID string               `json:"bundle_id"`
	Status   InflightBundleStatus `json:"status"`
	Slot     *int64               `json:"landed_slot"`
}

// curl https://mainnet.block-engine.jito.wtf/api/v1/getInflightBundleStatuses -X POST -H "Content-Type: application/json" -d '
//
//	{
//	  "jsonrpc": "2.0",
//	  "id": 1,
//	  "method": "getInflightBundleStatuses",
//	  "params": [
//	    [
//	      "b31e5fae4923f345218403ac1ab242b46a72d4f2a38d131f474255ae88f1ec9a",
//	      "e3c4d7933cf3210489b17307a14afbab2e4ae3c67c9e7157156f191f047aa6e8",
//	      "a7abecabd9a165bc73fd92c809da4dc25474e1227e61339f02b35ce91c9965e2",
//	      "e3934d2f81edbc161c2b8bb352523cc5f74d49e8d4db81b222c553de60a66514",
//	      "2cd515429ae99487dfac24b170248f6929e4fd849aa7957cccc1daf75f666b54"
//	    ]
//	  ]
//	}
//
// '
func (jito *Jito) GetInflightBundleStatuses(ctx context.Context, bundleIDs []string) (resp GetInflightBundleStatusesResp, err error) {
	err = (*rpc.Client)(jito).RPCCallForInto(ctx, &resp, "/api/v1/getInflightBundleStatuses", "getInflightBundleStatuses", append([]any{}, bundleIDs))
	return
}
