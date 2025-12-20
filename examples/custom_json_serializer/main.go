package main

import (
	"encoding/json"
	"fmt"

	"trpc.group/trpc-go/trpc-go/codec"
)

// https://github.com/trpc-group/trpc-go/issues/157

// StdJSONSerializer 1. 定义我们自己的 Serializer (使用官方标准库，或追求性能情况下使用sonic)
type StdJSONSerializer struct{}

func (s *StdJSONSerializer) Unmarshal(in []byte, body interface{}) error {
	fmt.Println("[StdJSONSerializer] Using official encoding/json Unmarshal")
	return json.Unmarshal(in, body)
}

func (s *StdJSONSerializer) Marshal(body interface{}) ([]byte, error) {
	fmt.Println("[StdJSONSerializer] Using official encoding/json Marshal")
	return json.Marshal(body)
}

func main() {
	data := []byte(`{"name": "trpc", "age": 18}`)
	var result map[string]interface{}

	fmt.Println("=== 1. Before Registration (Default Behavior) ===")
	// 获取默认的 JSON Serializer
	// 默认情况下，trpc-go 使用 json-iterator
	defaultSerializer := codec.GetSerializer(codec.SerializationTypeJSON)
	if defaultSerializer != nil {
		_ = defaultSerializer.Unmarshal(data, &result)
		fmt.Printf("Default Result: %+v\n", result)
	}

	fmt.Println("\n=== 2. Registering Custom Serializer ===")
	// 核心步骤：覆盖默认的 JSON Serializer
	// 任何后续调用 codec.GetSerializer(codec.SerializationTypeJSON) 的地方都会受到影响
	codec.RegisterSerializer(codec.SerializationTypeJSON, &StdJSONSerializer{})

	fmt.Println("\n=== 3. After Registration (Custom Behavior) ===")
	// 再次获取（此时已经是我们要的 StdJSONSerializer 了）
	newSerializer := codec.GetSerializer(codec.SerializationTypeJSON)
	if newSerializer != nil {
		// 这里的调用应该会打印 "[StdJSONSerializer] ..."
		_ = newSerializer.Unmarshal(data, &result)
		fmt.Printf("Custom Result: %+v\n", result)
	}
}
