package xwebsocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketClientOptions struct {
	Headers         http.Header
	Dialer          *websocket.Dialer
	ReadBufferSize  int
	WriteBufferSize int
	WriteMode       int
	WithPing        struct {
		Enabled       bool
		Interval      time.Duration
		CustomMessage []byte
	}
	OnClose     func()
	OnReconnect func()
}

type WebsocketClientOption func(*WebsocketClientOptions)

func WithHeaders(headers http.Header) WebsocketClientOption {
	return func(o *WebsocketClientOptions) {
		o.Headers = headers
	}
}

func WithDialer(dialer *websocket.Dialer) WebsocketClientOption {
	return func(o *WebsocketClientOptions) {
		o.Dialer = dialer
	}
}

func WithReadBufferSize(readBufferSize int) WebsocketClientOption {
	return func(o *WebsocketClientOptions) {
		o.ReadBufferSize = readBufferSize
	}
}

func WithWriteBufferSize(writeBufferSize int) WebsocketClientOption {
	return func(o *WebsocketClientOptions) {
		o.WriteBufferSize = writeBufferSize
	}
}

func WithPing(pingInterval time.Duration, customMessage []byte) WebsocketClientOption {
	return func(o *WebsocketClientOptions) {
		o.WithPing.Enabled = true
		o.WithPing.Interval = pingInterval
		o.WithPing.CustomMessage = customMessage
	}
}

func WithOnClose(onClose func()) WebsocketClientOption {
	return func(o *WebsocketClientOptions) {
		o.OnClose = onClose
	}
}

func WithOnReconnect(onReconnect func()) WebsocketClientOption {
	return func(o *WebsocketClientOptions) {
		o.OnReconnect = onReconnect
	}
}
