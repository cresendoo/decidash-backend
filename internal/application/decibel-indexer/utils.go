package decibelindexer

import (
	"bytes"
	"encoding/json"
)

func MapToStructJSON[T any](m map[string]any, out *T) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(out)
}
