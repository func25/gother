package gother

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

type IWorker interface {
	GetBlock(ctx context.Context) (uint64, error)
	UpdateBlock(ctx context.Context, block uint64) error
	ProcessLog(log types.Log) error
}
