package gother

import (
	"crypto/ecdsa"
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

func RecoverETHSig(msg, sig []byte) (*ecdsa.PublicKey, error) {
	if sig[64] != 27 && sig[64] != 28 {
		return nil, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[64] -= 27

	pubKey, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func RecoverHexSig(msg []byte, hexSig string) (*ecdsa.PublicKey, error) {
	sig, err := hexutil.Decode(hexSig)
	if err != nil {
		return nil, err
	}

	return RecoverETHSig(msg, sig)
}
