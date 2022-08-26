package gother

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type GotherClient struct {
	*ethclient.Client
}

var Client *GotherClient = &GotherClient{}

// DialCtx will connect to `url` and set the default client to client(url) if default client is nil
func DialCtx(ctx context.Context, url string) (*GotherClient, error) {
	var err error
	if Client.Client == nil {
		Client.Client, err = ethclient.DialContext(ctx, url)
	}

	return Client, err
}

// DialCtx will connect to `url` and set the default client to client(url) if default client is nil
func Dial(url string) (err error) {
	Client.Client, err = ethclient.Dial(url)

	return err
}

// ForceSetup is used in case you already have ethclient.Client and want to use features of gother
func ForceSetup(c *ethclient.Client) *GotherClient {
	Client = &GotherClient{
		Client: c,
	}

	return Client
}

func (c GotherClient) IsSmartContract(ctx context.Context, addr string) (bool, error) {
	address := common.HexToAddress(addr)
	bytecode, err := c.CodeAt(ctx, address, nil) // nil is latest block
	if err != nil {
		return false, err
	}

	return len(bytecode) > 0, nil
}

func (c GotherClient) Balance(ctx context.Context, address string) (*big.Int, error) {
	return c.Client.BalanceAt(ctx, common.HexToAddress(address), nil)
}
