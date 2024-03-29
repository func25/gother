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
	offset    uint64           // offset from latest block (unstable blocks), default: 1
	From      uint64           // scan from?
	Addresses []common.Address // smart contract address

	client *GotherClient
}

func (s Scanner) Clone() *Scanner {
	return &s
}

// NewScanner create a scanner which will scan `scanNum` blocks each time, it scans from `from`
// and only scans log from smart contracts which has address in `addresses` array
func NewScanner(scanNum uint64, from uint64, addresses ...common.Address) *Scanner {
	return &Scanner{
		ScanNum:   scanNum,
		From:      from,
		Addresses: addresses,
		client:    Client,
		offset:    1,
	}
}

// InjClient set the client for scanner
func (sc *Scanner) InjClient(client *GotherClient) *Scanner {
	sc.client = client

	return sc
}

// InjOffset sets the offset from latest block, this will scan all the blocks except last `offset` blocks
// because we need to scan the most stable blocks, latest blocks may be forked or replaced...
func (sc *Scanner) InjOffset(offset uint64) *Scanner {
	sc.offset = offset

	return sc
}

// InjAddresses changes the addresses of smart contracts that will be scanned
func (sc *Scanner) InjAddresses(addresses ...common.Address) *Scanner {
	sc.Addresses = addresses

	return sc
}

// Scan scans all the logs from `scanner.from` to `scanner.from + scanner.ScanNum` and retrieves logs to process
// return the current block after scanning (currentBlock = min(scanner.from + scanner.ScanNum, latestBlockNum - offset))
func (sc *Scanner) Scan(ctx context.Context) (logs []types.Log, scannedBlock uint64, err error) {
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
