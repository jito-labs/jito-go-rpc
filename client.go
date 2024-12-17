package jitorpc

import (
	"net/http"
)

// JitoJsonRpcClient represents a client for the Jito JSON-RPC API.
type JitoJsonRpcClient struct {
	BaseURL string
	UUID    string
	Client  *http.Client
	Debug   *bool
}

// NewJitoJsonRpcClient creates a new JitoJsonRpcClient.
func NewJitoJsonRpcClient(baseURL string, uuid string) *JitoJsonRpcClient {
	return &JitoJsonRpcClient{
		BaseURL: baseURL,
		UUID:    uuid,
		Client:  &http.Client{},
	}
}

func (c *JitoJsonRpcClient) isDebugEnabled() bool {
	return c.Debug != nil && *c.Debug
}
