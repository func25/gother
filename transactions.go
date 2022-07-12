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

type SmcTxData struct {
	Smc   Smc
	Value *big.Int
	Data  []byte
}

// functask: optimize this
// functask: add more option
func (c account) NewSmcTx(ctx context.Context, tx SmcTxData) (*bind.TransactOpts, error) {
	// functask: validate

	var err error

	pri := c.pri
	if len(pri) == 0 {
		pri = c.pri
	}

	priK, err := crypto.HexToECDSA(c.pri)
	if err != nil {
		return nil, err
	}

	pubK := priK.Public()
	pubKECDSA, ok := pubK.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not type of *ecdsa.PublicKey")
	}

	from := crypto.PubkeyToAddress(*pubKECDSA)
	to := common.HexToAddress(tx.Smc.Address)

	gasPrice, err := c.Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	gasLimit, err := c.Client.EstimateGas(ctx, ethereum.CallMsg{
		From:     from,
		To:       &to,
		Value:    tx.Value,
		GasPrice: gasPrice,
		Data:     tx.Data,
	})
	if err != nil {
		return nil, err
	}

	// functask: cache this
	chainID, err := c.Client.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	txOps, err := bind.NewKeyedTransactorWithChainID(priK, chainID)
	if err != nil {
		return nil, err
	}

	txOps.Value = tx.Value
	txOps.GasLimit = gasLimit + gasLimit*5/100 // 105% gas limit
	txOps.GasPrice = gasPrice

	return txOps, nil
}
