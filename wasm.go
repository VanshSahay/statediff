//go:build js && wasm

package main

import (
	"context"
	"fmt"
	"syscall/js"
	"time"
)

func jsAppend(html string) {
	js.Global().Call("stateDiffAppend", html)
}

func jsStatus(text string) {
	js.Global().Call("stateDiffStatus", text)
}

func jsError(text string) {
	js.Global().Call("stateDiffError", text)
}

func watch(ctx context.Context, rpcURL string) {
	var prev *Block
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			head, err := getBlockNumber(ctx, rpcURL)
			if err != nil {
				jsError("head poll: " + err.Error())
				continue
			}
			cur, err := getBlockByNumber(ctx, rpcURL, head, true)
			if err != nil {
				jsError("fetch block: " + err.Error())
				continue
			}
			if prev != nil && cur.Number == prev.Number {
				continue
			}
			bal, err := getBAL(ctx, rpcURL, head)
			if err != nil {
				jsError("getBAL: " + err.Error())
			}

			if prev == nil {
				prev = cur
				jsStatus(fmt.Sprintf("watching from block %d", hexToUint64(cur.Number)))
				continue
			}

			pb := process(prev, cur)
			jsAppend(renderDiffHTML(pb))
			if bal != nil {
				jsAppend(renderBALHTML(hexToUint64(cur.Number), bal))
			}
			prev = cur
		}
	}
}

var cancelWatch context.CancelFunc

func start(rpcURL string) {
	stop()
	ctx, cancel := context.WithCancel(context.Background())
	cancelWatch = cancel
	go watch(ctx, rpcURL)
}

func stop() {
	if cancelWatch != nil {
		cancelWatch()
		cancelWatch = nil
	}
	jsStatus("stopped")
}

func main() {
	api := map[string]any{
		"start": js.FuncOf(func(this js.Value, args []js.Value) any {
			url := ""
			if len(args) > 0 {
				url = args[0].String()
			}
			start(url)
			return nil
		}),
		"stop": js.FuncOf(func(this js.Value, args []js.Value) any {
			stop()
			return nil
		}),
	}
	js.Global().Set("statediff", js.ValueOf(api))

	select {}
}
