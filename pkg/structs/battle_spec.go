package structs

type Format string

const (
	// FormatGen8Random random SS singles battle
	//	FormatGen8Random Format = "[Gen 8] Random Battle"

	// FormatGen8 SS singles battle with given teams
	FormatGen8 Format = "[Gen 8] Anything Goes"

	// FormatGen8DoublesRandom a SS doubles battle with random teams
	//	FormatGen8DoublesRandom Format = "[Gen 8] Random Doubles Battle"

	// FormatGen8Doubles a SS doubles battle with given teams
	FormatGen8Doubles Format = "[Gen 8] Doubles Ubers"
)

// isDoubles holds if a format is a double battle.
// At present it's assumed not-doubles is singles.
var isDoubles = map[Format]bool{
	FormatGen8Doubles: true,
	//	FormatGen8DoublesRandom: true,
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

	// Players indicates who is in the battle & their team (if required).
	// `Random` battles don't require a team, but we still need the players
	// given.
	Players map[string][]*PokemonSpec

	// Seed for internal RNG
	Seed int
}
