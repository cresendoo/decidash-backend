package types

import "testing"

func TestUint64_ValueAndScan(t *testing.T) {
	original := Uint64(42)
	val, err := original.Value()
	if err != nil {
		t.Fatalf("Uint64.Value() error = %v", err)
	}
	if val != "42" {
		t.Fatalf("Uint64.Value() = %v, want %v", val, "42")
	}

	var parsed Uint64
	if err := parsed.Scan("42"); err != nil {
		t.Fatalf("Uint64.Scan(string) error = %v", err)
	}
	if parsed != 42 {
		t.Fatalf("Uint64.Scan(string) = %v, want %v", parsed, 42)
	}

	if err := parsed.Scan([]byte("17")); err != nil {
		t.Fatalf("Uint64.Scan([]byte) error = %v", err)
	}
	if parsed != 17 {
		t.Fatalf("Uint64.Scan([]byte) = %v, want %v", parsed, 17)
	}

	if err := parsed.Scan(int64(3)); err != nil {
		t.Fatalf("Uint64.Scan(int64) error = %v", err)
	}
	if parsed != 3 {
		t.Fatalf("Uint64.Scan(int64) = %v, want %v", parsed, 3)
	}

	if err := parsed.Scan(nil); err != nil {
		t.Fatalf("Uint64.Scan(nil) error = %v", err)
	}
	if parsed != 0 {
		t.Fatalf("Uint64.Scan(nil) = %v, want %v", parsed, 0)
	}

	if err := parsed.Scan(int64(-1)); err == nil {
		t.Fatalf("Uint64.Scan(negative) expected error")
	}

	var zero Uint64
	if got := zero.GormDataType(); got != "string" {
		t.Fatalf("Uint64.GormDataType() = %s, want string", got)
	}
}

func TestUint128_ValueAndScan(t *testing.T) {
	const sample = "340282366920938463463374607431768211455"
	var u Uint128
	if err := u.Scan(sample); err != nil {
		t.Fatalf("Uint128.Scan(string) error = %v", err)
	}
	if u.String() != sample {
		t.Fatalf("Uint128.Scan(string) = %s, want %s", u.String(), sample)
	}

	if val, err := u.Value(); err != nil || val != sample {
		if err != nil {
			t.Fatalf("Uint128.Value() error = %v", err)
		}
		t.Fatalf("Uint128.Value() = %v, want %s", val, sample)
	}

	if err := u.Scan(int64(12)); err != nil {
		t.Fatalf("Uint128.Scan(int64) error = %v", err)
	}
	if u.String() != "12" {
		t.Fatalf("Uint128.Scan(int64) = %s, want 12", u.String())
	}

	if err := u.Scan(nil); err != nil {
		t.Fatalf("Uint128.Scan(nil) error = %v", err)
	}
	if u.String() != "0" {
		t.Fatalf("Uint128.Scan(nil) = %s, want 0", u.String())
	}

	if err := u.Scan("-1"); err == nil {
		t.Fatalf("Uint128.Scan(negative) expected error")
	}

	var zero Uint128
	if got := zero.GormDataType(); got != "string" {
		t.Fatalf("Uint128.GormDataType() = %s, want string", got)
	}
}

func TestUint256_ValueAndScan(t *testing.T) {
	const sample = "115792089237316195423570985008687907853269984665640564039457584007913129639935"
	var u Uint256
	if err := u.Scan(sample); err != nil {
		t.Fatalf("Uint256.Scan(string) error = %v", err)
	}
	if u.String() != sample {
		t.Fatalf("Uint256.Scan(string) = %s, want %s", u.String(), sample)
	}

	if val, err := u.Value(); err != nil || val != sample {
		if err != nil {
			t.Fatalf("Uint256.Value() error = %v", err)
		}
		t.Fatalf("Uint256.Value() = %v, want %s", val, sample)
	}

	if err := u.Scan(nil); err != nil {
		t.Fatalf("Uint256.Scan(nil) error = %v", err)
	}
	if u.String() != "0" {
		t.Fatalf("Uint256.Scan(nil) = %s, want 0", u.String())
	}

	if err := u.Scan("-1"); err == nil {
		t.Fatalf("Uint256.Scan(negative) expected error")
	}

	var zero Uint256
	if got := zero.GormDataType(); got != "string" {
		t.Fatalf("Uint256.GormDataType() = %s, want string", got)
	}
}
