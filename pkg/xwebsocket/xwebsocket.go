package xwebsocket

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	ctx    context.Context
	cancel context.CancelFunc

	conn *websocket.Conn
	url  string

	received chan []byte
	send     chan []byte
	ping     chan []byte

	mu        sync.RWMutex
	closeOnce sync.Once

	logger  *slog.Logger
	options WebsocketClientOptions
}

func New(ctx context.Context, url string, options ...WebsocketClientOption) (*WebsocketClient, error) {
	defaultOptions := &WebsocketClientOptions{
		Dialer:          &websocket.Dialer{HandshakeTimeout: 10 * time.Second},
		Headers:         http.Header{},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		WriteMode:       websocket.TextMessage,
	}
	for _, option := range options {
		option(defaultOptions)
	}

	dialer := defaultOptions.Dialer

	conn, _, err := dialer.Dial(url, defaultOptions.Headers)
	if err != nil {
		return nil, errorx.Wrap(err)
	}

	c := &WebsocketClient{
		conn:     conn,
		url:      url,
		options:  *defaultOptions,
		received: make(chan []byte, defaultOptions.ReadBufferSize),
		send:     make(chan []byte, defaultOptions.WriteBufferSize),
		ping:     make(chan []byte, 1),
		logger:   slog.With("name", "xwebsocket"),
	}
	c.ctx, c.cancel = context.WithCancel(ctx)

	go c.readPump()
	go c.writePump()

	if c.options.WithPing.Enabled {
		go c.pingPump()
	}
	return c, nil
}

func (c *WebsocketClient) Send(message []byte) error {
	var err error
	defer func() {
		if raw := recover(); raw != nil {
			if captured, ok := raw.(error); ok {
				err = errorx.WrapDepth(captured, 6)
			} else {
				err = errorx.Errorf("send on closed channel")
			}
		}
	}()
	c.send <- message
	return err
}

func (c *WebsocketClient) Read() ([]byte, bool) {
	msg, ok := <-c.received
	return msg, ok
}

func (c *WebsocketClient) Close() error {
	c.closeOnce.Do(func() {
		close(c.received)
		close(c.send)
		close(c.ping)
		c.mu.Lock()
		defer c.mu.Unlock()
		c.conn.Close()
		c.cancel()
		if c.options.OnClose != nil {
			c.options.OnClose()
		}
	})
	return nil
}

func (c *WebsocketClient) readPump() {
	defer func() {
		c.Close()
		c.logger.Info("readPump closed")
	}()
	c.logger.Info("readPump started")

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) ||
				errors.Is(err, syscall.ECONNRESET) ||
				strings.Contains(err.Error(), "unexpected EOF") ||
				strings.Contains(err.Error(), "going away") {
				c.logger.Warn("connection closed", "error", err)

				maxRetries := 10
				baseDelay := time.Second
				maxDelay := 30 * time.Second

				for i := 0; i < maxRetries; i++ {
					if i == maxRetries-1 {
						c.logger.Error("failed to reconnect after max retries", "error", errorx.Wrap(err).With("body", string(message)))
						return
					}

					select {
					case <-c.ctx.Done():
						return
					default:
					}

					// exponential backoff with jitter
					delay := baseDelay * time.Duration(1<<uint(i))
					if delay > maxDelay {
						delay = maxDelay
					}
					jitter := time.Duration(rand.Int63n(int64(delay / 4)))
					delay = delay + jitter

					dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}
					conn, _, err := dialer.Dial(c.url, c.options.Headers)
					if err == nil {
						c.mu.Lock()
						c.conn.Close()
						c.conn = conn
						c.mu.Unlock()
						c.logger.Info("reconnected successfully", "attempt", i+1)
						if c.options.OnReconnect != nil {
							c.options.OnReconnect()
						}
						break
					}
					c.logger.Warn("reconnection attempt failed", "error", errorx.Wrap(err).With("body", string(message)), "attempt", i+1, "next_retry", delay)
					time.Sleep(delay)
				}
				continue
			} else {
				if errors.Is(err, net.ErrClosed) {
					c.logger.Info("connection closed by client")
					return
				}
				c.logger.Error("failed to read message", "error", errorx.Wrap(err).With("body", string(message)))
				return
			}
		}

		select {
		case c.received <- message:
		case <-c.ctx.Done():
			select {
			case <-c.received:
			default:
				return
			}
		}
	}
}

func (c *WebsocketClient) writePump() {
	defer func() {
		c.logger.Info("writePump closed")
	}()
	c.logger.Info("writePump started")

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}
			c.mu.RLock()
			c.logger.Debug("writing message", "body", string(message))
			if err := c.conn.WriteMessage(c.options.WriteMode, message); err != nil {
				c.logger.Error("failed to write message", "error", errorx.Wrap(err).With("req_body", string(message)))
			}
			c.mu.RUnlock()
		default:
			select {
			case <-c.ctx.Done():
				return
			default:
				select {
				case pingMessage := <-c.ping:
					c.mu.RLock()
					if len(pingMessage) == 0 {
						if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
							c.logger.Error("failed to write ping message", "error", errorx.Wrap(err))
						}
					} else {
						if err := c.conn.WriteMessage(c.options.WriteMode, pingMessage); err != nil {
							c.logger.Error("failed to write ping message", "error", errorx.Wrap(err))
						}
					}
					c.mu.RUnlock()
				default:
				}
			}
		}
	}
}

func (c *WebsocketClient) pingPump() {
	ticker := time.NewTicker(c.options.WithPing.Interval)
	defer func() {
		ticker.Stop()
		c.logger.Info("pingPump closed")
	}()
	c.logger.Info("pingPump started")

	var pingMessage []byte
	if c.options.WithPing.Enabled {
		pingMessage = c.options.WithPing.CustomMessage
	}

	for {
		select {
		case <-ticker.C:
			c.ping <- pingMessage
		case <-c.ctx.Done():
			return
		}
	}
}
