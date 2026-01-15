package main

import (
	"fmt"

	"github.com/AlexanderYAPPO/go-papa-carlo/parse"
)

func main() {
	res := parse.Parse("ClassA", "/Users/yappo/projects/go-test-project/pkg1/class_a.go")
	fmt.Println("Imports: --------------------------------")
	for _, i := range res.Imports {
		fmt.Println(i.Alias, i.Path)
	}
	fmt.Println("Fields: --------------------------------")
	for _, field := range res.Fields {
		fmt.Println(field.Name, field.Type)
	}
	fmt.Println("--------------------------------")
}
