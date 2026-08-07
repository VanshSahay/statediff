package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestVerifyBALRealDevnet(t *testing.T) {
	data, err := os.ReadFile("/tmp/bal-fixture.json")
	if err != nil {
		t.Skip("no live fixture captured")
	}
	want, err := os.ReadFile("/tmp/bal-hash.txt")
	if err != nil {
		t.Skip("no live hash captured")
	}
	var bal BlockAccessList
	if err := json.Unmarshal(data, &bal); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	ok, bytes := verifyBAL(bal, string(want))
	if !ok {
		t.Fatalf("BAL hash mismatch: header says %s but our RLP re-encoding hashes differently (encoded %d bytes)", want, bytes)
	}
	t.Logf("verified: %d accounts, %d RLP bytes, hash matches header commitment", len(bal), bytes)
}

func TestVerifyBALNegative(t *testing.T) {
	bal := BlockAccessList{{Address: "0x0000000000000000000000000000000000000001"}}
	if ok, _ := verifyBAL(bal, "0x"+string(make([]byte, 64))); ok {
		t.Fatal("expected mismatch for a bogus hash")
	}
}

func TestSlotLabelProbe(t *testing.T) {
	addr := "0x8943545177806ED17B9F23F0a21ee5948eCaa776"
	slot := balSlot0(addr)
	if slot == "" {
		t.Fatal("expected a probe slot for a valid address")
	}
	bal := BlockAccessList{{Address: "0x1111111111111111111111111111111111111111", StorageChanges: []SlotChanges{{Key: "0x" + slot}}}}
	txs := []Transaction{{From: addr}}
	labels := buildSlotLabels(bal, txs)
	if len(labels) != 1 {
		t.Fatalf("expected 1 label, got %d", len(labels))
	}
	if got := labels["0x1111111111111111111111111111111111111111:0x"+slot]; got == "" {
		t.Fatal("expected a balanceOf label on the matching slot")
	}
}
