package gother

import (
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var _addrRegex = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

var _signPrefix = "\x19Ethereum Signed Message:\n"

const (
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

func Uint(mul int, data []byte) []byte {
	return common.LeftPadBytes(data, mul/8)
}

func IsAddress(addr string) bool {
	return _addrRegex.MatchString(addr)
}
