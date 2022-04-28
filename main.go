package main

import (
	"rustgo"
)

func main() {
	wasm := rustgo.NewWasmLib("rustgo/id_based_self_encryption.wasm")
	defer wasm.Close()

	_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "test.txt")
}
