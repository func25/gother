package gother

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	SIG_PREFIX = "\x19Ethereum Signed Message:\n"
	SIG_V      = 27
)

type Wallet struct {
	Private      string
	PrivateECDSA *ecdsa.PrivateKey
	Public       string
	PublicECDSA  *ecdsa.PublicKey
	Address      string
}

// Keccak256Sign will hash data to 32 bytes (= keccak256) then signing it,
// return hexa-signature which is 132 bytes (0x...)
func Keccak256Sign(prv string, data ...[]byte) (signature string, err error) {
	dataHash := crypto.Keccak256Hash(data...)
	return Sign(prv, dataHash.Bytes())
}

// Keccak256SignBytes will hash data to 32 bytes (= keccak256) then signing it
// return bytes-signature which is 65 bytes (32 bytes r, 32 bytes s, 1 byte v)
func Keccak256SignBytes(prv string, data ...[]byte) (signature []byte, err error) {
	dataHash := crypto.Keccak256Hash(data...)

	sigHex, err := Sign(prv, dataHash.Bytes())
	if err != nil {
		return nil, err
	}

	return hexutil.Decode(sigHex[2:]) // remove "0x"
}

func Keccak256Hash(data ...[]byte) (hash string) {
	return crypto.Keccak256Hash(data...).Hex()
}

// JsonABI parses json abi to abi.ABI, panic if error
func JsonABI(reader io.Reader) abi.ABI {
	abi, err := abi.JSON(reader)
	if err != nil {
		panic(err)
	}

	return abi
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
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}

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

	return RecoverETHSig(msg, sig)
}

func NewWallet() (*Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	hexPrivate := hexutil.Encode(crypto.FromECDSA(privateKey))[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
	}
	hexPublic := hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA)[4:])

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return &Wallet{
		Private:      hexPrivate,
		PrivateECDSA: privateKey,
		Public:       hexPublic,
		PublicECDSA:  publicKeyECDSA,
		Address:      address,
	}, nil
}
