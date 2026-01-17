package main

import (
	"errors"
	"fmt"

	"github.com/AlexanderYAPPO/go-papa-carlo/fixtures/struct_with_many_fields/pkg2"
	"github.com/AlexanderYAPPO/go-papa-carlo/fixtures/struct_with_many_fields/pkg3"
	"github.com/AlexanderYAPPO/go-papa-carlo/generate"
	"github.com/AlexanderYAPPO/go-papa-carlo/parse"
	"github.com/AlexanderYAPPO/go-papa-carlo/test_env/struct_with_many_fields/pkg1"
)

func printBuilder() {
	res := parse.Parse("ClassA", "fixtures/struct_with_many_fields/pkg1/class_a.go")
	generatedCode := generate.Generate(res)
	fmt.Println(generatedCode)
}

func useBuilder() {
	obj := pkg1.New().
	WithPublicStringField("string value").
	WithprivateStringField("private string value").
	WithpublicIntField(1).
	WithpublicIntSliceField([]int{1, 2}).
	WithPublicMapField(map[string]int{"a": 5}).
	WithPublicMapClassB(map[string]pkg1.ClassB{"x": {}}).
	WithPublicClassX(pkg2.ClassX{}).
	WithPublicClassB(pkg3.ClassB{}).
	WithPublicMapClassX(map[string]pkg2.ClassX{"xb": {}}).
	WithPublicAnyField(5).
	WithPublicChannel(nil).
	WithPublicChannelIn(nil).
	WithPublicChannelOut(nil).
	WithPublicLambdaFunction(func(i int, s string) (int, error) {return 0, nil}).
	WithPublicErr(errors.New("")).
	WithPublicLambdaFunctionAlias(func(i int, s string) (int, error) {return 0, nil}).
	WithPublicInterface(nil).
	WithPublicDefinedInterface(nil).Build()
	fmt.Printf("%#v\n", obj)
}

func main() {
	useBuilder()
}
