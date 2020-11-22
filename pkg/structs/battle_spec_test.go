package structs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsDoubles(t *testing.T) {
	assert.True(t, IsDoubles(FormatGen8Doubles))
	assert.True(t, IsDoubles(FormatGen8DoublesRandom))
	assert.True(t, !IsDoubles(FormatGen8))
	assert.True(t, !IsDoubles(FormatGen8Random))
}
