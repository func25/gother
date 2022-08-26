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
	"github.com/ethereum/go-ethereum/core/types"
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

// NewSmcTx creates a transaction calling smart contract with suggest gas price and 110% gas limit
func (c account) NewSmcTx(ctx context.Context, tx SmcTxData) (*bind.TransactOpts, error) {
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
	txOps.GasLimit = gasLimit + gasLimit*10/100 // 110% gas limit
	txOps.GasPrice = gasPrice

	return txOps, nil
}

type TxData struct {
	To    string
	Value *big.Int
}

func (c account) NewTx(ctx context.Context, data TxData) (*types.Transaction, error) {
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
	to := common.HexToAddress(data.To)

	// nonce
	nonce, err := c.PendingNonceAt(ctx, from)
	if err != nil {
		return nil, err
	}

	// gas price
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	// gas limit
	gasLimit, err := c.EstimateGas(ctx, ethereum.CallMsg{
		From:     from,
		To:       &to,
		Value:    data.Value,
		GasPrice: gasPrice,
		Data:     nil,
	})
	if err != nil {
		return nil, err
	}

	// chainid
	// functask: cache this
	chainID, err := c.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nonce, to, data.Value, gasLimit, gasPrice, []byte{})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), priK)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// sign transaction and send
// return signed transaction
func (c account) SendTx(ctx context.Context, data TxData) (*types.Transaction, error) {
	tx, err := c.NewTx(ctx, data)
	if err != nil {
		return nil, err
	}

	if err := c.SendTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}
