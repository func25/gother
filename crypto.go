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

func SignRaw(prv string, data ...[]byte) (str string, err error) {
	dataHash := crypto.Keccak256Hash(data...)
	msg := fmt.Sprintf("%s%d%s", _signPrefix, len(dataHash), dataHash.Bytes())
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

func RecoverECDSA(data []byte, signature []byte) ([]byte, error) {
	content := crypto.Keccak256Hash(data)
	public, err := crypto.Ecrecover(content.Bytes(), signature)
	if err != nil {
		return []byte{}, err
	}
	return public, nil
}

func Uint(mul int, data []byte) []byte {
	return common.LeftPadBytes(data, mul/8)
}
