package gother

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

type IAgent interface {
	NextBlock(ctx context.Context) (uint64, error)
	UpdateBlock(ctx context.Context, block uint64) error
	ProcessLog(ctx context.Context, log types.Log) error
}

//Lazier is too lazy, he/she ignores errors when processing logs and updating block
type Lazier[T IAgent] struct {
	Agent    T
	Duration time.Duration
}

func (sol Lazier[T]) Scan(s Scanner) chan struct{} {
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-time.After(sol.Duration):
				ctx := context.Background()

				// get block
				if v, err := sol.Agent.NextBlock(ctx); err != nil {
					continue
				} else {
					s.From = v
				}

				// scan logs
				logs, currentBlock, err := s.Scan(ctx)
				if err != nil {
					continue
				}

				// process logs
				for i := range logs {
					_ = sol.Agent.ProcessLog(ctx, logs[i])
				}
				s.From = currentBlock + 1

				// update block
				_ = sol.Agent.UpdateBlock(ctx, currentBlock)
			case <-stop:
				return
			}
		}
	}()

	return stop
}
