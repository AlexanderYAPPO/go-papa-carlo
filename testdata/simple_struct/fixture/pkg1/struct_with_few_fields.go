package pkg1

import (
	"time"
)

// StructWithFewFields ...
type StructWithFewFields struct {
	FieldInt       int
	FieldBool      bool
	FieldDatetime time.Time
	FieldMapField map[string]int
}
