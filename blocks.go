package gother

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

func (c *GotherClient) HeaderLatest(ctx context.Context) (*types.Header, error) {
	return c.HeaderByNumber(ctx, nil)
}

func (c *GotherClient) HeaderNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return c.HeaderByNumber(ctx, number)
}
