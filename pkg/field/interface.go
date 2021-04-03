package field

import (
	"github.com/voidshard/poke-showdown-go/pkg/sim"
)

// Watcher keeps track of the field based on updates.
type Watcher interface {
	// Update Watcher info based on update event
	Update(*sim.Update)

	// Turn returns the current game turn
	Turn() int

	// WhoIs returns the pokemon at the given slot (if any)
	WhoIs(string) *sim.Pokemon
}

// NewWatcher returns a default field watcher
func NewWatcher() Watcher {
	return &Field{
		sides: map[string]*sim.Side{},
		slots: map[string]*sim.Pokemon{},
	}
}
