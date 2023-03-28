package internal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/func25/gother"
)

func init() {
	gother.Dial("https://data-seed-prebsc-1-s3.binance.org:8545/")
}

func TestScan(t *testing.T) {
	ctx := context.Background()

	scan := gother.NewScanner(100, 21600030, common.HexToAddress("0xbA01E92eA9B940745f89785fC9cED4DDc17Da450"))
	logs, scannedBlock, err := scan.Scan(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	for _, l := range logs {
		fmt.Println(l.TxHash)
	}

	fmt.Printf("I have scanned to block %d\n", scannedBlock)
}

func TestCurrency(t *testing.T) {
	acc := gother.NewAccount("")

	ctx := context.Background()

	fmt.Println(acc.Balance(ctx))
}

type Agent struct {
	Block uint64
}

// GetBlock implements gother.IWorker
func (s *Agent) FromBlock(ctx context.Context) (uint64, error) {
	return s.Block + 1, nil
}

func (s *Agent) ProcessLogs(ctx context.Context, from, to uint64, logs []types.Log) error {
	for _, log := range logs {
		if log.Removed {
			return nil
		}

		fmt.Println(log.TxHash.Hex())
	}

	return nil
}

func (s *Agent) UpdateBlock(ctx context.Context, block uint64) error {
	s.Block = block
	return nil
}

func TestMe(t *testing.T) {
	agent := Agent{Block: 21600130}

	lazier := gother.Lazier[*Agent]{Agent: &agent, Duration: time.Second * 3}
	lazier.Scan(*gother.NewScanner(100, 0))

	for {
	}
}
