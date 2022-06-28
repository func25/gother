package gother

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

func HeaderLatest(ctx context.Context) (*types.Header, error) {
	return Client.HeaderByNumber(ctx, nil)
}

func HeaderNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return Client.HeaderByNumber(ctx, number)
}

func (a GotherClient) HeaderLatest(ctx context.Context) (*types.Header, error) {
	return a.Client.HeaderByNumber(ctx, nil)
}

func (a GotherClient) HeaderNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return a.Client.HeaderByNumber(ctx, number)
}
