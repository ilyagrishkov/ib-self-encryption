package rustgo

import (
	wasm "github.com/wasmerio/wasmer-go/wasmer"
	"io/ioutil"
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
	Instance   *wasm.Instance
	References map[int32]int
}

func NewWasmLib(filepath string) WasmLib {
	wasmBytes, _ := ioutil.ReadFile(filepath)

	store := wasm.NewStore(wasm.NewEngine())
	module, _ := wasm.NewModule(store, wasmBytes)

	wasiEnv, _ := wasm.NewWasiStateBuilder("wasi-program").InheritStdout().MapDirectory("./", ".").Finalize()
	importObject, _ := wasiEnv.GenerateImportObject(store, module)

	instance, _ := wasm.NewInstance(module, importObject)

	return WasmLib{
		Instance:   instance,
		References: make(map[int32]int),
	}
}

func (wasmLib WasmLib) Invoke(function string, returnType Type, arguments ...interface{}) (interface{}, error) {
	f, _ := wasmLib.Instance.Exports.GetFunction(function)

	var args []interface{}
	var pointers []Pair[*uint8, int32]
	for _, argument := range arguments {
		isPointer := reflect.ValueOf(argument).Kind() == reflect.Ptr
		if isPointer {
			allocate, _ := wasmLib.Instance.Exports.GetFunction("allocate")
			allocateResult, _ := allocate(1)
			inputPointer := allocateResult.(int32)

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
				allocate, _ := wasmLib.Instance.Exports.GetFunction("allocate")
				allocateResult, _ := allocate(lengthOfSubject)
				inputPointer := allocateResult.(int32)

				// Write the subject into the memory.
				memory, _ := wasmLib.Instance.Exports.GetMemory("memory")
				memorySlice := memory.Data()[inputPointer:]

				for nth := 0; nth < lengthOfSubject; nth++ {
					memorySlice[nth] = subject[nth]
				}

				// C-string terminates by NULL.
				memorySlice[lengthOfSubject] = 0
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
		return result.(int32), nil
	case Int64:
		return result.(int64), nil
	case Float32:
		return result.(float32), nil
	case Float64:
		return result.(float64), nil
	case String:
		{
			outputPointer := result.(int32)
			memory, _ := wasmLib.Instance.Exports.GetMemory("memory")
			memorySlice := memory.Data()[outputPointer:]
			nth := 0
			var output strings.Builder

			for {
				if memorySlice[nth] == 0 {
					break
				}

				output.WriteByte(memorySlice[nth])
				nth++
			}
			wasmLib.References[outputPointer] = nth
			return output.String(), nil
		}
	case Array:
		{
			memory, _ := wasmLib.Instance.Exports.GetMemory("memory")
			for _, pointer := range pointers {

				memorySlice := memory.Data()[pointer.Second]
				*pointer.First = memorySlice
			}
			outputPointer := result.(int32)
			memorySlice := memory.Data()[outputPointer:]
			size := int(memory.Data()[pointers[len(pointers)-1].Second])
			nth := 0
			var output []byte

			for {
				if nth == size {
					break
				}
				output = append(output, memorySlice[nth])
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
	deallocate, _ := wasmLib.Instance.Exports.GetFunction("deallocate")
	for ptr, length := range wasmLib.References {
		_, err := deallocate(ptr, length)
		if err != nil {
			log.Printf("Failed to deallocate %d of length %d\n", ptr, length)
		}
	}
	wasmLib.Instance.Close()
}
