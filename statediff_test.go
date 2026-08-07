package main

import (
	"testing"
)

func TestProcess(t *testing.T) {
	prev := &Block{Number: "0x1", GasUsed: "0x5208", Timestamp: "0x0", Transactions: []Transaction{{Hash: "0xa"}}}
	cur := &Block{Number: "0x2", GasUsed: "0xa410", Timestamp: "0x4", Transactions: []Transaction{{Hash: "0xb", To: "0x1"}, {Hash: "0xc", To: ""}}}

	pb := process(prev, cur)
	if pb.GasDelta != 21000 {
		t.Fatalf("gas delta: got %d, want 21000", pb.GasDelta)
	}
	if pb.TimeDelta != 4 {
		t.Fatalf("time delta: got %d, want 4", pb.TimeDelta)
	}
	if len(pb.TxAdded) != 2 || len(pb.TxRemoved) != 1 {
		t.Fatalf("tx counts: added %d removed %d, want 2/1", len(pb.TxAdded), len(pb.TxRemoved))
	}
	if len(pb.NewContracts) != 1 {
		t.Fatalf("new contracts: got %d, want 1", len(pb.NewContracts))
	}
}
