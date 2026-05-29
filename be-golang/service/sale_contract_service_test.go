package service

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateEventSignature(t *testing.T) {
	saleContractService := NewSaleContractService(nil)

	signature := saleContractService.CalculateEventSignature("SaleDeployed(address)")

	assert.Equal(t, "0x65c0ac3f6aa97317ad1e9f6c73af709aad47dc12a97239e1b08a43a73195f7e0", signature.Hex())
}

func TestBytesToAddressFrom20Bytes(t *testing.T) {
	input := common.Hex2Bytes("0000000000000000000000008acd85898458400f7db866d53fcff6f0d49741ff")

	address, err := bytesToAddressFrom20Bytes(input)

	require.NoError(t, err)
	assert.Equal(t, "0x8aCd85898458400f7Db866d53FCFF6f0D49741FF", address.Hex())
}
