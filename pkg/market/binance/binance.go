package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/cresendoo/decidash-backend/pkg/xwebsocket"
)

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	logger *slog.Logger
	
	conn   *xwebsocket.WebsocketClient
	subscribed map[string]bool
	requestID atomic.Int64
	orderbook OrderbookMap
}

func NewClient(rootCtx context.Context) (*Client, error) {
	ctx, cancel := context.WithCancel(rootCtx)	
	client := Client{
		ctx:    ctx,
		cancel: cancel,
		logger: slog.With("name", "market.binance"),
		subscribed: make(map[string]bool),
		orderbook: NewOrderbookMap(),
	}
	return &client, nil
}

func (c *Client) Start() error {
	if err := c.connect(); err != nil {
		return err
	}
	go c.readLoop()
	return nil
}

func (c *Client) Stop() error {
	c.cancel()
	c.wg.Wait()
	return nil
}

func (c *Client) Subscribe(param string) error {
	if c.subscribed[param] {
		return nil
	}
	c.subscribed[param] = true
	requestID := c.requestID.Add(1)
	return c.conn.Send(
		[]byte(fmt.Sprintf(`{"method":"SUBSCRIBE","params":["%s"],"id":%d}`, param, requestID)),
	)
}

func (c *Client) GetOrderbook(symbol string) (Orderbook, bool) {
	mu.RLock()
	orderbook, ok := c.orderbook[symbol]
	orderbook = orderbook.DeepCopy()
	mu.RUnlock()
	if !ok {
		return Orderbook{}, false
	}
	return orderbook, true
}

func (c *Client) connect() error {
	conn, err := xwebsocket.New(
		c.ctx,
		websocketURL,
		xwebsocket.WithOnReconnect(c.onReconnect),
	)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) readLoop() {
	c.wg.Add(1)
	defer c.wg.Done()
	defer c.logger.Debug("readLoop closed")

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			message, ok := c.conn.Read()
			if !ok {
				return
			}
			c.onMessage(message)
		}
	}
}

func (c *Client) onMessage(message []byte) {
	var response BookDepthResponse
	if err := json.Unmarshal(message, &response); err != nil {
		c.logger.Error("failed to unmarshal message", "error", err)
		return
	}
	c.orderbook.Update(response)
}

func (c *Client) onReconnect() {
	for param := range c.subscribed {
		if err := c.Subscribe(param); err != nil {
			c.logger.Error("failed to subscribe", "error", err)
		}
	}
	c.logger.Info("reconnected")
}