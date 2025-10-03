package client

import (
	"fmt"
	"io"
	"testing"
	"time"
)

const (
	// testServerAddr = "http://192.168.0.2:8080/v1"
	testServerAddr = "https://api.netna.staging.aptoslabs.com/v1"
)

func TestFullnodeRpcClient(t *testing.T) {
	client, err := NewFullnodeRpcClient(testServerAddr, "")
	if err != nil {
		t.Fatalf("failed to create fullnode rpc client: %v", err)
	}

	txs, _, err := client.getTransactions(nil, nil)
	if err != nil {
		t.Fatalf("failed to get fullnode rpc client: %v", err)
	}

	t.Log(len(txs))
}

func TestFullnodeRpcStream(t *testing.T) {
	client, err := NewFullnodeRpcClient(testServerAddr, "")
	if err != nil {
		t.Fatalf("failed to create fullnode rpc client: %v", err)
	}

	stream, err := client.NewStream(0, 100)
	if err != nil {
		t.Fatalf("failed to create fullnode rpc stream: %v", err)
	}
	defer stream.Close()

	go func() {
		time.Sleep(time.Second * 10)
		stream.Close()
	}()

	for {
		txs, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return
			}
			t.Fatalf("failed to read txs: %v", err)
			return
		}
		fmt.Println(len(txs))
	}
}
