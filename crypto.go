package gother

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	_signPrefix    = "\x19Ethereum Signed Message:\n"
	_signV         = 27
	_signLastBytes = 64
)

func Keccak256Sign(prv string, data ...[]byte) (str string, err error) {
	dataHash := crypto.Keccak256Hash(data...)
	return Sign(prv, dataHash.Bytes())
}

func Sign(prv string, data []byte) (str string, err error) {
	msg := fmt.Sprintf("%s%d%s", _signPrefix, len(data), data)
	ethHash := crypto.Keccak256Hash([]byte(msg))

	privateKey, err := crypto.HexToECDSA(prv)
	if err != nil {
		return
	}

	signature, err := crypto.Sign(ethHash.Bytes(), privateKey)
	if err != nil {
		return
	}

	signature[_signLastBytes] += _signV

	return hexutil.Encode(signature), nil
}

func Uint(mul int, data []byte) []byte {
	return common.LeftPadBytes(data, mul/8)
}
