#!/bin/sh

TEMP_DIR=$(mktemp -d)
WASM_SRC=$(pwd)/wasm-src
EXCLUDE_FILE=$WASM_SRC/.update-exclude 

cd $TEMP_DIR
git clone --depth 1 https://github.com/CyberL1/html-cli .

rsync -av --exclude-from $EXCLUDE_FILE $TEMP_DIR/* $WASM_SRC
rm -rf $TEMP_DIR
