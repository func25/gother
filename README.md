# gother

gother is the library supporting interaction with smart contract, fetch contract's logs from blocks,...

## Content:

- [Installation](#installation)
- [Get started](#get-started)
  - [Dial to node](#dial-to-node)
  - [Worker workflow](#worker-workflow)
    - [Scanner: collect the logs](#scanner-scan-blocks-to-get-the-logs)
    - [Agent & Lazier: collect the logs, but in loop](#agent--lazier-scan-but-in-the-loop)

## Installation

1. You first need Go installed (version 1.18+ is required):

```sh
$ go get -u github.com/func25/gother
```

2. Import it in your code:

```go
import "github.com/func25/gother"
```

## Get started

### Dial to node

You need to dial the RPC node first:

```go
gother.DialCtx(context.Background(), "https://data-seed-prebsc-1-s3.binance.org:8545/")

// or gother.Dial("https://data-seed-prebsc-1-s3.binance.org:8545/")
```

### Worker workflow

In some cases, we need to crawl contract logs from the block; you can collect with the support of a scanner and lazier.

#### Scanner: scan blocks to get the logs

```go
// scan 100 blocks from block 2160030
scan := gother.NewScanner(100, 21600030)

logs, scannedBlock, err := scan.Scan(ctx)
if err != nil {
  return err
}

// print the transaction hash of logs
for _, l := range logs {
  fmt.Println(l.TxHash)
}

fmt.Printf("I have scanned to block %d\n", scannedBlock)
```

But how can you scan logs emitted from a specific smart contract?
```go
// scan 100 blocks from block 2160030 and get the logs 
// which emitted from contract: 0xbA01E92eA9B940745f89785fC9cED4DDc17Da450
scan := gother.NewScanner(100, 21600030, common.HexToAddress("0xbA01E92eA9B940745f89785fC9cED4DDc17Da450"))
```

#### Agent & Lazier: scan but in the loop

This lib supports you define an agent to do the loop: scan -> wait -> scan -> wait -> scan..., and the agent must meet the interface IAgent
```go
// the target interface 
type IAgent interface {
	FromBlock(ctx context.Context) (uint64, error)       // get the next block which want to scan from
	ProcessLog(ctx context.Context, log types.Log) error // process the logs that agent collects
	UpdateBlock(ctx context.Context, block uint64) error // update the scanned block after scanning
}

// create an agent
type Agent struct {
  Block uint64
}

// Scanner will scan from this block
func (s *Agent) FromBlock(ctx context.Context) (uint64, error) {
	return s.Block + 1, nil
}

// Process the crawled log
func (s *Agent) ProcessLog(ctx context.Context, log types.Log) error {
	if log.Removed {
		return nil
	}

	fmt.Println(log.TxHash.Hex())

	return nil
}

// Save the scanned block after scanning
func (s *Agent) UpdateBlock(ctx context.Context, block uint64) error {
	s.Block = block
	return nil
}
```

After creating the agent, wrap it with the lazier and do scanning
```go
// create an agent
agent := Agent{Block: 21600130}

// wrap the agent with lazier and wait duration of lazier is 3 seconds
lazier := gother.Lazier[*Agent]{Agent: &agent, Duration: time.Second * 3}

// scan 100 blocks each time, the `from` of scanner will be replaced with FromBlock(ctx) of agent
lazier.Scan(*gother.NewScanner(100, 0))

for {} // block the thread
```
