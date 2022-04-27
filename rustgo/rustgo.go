package rustgo

import (
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"log"
	"reflect"
	"strings"
)

type Type int64

const (
	Int32 Type = iota
	Int64
	Float32
	Float64
	String
	Array
	Void
)

type Pair[T, U any] struct {
	First  T
	Second U
}

type WasmLib struct {
	Instance   wasm.Instance
	References map[int32]int
}

func NewWasmLib(filepath string) WasmLib {
	bytes, _ := wasm.ReadBytes(filepath)
	instance, _ := wasm.NewInstance(bytes)

	return WasmLib{
		Instance:   instance,
		References: make(map[int32]int),
	}
}

func (wasmLib WasmLib) Invoke(function string, returnType Type, arguments ...interface{}) (interface{}, error) {
	f := wasmLib.Instance.Exports[function]

	var args []interface{}
	var pointers []Pair[*uint8, int32]
	for _, argument := range arguments {
		isPointer := reflect.ValueOf(argument).Kind() == reflect.Ptr
		if isPointer {
			allocateResult, _ := wasmLib.Instance.Exports["allocate"](1)
			inputPointer := allocateResult.ToI32()

			wasmLib.References[inputPointer] = 1
			args = append(args, inputPointer)
			pointers = append(pointers, Pair[*uint8, int32]{argument.(*uint8), inputPointer})
			continue
		}

		switch argument.(type) {
		case string:
			{
				subject := argument.(string)
				lengthOfSubject := len(subject)
				allocateResult, _ := wasmLib.Instance.Exports["allocate"](lengthOfSubject)
				inputPointer := allocateResult.ToI32()

				// Write the subject into the memory.
				memory := wasmLib.Instance.Memory.Data()[inputPointer:]

				for nth := 0; nth < lengthOfSubject; nth++ {
					memory[nth] = subject[nth]
				}

				// C-string terminates by NULL.
				memory[lengthOfSubject] = 0
				wasmLib.References[inputPointer] = lengthOfSubject
				args = append(args, inputPointer)
			}
		//case uint32:
		//	{
		//		allocateResult, _ := wasmLib.Instance.Exports["allocate"](1)
		//		inputPointer := allocateResult.ToI32()
		//
		//		wasmLib.References[inputPointer] = 1
		//		args = append(args, inputPointer)
		//		pointers = append(pointers, inputPointer)
		//	}
		default:
			args = append(args, argument)
		}
	}
	result, err := f(args...)
	if err != nil {
		return nil, err
	}

	switch returnType {
	case Int32:
		return result.ToI32(), nil
	case Int64:
		return result.ToI64(), nil
	case Float32:
		return result.ToF32(), nil
	case Float64:
		return result.ToF64(), nil
	case String:
		{
			outputPointer := result.ToI32()
			memory := wasmLib.Instance.Memory.Data()[outputPointer:]
			nth := 0
			var output strings.Builder

			for {
				if memory[nth] == 0 {
					break
				}

				output.WriteByte(memory[nth])
				nth++
			}
			wasmLib.References[outputPointer] = nth
			return output.String(), nil
		}
	case Array:
		{
			for _, pointer := range pointers {
				memory := wasmLib.Instance.Memory.Data()[pointer.Second]
				*pointer.First = memory
			}
			outputPointer := result.ToI32()
			memory := wasmLib.Instance.Memory.Data()[outputPointer:]
			size := int(wasmLib.Instance.Memory.Data()[pointers[len(pointers)-1].Second])
			nth := 0
			var output []byte

			for {
				if nth == size {
					break
				}
				output = append(output, memory[nth])
				nth++
			}
			wasmLib.References[outputPointer] = size
			return output, nil
		}
	case Void:
		return nil, nil
	}
	panic("Couldn't invoke the function")
}

func (wasmLib WasmLib) Close() {
	deallocate := wasmLib.Instance.Exports["deallocate"]
	for ptr, length := range wasmLib.References {
		_, err := deallocate(ptr, length)
		if err != nil {
			log.Printf("Failed to deallocate %d of length %d\n", ptr, length)
		}
	}
	wasmLib.Instance.Close()
}
