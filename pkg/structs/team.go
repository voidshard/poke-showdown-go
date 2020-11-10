package structs

import (
	"fmt"
	"strings"
)

// Battlemon represents a pokemon with battle relevant stats, items, moves etc
type Battlemon struct {
	Name             string
	Species          string
	Item             string
	Ability          string
	Moves            []string // max len 4
	Nature           string
	EffortValues     *StatValues // 0-255
	Gender           string
	IndividualValues *StatValues // 0-31
	Shiny            bool
	Level            int // 1-100

	// four bonus values whacked in at the end
	Happiness    int
	HPType       string
	PokeballType string
	GigantaMax   bool
}

// Pack a battlemon into the pokemon simulator format
func (b *Battlemon) Pack() string {
	return strings.Join(
		[]string{
			b.Name,
			b.Species,
			b.Item,
			b.Ability,
			strings.Join(b.Moves, ","),
			b.Nature,
			fmt.Sprintf(
				"%d,%d,%d,%d,%d,%d",
				b.EffortValues.HP,
				b.EffortValues.Attack,
				b.EffortValues.Defense,
				b.EffortValues.SpecialAttack,
				b.EffortValues.SpecialDefense,
				b.EffortValues.Speed,
			),
			b.Gender,
			fmt.Sprintf(
				"%d,%d,%d,%d,%d,%d",
				b.IndividualValues.HP,
				b.IndividualValues.Attack,
				b.IndividualValues.Defense,
				b.IndividualValues.SpecialAttack,
				b.IndividualValues.SpecialDefense,
				b.IndividualValues.Speed,
			),
			packbool(b.Shiny, "S"),
			b.Level,
			fmt.Sprintf(
				"%d,%s,%s,%s",
				b.Happiness,
				b.HPType,
				b.PokeballType,
				packbool(b.GigantaMax, "G"),
			),
		},
		"|",
	)
}

// packbool turns a given bool in a blank or the given value
func packbool(b bool, value string) string {
	if b {
		return value
	}
	return ""
}

// StatValues holds values for the 6 pokemon stats
type StatValues struct {
	HP             int
	Attack         int
	Defense        int
	SpecialAttack  int
	SpecialDefense int
	Speed          int
}
