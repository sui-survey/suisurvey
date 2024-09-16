package utils

import (
	"encoding/base64"
	"math/big"
)

func NumToBlobId(num *big.Int) string {
	extractedBytes := make([]byte, 32)

	for i := uint(0); i < 32; i++ {
		extractedBytes[i] = byte(num.Uint64() & 0xff)
		num.Rsh(num, 8)
	}
	encoded := base64.RawURLEncoding.EncodeToString(extractedBytes)
	return encoded
}
