package pkg1

import (
	"time"

	pkg2alias "github.com/AlexanderYAPPO/go-papa-carlo/test_scenarios/struct_with_many_fields/pkg2"
	"github.com/AlexanderYAPPO/go-papa-carlo/test_scenarios/struct_with_many_fields/pkg3"
)

type LambdaFunction func(int, string) (int, error)

type StructA struct {
	StringField         string
	IntField            int
	IntSliceField       []int
	MapField            map[string]int
	MapStructB          map[string]StructB
	AliasedStruct       pkg2alias.StructX
	NonaliasedStruct    pkg3.StructB
	MapStructX          map[string]pkg2alias.StructX
	AnyField            any
	Channel             chan int
	ChannelIn           chan<- int
	ChannelOut          <-chan int
	LambdaFunction      func(int, string) (int, error)
	Err                 error
	LambdaFunctionAlias LambdaFunction
	VariadicFunction    func(...int) (int, error)
	Interface           interface{}
	DefinedInterface    InterfaceA
	Datetime            time.Time
	ParenthesizedStruct (StructB)
}

type InterfaceA interface{}

type StructB struct {
	FieldOne int
}

type ClassC[T any] struct {
	FieldGen T
}

type privateClassC struct {
	fieldOne int
}
