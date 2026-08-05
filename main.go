package main

import (
	"context"
	"strconv"
	"strings"
)

func hexToUint64(s string) uint64 {
	n, _ := strconv.ParseUint(strings.TrimPrefix(s, "0x"), 16, 64)
	return n
}

func process(prev, cur *Block) ProcessedBlock {
	pb := ProcessedBlock{Block: *cur}
	pb.GasDelta = int64(hexToUint64(cur.GasUsed)) - int64(hexToUint64(prev.GasUsed))
	pb.TimeDelta = int64(hexToUint64(cur.Timestamp)) - int64(hexToUint64(prev.Timestamp))

	prevHashes := make(map[string]bool)
	for _, tx := range prev.Transactions {
		prevHashes[tx.Hash] = true
	}
	for _, tx := range cur.Transactions {
		if !prevHashes[tx.Hash] {
			pb.TxAdded = append(pb.TxAdded, tx)
		}
		if tx.To == "" {
			pb.NewContracts = append(pb.NewContracts, tx)
		}
	}

	curHashes := make(map[string]bool)
	for _, tx := range cur.Transactions {
		curHashes[tx.Hash] = true
	}
	for _, tx := range prev.Transactions {
		if !curHashes[tx.Hash] {
			pb.TxRemoved = append(pb.TxRemoved, tx)
		}
	}

	return pb
}

func main() {
	ctx := context.Background()
	livews(ctx)
}
