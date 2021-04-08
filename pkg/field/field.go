package field

import (
	"github.com/voidshard/poke-showdown-go/pkg/event"
	"github.com/voidshard/poke-showdown-go/pkg/sim"
)

// Field implements some basic event watching.
type Field struct {
	turn  int
	sides map[string]*sim.Side
	slots map[string]*sim.Pokemon
}

// Turn returns the current turn of the game.
// Nb. compulsory switches happen before the start of the next turn
// where applicable.
func (f *Field) Turn() int {
	return f.turn
}

// WhoIs returns the pokemon in the given slot or nil if there is none.
func (f *Field) WhoIs(slotID string) *sim.Pokemon {
	p, ok := f.slots[slotID]
	if !ok {
		return nil
	}
	return p
}

// Update sets internal data based on the given update struct
func (f *Field) Update(ud *sim.Update) {
	if ud.Side != nil {
		f.sides[ud.Side.Player] = ud.Side
		for i, s := range ud.Side.Field {
			f.slots[s.ID] = ud.Side.Pokemon[i]
		}
		return
	}

	if ud.Event != nil {
		switch ud.Event.Type {
		case event.Turn:
			f.turn = ud.Event.Magnitude
		}
	}
}
