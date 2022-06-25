package gother

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type scanner struct {
	GotherClient
	BlockNum uint64

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

func (sc *scanner) ScanNext(ctx context.Context) (logs []types.Log, currentBlock uint64, err error) {
	latest, err := sc.GotherClient.HeaderLatest(ctx)
	if err != nil {
		return
	}

	latestNum := latest.Number.Uint64()

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
