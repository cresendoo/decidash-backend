package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"io"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

type AESGCM struct {
	secret []byte
	key    []byte
}

func NewAESGCM(secret string, salt string) (*AESGCM, error) {
	s, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, errorx.Wrap(err)
	}
	k := pbkdf2.Key(s, []byte(salt), 2048, 64, sha1.New)
	return &AESGCM{s, k}, nil
}

func (e *AESGCM) deriveKey(id int64) (aesKey, nonce []byte) {
	aesKey = e.key[:32]

	buf := new(bytes.Buffer)
	buf.WriteString("DECIDASH!@#$%^&*()")
	_ = binary.Write(buf, binary.LittleEndian, id)
	kdf := hkdf.New(sha256.New, e.secret, e.key[32:], buf.Bytes())
	nonce = make([]byte, 12)
	_, _ = io.ReadFull(kdf, nonce)
	return
}

func (e *AESGCM) Encrypt(id int64, plainText []byte) (cipherText string, err error) {
	aesKey, nonce := e.deriveKey(id)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		err = errorx.Wrap(err)
		return
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		err = errorx.Wrap(err)
		return
	}
	enc := aesgcm.Seal(nil, nonce, plainText, nil)
	cipherText = base64.StdEncoding.EncodeToString(enc)
	return
}

func (e *AESGCM) Decrypt(id int64, cipherText string) (msg []byte, err error) {
	encrypted, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		err = errorx.Wrap(err)
		return
	}

	aesKey, nonce := e.deriveKey(id)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		err = errorx.Wrap(err)
		return
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		err = errorx.Wrap(err)
		return
	}

	msg, err = aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		err = errorx.Wrap(err)
		return
	}

	return
}
