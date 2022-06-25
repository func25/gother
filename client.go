package gother

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type GotherClient struct {
	*ethclient.Client
	prv string
}

var Client *GotherClient

func DialCtx(ctx context.Context, url string) (err error) {
	Client.Client, err = ethclient.DialContext(ctx, url)

	return err
}

func Dial(url string) (err error) {
	Client.Client, err = ethclient.Dial(url)

	return err
}

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

func (c GotherClient) SignRaw(data ...[]byte) (str string, err error) {
	return SignRaw(c.prv, data...)
}

func (c *GotherClient) InjectPrivate(prv string) *GotherClient {
	c.prv = prv
	return c
}
