package gother

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type scanner struct {
	*GotherClient
	BlockNum uint64

	stable    uint64
	from      uint64
	addresses []common.Address
}

func NewScanner(blocks uint64) *scanner {
	return &scanner{
		GotherClient: Client,
		BlockNum:     blocks,
	}
}

func (sc *scanner) From(from uint64) *scanner {
	sc.from = from

	return sc
}

func (sc *scanner) Addresses(addrs ...common.Address) *scanner {
	sc.addresses = append(sc.addresses, addrs...)

	return sc
}

func (sc *scanner) LatestStable(block uint64) {
	sc.stable = block
}

func (sc *scanner) ScanNext(ctx context.Context) (logs []types.Log, currentBlock uint64, err error) {
	latestBlock, err := sc.GotherClient.HeaderLatest(ctx)
	if err != nil {
		return
	}

	latestNum := latestBlock.Number.Uint64() - sc.stable

	to := sc.from + sc.BlockNum
	if to > latestNum {
		to = latestNum
	}

	logs, err = sc.Client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(sc.from),
		ToBlock:   new(big.Int).SetUint64(to),
		Addresses: sc.addresses,
	})

	if err != nil {
		return
	}

	return logs, to, err
}

// functask: retry
type LoopConfig struct {
	Duration  time.Duration                                         // duration gap for each call
	Emit      func(logs types.Log)                                  // called for each log collected
	AfterHook func(ogs []types.Log, currentBlock uint64, err error) // called after 1 round of scan
}

func (sc *scanner) ScanLogsLoop(ctx context.Context, cfg LoopConfig) chan struct{} {
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(cfg.Duration):
				logs, currentBlock, err := sc.ScanNext(ctx)

				if cfg.Emit != nil {
					for i := range logs {
						cfg.Emit(logs[i])
					}
				}

				sc.from = currentBlock + 1

				if cfg.AfterHook != nil {
					cfg.AfterHook(logs, currentBlock, err)
				}
			case <-stop:
				return
			}
		}
	}()

	return stop
}
