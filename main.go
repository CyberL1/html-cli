//go:build !js || !wasm
// +build !js !wasm

package main

import "html-cli/cmd"

func main() {
	cmd.Execute()
}
