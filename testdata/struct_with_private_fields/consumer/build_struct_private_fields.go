package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"example.com/fixture/pkg"
)

func TestBuilderWorks(t *testing.T) {
	got := pkg.NewStructWithPrivateFieldsThatWorksBuilder().
		WithPublicField(1).
		Build()
	want := pkg.StructWithPrivateFieldsThatWorks{
		PublicField: 1,
	}
	assert.Equal(t, want, got)
}
