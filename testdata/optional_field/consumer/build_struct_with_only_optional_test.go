package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"example.com/fixture/pkg1"
)

// TestBuilderOnlyOptionalFieldsWithoutOptionalSetters tests that the builder works for a struct with only optional
// when none of the fields are set.
func TestBuilderOnlyOptionalFieldsWithoutOptionalSetters(t *testing.T) {
	got := pkg1.NewStructWithOnlyOptionalBuilder().
		Build()
	want := pkg1.StructWithOnlyOptional{}
	assert.Equal(t, want, got)
}

// TestBuilderOnlyOptionalFieldsWithSomeOptionalSetters tests that the builder works for a struct with only optional
// when some of the fields are set.
func TestBuilderOnlyOptionalFieldsWithSomeOptionalSetters(t *testing.T) {
	got := pkg1.NewStructWithOnlyOptionalBuilder().
		WithOptionalName("alice").
		WithOptionalEnabled(true).
		Build()
	want := pkg1.StructWithOnlyOptional{
		Name:    "alice",
		Enabled: true,
	}
	assert.Equal(t, want, got)
}
