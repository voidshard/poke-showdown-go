package sim

import (
	"log"

	"github.com/voidshard/poke-showdown-go/pkg/internal/structs"
	data "github.com/voidshard/poke-showdown-go/pkg/pokedata"
)

const (
	// Target details what pokemon can be targeted in doubles+ battles
	// github.com/smogon/pokemon-showdown/blob/master/sim/dex-moves.ts
	// play.pokemonshowdown.com/data/moves.json
	/**
	 * adjacentAlly - Only relevant to Doubles or Triples, the move only targets an ally of the user.
	 * adjacentAllyOrSelf - The move can target the user or its ally.
	 * adjacentFoe - The move can target a foe, but not (in Triples) a distant foe.
	 * all - The move targets the field or all Pokémon at once.
	 * allAdjacent - The move is a spread move that also hits the user's ally.
	 * allAdjacentFoes - The move is a spread move.
	 * allies - The move affects all active Pokémon on the user's team.
	 * allySide - The move adds a side condition on the user's side.
	 * allyTeam - The move affects all unfainted Pokémon on the user's team.
	 * any - The move can hit any other active Pokémon, not just those adjacent.
	 * foeSide - The move adds a side condition on the foe's side.
	 * normal - The move can hit one adjacent Pokémon of your choice.
	 * randomNormal - The move targets an adjacent foe at random.
	 * scripted - The move targets the foe that damaged the user.
	 * self - The move affects the user of the move.
	 */
	AdjAlly       = "adjacentAlly" // helping hand
	AdjAllyOrSelf = "adjacentAllyOrSelf"
	AdjFoe        = "adjacentFoe"
	All           = "all"
	Adj           = "allAdjacent" // earthquake
	AdjFoes       = "allAdjacentFoes"
	Allies        = "allies"   // life dew
	AllySide      = "allySide" // tailwind
	AllyTeam      = "allyTeam" // heal bell
	Any           = "any"
	FoeSide       = "foeSide"      // electroweb
	Normal        = "normal"       // waterfall
	RandomNormal  = "randomNormal" // metronome
	Scripted      = "scripted"     // counter
	Self          = "self"         // recover
)

// RequiresDoublesTarget returns if a target must be specified in doubles
func RequiresDoublesTarget(target string) bool {
	switch target {
	case AdjAllyOrSelf, AdjFoe, Any, Normal:
		return true
	}
	return false
}

// ValidDoublesTargets returns who a move may target in doubles.
// github.com/smogon/pokemon-showdown/blob/master/sim/SIM-PROTOCOL.md
// +2 +1  [enemies]
// -1 -2  [allies]
func ValidDoublesTargets(target string) []int {
	switch target {
	case AdjAllyOrSelf:
		return []int{-1, -2}
	case AdjFoe, Normal:
		return []int{1, 2}
	case Any:
		return []int{1, 2, -1, -2}
	}
	return nil
}

// Move struct expands on move data to add battle info as well
type Move struct {
	// Id is the name string, but no symbols & lowercase
	ID string `json:"id"`

	// Fields only supplied during battle
	PP       int  `json:"pp"`
	Disabled bool `json:"disabled"`

	// Fields that can be parsed from Showdown or copied from data
	Name           string `json:"move"`
	MaxPP          int    `json:"maxpp"`
	Target         string `json:"target"`
	Number         int
	Accuracy       int
	Power          int
	Category       string
	Priority       int
	Flags          map[string]int
	Secondary      data.Secondary
	Type           string
	Description    string
	VolatileStatus string
	Boosts         data.Stats
}

// NewMove builds a Move from it's ID & then calls Inflate
func NewMove(id string) (*Move, error) {
	m := &Move{ID: id}
	return m, m.Inflate()
}

// Inflate attempts to set data using the ID field, pulling
// from the Pokemon-Showdown dataset
func (m *Move) Inflate() error {
	dex, err := data.MoveDex(m.ID)
	if err != nil {
		return err
	}

	m.Name = dex.Name
	m.Target = dex.Target
	m.MaxPP = dex.MaxPP
	m.Number = dex.Number
	m.Accuracy = dex.Accuracy
	m.Power = dex.Power
	m.Category = dex.Category
	m.Priority = dex.Priority
	m.Flags = dex.Flags
	m.Secondary = dex.Secondary
	m.Type = dex.Type
	m.Description = dex.Description
	m.VolatileStatus = dex.VolatileStatus
	m.Boosts = dex.Boosts

	return nil
}

// toMove converts an internal parsed move to sim.Move
func toMove(in *structs.Move) *Move {
	if in == nil {
		return nil
	}

	id := in.ID
	if in.ID == "" {
		id = in.Name
	}

	mov, err := NewMove(id)
	if err != nil {
		log.Printf("[move.go] failed to find move '%s'\n", in.ID)
		return nil
	}

	mov.PP = in.PP
	mov.Disabled = in.Disabled

	return mov
}
