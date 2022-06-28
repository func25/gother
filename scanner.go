package gother

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Scanner struct {
	ScanNum   uint64           // how many blocks scanning each time
	offset    uint64           // offset from latest block (unstable blocks)
	From      uint64           // scan from?
	Addresses []common.Address // smart contract address

	client *GotherClient
}

func (s Scanner) Clone() *Scanner {
	return &s
}

// default offset is 1
func NewScanner(scanNum uint64, from uint64, addresses ...common.Address) *Scanner {
	return &Scanner{
		ScanNum:   scanNum,
		From:      from,
		Addresses: addresses,
		client:    Client,
		offset:    1,
	}
}

func (sc *Scanner) InjClient(client *GotherClient) *Scanner {
	sc.client = client

	return sc
}

func (sc *Scanner) InjOffset(offset uint64) *Scanner {
	sc.offset = offset

	return sc
}

func (sc *Scanner) Scan(ctx context.Context) (logs []types.Log, currentBlock uint64, err error) {
	latestBlock, err := sc.client.HeaderLatest(ctx)
	if err != nil {
		return
	}

	latestNum := latestBlock.Number.Uint64() - sc.offset

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
