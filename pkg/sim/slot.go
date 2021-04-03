package sim

import (
	"fmt"
)

// Slot is a place on the field that can hold a pokemon.
// Note that it may not necessarily hold one (in the case of doubles).
type Slot struct {
	// Showdown uses references like 'p1a' 'p2a' 'p2b' to give player & slot addresses
	ID string

	// This slot is required to switch a new pokemon in
	Switch bool

	// fields set only if slot is occupied by a pokemon
	// Pokemon identity string
	Ident string

	// Pokemon index (in team)
	Index int

	// This pokemon is unable to switch out
	Trapped bool

	// Options of pokemon in this slot
	Options *Options
}

// slotID makes a showdown slot ID from a player & position
func slotID(i int, player string) string {
	name := ""
	switch i {
	case 1:
		name = "b"
	case 2:
		name = "c"
	case 3:
		name = "d"
	default:
		name = "a"
	}

	return fmt.Sprintf("%s%s", player, name)
}
