package types

import (
	"encoding/json"
	"math/big"
	"testing"
)

func TestUint64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Uint64
		wantErr bool
	}{
		{
			name:    "valid string number",
			input:   `"123"`,
			want:    Uint64(123),
			wantErr: false,
		},
		{
			name:    "valid number",
			input:   `456`,
			want:    Uint64(456),
			wantErr: false,
		},
		{
			name:    "zero",
			input:   `"0"`,
			want:    Uint64(0),
			wantErr: false,
		},
		{
			name:    "max uint64",
			input:   `"18446744073709551615"`,
			want:    Uint64(18446744073709551615),
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   `""`,
			want:    Uint64(0),
			wantErr: true,
		},
		{
			name:    "null",
			input:   `null`,
			want:    Uint64(0),
			wantErr: true,
		},
		{
			name:    "invalid string",
			input:   `"abc"`,
			want:    Uint64(0),
			wantErr: true,
		},
		{
			name:    "negative number",
			input:   `"-1"`,
			want:    Uint64(0),
			wantErr: true,
		},
		{
			name:    "overflow",
			input:   `"18446744073709551616"`,
			want:    Uint64(0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u Uint64
			err := u.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && u != tt.want {
				t.Errorf("Uint64.UnmarshalJSON() = %v, want %v", u, tt.want)
			}
		})
	}
}

