package main

import (
	"rustgo"
)

func main() {
	wasm := rustgo.NewWasmLib("rustgo/ib_self_encryption_rust.wasm")
	defer wasm.Close()

	_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "test.txt")

	_, _ = wasm.Invoke("self_decrypt", rustgo.Void)
}
