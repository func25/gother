package gother

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

type IAgent interface {
	FromBlock(ctx context.Context) (uint64, error)                            // get the next block which want to scan from
	ProcessLogs(ctx context.Context, from, to uint64, logs []types.Log) error // process the logs that agent collects
}

// Lazier is too lazy, he/she ignores errors when processing logs and updating block
type Lazier[T IAgent] struct {
	Agent    T
	Duration time.Duration
}

// Lazier will scan the blocks to get the logs with 3 steps with IAgent:
//   - FromBlock -> ProcessLog -> UpdateBlock
//
// Read IAgent interface for more information
func (sol Lazier[T]) Scan(s Scanner) chan struct{} {
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-time.After(sol.Duration):
				ctx := context.Background()

				// get block
				if v, err := sol.Agent.FromBlock(ctx); err != nil {
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
				sol.Agent.ProcessLogs(ctx, s.From, currentBlock, logs)
				s.From = currentBlock + 1
			case <-stop:
				return
			}
		}
	}()

	return stop
}
