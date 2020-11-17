package structs

type Format string

const (
	FormatGen8Random        Format = "gen8randombattle"
	FormatGen8              Format = "[Gen 8] Anything Goes"
	FormatGen8DoublesRandom Format = "[Gen 8] Random Doubles Battle"
	FormatGen8Doubles       Format = "[Gen 8] Doubles Ubers"
)

var isDoubles = map[Format]bool{
	FormatGen8Doubles:       true,
	FormatGen8DoublesRandom: true,
}

func IsDoubles(f Format) bool {
	ok := isDoubles[f]
	return ok
}

type BattleSpec struct {
	Format Format

	Players map[string][]*PokemonSpec
}
