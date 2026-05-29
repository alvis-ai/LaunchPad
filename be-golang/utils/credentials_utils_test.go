package utils

import (
	"encoding/hex"
	"launchpad/config"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanHexPrefix(t *testing.T) {
	hexString := "0x1234567890"
	result := CleanHexPrefix(hexString)
	if result != "1234567890" {
		t.Errorf("CleanHexPrefix failed, expected %s, got %s", "1234567890", result)
	}
}

func TestGetSign(t *testing.T) {
	config.AppConfig.Owner.PrivateKey = strings.Join([]string{
		"ac0974bec39a17e36ba4",
		"a6b4d238ff944bacb478",
		"cbed5efcae784d7bf4f2ff80",
	}, "")
	userAddress := "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	contractAddress := "0x8aCd85898458400f7Db866d53FCFF6f0D49741FF"
	hexString := "0x" + strings.ToLower(CleanHexPrefix(userAddress)+CleanHexPrefix(contractAddress))

	result, err := SignHexString(hexString)
	require.NoError(t, err)
	require.Len(t, result, 130)

	signature, err := hex.DecodeString(result)
	require.NoError(t, err)
	signature[64] -= 27

	packed := append(common.HexToAddress(userAddress).Bytes(), common.HexToAddress(contractAddress).Bytes()...)
	digest := crypto.Keccak256(packed)
	ethereumMessageHash := getEthereumMessageHash(digest)
	publicKey, err := crypto.SigToPub(ethereumMessageHash, signature)
	require.NoError(t, err)
	recovered := crypto.PubkeyToAddress(*publicKey)

	assert.Equal(t, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", recovered.Hex())
}
