package structs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBattleState(t *testing.T) {
	result := NewBattleState()

	assert.NotNil(t, result)
	assert.NotNil(t, result.Field)
	assert.NotNil(t, result.Events)
}
