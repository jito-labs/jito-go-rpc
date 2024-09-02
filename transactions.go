package jitorpc

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *JitoJsonRpcClient) SendTxn(params interface{}, bundleOnly bool) (json.RawMessage, error) {
	endpoint := "/transactions"
	queryParams := []string{}

	if bundleOnly {
		queryParams = append(queryParams, "bundleOnly=true")
	}
	if c.UUID != "" {
		queryParams = append(queryParams, fmt.Sprintf("uuid=%s", c.UUID))
	}

	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, strings.Join(queryParams, "&"))
	}

	return c.sendRequest(endpoint, "sendTransaction", params)
}
