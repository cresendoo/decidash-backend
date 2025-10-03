package binance

import (
	"encoding/json"
	"fmt"
	"time"
)

type SubscribeParams struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int64    `json:"id"`
}

func BookDepthStreamsParam(symbol string, interval string) string {
	return fmt.Sprintf("%s@depth@%s", symbol, interval) // ethusdc@depth@500ms
}

type Response[T any] struct {
	Stream string `json:"stream"`
	Data   T      `json:"data"`
}

type BookDepthResponse struct {
	Event string `json:"e"` // Event type
	EventTime time.Time `json:"E"` // Event time
	TransactionTime time.Time `json:"T"` // Transaction time
	Symbol string `json:"s"` // Symbol
	FirstUpdateID int64 `json:"U"` // First update ID in event
	FinalUpdateID int64 `json:"u"` // Final update ID in event
	FinalUpdateIDInLastStream int64 `json:"pu"` // Final update Id in last stream(ie `u` in last stream)
	Bids [][]string `json:"b"` // Bids to be updated, [price, quantity]
	Asks [][]string `json:"a"` // Asks to be updated, [price, quantity]
}

func (r *BookDepthResponse) UnmarshalJSON(data []byte) error {
	var response Response[RawBookDepthResponse]
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	r.Event = response.Data.Event
	r.EventTime = time.UnixMilli(response.Data.EventTime)
	r.TransactionTime = time.UnixMilli(response.Data.TransactionTime)
	r.Symbol = response.Data.Symbol
	r.FirstUpdateID = response.Data.FirstUpdateID
	r.FinalUpdateID = response.Data.FinalUpdateID
	r.FinalUpdateIDInLastStream = response.Data.FinalUpdateIDInLastStream
	r.Bids = response.Data.Bids
	r.Asks = response.Data.Asks
	return nil
}

type RawBookDepthResponse struct {
	Event string `json:"e"` // Event type
	EventTime int64 `json:"E"` // Event time
	TransactionTime int64 `json:"T"` // Transaction time
	Symbol string `json:"s"` // Symbol
	FirstUpdateID int64 `json:"U"` // First update ID in event
	FinalUpdateID int64 `json:"u"` // Final update ID in event
	FinalUpdateIDInLastStream int64 `json:"pu"` // Final update Id in last stream(ie `u` in last stream)
	Bids [][]string `json:"b"` // Bids to be updated, [price, quantity]
	Asks [][]string `json:"a"` // Asks to be updated, [price, quantity]
}