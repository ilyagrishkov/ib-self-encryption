package main

import (
	"encoding/json"
	"fmt"
	"ibse/internal"
	"io/ioutil"
	"os"
	"rustgo"
	"testing"
	"time"
)

func TestBench100kb(t *testing.T) {
	internal.Init()
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, "./resources")

	n := 100

	var results []float32

	for i := 0; i < n; i++ {
		start := time.Now()
		_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "100KB", "randomidentity")
		elapsed := time.Since(start)
		results = append(results, float32(elapsed.Seconds()))
	}

	file, _ := json.Marshal(results)
	_ = os.MkdirAll("output/", os.ModePerm)
	_ = ioutil.WriteFile("output/100KB.json", file, 0644)
	_ = os.RemoveAll("resources/chunk_store")
}

func TestBench250kb(t *testing.T) {
	internal.Init()
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, "./resources")

	n := 100

	var results []float32

	for i := 0; i < n; i++ {
		start := time.Now()
		_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "250KB", "randomidentity")
		elapsed := time.Since(start)
		results = append(results, float32(elapsed.Seconds()))
	}

	file, _ := json.Marshal(results)
	_ = os.MkdirAll("output/", os.ModePerm)
	_ = ioutil.WriteFile("output/250KB.json", file, 0644)
	_ = os.RemoveAll("resources/chunk_store")
}

func TestBench500kb(t *testing.T) {
	internal.Init()
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, "./resources")

	n := 100

	var results []float32

	for i := 0; i < n; i++ {
		start := time.Now()
		_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "500KB", "randomidentity")
		elapsed := time.Since(start)
		results = append(results, float32(elapsed.Seconds()))
	}

	file, _ := json.Marshal(results)
	_ = os.MkdirAll("output/", os.ModePerm)
	_ = ioutil.WriteFile("output/500KB.json", file, 0644)
	_ = os.RemoveAll("resources/chunk_store")
}

func TestBench750kb(t *testing.T) {
	internal.Init()
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, "./resources")

	n := 100

	var results []float32

	for i := 0; i < n; i++ {
		start := time.Now()
		_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "750KB", "randomidentity")
		elapsed := time.Since(start)
		results = append(results, float32(elapsed.Seconds()))
	}

	file, _ := json.Marshal(results)
	_ = os.MkdirAll("output/", os.ModePerm)
	_ = ioutil.WriteFile("output/100KB.json", file, 0644)
	_ = os.RemoveAll("resources/chunk_store")
}

func TestBench1mb(t *testing.T) {
	internal.Init()
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, "./resources")

	n := 100

	var results []float32

	for i := 0; i < n; i++ {
		start := time.Now()
		_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "1MB", "randomidentity")
		elapsed := time.Since(start)
		results = append(results, float32(elapsed.Seconds()))
	}

	file, _ := json.Marshal(results)
	_ = os.MkdirAll("output/", os.ModePerm)
	_ = ioutil.WriteFile("output/1MB.json", file, 0644)
	_ = os.RemoveAll("resources/chunk_store")
}

func TestBench10mb(t *testing.T) {
	internal.Init()
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, "./resources")

	n := 100

	var results []float32

	for i := 0; i < n; i++ {
		start := time.Now()
		_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "10MB", "randomidentity")
		elapsed := time.Since(start)
		results = append(results, float32(elapsed.Seconds()))
	}

	file, _ := json.Marshal(results)
	_ = os.MkdirAll("output/", os.ModePerm)
	_ = ioutil.WriteFile("output/10MB.json", file, 0644)
	_ = os.RemoveAll("resources/chunk_store")
}

func TestBench25mb(t *testing.T) {
	internal.Init()
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, "./resources")

	n := 100

	var results []float32

	for i := 0; i < n; i++ {
		start := time.Now()
		_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "25MB", "randomidentity")
		elapsed := time.Since(start)
		results = append(results, float32(elapsed.Seconds()))
	}

	file, _ := json.Marshal(results)
	_ = os.MkdirAll("output/", os.ModePerm)
	_ = ioutil.WriteFile("output/25MB.json", file, 0644)
	_ = os.RemoveAll("resources/chunk_store")
}

func TestBench50mb(t *testing.T) {
	internal.Init()
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, "./resources")

	n := 100

	var results []float32

	for i := 0; i < n; i++ {
		start := time.Now()
		_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "50MB", "randomidentity")
		elapsed := time.Since(start)
		results = append(results, float32(elapsed.Seconds()))
	}

	file, _ := json.Marshal(results)
	_ = os.MkdirAll("output/", os.ModePerm)
	_ = ioutil.WriteFile("output/50MB.json", file, 0644)
	_ = os.RemoveAll("resources/chunk_store")
}

func TestBench75mb(t *testing.T) {
	internal.Init()
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, "./resources")

	n := 100

	var results []float32

	for i := 0; i < n; i++ {
		start := time.Now()
		_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "75MB", "randomidentity")
		elapsed := time.Since(start)
		results = append(results, float32(elapsed.Seconds()))
	}

	file, _ := json.Marshal(results)
	_ = os.MkdirAll("output/", os.ModePerm)
	_ = ioutil.WriteFile("output/75MB.json", file, 0644)
	_ = os.RemoveAll("resources/chunk_store")
}

func TestBench100mb(t *testing.T) {
	internal.Init()
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, "./resources")

	n := 100

	var results []float32

	for i := 0; i < n; i++ {
		start := time.Now()
		_, _ = wasm.Invoke("self_encrypt", rustgo.Void, "100MB", "randomidentity")
		elapsed := time.Since(start)
		results = append(results, float32(elapsed.Seconds()))
	}

	file, _ := json.Marshal(results)
	_ = os.MkdirAll("output/", os.ModePerm)
	_ = ioutil.WriteFile("output/100MB.json", file, 0644)
	_ = os.RemoveAll("resources/chunk_store")
}
