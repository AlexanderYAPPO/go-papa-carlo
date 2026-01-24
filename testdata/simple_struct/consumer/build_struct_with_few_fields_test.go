package consumer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"example.com/fixture/pkg1"
)

func TestBuilderWorks(t *testing.T) {
	got := pkg1.NewStructWithFewFieldsBuilder().
		WithFieldInt(1).
		WithFieldBool(true).
		WithFieldDatetime(time.Date(2026, 1, 17, 15, 42, 10, 0, time.UTC)).
		WithFieldMapField(map[string]int{"two": 2, "three": 3}).
		Build()
	want := pkg1.StructWithFewFields{
		FieldInt:       1,
		FieldBool:      true,
		FieldDatetime:  time.Date(2026, 1, 17, 15, 42, 10, 0, time.UTC),
		FieldMapField:  map[string]int{"two": 2, "three": 3},
	}
	assert.Equal(t, want, got)
}
