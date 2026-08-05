package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const rpcURL = "https://ethereum-rpc.publicnode.com"
const wsURL = "wss://ethereum-rpc.publicnode.com"

func (e *RPCError) Error() string {
	return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
}

func livews(ctx context.Context) {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	subReq := RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_subscribe",
		Params:  []any{"newHeads"},
		ID:      1,
	}

	if err := conn.WriteJSON(subReq); err != nil {
		log.Fatal("subscribe:", err)
	}

	_, _, err = conn.ReadMessage()
	if err != nil {
		log.Fatal("sub ack:", err)
	}

	var prev *Block
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		var msg SubscriptionMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("unmarshal:", err)
			continue
		}

		cur, err := getBlockByNumber(ctx, msg.Params.Result.Number, true)
		if err != nil {
			log.Println("fetch block:", err)
			continue
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

	resp, err := http.DefaultClient.Do(httpReq)
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
