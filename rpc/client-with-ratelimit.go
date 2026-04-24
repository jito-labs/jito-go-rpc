package rpc

import (
	"context"
	"io"
	"net/http"

	"github.com/jito-labs/jito-go-rpc/rpc/jsonrpc"
	"go.uber.org/ratelimit"
)

var _ JSONRPCClient = &clientWithRateLimiting{}

type clientWithRateLimiting struct {
	rpcClient   jsonrpc.RPCClient
	rateLimiter ratelimit.Limiter
}

// NewWithRateLimit creates a new rate-limitted Solana RPC client.
func NewWithRateLimit(
	rpcEndpoint string,
	rps int, // requests per second
	headers map[string]string,
) JSONRPCClient {
	opts := &jsonrpc.RPCClientOpts{
		HTTPClient:    newHTTP(),
		CustomHeaders: headers,
	}

	rpcClient := jsonrpc.NewClientWithOpts(rpcEndpoint, opts)

	return &clientWithRateLimiting{
		rpcClient:   rpcClient,
		rateLimiter: ratelimit.New(rps),
	}
}

func (wr *clientWithRateLimiting) CallForInto(ctx context.Context, out any, path string, method string, params []any) error {
	wr.rateLimiter.Take()
	return wr.rpcClient.CallForInto(ctx, &out, path, method, params)
}

func (wr *clientWithRateLimiting) CallWithCallback(
	ctx context.Context,
	path string,
	method string,
	params []any,
	callback func(*http.Request, *http.Response) error,
) error {
	wr.rateLimiter.Take()
	return wr.rpcClient.CallWithCallback(ctx, path, method, params, callback)
}

func (wr *clientWithRateLimiting) CallBatch(
	ctx context.Context,
	path string,
	requests jsonrpc.RPCRequests,
) (jsonrpc.RPCResponses, error) {
	wr.rateLimiter.Take()
	return wr.rpcClient.CallBatch(ctx, path, requests)
}

// Close closes clientWithRateLimiting.
func (cl *clientWithRateLimiting) Close() error {
	if c, ok := cl.rpcClient.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
