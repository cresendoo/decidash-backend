package client

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/aptos-labs/aptos-go-sdk/api"
)

type fullnodeRpcClient struct {
	httpClient *http.Client

	baseUrl *url.URL
	apiKey  string
}

type FullnodeRpcStream struct {
	ctx    context.Context
	cancel context.CancelFunc
	txChan chan []*api.UserTransaction
	wg     sync.WaitGroup
}

func (s *FullnodeRpcStream) Recv() ([]*api.UserTransaction, error) {
	txs, ok := <-s.txChan
	if !ok {
		return nil, io.EOF
	}
	return txs, nil
}

func (s *FullnodeRpcStream) Close() {
	s.cancel()
	s.wg.Wait()
}

// Using FullNode RPC Client.
// @param serverAddr: http://localhost:8080/v1
func NewFullnodeRpcClient(serverAddr string, apiKey string) (*fullnodeRpcClient, error) {
	baseUrl, err := url.Parse(serverAddr)
	if err != nil {
		return nil, err
	}
	return &fullnodeRpcClient{
		apiKey:     apiKey,
		baseUrl:    baseUrl,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (c *fullnodeRpcClient) NewStream(
	startVersion uint64,
	limit uint64,
) (*FullnodeRpcStream, error) {
	ctx, cancel := context.WithCancel(context.Background())
	stream := &FullnodeRpcStream{
		ctx:    ctx,
		cancel: cancel,
		txChan: make(chan []*api.UserTransaction, 128),
	}

	limit = uint64(math.Min(float64(limit), 25))

	delay := time.Millisecond * 25

	stream.wg.Add(1)
	go func() {
		defer stream.wg.Done()
		defer close(stream.txChan)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				txs, count, err := c.getTransactions(&startVersion, &limit)
				if err != nil {
					slog.Error("failed to get transactions", "error", err)
					return
				}
				if count == 0 || count < int(limit) {
					time.Sleep(delay)
					continue
				}

				lastTx := txs[len(txs)-1]
				startVersion = lastTx.Version + 1
				select {
				case <-ctx.Done():
					return
				case stream.txChan <- txs:
				}
			}
		}
	}()

	return stream, nil
}

func (c *fullnodeRpcClient) getTransactions(
	startVersion *uint64,
	limit *uint64,
) ([]*api.UserTransaction, int, error) {
	requestURI := c.baseUrl.JoinPath("/transactions")
	params := url.Values{}
	if startVersion != nil {
		params.Set("start", strconv.FormatUint(*startVersion, 10))
	}
	if limit != nil {
		params.Set("limit", strconv.FormatUint(*limit, 10))
	}

	req, err := http.NewRequest("GET", requestURI.String(), nil)
	if err != nil {
		return nil, 0, err
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var txs []*api.CommittedTransaction
	if err := json.Unmarshal(body, &txs); err != nil {
		return nil, 0, err
	}

	var userTxs []*api.UserTransaction
	for _, tx := range txs {
		if tx.Type == api.TransactionVariantUser {
			userTx, err := tx.UserTransaction()
			if err != nil {
				continue
			}
			userTxs = append(userTxs, userTx)
		}
	}
	return userTxs, len(txs), nil
}
