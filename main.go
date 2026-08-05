package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

func hexToUint64(s string) uint64 {
	n, _ := strconv.ParseUint(strings.TrimPrefix(s, "0x"), 16, 64)
	return n
}

func main() {
	cur, _ := getBlockByNumber(context.Background(), "latest", true)
	n := hexToUint64(cur.Number)
	prev, _ := getBlockByNumber(context.Background(), fmt.Sprintf("0x%x", n-1), true)

	fmt.Printf("block delta:     %d\n", int64(n)-int64(hexToUint64(prev.Number)))
	fmt.Printf("time delta:      %ds\n", int64(hexToUint64(cur.Timestamp))-int64(hexToUint64(prev.Timestamp)))
	fmt.Printf("gas delta:       %d\n", int64(hexToUint64(cur.GasUsed))-int64(hexToUint64(prev.GasUsed)))
	fmt.Printf("stateRoot same:  %t\n", prev.StateRoot == cur.StateRoot)
}
