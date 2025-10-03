package encrypt

import (
	"database/sql/driver"
	"fmt"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
)

func NewColumnCipher(name, col string) Cipher {
	c, err := NewAESGCM(secret, fmt.Sprintf("%s.%s", name, col))
	if err != nil {
		panic(err)
	}
	return c
}

type EncryptID interface {
	EncryptID() int64
}

type Column struct {
	ctx    EncryptID
	cipher Cipher
	c      interface{}
	v      interface{}
}

func NewColumn(ctx EncryptID, c interface{}, cipher Cipher) *Column {
	return &Column{ctx, cipher, c, nil}
}

func (e *Column) Assign() (err error) {
	var dec []byte
	if e.v != nil {
		var encString string
		switch v := e.v.(type) {
		case string:
			encString = v
		case []byte:
			encString = string(v)
		default:
			return errorx.New("Invalid value type for decrypt").With("type", e.v)
		}

		if len(encString) > 0 {
			dec, err = e.cipher.Decrypt(e.ctx.EncryptID(), encString)
			if err != nil {
				return errorx.Wrap(err)
			}
		}
	} else {
		return errorx.Wrap(convertAssign(e.c, e.v))
	}

	return errorx.Wrap(convertAssign(e.c, dec))
}

func (e *Column) Scan(value interface{}) (err error) {
	e.v = value
	return nil
}

func (e Column) Value() (driver.Value, error) {
	var data []byte
	var val interface{}
	var err error

	if valuer, ok := e.c.(driver.Valuer); ok {
		val, err = valuer.Value()
		if err != nil {
			return nil, errorx.Wrap(err)
		}
	} else {
		val = e.c
	}

	switch v := val.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	case nil:
		return nil, nil
	default:
		return nil, errorx.New("Invalid value type for encrypt").With("v", v).With("c", e.c)
	}

	if len(data) > 0 {
		return e.cipher.Encrypt(e.ctx.EncryptID(), data)
	}

	return data, nil
}

// TruncateColumnValue 암호화 컬럼의 경우 암호화 된 글자의 길이를 원문에서 알 수 없기 때문에
// column size 에서 자른다.
func TruncateColumnValue(e Cipher, id int64, plainText string, columnSize int) (string, error) {
	runes := []rune(plainText)
	for i := 0; i < len(runes); i++ {
		text := string(runes[:len(runes)-i])
		cipherText, err := e.Encrypt(id, []byte(text))
		if err != nil {
			return "", errorx.Wrap(err).
				With("msg_id", id).
				With("plain_text", plainText).
				With("column_size", columnSize)
		} else if len(cipherText) <= columnSize {
			return text, nil
		}
	}
	return "", nil
}
