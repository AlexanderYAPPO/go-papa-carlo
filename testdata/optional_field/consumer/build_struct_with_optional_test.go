package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"example.com/fixture/pkg1"
)

func TestBuilderOptionalFieldCanBeSkipped(t *testing.T) {
	got := pkg1.NewStructWithOptionalFieldBuilder().
		WithRequiredInt(7).
		WithRequiredMap(map[string]int{"x": 1}).
		Build()
	want := pkg1.StructWithOptionalField{
		RequiredInt: 7,
		RequiredMap: map[string]int{"x": 1},
	}
	assert.Equal(t, want, got)
}

func TestBuilderOptionalFieldCanBeSet(t *testing.T) {
	got := pkg1.NewStructWithOptionalFieldBuilder().
		WithRequiredInt(7).
		WithRequiredMap(map[string]int{"x": 1}).
		WithOptionalOptionalString("hello").
		WithOptionalOptionalBool(true).
		Build()
	want := pkg1.StructWithOptionalField{
		RequiredInt:    7,
		OptionalString: "hello",
		OptionalBool:   true,
		RequiredMap:    map[string]int{"x": 1},
	}
	assert.Equal(t, want, got)
}
