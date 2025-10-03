package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

var _processUID string

func init() {
	var buf [12]byte
	_, _ = rand.Read(buf[:])
	b64 := base64.StdEncoding.EncodeToString(buf[:])
	b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)

	_processUID = b64[0:10]
}

func ProcessUID() string {
	return _processUID
}
