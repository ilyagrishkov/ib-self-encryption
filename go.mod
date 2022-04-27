module fabric_app

go 1.18

require rustgo v0.0.1

require (
	github.com/wasmerio/go-ext-wasm v0.3.1 // indirect
	github.com/wasmerio/wasmer-go v1.0.4 // indirect
)

replace rustgo => ./rustgo
