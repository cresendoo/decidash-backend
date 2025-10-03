package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

func MD5Hash(data string) string {
	hash := md5.Sum([]byte(data))
	result := hex.EncodeToString(hash[:])

	for len(result) < 32 {
		result = "0" + result
	}
	return result
}

func MD5HashWithReader(reader io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, reader); err != nil {
		return "", errorx.Wrap(err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func MD5HashWithFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", errorx.Wrap(err)
	}
	return MD5HashWithReader(file)
}
