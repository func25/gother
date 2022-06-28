package gother

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Smc struct {
	ABI     abi.ABI
	Address string
}

// functask: optimize this
// functask: add more option
func (c GotherClient) NewSmcTx(ctx context.Context, smc Smc, value *big.Int, data []byte) (*bind.TransactOpts, error) {
	var err error

	priK, err := crypto.HexToECDSA(c.prv)
	if err != nil {
		return nil, err
	}

	pubK := priK.Public()
	pubKECDSA, ok := pubK.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not type of *ecdsa.PublicKey")
	}

	from := crypto.PubkeyToAddress(*pubKECDSA)
	to := common.HexToAddress(smc.Address)

	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	gasLimit, err := c.EstimateGas(ctx, ethereum.CallMsg{
		From:     from,
		To:       &to,
		Value:    value,
		GasPrice: gasPrice,
		Data:     data,
	})
	if err != nil {
		return nil, err
	}

	// functask: cache this
	chainID, err := c.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	txOps, err := bind.NewKeyedTransactorWithChainID(priK, chainID)
	if err != nil {
		return nil, err
	}

	txOps.Value = value
	txOps.GasLimit = gasLimit
	txOps.GasPrice = gasPrice

	return txOps, nil
}
