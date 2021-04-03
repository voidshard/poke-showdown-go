package sim

import (
	"github.com/voidshard/poke-showdown-go/pkg/internal/structs"
)

// Side represents the current status & options of one side in a battle.
type Side struct {
	// Player id
	Player string `json:"player"`

	// True if this player's side doesn't need to make a decision for the battle
	// to progress (ie. a pokemon has been switched / fainted on the other team)
	Wait bool `json:"wait"`

	// Field is the players field
	Field []*Slot

	// Pokemon is the list of pokemon on this player's team.
	// Order here is important (for "switch" instructions)
	Pokemon []*Pokemon `json:"pokemon"`
}

// toSide mutates the raw parsed showdown update in to our nicer style
// of update & fleshes out everything we can
func toSide(in *structs.Update) *Side {
	out := &Side{
		Player:  in.Team.Player,
		Wait:    in.Wait,
		Field:   []*Slot{},
		Pokemon: []*Pokemon{},
	}

	// Rather than return active slot info depending on force switch
	// we'll always return it, but tag slots with "this must switch"
	// along with no other options entries.
	//
	// Observations:
	// - according to Zarel the order of pokemon in the 'Active'
	//   block matches the pokemon in 'Side'
	// - if a pokemon must switch then ForceSwitch []bool is given, which
	//   will have an entry per slot
	// - if under a force switch no 'Active' info is given
	if len(in.ForceSwitch) > 0 {
		for i, mustSwitch := range in.ForceSwitch {
			lot := &Slot{
				ID:     slotID(i, in.Team.Player),
				Switch: mustSwitch,
				Ident:  in.Team.Pokemon[i].Ident,
				Index:  i,
			}
			out.Field = append(out.Field, lot)
		}
	} else {
		for i, data := range in.Active {
			lot := &Slot{
				ID:      slotID(i, in.Team.Player),
				Trapped: data.Trapped, // pokemon cannot switch
				Options: toOptions(data),
				Ident:   in.Team.Pokemon[i].Ident,
				Index:   i,
			}
			out.Field = append(out.Field, lot)
		}
	}
	for _, p := range in.Team.Pokemon {
		out.Pokemon = append(out.Pokemon, toPokemon(p))
	}

	return out
}
