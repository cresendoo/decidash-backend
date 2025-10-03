package binance

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestBookDepthStreamsParam(t *testing.T) {
	param := BookDepthStreamsParam("ethusdc", "500ms")
	if param != "ethusdc@depth@500ms" {
		t.Errorf("expected ethusdc@depth@500ms, got %s", param)
	}
}

func TestSubscribe(t *testing.T) {
	binance, err := NewClient(context.Background())
	if err != nil {
		t.Fatalf("failed to create binance client: %v", err)
	}
	if err := binance.Start(); err != nil {
		t.Fatalf("failed to start binance client: %v", err)
	}
	if err := binance.Subscribe(BookDepthStreamsParam("ethusdc", "500ms")); err != nil {
		t.Fatalf("failed to subscribe: %v", err)
	}
	time.Sleep(10 * time.Second)
	if err := binance.Stop(); err != nil {
		t.Fatalf("failed to stop binance client: %v", err)
	}
	orderbook, ok := binance.GetOrderbook("ETHUSDC")
	if !ok {
		t.Fatalf("failed to get orderbook")
	}
	slog.Info("orderbook", "orderbook", orderbook)
	slog.Info("max bid", "max bid", orderbook.MaxBid())
	slog.Info("min ask", "min ask", orderbook.MinAsk())
}