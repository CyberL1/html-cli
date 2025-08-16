#!/bin/sh

PLAYGROUND_DIR=$(dirname "$0")

cp $(go env GOROOT)/lib/wasm/wasm_exec.js $PLAYGROUND_DIR
GOOS=js GOARCH=wasm go build -o $PLAYGROUND_DIR/html-cli.wasm
