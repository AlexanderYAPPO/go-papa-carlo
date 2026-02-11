package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"example.com/fixture/pkg1"
)

func TestBuilderOnlyOptionalFieldsWithoutOptionalSetters(t *testing.T) {
	got := pkg1.NewStructWithOnlyOptionalBuilder().
		Build()
	want := pkg1.StructWithOnlyOptional{}
	assert.Equal(t, want, got)
}

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
