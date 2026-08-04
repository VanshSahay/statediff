package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	body := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`
	resp, err := http.Post("https://ethereum-rpc.publicnode.com", "application/json", strings.NewReader(body))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	out, _ := io.ReadAll(resp.Body)
	fmt.Println(string(out))
}
