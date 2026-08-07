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

func renderDiffHTML(pb ProcessedBlock) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<div class="block">block %d <span class="dim">|</span> +%ds <span class="dim">|</span> gas %+d <span class="dim">|</span> txs <span class="pos">+%d</span> <span class="neg">-%d</span> <span class="dim">|</span> new contracts %d</div>`,
		hexToUint64(pb.Number),
		pb.TimeDelta,
		pb.GasDelta,
		len(pb.TxAdded), len(pb.TxRemoved),
		len(pb.NewContracts),
	)
	return b.String()
}

func renderBALHTML(blockNum uint64, bal BlockAccessList) string {
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

		if _, ok := systemContracts[acct.Address]; ok {
			systemCount++
			continue
		}
		userCount++

		fmt.Fprintf(&b, `<div class="acct">%s</div>`, short(acct.Address))

		for _, bc := range acct.BalanceChanges {
			fmt.Fprintf(&b, `<div class="chg">balance → <span class="val">%s ETH</span> <span class="dim">(tx %d)</span></div>`,
				weiToEth(bc.Value), hexToUint64(bc.Index))
		}
		for _, nc := range acct.NonceChanges {
			fmt.Fprintf(&b, `<div class="chg">nonce → <span class="val">%d</span> <span class="dim">(tx %d)</span></div>`,
				hexToUint64(nc.Value), hexToUint64(nc.Index))
		}
		for _, sc := range acct.StorageChanges {
			for _, chg := range sc.Changes {
				fmt.Fprintf(&b, `<div class="chg">slot %s = <span class="val">%s</span> <span class="dim">(tx %d)</span></div>`,
					short(sc.Key), short(chg.Value), hexToUint64(chg.Index))
			}
		}
		if len(acct.CodeChanges) > 0 {
			fmt.Fprintf(&b, `<div class="chg">code deployed <span class="dim">(%d changes)</span></div>`, len(acct.CodeChanges))
		}
	}

	var out strings.Builder
	fmt.Fprintf(&out, `<div class="bal">── block %d state changes: %d user accounts, %d system contracts ──</div>`,
		blockNum, userCount, systemCount)
	out.WriteString(b.String())
	return out.String()
}
