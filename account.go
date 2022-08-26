package gother

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type account struct {
	*GotherClient
	pri string
}

func NewAccount(privateKey string) *account {
	return &account{
		pri:          privateKey,
		GotherClient: Client,
	}
}

func (a *account) SetClient(client *GotherClient) *account {
	a.GotherClient = client
	return a
}

func (a *account) SetPrivate(privateKey string) *account {
	a.pri = privateKey
	return a
}

func (c account) Keccak256Sign(data ...[]byte) (str string, err error) {
	return Keccak256Sign(c.pri, data...)
}

func (c account) Sign(data []byte) (str string, err error) {
	return Sign(c.pri, data)
}

func (c account) PublicKey() (string, error) {
	priv, err := crypto.HexToECDSA(c.pri)
	if err != nil {
		return "", err
	}

	return Util.GetPublicHex(priv)
}

func (c account) Address() (string, error) {
	priv, err := crypto.HexToECDSA(c.pri)
	if err != nil {
		return "", err
	}

	pub, err := Util.GetPublic(priv)
	if err != nil {
		return "", err
	}

	return Util.GetAddressHex(pub), nil
}

func (c account) Balance(ctx context.Context) (*big.Int, error) {
	address, err := c.Address()
	if err != nil {
		return nil, err
	}

	balance, err := c.Client.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
