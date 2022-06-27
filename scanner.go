package gother

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Scanner struct {
	ScanNum uint64

	stable    uint64
	from      uint64
	addresses []common.Address
	client    *GotherClient
}

func NewScanner(scanNum uint64) *Scanner {
	return &Scanner{
		client:  Client,
		ScanNum: scanNum,
	}
}

func (sc *Scanner) From(from uint64) *Scanner {
	sc.from = from

	return sc
}

func (sc *Scanner) AddAddresses(addrs ...common.Address) *Scanner {
	sc.addresses = append(sc.addresses, addrs...)

	return sc
}

func (sc *Scanner) LatestStable(block uint64) {
	sc.stable = block
}

func (sc *Scanner) ScanNext(ctx context.Context) (logs []types.Log, currentBlock uint64, err error) {
	latestBlock, err := sc.client.HeaderLatest(ctx)
	if err != nil {
		return
	}

	latestNum := latestBlock.Number.Uint64() - sc.stable

	to := sc.from + sc.ScanNum
	if to > latestNum {
		to = latestNum
	}

	logs, err = sc.client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(sc.from),
		ToBlock:   new(big.Int).SetUint64(to),
		Addresses: sc.addresses,
	})

	if err != nil {
		return
	}

	return logs, to, err
}

//SoldierScan ignore error when process logs and update block
func (sc *Scanner) SoldierScan(s IWorker, dur time.Duration) chan struct{} {
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-time.After(dur):
				ctx := context.Background()

				// get block
				if v, err := s.GetBlock(ctx); err != nil {
					continue
				} else {
					sc.from = v + 1
				}

				// scan logs
				logs, currentBlock, err := sc.ScanNext(ctx)
				if err != nil {
					continue
				}

				// process logs
				for i := range logs {
					_ = s.ProcessLog(logs[i])
				}
				sc.from = currentBlock + 1

				// update block
				_ = s.UpdateBlock(ctx, currentBlock)
			case <-stop:
				return
			}
		}
	}()

	return stop
}
