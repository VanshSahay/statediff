package main

import (
	"fmt"
	"math/big"
	"strings"
)

var systemContracts = map[string]string{
	"0x000f3df6d732807ef1319fb7b8bb8522d0beac02": "beacon-roots (EIP-4788)",
	"0x0000f90827f1c53a10cb7a02335b175320002935": "blockhash-history (EIP-2935)",
}

func short(hexStr string) string {
	if len(hexStr) <= 12 {
		return hexStr
	}
	return hexStr[:8] + "…" + hexStr[len(hexStr)-4:]
}

func weiToEth(hexWei string) string {
	wei := hexToBig(hexWei)
	f := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
	return f.Text('f', 6)
}

func renderBAL(blockNum uint64, bal BlockAccessList) {
	var systemCount, userCount int
	var b strings.Builder

	for _, acct := range bal {
		hasWrites := len(acct.StorageChanges) > 0 ||
			len(acct.BalanceChanges) > 0 ||
			len(acct.NonceChanges) > 0 ||
			len(acct.CodeChanges) > 0

		if !hasWrites {
			continue
		}

		if label, ok := systemContracts[acct.Address]; ok {
			systemCount++
			_ = label
			continue
		}
		userCount++

		fmt.Fprintf(&b, "  %s\n", short(acct.Address))

		for _, bc := range acct.BalanceChanges {
			fmt.Fprintf(&b, "    balance → %s ETH  (tx %d)\n",
				weiToEth(bc.Value), hexToUint64(bc.Index))
		}
		for _, nc := range acct.NonceChanges {
			fmt.Fprintf(&b, "    nonce → %d  (tx %d)\n",
				hexToUint64(nc.Value), hexToUint64(nc.Index))
		}
		for _, sc := range acct.StorageChanges {
			for _, chg := range sc.Changes {
				fmt.Fprintf(&b, "    slot %s = %s  (tx %d)\n",
					short(sc.Key), short(chg.Value), hexToUint64(chg.Index))
			}
		}
		if len(acct.CodeChanges) > 0 {
			fmt.Fprintf(&b, "    code deployed (%d changes)\n", len(acct.CodeChanges))
		}
	}

	fmt.Printf("── block %d state changes: %d user accounts, %d system contracts ──\n",
		blockNum, userCount, systemCount)
	if b.Len() > 0 {
		fmt.Print(b.String())
	}
}
