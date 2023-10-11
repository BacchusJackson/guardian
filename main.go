// Package main
package main

import "github.com/bytecodealliance/wasmtime-go"

const kanikoWasmFile = "~/.local/wasm/kaniko/kaniko-executor.wasm"

func main() {
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	module, err := wasmtime.NewModuleFromFile(engine, kanikoWasmFile)
	if err != nil {
		panic(err)
	}
	instance, err := wasmtime.NewInstance(store, module, []wasmtime.AsExtern{})
	if err != nil {
		panic(err)
	}
	instance.GetFunc()
}
