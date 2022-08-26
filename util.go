package gother

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

type gotherUtil struct{}

var Util gotherUtil

// IsValidAddress validate hex address
func (gotherUtil) IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// IsZeroAddress validate if it's a 0 address
func (gotherUtil) IsZeroAddress(iaddress interface{}) bool {
	var address common.Address
	switch v := iaddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}

// ToDecimal wei to decimals
func (gotherUtil) ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ToWei decimals to wei
func (gotherUtil) ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

// CalcGasCost calculate gas cost given gas limit (units) and gas price (wei)
func (gotherUtil) CalcGasCost(gasLimit uint64, gasPrice *big.Int) *big.Int {
	gasLimitBig := big.NewInt(int64(gasLimit))
	return gasLimitBig.Mul(gasLimitBig, gasPrice)
}

// SigRSV signatures R S V returned as arrays
func (gotherUtil) SigRSV(isig interface{}) ([32]byte, [32]byte, uint8) {
	var sig []byte
	switch v := isig.(type) {
	case []byte:
		sig = v
	case string:
		sig, _ = hexutil.Decode(v)
	}

	sigstr := common.Bytes2Hex(sig)
	rS := sigstr[0:64]
	sS := sigstr[64:128]
	R := [32]byte{}
	S := [32]byte{}
	copy(R[:], common.FromHex(rS))
	copy(S[:], common.FromHex(sS))
	vStr := sigstr[128:130]
	vI, _ := strconv.Atoi(vStr)
	V := uint8(vI + 27)

	return R, S, V
}

func (gotherUtil) ToEthSignedMessageHash(data []byte) []byte {
	return crypto.Keccak256([]byte(fmt.Sprintf("%s%d%s", SIG_PREFIX, len(data), data)))
}

// func (gotherUtil) HexToECDSAPrv(priv string) (*ecdsa.PrivateKey, error) {
// 	pk := new(ecdsa.PrivateKey)
// 	pk.D, _ = new(big.Int).SetString(priv, 16)
// 	pk.PublicKey.Curve = elliptic.P256()
// 	pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())
// 	return pk, nil
// }

func (gotherUtil) GetPublicHex(priv *ecdsa.PrivateKey) (string, error) {
	publicKey := priv.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("cannot parse public key to *edsa.PublicKey")
	}
	return hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA))[4:], nil
}

func (gotherUtil) GetPublic(priv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	publicKey := priv.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot parse public key to *edsa.PublicKey")
	}
	return publicKeyECDSA, nil
}

func (gotherUtil) GetAddressHex(pub *ecdsa.PublicKey) string {
	return crypto.PubkeyToAddress(*pub).Hex()
}
