package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	block, err := getBlockByNumber(context.Background(), "latest", true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("block %s  txs=%d  stateRoot=%s\n",
		block.Number, len(block.Transactions), block.StateRoot)
}
