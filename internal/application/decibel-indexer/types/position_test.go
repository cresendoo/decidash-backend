package types

import (
	"encoding/json"
	"os"
	"testing"
)

type transaction struct {
	Changes []struct {
		Data struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		} `json:"data"`
	} `json:"changes"`
}

func TestParsePerpPositionsFromTransaction(t *testing.T) {
	raw, err := os.ReadFile("testdata/tx_32667225.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	var tx transaction
	if err := json.Unmarshal(raw, &tx); err != nil {
		t.Fatalf("failed to unmarshal transaction: %v", err)
	}

	const module = "0xb8a5788314451ce4d2fbbad32e1bad88d4184b73943b7fe5166eab93cf1a5a95::perp_positions::"

	var crossed *CrossedPosition
	var isolated *IsolatedPosition

	for _, change := range tx.Changes {
		switch change.Data.Type {
		case module + "CrossedPosition":
			var cp CrossedPosition
			if err := json.Unmarshal(change.Data.Data, &cp); err != nil {
				t.Fatalf("failed to parse CrossedPosition: %v", err)
			}
			crossed = &cp
		case module + "IsolatedPosition":
			var ip IsolatedPosition
			if err := json.Unmarshal(change.Data.Data, &ip); err != nil {
				t.Fatalf("failed to parse IsolatedPosition: %v", err)
			}
			isolated = &ip
		}
	}

	if crossed == nil {
		t.Fatalf("CrossedPosition resource not found in fixture")
	}
	if isolated == nil {
		t.Fatalf("IsolatedPosition resource not found in fixture")
	}

	if len(crossed.Positions) == 0 {
		t.Fatalf("expected crossed positions to be populated")
	}

	first := crossed.Positions[0]
	if first.Size.String() != "1357042816" {
		t.Errorf("unexpected first position size: %s", first.Size.String())
	}
	if first.FundingIndexAtLastUpdate.Index.String() != "170141183460469231731684588685638315128" {
		t.Errorf("unexpected funding index: %s", first.FundingIndexAtLastUpdate.Index.String())
	}

	if isolated.Position.UserLeverage != 10 {
		t.Errorf("unexpected isolated leverage: %d", isolated.Position.UserLeverage)
	}
	if !isolated.Position.IsLong {
		t.Errorf("isolated position should be long")
	}

	if isolated.Position.Size.Uint64() != 0 {
		t.Errorf("expected isolated position size to be zero, got %d", isolated.Position.Size.Uint64())
	}
}
