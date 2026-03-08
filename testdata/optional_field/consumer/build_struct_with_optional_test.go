package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"example.com/fixture/pkg1"
)

// TestBuilderOptionalFieldCanBeSkipped tests that the builder works for a struct with optional fields
// when none of the optional fields are set.
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

// TestBuilderOptionalFieldCanBeSet tests that the builder works for a struct with optional fields
// when some of the optional fields are set.
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
