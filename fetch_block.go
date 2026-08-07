package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (e *RPCError) Error() string {
	return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
}

func getBlockNumber(ctx context.Context, rpcURL string) (string, error) {
	reqBody, err := json.Marshal(RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []any{},
		ID:      1,
	})
	if err != nil {
		return "", err
	}
	resp, err := doRPC(ctx, rpcURL, reqBody)
	if err != nil {
		return "", err
	}
	var r struct {
		Result string    `json:"result"`
		Error  *RPCError `json:"error"`
	}
	if err := json.Unmarshal(resp, &r); err != nil {
		return "", err
	}
	if r.Error != nil {
		return "", r.Error
	}
	return r.Result, nil
}

func getBAL(ctx context.Context, rpcURL, tag string) (BlockAccessList, error) {
	reqBody, err := json.Marshal(RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockAccessList",
		Params:  []any{tag},
		ID:      1,
	})
	if err != nil {
		return nil, err
	}
	resp, err := doRPC(ctx, rpcURL, reqBody)
	if err != nil {
		return nil, err
	}
	var r BALResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}
	if r.Error != nil {
		return nil, r.Error
	}
	return r.Result, nil
}

func getBlockByNumber(ctx context.Context, rpcURL, tag string, fullTx bool) (*Block, error) {
	reqBody, err := json.Marshal(RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []any{tag, fullTx},
		ID:      1,
	})
	if err != nil {
		return nil, err
	}
	resp, err := doRPC(ctx, rpcURL, reqBody)
	if err != nil {
		return nil, err
	}
	var r RPCResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}
	if r.Error != nil {
		return nil, r.Error
	}
	return &r.Result, nil
}

func doRPC(ctx context.Context, rpcURL string, reqBody []byte) ([]byte, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimSpace(rpcURL), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	var httpClient = &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
