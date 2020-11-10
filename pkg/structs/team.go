package structs

import (
	"fmt"
	"strings"
)

// Battlemon represents a pokemon with battle relevant stats, items, moves etc
type Battlemon struct {
	Name             string      `json:"name"`
	Species          string      `json:"species"`
	Item             string      `json:"item"`
	Ability          string      `json:"ability"`
	Moves            []string    `json:"moves"` // max len 4
	Nature           string      `json:"nature"`
	EffortValues     *StatValues `json:"evs"`    // 0-255
	Gender           string      `json:"gender"` // one of M F N
	IndividualValues *StatValues `json:"ivs"`    // 0-31
	Shiny            bool        `json:"shiny"`
	Level            int         `json:"level"` // 1-100
	Happiness        int         `json:"happiness"`
	HPType           string      `json:"hpType"`
	PokeballType     string      `json:"pokeball"`
	GigantaMax       bool        `json:"gigantamax"`
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
	HP             int `json:"hp"`
	Attack         int `json:"atk"`
	Defense        int `json:"def"`
	SpecialAttack  int `json:"spa"`
	SpecialDefense int `json:"spd"`
	Speed          int `json:"spe"`
}
