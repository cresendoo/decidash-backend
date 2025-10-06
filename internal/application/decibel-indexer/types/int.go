package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type Uint64 uint64

type Uint128 struct {
	value big.Int
}

type Uint256 struct {
	value big.Int
}

func (u *Uint64) UnmarshalJSON(data []byte) error {
	s, err := parseNumericString(data)
	if err != nil {
		return err
	}
	if s == "" {
		return errors.New("types: empty string for uint64")
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return fmt.Errorf("types: invalid uint64 %q: %w", s, err)
	}
	*u = Uint64(v)
	return nil
}

func (u Uint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u Uint64) Uint64() uint64 {
	return uint64(u)
}

func (u Uint64) String() string {
	return strconv.FormatUint(uint64(u), 10)
}

func (u Uint64) BigInt() *big.Int {
	return new(big.Int).SetUint64(uint64(u))
}

func (u Uint64) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u *Uint64) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		*u = 0
		return nil
	case int64:
		if v < 0 {
			return fmt.Errorf("types: negative value for uint64: %d", v)
		}
		*u = Uint64(v)
		return nil
	case uint64:
		*u = Uint64(v)
		return nil
	case []byte:
		return u.fromString(string(v))
	case string:
		return u.fromString(v)
	default:
		return fmt.Errorf("types: cannot scan Uint64 from %T", value)
	}
}

func (u *Uint64) fromString(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		*u = 0
		return nil
	}
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return fmt.Errorf("types: invalid uint64 %q: %w", s, err)
	}
	*u = Uint64(val)
	return nil
}

func (Uint64) GormDataType() string {
	return "string"
}

func (u *Uint128) UnmarshalJSON(data []byte) error {
	s, err := parseNumericString(data)
	if err != nil {
		return err
	}
	if s == "" {
		return errors.New("types: empty string for uint128")
	}
	bi, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return fmt.Errorf("types: invalid uint128 %q", s)
	}
	if bi.Sign() < 0 {
		return fmt.Errorf("types: negative uint128 %q", s)
	}
	if bi.BitLen() > 128 {
		return fmt.Errorf("types: uint128 overflow %q", s)
	}
	u.value.Set(bi)
	return nil
}

func (u Uint128) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u Uint128) String() string {
	return u.value.String()
}

func (u Uint128) BigInt() *big.Int {
	return new(big.Int).Set(&u.value)
}

func (u Uint128) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u *Uint128) SetBigInt(b *big.Int) error {
	if b.Sign() < 0 {
		return errors.New("types: negative value for uint128")
	}
	if b.BitLen() > 128 {
		return errors.New("types: uint128 overflow")
	}
	u.value.Set(b)
	return nil
}

func (u *Uint128) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		u.value.SetUint64(0)
		return nil
	case int64:
		if v < 0 {
			return fmt.Errorf("types: negative value for uint128: %d", v)
		}
		u.value.SetUint64(uint64(v))
		return nil
	case uint64:
		u.value.SetUint64(v)
		return nil
	case []byte:
		return u.setFromString(string(v))
	case string:
		return u.setFromString(v)
	default:
		return fmt.Errorf("types: cannot scan Uint128 from %T", value)
	}
}

func (u *Uint128) setFromString(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		u.value.SetUint64(0)
		return nil
	}
	bi, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return fmt.Errorf("types: invalid uint128 %q", s)
	}
	if bi.Sign() < 0 {
		return fmt.Errorf("types: negative uint128 %q", s)
	}
	if bi.BitLen() > 128 {
		return fmt.Errorf("types: uint128 overflow %q", s)
	}
	u.value.Set(bi)
	return nil
}

func (Uint128) GormDataType() string {
	return "string"
}

func (u *Uint256) UnmarshalJSON(data []byte) error {
	s, err := parseNumericString(data)
	if err != nil {
		return err
	}
	if s == "" {
		return errors.New("types: empty string for uint256")
	}
	bi, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return fmt.Errorf("types: invalid uint256 %q", s)
	}
	if bi.Sign() < 0 {
		return fmt.Errorf("types: negative uint256 %q", s)
	}
	if bi.BitLen() > 256 {
		return fmt.Errorf("types: uint256 overflow %q", s)
	}
	u.value.Set(bi)
	return nil
}

func (u Uint256) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u Uint256) String() string {
	return u.value.String()
}

func (u Uint256) BigInt() *big.Int {
	return new(big.Int).Set(&u.value)
}

func (u *Uint256) SetBigInt(b *big.Int) error {
	if b.Sign() < 0 {
		return errors.New("types: negative value for uint256")
	}
	if b.BitLen() > 256 {
		return errors.New("types: uint256 overflow")
	}
	u.value.Set(b)
	return nil
}

func (u Uint256) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u *Uint256) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		u.value.SetUint64(0)
		return nil
	case int64:
		if v < 0 {
			return fmt.Errorf("types: negative value for uint256: %d", v)
		}
		u.value.SetUint64(uint64(v))
		return nil
	case uint64:
		u.value.SetUint64(v)
		return nil
	case []byte:
		return u.setFromString(string(v))
	case string:
		return u.setFromString(v)
	default:
		return fmt.Errorf("types: cannot scan Uint256 from %T", value)
	}
}

func (u *Uint256) setFromString(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		u.value.SetUint64(0)
		return nil
	}
	bi, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return fmt.Errorf("types: invalid uint256 %q", s)
	}
	if bi.Sign() < 0 {
		return fmt.Errorf("types: negative uint256 %q", s)
	}
	if bi.BitLen() > 256 {
		return fmt.Errorf("types: uint256 overflow %q", s)
	}
	u.value.Set(bi)
	return nil
}

func (Uint256) GormDataType() string {
	return "string"
}

func parseNumericString(data []byte) (string, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return "", errors.New("types: empty json value for number")
	}
	if bytes.Equal(trimmed, []byte("null")) {
		return "", errors.New("types: null value for number")
	}
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return "", err
		}
		return s, nil
	}
	return string(trimmed), nil
}
