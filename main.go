package main

import (
	"fmt"

	"github.com/AlexanderYAPPO/go-papa-carlo/generate"
	"github.com/AlexanderYAPPO/go-papa-carlo/parse"
)

func main() {
	res := parse.Parse("ClassX", "fixtures/struct_with_many_fields/pkg2/class_x.go")
	fmt.Println("Imports: --------------------------------")
	for _, i := range res.Imports {
		fmt.Println(i.ReferenceName, i.Path)
	}
	fmt.Println("Fields: --------------------------------")
	for _, field := range res.Fields {
		fmt.Println(field.Name, field.Type)
	}
	fmt.Println("--------------------------------")

	fmt.Println("GENERATION PART: --------------------------------")
	generatedCode := generate.Generate(res)
	fmt.Println(generatedCode)
}