func TestUint64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		u    Uint64
		want string
	}{
		{
			name: "zero",
			u:    Uint64(0),
			want: `"0"`,
		},
		{
			name: "positive number",
			u:    Uint64(123),
			want: `"123"`,
		},
		{
			name: "max uint64",
			u:    Uint64(18446744073709551615),
			want: `"18446744073709551615"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.MarshalJSON()
			if err != nil {
				t.Errorf("Uint64.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Uint64.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestUint64_Uint64(t *testing.T) {
	tests := []struct {
		name string
		u    Uint64
		want uint64
	}{
		{
			name: "zero",
			u:    Uint64(0),
			want: 0,
		},
		{
			name: "positive number",
			u:    Uint64(123),
			want: 123,
		},
		{
			name: "max uint64",
			u:    Uint64(18446744073709551615),
			want: 18446744073709551615,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Uint64(); got != tt.want {
				t.Errorf("Uint64.Uint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_String(t *testing.T) {
	tests := []struct {
		name string
		u    Uint64
		want string
	}{
		{
			name: "zero",
			u:    Uint64(0),
			want: "0",
		},
		{
			name: "positive number",
			u:    Uint64(123),
			want: "123",
		},
		{
			name: "max uint64",
			u:    Uint64(18446744073709551615),
			want: "18446744073709551615",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); got != tt.want {
				t.Errorf("Uint64.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_BigInt(t *testing.T) {
	tests := []struct {
		name string
		u    Uint64
		want *big.Int
	}{
		{
			name: "zero",
			u:    Uint64(0),
			want: big.NewInt(0),
		},
		{
			name: "positive number",
			u:    Uint64(123),
			want: big.NewInt(123),
		},
		{
			name: "max uint64",
			u:    Uint64(18446744073709551615),
			want: new(big.Int).SetUint64(18446744073709551615),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.u.BigInt()
			if got.Cmp(tt.want) != 0 {
				t.Errorf("Uint64.BigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint128_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid string number",
			input:   `"123"`,
			want:    "123",
			wantErr: false,
		},
		{
			name:    "valid number",
			input:   `456`,
			want:    "456",
			wantErr: false,
		},
		{
			name:    "zero",
			input:   `"0"`,
			want:    "0",
			wantErr: false,
		},
		{
			name:    "max uint128",
			input:   `"340282366920938463463374607431768211455"`,
			want:    "340282366920938463463374607431768211455",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   `""`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "null",
			input:   `null`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid string",
			input:   `"abc"`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "negative number",
			input:   `"-1"`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "overflow",
			input:   `"340282366920938463463374607431768211456"`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u Uint128
			err := u.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint128.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && u.String() != tt.want {
				t.Errorf("Uint128.UnmarshalJSON() = %v, want %v", u.String(), tt.want)
			}
		})
	}
}

func TestUint128_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		u    Uint128
		want string
	}{
		{
			name: "zero",
			u:    Uint128{value: *big.NewInt(0)},
			want: `"0"`,
		},
		{
			name: "positive number",
			u:    Uint128{value: *big.NewInt(123)},
			want: `"123"`,
		},
		{
			name: "max uint128",
			u: func() Uint128 {
				bi, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10)
				return Uint128{value: *bi}
			}(),
			want: `"340282366920938463463374607431768211455"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.MarshalJSON()
			if err != nil {
				t.Errorf("Uint128.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Uint128.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestUint128_String(t *testing.T) {
	tests := []struct {
		name string
		u    Uint128
		want string
	}{
		{
			name: "zero",
			u:    Uint128{value: *big.NewInt(0)},
			want: "0",
		},
		{
			name: "positive number",
			u:    Uint128{value: *big.NewInt(123)},
			want: "123",
		},
		{
			name: "max uint128",
			u: func() Uint128 {
				bi, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10)
				return Uint128{value: *bi}
			}(),
			want: "340282366920938463463374607431768211455",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); got != tt.want {
				t.Errorf("Uint128.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint128_BigInt(t *testing.T) {
	tests := []struct {
		name string
		u    Uint128
		want *big.Int
	}{
		{
			name: "zero",
			u:    Uint128{value: *big.NewInt(0)},
			want: big.NewInt(0),
		},
		{
			name: "positive number",
			u:    Uint128{value: *big.NewInt(123)},
			want: big.NewInt(123),
		},
		{
			name: "max uint128",
			u: func() Uint128 {
				bi, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10)
				return Uint128{value: *bi}
			}(),
			want: func() *big.Int {
				bi, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10)
				return bi
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.u.BigInt()
			if got.Cmp(tt.want) != 0 {
				t.Errorf("Uint128.BigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint128_SetBigInt(t *testing.T) {
	tests := []struct {
		name    string
		input   *big.Int
		wantErr bool
	}{
		{
			name:    "valid positive number",
			input:   big.NewInt(123),
			wantErr: false,
		},
		{
			name:    "zero",
			input:   big.NewInt(0),
			wantErr: false,
		},
		{
			name: "max uint128",
			input: func() *big.Int {
				bi, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10)
				return bi
			}(),
			wantErr: false,
		},
		{
			name:    "negative number",
			input:   big.NewInt(-1),
			wantErr: true,
		},
		{
			name: "overflow",
			input: func() *big.Int {
				bi, _ := new(big.Int).SetString("340282366920938463463374607431768211456", 10)
				return bi
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u Uint128
			err := u.SetBigInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint128.SetBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && u.BigInt().Cmp(tt.input) != 0 {
				t.Errorf("Uint128.SetBigInt() = %v, want %v", u.BigInt(), tt.input)
			}
		})
	}
}

func TestUint256_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid string number",
			input:   `"123"`,
			want:    "123",
			wantErr: false,
		},
		{
			name:    "valid number",
			input:   `456`,
			want:    "456",
			wantErr: false,
		},
		{
			name:    "zero",
			input:   `"0"`,
			want:    "0",
			wantErr: false,
		},
		{
			name:    "max uint256",
			input:   `"115792089237316195423570985008687907853269984665640564039457584007913129639935"`,
			want:    "115792089237316195423570985008687907853269984665640564039457584007913129639935",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   `""`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "null",
			input:   `null`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid string",
			input:   `"abc"`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "negative number",
			input:   `"-1"`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "overflow",
			input:   `"115792089237316195423570985008687907853269984665640564039457584007913129639936"`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u Uint256
			err := u.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint256.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && u.String() != tt.want {
				t.Errorf("Uint256.UnmarshalJSON() = %v, want %v", u.String(), tt.want)
			}
		})
	}
}

func TestUint256_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		u    Uint256
		want string
	}{
		{
			name: "zero",
			u:    Uint256{value: *big.NewInt(0)},
			want: `"0"`,
		},
		{
			name: "positive number",
			u:    Uint256{value: *big.NewInt(123)},
			want: `"123"`,
		},
		{
			name: "max uint256",
			u: func() Uint256 {
				bi, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
				return Uint256{value: *bi}
			}(),
			want: `"115792089237316195423570985008687907853269984665640564039457584007913129639935"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.MarshalJSON()
			if err != nil {
				t.Errorf("Uint256.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Uint256.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestUint256_String(t *testing.T) {
	tests := []struct {
		name string
		u    Uint256
		want string
	}{
		{
			name: "zero",
			u:    Uint256{value: *big.NewInt(0)},
			want: "0",
		},
		{
			name: "positive number",
			u:    Uint256{value: *big.NewInt(123)},
			want: "123",
		},
		{
			name: "max uint256",
			u: func() Uint256 {
				bi, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
				return Uint256{value: *bi}
			}(),
			want: "115792089237316195423570985008687907853269984665640564039457584007913129639935",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); got != tt.want {
				t.Errorf("Uint256.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint256_BigInt(t *testing.T) {
	tests := []struct {
		name string
		u    Uint256
		want *big.Int
	}{
		{
			name: "zero",
			u:    Uint256{value: *big.NewInt(0)},
			want: big.NewInt(0),
		},
		{
			name: "positive number",
			u:    Uint256{value: *big.NewInt(123)},
			want: big.NewInt(123),
		},
		{
			name: "max uint256",
			u: func() Uint256 {
				bi, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
				return Uint256{value: *bi}
			}(),
			want: func() *big.Int {
				bi, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
				return bi
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.u.BigInt()
			if got.Cmp(tt.want) != 0 {
				t.Errorf("Uint256.BigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint256_SetBigInt(t *testing.T) {
	tests := []struct {
		name    string
		input   *big.Int
		wantErr bool
	}{
		{
			name:    "valid positive number",
			input:   big.NewInt(123),
			wantErr: false,
		},
		{
			name:    "zero",
			input:   big.NewInt(0),
			wantErr: false,
		},
		{
			name: "max uint256",
			input: func() *big.Int {
				bi, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
				return bi
			}(),
			wantErr: false,
		},
		{
			name:    "negative number",
			input:   big.NewInt(-1),
			wantErr: true,
		},
		{
			name: "overflow",
			input: func() *big.Int {
				bi, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639936", 10)
				return bi
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u Uint256
			err := u.SetBigInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint256.SetBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && u.BigInt().Cmp(tt.input) != 0 {
				t.Errorf("Uint256.SetBigInt() = %v, want %v", u.BigInt(), tt.input)
			}
		})
	}
}

func TestParseNumericString(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr bool
	}{
		{
			name:    "quoted string",
			input:   []byte(`"123"`),
			want:    "123",
			wantErr: false,
		},
		{
			name:    "unquoted number",
			input:   []byte(`456`),
			want:    "456",
			wantErr: false,
		},
		{
			name:    "zero",
			input:   []byte(`"0"`),
			want:    "0",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   []byte(`""`),
			want:    "",
			wantErr: false,
		},
		{
			name:    "null",
			input:   []byte(`null`),
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty bytes",
			input:   []byte(``),
			want:    "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   []byte(`   `),
			want:    "",
			wantErr: true,
		},
		{
			name:    "whitespace around number",
			input:   []byte(`  123  `),
			want:    "123",
			wantErr: false,
		},
		{
			name:    "whitespace around quoted string",
			input:   []byte(`  "123"  `),
			want:    "123",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNumericString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNumericString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseNumericString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// JSON round-trip tests
func TestUint64_JSONRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		u    Uint64
	}{
		{
			name: "zero",
			u:    Uint64(0),
		},
		{
			name: "positive number",
			u:    Uint64(123),
		},
		{
			name: "max uint64",
			u:    Uint64(18446744073709551615),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.u)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			// Unmarshal back
			var u Uint64
			err = json.Unmarshal(jsonData, &u)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Check equality
			if u != tt.u {
				t.Errorf("JSON round-trip failed: got %v, want %v", u, tt.u)
			}
		})
	}
}

func TestUint128_JSONRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		u    Uint128
	}{
		{
			name: "zero",
			u:    Uint128{value: *big.NewInt(0)},
		},
		{
			name: "positive number",
			u:    Uint128{value: *big.NewInt(123)},
		},
		{
			name: "max uint128",
			u: func() Uint128 {
				bi, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10)
				return Uint128{value: *bi}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.u)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			// Unmarshal back
			var u Uint128
			err = json.Unmarshal(jsonData, &u)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Check equality
			if u.String() != tt.u.String() {
				t.Errorf("JSON round-trip failed: got %v, want %v", u.String(), tt.u.String())
			}
		})
	}
}

func TestUint256_JSONRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		u    Uint256
	}{
		{
			name: "zero",
			u:    Uint256{value: *big.NewInt(0)},
		},
		{
			name: "positive number",
			u:    Uint256{value: *big.NewInt(123)},
		},
		{
			name: "max uint256",
			u: func() Uint256 {
				bi, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
				return Uint256{value: *bi}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.u)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			// Unmarshal back
			var u Uint256
			err = json.Unmarshal(jsonData, &u)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Check equality
			if u.String() != tt.u.String() {
				t.Errorf("JSON round-trip failed: got %v, want %v", u.String(), tt.u.String())
			}
		})
	}
}
