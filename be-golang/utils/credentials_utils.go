package utils

import (
	"encoding/hex"
	"fmt"
	"launchpad/config"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	sha3 "golang.org/x/crypto/sha3"
)

const MESSAGE_PREFIX = "\x19Ethereum Signed Message:\n"
const HEX_PREFIX = "0x"

func GetSign(hexString string) string {
	signature, err := SignHexString(hexString)
	if err != nil {
		return ""
	}
	return signature
}

func SignHexString(hexString string) (string, error) {
	bytes, err := hexStringToByteArray(hexString)
	if err != nil {
		return "", err
	}
	messageHash := sha3Hash(bytes)
	ethereumMessageHash := getEthereumMessageHash(messageHash)
	return doSign(ethereumMessageHash)
}

func hexStringToByteArray(input string) ([]byte, error) {
	// Remove the "0x" prefix if it exists
	cleanInput := CleanHexPrefix(input)

	length := len(cleanInput)

	if length == 0 {
		return []byte{}, nil
	}

	// Prepare byte array for the result
	var data []byte
	var startIdx int

	// If the length is odd, process the first character separately
	if length%2 != 0 {
		// Create byte array of size (length / 2 + 1)
		data = make([]byte, (length/2)+1)

		// Convert the first single hex character into a byte and store it
		firstByte, err := strconv.ParseUint(string(cleanInput[0]), 16, 8)
		if err != nil {
			return nil, err
		}
		data[0] = byte(firstByte)
		startIdx = 1
	} else {
		// Create byte array of size (length / 2)
		data = make([]byte, length/2)
		startIdx = 0
	}

	// Iterate through the rest of the string and process two hex characters at a time
	for i := startIdx; i < length; i += 2 {
		highNibble, err := strconv.ParseUint(string(cleanInput[i]), 16, 8)
		if err != nil {
			return nil, err
		}
		lowNibble, err := strconv.ParseUint(string(cleanInput[i+1]), 16, 8)
		if err != nil {
			return nil, err
		}
		data[(i+1)/2] = byte((highNibble << 4) + lowNibble)
	}

	return data, nil
}

func doSign(message []byte) (string, error) {
	if len(message) != 32 {
		return "", fmt.Errorf("message hash must be 32 bytes, got %d", len(message))
	}
	prv, err := crypto.HexToECDSA(config.AppConfig.Owner.PrivateKey)
	if err != nil {
		return "", err
	}
	signature, err := secp256k1.Sign(message, crypto.FromECDSA(prv))
	if err != nil {
		return "", err
	}
	// Add the recovery id to the signature
	signature[64] += 27
	return hex.EncodeToString(signature), nil
}

func getEthereumMessagePrefix(message []byte) []byte {
	return []byte(fmt.Sprintf("%s%d", MESSAGE_PREFIX, len(message)))
}

func getEthereumMessageHash(message []byte) []byte {
	prefix := getEthereumMessagePrefix(message)
	result := append(prefix, message...)
	return sha3Hash(result)
}

func sha3Hash(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}

func CleanHexPrefix(hexString string) string {
	if strings.HasPrefix(hexString, HEX_PREFIX) {
		return hexString[2:] // 返回去掉前缀后的字符串
	}
	return hexString // 如果没有前缀，则返回原始字符串
}

func ToHexStringNoPrefixZeroPadded(value int64, size int) string {
	return ToHexStringZeroPadded(value, size, false)
}

func ToHexStringZeroPadded(value int64, size int, withPrefix bool) string {
	result := strconv.FormatInt(value, 16)

	len := len(result) // 获取字符串长度
	if len > size {
		panic("value " + result + "is large then size " + strconv.Itoa(size))
	} else if value < 0 {
		panic("value " + result + "cannot be negative")
	}

	if len < size {
		result = strings.Repeat("0", size-len) + result
	}
	if withPrefix {
		return HEX_PREFIX + result
	} else {
		return result
	}
}

// ToHexStringNoPrefixZeroPaddedFromString 将字符串形式的大数字转换为十六进制字符串（无前缀，零填充）
func ToHexStringNoPrefixZeroPaddedFromString(valueStr string, size int) (string, error) {
	value := new(big.Int)
	_, ok := value.SetString(valueStr, 10)
	if !ok {
		return "", fmt.Errorf("invalid number string: %s", valueStr)
	}

	if value.Sign() < 0 {
		return "", fmt.Errorf("value cannot be negative: %s", valueStr)
	}

	result := fmt.Sprintf("%x", value)
	len := len(result)
	if len > size {
		return "", fmt.Errorf("value %s is larger than size %d", result, size)
	}

	if len < size {
		result = strings.Repeat("0", size-len) + result
	}

	return result, nil
}
