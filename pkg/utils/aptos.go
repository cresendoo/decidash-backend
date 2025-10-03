package utils

import (
	"encoding/hex"
)

func BytesToHex(bytes []byte) string {
	return "0x" + hex.EncodeToString(bytes)
}
