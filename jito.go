package jitorpc

import (
	"math/rand"
	"time"

	"github.com/jito-labs/jito-go-rpc/rpc"
)

var (
	endpoints = map[Region]string{
		regions[0]: "https://mainnet.block-engine.jito.wtf",
		regions[1]: "https://amsterdam.mainnet.block-engine.jito.wtf",
		regions[2]: "https://dublin.mainnet.block-engine.jito.wtf",
		regions[3]: "https://frankfurt.mainnet.block-engine.jito.wtf",
		regions[4]: "https://london.mainnet.block-engine.jito.wtf",
		regions[5]: "https://ny.mainnet.block-engine.jito.wtf",
		regions[6]: "https://slc.mainnet.block-engine.jito.wtf",
		regions[7]: "https://singapore.mainnet.block-engine.jito.wtf",
		regions[8]: "https://tokyo.mainnet.block-engine.jito.wtf",
	}
)

type Jito rpc.Client

func NewJito(opts ...JitoOption) *Jito {
	options := &jitoOptions{}

	for _, opt := range opts {
		opt(options)
	}

	uri, ok := endpoints[options.region]
	if !ok {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		uri = endpoints[regions[rnd.Intn(len(regions))]]
	}

	if options.rps == 0 {
		options.rps = 1
	}

	header := make(map[string]string)
	if options.xJitoAuth != "" {
		header["x-jito-auth"] = options.xJitoAuth
	}

	return (*Jito)(rpc.NewWithCustomRPCClient(rpc.NewWithRateLimit(uri, options.rps, header)))
}

func (jito *Jito) Close() error {
	return (*rpc.Client)(jito).Close()
}

type jitoOptions struct {
	region    Region
	rps       int
	xJitoAuth string
}

type JitoOption func(opts *jitoOptions)

func WithRegion(region Region) JitoOption {
	return func(opts *jitoOptions) {
		opts.region = region
	}
}

func WithRPS(rps int) JitoOption {
	return func(opts *jitoOptions) {
		opts.rps = rps
	}
}

func WithXJitoAuth(xJitoAuth string) JitoOption {
	return func(opts *jitoOptions) {
		opts.xJitoAuth = xJitoAuth
	}
}
