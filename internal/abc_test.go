package internal

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/func25/gother"
)

func TestMe(t *testing.T) {
	err := gother.Dial("https://data-seed-prebsc-1-s1.binance.org:8545")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSign(t *testing.T) {
	var uniqueID uint64 = 71046
	var tokenID1 int64 = 86401635654417
	var tokenID2 int64 = 73691635656399
	var tokenID3 int64 = 16401635657790
	var rawAmount string = "12397129371289371"
	_amount, _ := new(big.Int).SetString(rawAmount, 10)

	var _expiredAt int64 = 1655992446
	_id := new(big.Int).SetUint64(uint64(uniqueID))
	_token0 := new(big.Int).SetInt64(tokenID1)
	_token1 := new(big.Int).SetInt64(tokenID2)
	_token2 := new(big.Int).SetInt64(tokenID3)

	id := common.LeftPadBytes(_id.Bytes(), 32)
	amount := common.LeftPadBytes(_amount.Bytes(), 32)
	token0 := common.LeftPadBytes(_token0.Bytes(), 32)
	token1 := common.LeftPadBytes(_token1.Bytes(), 32)
	token2 := common.LeftPadBytes(_token2.Bytes(), 32)
	expiredAt := common.LeftPadBytes(new(big.Int).SetInt64(_expiredAt).Bytes(), 32)
	sender := common.HexToAddress(address).Bytes()
	signature2, err := gother.Keccak256Sign(privateKey, sender, id, amount, token0, token1, token2, expiredAt)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(signature2)
}

var privateKey string = "32bef9ccee825377c9b0a213a5e42d45620b6287605cda880ba689b3a24daddc"
var publicKey string = "047fb650d77269b242125aad7a971b5178ae54fa4087ffedf862e3a7c683a9cd0179dfa17a6751dfe8d3d548106dc3c54735b81649521226b571c307d1c1e47af5"
var address string = "0x95EA3cC6b50a42a7b1A33bfc6F5a3a5A8bFc07f0"
