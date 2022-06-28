package gother

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type scanner struct {
	ScanNum   uint64
	Offset    uint64
	From      uint64
	Addresses []common.Address

	client *GotherClient
}

func (s scanner) Clone() *scanner {
	return &s
}

func NewScanner(scanNum uint64) *scanner {
	return &scanner{
		client:  Client,
		ScanNum: scanNum,
	}
}

func (sc *scanner) InjClient(client *GotherClient) *scanner {
	sc.client = client

	return sc
}

func (sc *scanner) Scan(ctx context.Context) (logs []types.Log, currentBlock uint64, err error) {
	latestBlock, err := sc.client.HeaderLatest(ctx)
	if err != nil {
		return
	}

	latestNum := latestBlock.Number.Uint64() - sc.Offset

	to := sc.From + sc.ScanNum
	if to > latestNum {
		to = latestNum
	}

	if sc.From > to {
		return nil, sc.From - 1, nil
	}

	logs, err = sc.client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(sc.From),
		ToBlock:   new(big.Int).SetUint64(to),
		Addresses: sc.Addresses,
	})

	if err != nil {
		return
	}

	return logs, to, err
}
