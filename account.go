package gother

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	*GotherClient
	pri string
}

func NewAccount(privateKey string) *Account {
	return &Account{
		pri:          privateKey,
		GotherClient: Client,
	}
}

func (a *Account) SetClient(client *GotherClient) *Account {
	a.GotherClient = client
	return a
}

func (a *Account) SetPrivate(privateKey string) *Account {
	a.pri = privateKey
	return a
}

func (c Account) Keccak256Sign(data ...[]byte) (str string, err error) {
	return Keccak256Sign(c.pri, data...)
}

func (c Account) Sign(data []byte) (str string, err error) {
	return Sign(c.pri, data)
}

func (c Account) PublicKey() (string, error) {
	priv, err := crypto.HexToECDSA(c.pri)
	if err != nil {
		return "", err
	}

	return Util.GetPublicHex(priv)
}

func (c Account) Address() (string, error) {
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

func (c Account) Balance(ctx context.Context) (*big.Int, error) {
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
