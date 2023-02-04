package main

import (
	"di-demo/demo"
	"di-demo/di"
	"fmt"
)

func main() {
	// 构造全量字段
	c := di.NewContainer()
	c.Provide(&demo.StructA{}, func() (di.Dependency, error) {
		return &demo.StructA{}, nil
	})
	c.Provide(&demo.StructB{}, func() (di.Dependency, error) {
		return &demo.StructB{}, nil
	})
	c.Provide(&demo.StructC{}, func() (di.Dependency, error) {
		return &demo.StructC{
			Foo: "bar",
		}, nil
	})
	// end

	// 业务
	result, err := c.Get(&demo.StructC{})
	if err != nil {
		fmt.Printf("Failed to factory: %v\n", err)
		return
	}
	fmt.Printf("Obj: %+v\n", result)
	fmt.Printf("Container: %+v\n", c)
}
