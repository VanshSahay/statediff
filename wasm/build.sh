#!/bin/sh
set -e
cd "$(dirname "$0")/.."
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" wasm/wasm_exec.js
GOOS=js GOARCH=wasm go build -o wasm/statediff.wasm .
ls -lh wasm/statediff.wasm wasm/wasm_exec.js
