package pkg1

import (
	pkg2alias "github.com/AlexanderYAPPO/go-papa-carlo/fixtures/struct_with_many_fields/pkg2"
	"github.com/AlexanderYAPPO/go-papa-carlo/fixtures/struct_with_many_fields/pkg3"
)

type lambdaFunction func(int, string) (int, error)

type ClassA struct {
	PublicStringField         string
	privateStringField        string
	publicIntField            int
	publicIntSliceField       []int
	PublicMapField            map[string]int
	PublicMapClassB           map[string]ClassB
	PublicClassX              pkg2alias.ClassX
	PublicClassB              pkg3.ClassB
	PublicMapClassX           map[string]pkg2alias.ClassX
	PublicAnyField            any
	PublicChannel             chan int
	PublicChannelIn           chan<- int
	PublicChannelOut          <-chan int
	PublicLambdaFunction      func(int, string) (int, error)
	PublicErr                 error
	PublicLambdaFunctionAlias lambdaFunction
	PublicInterface           interface{}
	PublicDefinedInterface    InterfaceA
}

type InterfaceA interface{}

type ClassB struct {
	FieldOne int
}

type ClassC[T any] struct {
	FieldGen T
}

type privateClassC struct {
	fieldOne int
}
