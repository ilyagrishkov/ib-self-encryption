package main

import (
	"fmt"
	"rustgo"
)

func main() {
	wasm := rustgo.NewWasmLib("rustgo/id_based_self_encryption.wasm")
	defer wasm.Close()

	res, _ := wasm.Invoke("greet", rustgo.String, "Testhhfjhffgfgjhfghjfghfghfjgfhgfgf")
	fmt.Println(res)
}
