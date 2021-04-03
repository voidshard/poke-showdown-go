package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var dataTestGender = []struct {
	In     string
	Expect string
}{
	{"Pincurchin, L88, M", "M"},
	{"Froslass, L86, F", "F"},
	{"Glastrier, L84", ""},
}

func TestGender(t *testing.T) {
	for _, tt := range dataTestGender {
		result := (&Pokemon{Details: tt.In}).gender()

		assert.Equal(t, tt.Expect, result)
	}
}
