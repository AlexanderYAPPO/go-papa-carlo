package consumer

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"example.com/fixture/pkg1"
	"example.com/fixture/pkg2"
	"example.com/fixture/pkg3"
)

// TestBuilderWorks tests that the builder works for a struct with many fields 
// and it's expected that it covers all supported field types.
func TestBuilderWorks(t *testing.T) {
	channel := make(chan int)
	channelInSource := make(chan int)
	channelOutSource := make(chan int)
	var channelIn chan<- int = channelInSource
	var channelOut <-chan int = channelOutSource

	lambdaFunctionSum := func(i int, s string) (int, error) {
		return i + len(s), nil
	}
	lambdaFunctionAliasMult := func(i int, s string) (int, error) {
		return i * len(s), nil
	}

	expectedErr := errors.New("boom")
	expectedDatetime := time.Date(2026, 1, 17, 15, 42, 10, 0, time.UTC)

	got := pkg1.NewStructABuilder().
		WithStringField("hello").
		WithIntField(42).
		WithIntSliceField([]int{1, 2, 3}).
		WithMapField(map[string]int{"one": 1, "two": 2}).
		WithMapStructB(map[string]pkg1.StructB{
			"first":  {FieldOne: 10},
			"second": {FieldOne: 20},
		}).
		WithAliasedStruct(pkg2.StructX{FieldInt: 7, FieldBool: true}).
		WithNonaliasedStruct(pkg3.StructB{FieldStr: "value"}).
		WithMapStructX(map[string]pkg2.StructX{
			"entry": {FieldInt: 9, FieldBool: false},
		}).
		WithAnyField(map[string]bool{"ok": true}).
		WithChannel(channel).
		WithChannelIn(channelIn).
		WithChannelOut(channelOut).
		WithLambdaFunction(lambdaFunctionSum).
		WithErr(expectedErr).
		WithLambdaFunctionAlias(lambdaFunctionAliasMult).
		WithInterface([]string{"a", "b"}).
		WithDefinedInterface(map[string]string{"kind": "demo"}).
		WithDatetime(expectedDatetime).
		Build()

	want := pkg1.StructA{
		StringField:   "hello",
		IntField:      42,
		IntSliceField: []int{1, 2, 3},
		MapField:      map[string]int{"one": 1, "two": 2},
		MapStructB: map[string]pkg1.StructB{
			"first":  {FieldOne: 10},
			"second": {FieldOne: 20},
		},
		AliasedStruct:    pkg2.StructX{FieldInt: 7, FieldBool: true},
		NonaliasedStruct: pkg3.StructB{FieldStr: "value"},
		MapStructX: map[string]pkg2.StructX{
			"entry": {FieldInt: 9, FieldBool: false},
		},
		AnyField:            map[string]bool{"ok": true},
		Channel:             channel,
		ChannelIn:           channelIn,
		ChannelOut:          channelOut,
		LambdaFunction:      nil,
		Err:                 expectedErr,
		LambdaFunctionAlias: nil,
		Interface:           []string{"a", "b"},
		DefinedInterface:    map[string]string{"kind": "demo"},
		Datetime:            expectedDatetime,
	}

	gotComparable := got
	// lambda functions cannot be compared as is since they are stored as pointers.
	// so we set them to nil to make the comparison work and test the functions separately later.
	gotComparable.LambdaFunction = nil
	gotComparable.LambdaFunctionAlias = nil

	assert.Equal(t, want, gotComparable)

	// sum function variant expected
	lambdaResult, lambdaErr := got.LambdaFunction(5, "abc")
	assert.NoError(t, lambdaErr)
	assert.Equal(t, 8, lambdaResult)

	// mult function variant expected
	lambdaAliasResult, lambdaAliasErr := got.LambdaFunctionAlias(5, "abc")
	assert.NoError(t, lambdaAliasErr)
	assert.Equal(t, 15, lambdaAliasResult)
}
