package consumer

import (
	"fmt"
	"testing"
	"time"

	"example.com/fixture/pkg1"
)

func TestBuilderWorks(t *testing.T) {
	obj := pkg1.NewStructWithFewFieldsBuilder().
		WithFieldInt(1).
		WithFieldBool(true).
		WithFieldDatetime(time.Date(2026, 1, 17, 15, 42, 10, 0, time.UTC)).
		WithFieldMapField(map[string]int{"two": 2, "three": 3}).
		Build()
	fmt.Printf("%#v\n", obj)
}
