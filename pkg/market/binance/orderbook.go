package binance

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"sync"
)

type Orderbook struct {
	Symbol string
	IndexPrice float64
	Bids map[float64]float64
	Asks map[float64]float64
}
type OrderbookResponse struct {
	Symbol string `json:"s"`
	Bids [][]string `json:"b"`
	Asks [][]string `json:"a"`
}

func (o Orderbook) MarshalJSON() ([]byte, error) {
	bids := make([][]string, 0, len(o.Bids))
	for price, quantity := range o.Bids {
		bids = append(bids, []string{strconv.FormatFloat(price, 'f', -1, 64), strconv.FormatFloat(quantity, 'f', -1, 64)})
	}
	asks := make([][]string, 0, len(o.Asks))
	for price, quantity := range o.Asks {
		asks = append(asks, []string{strconv.FormatFloat(price, 'f', -1, 64), strconv.FormatFloat(quantity, 'f', -1, 64)})
	}
	sort.Slice(bids, func(i, j int) bool {
		return bids[i][0] < bids[j][0]
	})
	sort.Slice(asks, func(i, j int) bool {
		return asks[i][0] < asks[j][0]
	})
	return json.Marshal(OrderbookResponse{
		Symbol: o.Symbol,
		Bids: bids,
		Asks: asks,
	})
}

func (o Orderbook) DeepCopy() Orderbook {
	copied := Orderbook{
		Symbol:     o.Symbol,
		IndexPrice: o.IndexPrice,
		Bids:       make(map[float64]float64, len(o.Bids)),
		Asks:       make(map[float64]float64, len(o.Asks)),
	}
	for price, quantity := range o.Bids {
		copied.Bids[price] = quantity
	}
	for price, quantity := range o.Asks {
		copied.Asks[price] = quantity
	}
	return copied
}

var mu sync.RWMutex
type OrderbookMap map[string]Orderbook

func NewOrderbookMap() OrderbookMap {
	return make(OrderbookMap)
}

func (m OrderbookMap) Update(response BookDepthResponse) {
	mu.Lock()
	defer mu.Unlock()
	orderbook, ok := m[response.Symbol]
	if !ok {
		orderbook = Orderbook{
			Symbol: response.Symbol,
			Bids: make(map[float64]float64),
			Asks: make(map[float64]float64),
		}
	}
	for _, bid := range response.Bids {
		price, err := strconv.ParseFloat(bid[0], 64)
		if err != nil {
			continue
		}
		quantity, err := strconv.ParseFloat(bid[1], 64)
		if err != nil {
			continue
		}
		if quantity == 0 {
			delete(orderbook.Bids, price)
			continue
		}
		orderbook.Bids[price] = quantity
	}
	for _, ask := range response.Asks {
		price, err := strconv.ParseFloat(ask[0], 64)
		if err != nil {
			continue
		}
		quantity, err := strconv.ParseFloat(ask[1], 64)
		if err != nil {
			continue
		}
		if quantity == 0 {
			delete(orderbook.Asks, price)
			continue
		}
		orderbook.Asks[price] = quantity
	}
	orderbook.IndexPrice = (orderbook.MaxBid() + orderbook.MinAsk()) / 2
	m[response.Symbol] = orderbook
}

func (o *Orderbook) MinBid() float64 {
	minBid := math.MaxFloat64
	for price := range o.Bids {
		if price < minBid {
			minBid = price
		}
	}
	return minBid
}

func (o *Orderbook) MaxBid() float64 {
	maxBid := math.SmallestNonzeroFloat64
	for price := range o.Bids {
		if price > maxBid {
			maxBid = price
		}
	}
	return maxBid
}

func (o *Orderbook) MinAsk() float64 {
	minAsk := math.MaxFloat64
	for price := range o.Asks {
		if price < minAsk {
			minAsk = price
		}
	}
	return minAsk
}

func (o *Orderbook) MaxAsk() float64 {
	maxAsk := math.SmallestNonzeroFloat64
	for price := range o.Asks {
		if price > maxAsk {
			maxAsk = price
		}
	}
	return maxAsk
}