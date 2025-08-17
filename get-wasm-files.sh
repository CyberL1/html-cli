#!/bin/sh

PLAYGROUND_DIR=$(dirname "$0")

cp $(go env GOROOT)/lib/wasm/wasm_exec.js $PLAYGROUND_DIR

cd wasm-src
GOOS=js GOARCH=wasm go build -o ../html-cli.wasm
