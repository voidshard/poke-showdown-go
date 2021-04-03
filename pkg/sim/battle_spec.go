package sim

import (
	"fmt"
)

// Format is some battle format that our simulatr understands
type Format string

const (
	// FormatGen8 SS singles battle with given teams
	FormatGen8 Format = "[Gen 8] Anything Goes"

	// FormatGen8Doubles a SS doubles battle with given teams
	FormatGen8Doubles Format = "[Gen 8] Doubles Ubers"
)

// isDoubles holds if a format is a double battle.
// At present it's assumed not-doubles is singles.
var isDoubles = map[Format]bool{
	FormatGen8Doubles: true,
}

// IsDoubles returns if the given format is a `doubles` battle
func IsDoubles(f Format) bool {
	ok := isDoubles[f]
	return ok
}

// BattleSpec is all the required data to start a new battle
type BattleSpec struct {
	// Format is the battle style (singles/doubles)
	Format Format

	// Players indicates player team(s). Random battles can use nil
	// values to represent how many players there are.
	// The simulator refers to the players in order as p1 p2 p3 etc...
	Players [][]*PokemonSpec

	// Seed for internal RNG
	Seed int
}

// validate does some simple checks
func (b *BattleSpec) validate() error {
	if b.Players == nil {
		return fmt.Errorf("no player team data")
	}

	if len(b.Players) != 2 {
		return fmt.Errorf("two players are required")
	}

	for _, p := range b.Players {
		if len(p) < 1 {
			return fmt.Errorf("all players must have at least one pokemon")
		} else if len(p) > 6 {
			return fmt.Errorf("players cannot have more than six pokemon")
		}
		if IsDoubles(b.Format) && len(p) < 2 {
			return fmt.Errorf("doubles players must have at least two pokemon")
		}
	}

	return nil
}
