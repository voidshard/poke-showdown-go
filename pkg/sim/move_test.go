package sim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInflate(t *testing.T) {
	m := &Move{ID: "thunderbolt"}

	err := m.Inflate()

	assert.Nil(t, err)
	assert.Equal(t, "Thunderbolt", m.Name)
	assert.Equal(t, 90, m.Power)
	assert.Equal(t, 85, m.Number)
	assert.Equal(t, 100, m.Accuracy)
}

func TestNewMove(t *testing.T) {
	m, err := NewMove("thunderbolt")

	assert.Nil(t, err)
	assert.Equal(t, "Thunderbolt", m.Name)
	assert.Equal(t, 90, m.Power)
	assert.Equal(t, 85, m.Number)
	assert.Equal(t, 100, m.Accuracy)
}
