package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"example.com/fixture/pkg1"
)

// TestBuilderOmitsTaggedField tests that the builder works for a struct with `papa-carlo:"omit"` tag.
func TestBuilderOmitsTaggedField(t *testing.T) {
	got := pkg1.NewStructWithOmittedFieldBuilder().
		WithRequiredInt(11).
		WithFieldWithOtherPapaCarloTag("xyz").
		WithRequiredBool(true).
		Build()
	want := pkg1.StructWithOmittedField{
		RequiredInt:  11,
		FieldWithOtherPapaCarloTag: "xyz",
		RequiredBool: true,
	}
	assert.Equal(t, want, got)
}
