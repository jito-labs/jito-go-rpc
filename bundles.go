package jitorpc

import (
	"encoding/json"
	"fmt"
)

func (c *JitoJsonRpcClient) GetTipAccounts() (json.RawMessage, error) {
	endpoint := "/bundles"
	if c.UUID != "" {
		endpoint = fmt.Sprintf("%s?uuid=%s", endpoint, c.UUID)
	}
	return c.sendRequest(endpoint, "getTipAccounts", nil)
}

func (c *JitoJsonRpcClient) GetBundleStatuses(params interface{}) (json.RawMessage, error) {
	endpoint := "/bundles"
	if c.UUID != "" {
		endpoint = fmt.Sprintf("%s?uuid=%s", endpoint, c.UUID)
	}
	return c.sendRequest(endpoint, "getBundleStatuses", params)
}

func (c *JitoJsonRpcClient) SendBundle(params interface{}) (json.RawMessage, error) {
	endpoint := "/bundles"
	if c.UUID != "" {
		endpoint = fmt.Sprintf("%s?uuid=%s", endpoint, c.UUID)
	}
	return c.sendRequest(endpoint, "sendBundle", params)
}
