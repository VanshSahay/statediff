package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const rpcURL = "https://sensors-raises-edwards-task.trycloudflare.com"
const wsURL = "wss://received-reproduction-compatible-recreation.trycloudflare.com"

func (e *RPCError) Error() string {
	return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
}

func livews(ctx context.Context) {
	for {
		if err := runSubscription(ctx); err != nil {
			log.Println("subscription ended: ", err, "reconnecting in 3s")
			time.Sleep(3 * time.Second)
		}
	}
}

func runSubscription(ctx context.Context) error {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	subReq := RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_subscribe",
		Params:  []any{"newHeads"},
		ID:      1,
	}

	if err := conn.WriteJSON(subReq); err != nil {
		return err
	}
	var ack struct {
		ID    int       `json:"id"`
		Error *RPCError `json:"error"`
	}

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		if json.Unmarshal(raw, &ack) == nil && ack.ID == 1 && ack.Error == nil {
			break
		}
	}
	conn.SetReadDeadline(time.Time{})

	var prev *Block
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var msg SubscriptionMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}
		if msg.Params == nil {
			continue
		}

		cur, err := getBlockByNumber(ctx, msg.Params.Result.Number, true)
		if err != nil {
			log.Println("fetch block:", err)
			continue
		}

		bal, err := getBAL(ctx, msg.Params.Result.Number)
		if err != nil {
			log.Println("getBAL: ", err)
		} else {
			renderBAL(hexToUint64(cur.Number), bal)
		}

		if prev == nil {
			prev = cur
			log.Printf("watching from block %d\n", hexToUint64(cur.Number))
			continue
		}

		pb := process(prev, cur)
		printDiff(pb)
		prev = cur
	}

}

func printDiff(pb ProcessedBlock) {
	fmt.Printf("block %d | +%ds | gas %+d | txs +%d -%d | new contracts %d\n",
		hexToUint64(pb.Number),
		pb.TimeDelta,
		pb.GasDelta,
		len(pb.TxAdded), len(pb.TxRemoved),
		len(pb.NewContracts),
	)
}

func getBAL(ctx context.Context, tag string) (BlockAccessList, error) {
	reqBody, err := json.Marshal(RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockAccessList",
		Params:  []any{tag},
		ID:      1,
	})
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, rpcURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var r BALResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Error != nil {
		return nil, r.Error
	}

	return r.Result, nil
}

func getBlockByNumber(ctx context.Context, tag string, fullTx bool) (*Block, error) {
	reqBody, err := json.Marshal(RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []any{tag, fullTx},
		ID:      1,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, rpcURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	var httpClient = &http.Client{Timeout: 5 * time.Second}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Error != nil {
		return nil, r.Error
	}
	return &r.Result, nil
}
