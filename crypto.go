package gother

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	SIG_PREFIX = "\x19Ethereum Signed Message:\n"
	SIG_V      = 27
)

// Keccak256Sign will hash data to 32 bytes (= keccak256) then signing it
func Keccak256Sign(prv string, data ...[]byte) (str string, err error) {
	dataHash := crypto.Keccak256Hash(data...)
	return Sign(prv, dataHash.Bytes())
}

// Sign signs the data with prefix `\x19Ethereum Signed Message:\n${len(data)}`
func Sign(prv string, data []byte) (str string, err error) {
	msg := fmt.Sprintf("%s%d%s", SIG_PREFIX, len(data), data)
	ethHash := crypto.Keccak256Hash([]byte(msg))

	privateKey, err := crypto.HexToECDSA(prv)
	if err != nil {
		return
	}

	signature, err := crypto.Sign(ethHash.Bytes(), privateKey)
	if err != nil {
		return
	}

	signature[len(signature)-1] += SIG_V

	return hexutil.Encode(signature), nil
}

func Uint(mul int, data []byte) []byte {
	return common.LeftPadBytes(data, mul/8)
}

// RecoverHexSig will recover public key from [65]byte signature
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

// RecoverHexSig will recover public key from hex signature
func RecoverHexSig(msg []byte, hexSig string) (*ecdsa.PublicKey, error) {
	sig, err := hexutil.Decode(hexSig)
	if err != nil {
		return nil, err
	}

	if len(sig) < 65 {
		return nil, fmt.Errorf("signature invalid, len(sig) < 65")
	}

	if sig[64] == 0 || sig[64] == 1 {
		sig[64] += SIG_V
	}

	return RecoverETHSig(msg, sig)
}
