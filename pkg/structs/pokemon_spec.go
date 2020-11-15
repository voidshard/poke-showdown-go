package structs

import (
	"fmt"
	"strings"
)

// PokemonSpec represents a pokemon with battle relevant stats, items, moves etc as specified in order to *start* a battle.
// Pokemon returned during battle have different derived fields.
type PokemonSpec struct {
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

// PackTeam turns a list of PokemonSpec into a pokemon-showdown compliant `packed` string.
func PackTeam(in []*PokemonSpec) string {
	members := []string{}
	for _, p := range in {
		members = append(members, p.Pack())
	}
	return strings.Join(members, "]")
}

// Pack a battlemon into the pokemon simulator format
func (b *PokemonSpec) Pack() string {
	// nb. the battle simulator packs things more tightly than we do.
	// Essentially if an EV or IV is a default value (0 or 31) it can be
	// omitted (simply "").
	// If all values are "" then even the ',' symbols can be omitted,
	// rendering an empty "" for all six values (meaning "all are default".
	// This isn't really required for us so I'm omitting this.
	// Note however that this means *not* giving any IVs means they're all
	// considered to be 31 NOT 0.
	// Additionally the final values (happiness, hptype, ball, gigantamax)
	// are also entirely omitted if they're the default values.
	// We always include these values because it's much easier and we're
	// not sending packed data over the network.

	packedIvs := "31,31,31,31,31,31"
	if b.IndividualValues != nil {
		packedIvs = fmt.Sprintf(
			"%d,%d,%d,%d,%d,%d",
			b.IndividualValues.HP,
			b.IndividualValues.Attack,
			b.IndividualValues.Defense,
			b.IndividualValues.SpecialAttack,
			b.IndividualValues.SpecialDefense,
			b.IndividualValues.Speed,
		)
	}

	packedEvs := "85,85,85,85,85,85"
	if b.EffortValues != nil {
		packedEvs = fmt.Sprintf(
			"%d,%d,%d,%d,%d,%d",
			b.EffortValues.HP,
			b.EffortValues.Attack,
			b.EffortValues.Defense,
			b.EffortValues.SpecialAttack,
			b.EffortValues.SpecialDefense,
			b.EffortValues.Speed,
		)
	}

	packedSpecies := b.Species
	if b.Species == b.Name {
		packedSpecies = ""
	}

	return strings.Join(
		[]string{
			b.Name,
			packedSpecies,
			b.Item,
			b.Ability,
			strings.Join(b.Moves, ","),
			b.Nature,
			packedEvs,
			b.Gender,
			packedIvs,
			packbool(b.Shiny, "S"),
			fmt.Sprintf("%d", b.Level),
			fmt.Sprintf(
				"%d,%s,%s,%s",
				b.Happiness, // max is 255
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
