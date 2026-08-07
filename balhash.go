package main

import (
	"encoding/hex"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

func rlpBytes(b []byte) []byte {
	if len(b) == 1 && b[0] < 0x80 {
		return b
	}
	var out []byte
	if len(b) < 56 {
		out = append(out, 0x80+byte(len(b)))
	} else {
		l := uint64(len(b))
		var lenb []byte
		for l > 0 {
			lenb = append([]byte{byte(l & 0xff)}, lenb...)
			l >>= 8
		}
		out = append(out, 0xB7+byte(len(lenb)))
		out = append(out, lenb...)
	}
	return append(out, b...)
}

func rlpUint(v uint64) []byte {
	if v == 0 {
		return []byte{0x80}
	}
	var b []byte
	for v > 0 {
		b = append([]byte{byte(v & 0xff)}, b...)
		v >>= 8
	}
	return rlpBytes(b)
}

func rlpBig(n *big.Int) []byte {
	if n.Sign() == 0 {
		return []byte{0x80}
	}
	return rlpBytes(n.Bytes())
}

func rlpList(items ...[]byte) []byte {
	var payload []byte
	for _, it := range items {
		payload = append(payload, it...)
	}
	if len(payload) < 56 {
		return append([]byte{0xC0 + byte(len(payload))}, payload...)
	}
	l := uint64(len(payload))
	var lenb []byte
	for l > 0 {
		lenb = append([]byte{byte(l & 0xff)}, lenb...)
		l >>= 8
	}
	out := []byte{0xF7 + byte(len(lenb))}
	out = append(out, lenb...)
	return append(out, payload...)
}

func hexBig(s string) *big.Int {
	n := new(big.Int)
	n.SetString(strings.TrimPrefix(s, "0x"), 16)
	return n
}

func hexBytes(s string) []byte {
	b, _ := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	return b
}

func accountAccessRLP(a AccountAccess) []byte {
	var slots [][]byte
	for _, sc := range a.StorageChanges {
		var writes [][]byte
		for _, w := range sc.Changes {
			writes = append(writes, rlpList(rlpUint(hexToUint64(w.Index)), rlpBig(hexBig(w.Value))))
		}
		slots = append(slots, rlpList(rlpBig(hexBig(sc.Key)), rlpList(writes...)))
	}
	var reads [][]byte
	for _, r := range a.StorageReads {
		reads = append(reads, rlpBig(hexBig(r)))
	}
	var bals [][]byte
	for _, c := range a.BalanceChanges {
		bals = append(bals, rlpList(rlpUint(hexToUint64(c.Index)), rlpBig(hexBig(c.Value))))
	}
	var nonces [][]byte
	for _, c := range a.NonceChanges {
		nonces = append(nonces, rlpList(rlpUint(hexToUint64(c.Index)), rlpUint(hexToUint64(c.Value))))
	}
	var codes [][]byte
	for _, c := range a.CodeChanges {
		codes = append(codes, rlpList(rlpUint(hexToUint64(c.Index)), rlpBytes(hexBytes(c.Code))))
	}
	return rlpList(
		rlpBytes(hexBytes(a.Address)),
		rlpList(slots...),
		rlpList(reads...),
		rlpList(bals...),
		rlpList(nonces...),
		rlpList(codes...),
	)
}

func balRLP(list BlockAccessList) []byte {
	var accounts [][]byte
	for _, a := range list {
		accounts = append(accounts, accountAccessRLP(a))
	}
	return rlpList(accounts...)
}

func verifyBAL(list BlockAccessList, wantHash string) (bool, int) {
	enc := balRLP(list)
	h := sha3.NewLegacyKeccak256()
	h.Write(enc)
	got := hex.EncodeToString(h.Sum(nil))
	return got == strings.TrimPrefix(strings.TrimSpace(wantHash), "0x"), len(enc)
}

func buildSlotLabels(bal BlockAccessList, txs []Transaction) map[string]string {
	labels := make(map[string]string)
	probes := make(map[string]string)
	for _, tx := range txs {
		if tx.From != "" {
			probes[balSlot0(tx.From)] = "balanceOf(" + shortAddr(tx.From) + ")"
		}
		if tx.To != "" {
			probes[balSlot0(tx.To)] = "balanceOf(" + shortAddr(tx.To) + ")"
		}
	}
	if len(probes) == 0 {
		return labels
	}
	for _, acct := range bal {
		for _, sc := range acct.StorageChanges {
			if label, ok := probes[strings.TrimPrefix(sc.Key, "0x")]; ok {
				labels[acct.Address+":"+sc.Key] = label
			}
		}
	}
	return labels
}

func shortAddr(a string) string {
	if len(a) <= 12 {
		return a
	}
	return a[:8] + "…" + a[len(a)-4:]
}

func balSlot0(addr string) string {
	b := hexBytes(addr)
	if len(b) != 20 {
		return ""
	}
	var buf []byte
	buf = append(buf, make([]byte, 12)...)
	buf = append(buf, b...)
	buf = append(buf, make([]byte, 32)...)
	h := sha3.NewLegacyKeccak256()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}
