package sim

import (
	"fmt"
	"strings"

	data "github.com/voidshard/poke-showdown-go/pkg/pokedata"
)

// PokemonSpec represents a pokemon with battle relevant stats, items, moves etc as specified in order to *start* a battle.
// Pokemon returned during battle have different derived fields.
type PokemonSpec struct {
	Name             string   `json:"name"`
	Species          string   `json:"species"`
	Item             string   `json:"item"`
	Ability          string   `json:"ability"`
	Moves            []string `json:"moves"` // max len 4
	Nature           string   `json:"nature"`
	EffortValues     *Stats   `json:"evs"`    // 0-255
	Gender           string   `json:"gender"` // one of M F N
	IndividualValues *Stats   `json:"ivs"`    // 0-31
	Level            int      `json:"level"`  // 1-100
	Happiness        int      `json:"happiness"`
	HPType           string   `json:"hpType"`
	PokeballType     string   `json:"pokeball"`
	GigantaMax       bool     `json:"gigantamax"`
}

// enforceLimits clamps down on int values so they're within acceptable ranges
func (b *PokemonSpec) enforceLimits() {
	if b.Moves == nil || len(b.Moves) == 0 {
		b.Moves = []string{}
	} else if len(b.Moves) > 4 {
		b.Moves = b.Moves[:4]
	}

	if b.EffortValues != nil {
		b.EffortValues.HP = clamp(0, 255, b.EffortValues.HP)
		b.EffortValues.Attack = clamp(0, 255, b.EffortValues.Attack)
		b.EffortValues.Defense = clamp(0, 255, b.EffortValues.Defense)
		b.EffortValues.SpecialAttack = clamp(0, 255, b.EffortValues.SpecialAttack)
		b.EffortValues.SpecialDefense = clamp(0, 255, b.EffortValues.SpecialDefense)
		b.EffortValues.Speed = clamp(0, 255, b.EffortValues.Speed)

		if b.EffortValues.Sum() > 510 {
			b.EffortValues.HP = 85
			b.EffortValues.Attack = 85
			b.EffortValues.Defense = 85
			b.EffortValues.SpecialAttack = 85
			b.EffortValues.SpecialDefense = 85
			b.EffortValues.Speed = 85
		}
	}
	if b.IndividualValues != nil {
		b.IndividualValues.HP = clamp(0, 31, b.IndividualValues.HP)
		b.IndividualValues.Attack = clamp(0, 31, b.IndividualValues.Attack)
		b.IndividualValues.Defense = clamp(0, 31, b.IndividualValues.Defense)
		b.IndividualValues.SpecialAttack = clamp(0, 31, b.IndividualValues.SpecialAttack)
		b.IndividualValues.SpecialDefense = clamp(0, 31, b.IndividualValues.SpecialDefense)
		b.IndividualValues.Speed = clamp(0, 31, b.IndividualValues.Speed)
	}

	switch b.Gender {
	case "M", "F", "N":
		break
	default:
		b.Gender = ""
	}

	b.Level = clamp(1, 100, b.Level)
	b.Happiness = clamp(0, 255, b.Happiness)
}

// clamp makes an int between two given min, max values
func clamp(min, max, value int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// PackTeam turns a list of PokemonSpec into a pokemon-showdown compliant `packed` string.
func PackTeam(in []*PokemonSpec) (string, error) {
	members := []string{}
	for _, p := range in {
		pokestring, err := p.Pack()
		if err != nil {
			return "", err
		}
		members = append(members, pokestring)
	}
	return strings.Join(members, "]"), nil
}

// Pack a battlemon into the pokemon simulator format
func (b *PokemonSpec) Pack() (string, error) {
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
	b.enforceLimits()

	// --- basic checks
	if len(b.Moves) == 0 {
		return "", fmt.Errorf("at least one move required")
	} else if b.Species == "" && b.Name == "" {
		return "", fmt.Errorf("one of name/species is required")
	}

	// --- check pokemon name / species
	var dex *data.PokeDexItem
	sdex, serr := data.PokeDex(b.Species)
	ndex, nerr := data.PokeDex(b.Name)
	if serr != nil && nerr != nil {
		// neither name or species was found
		return "", fmt.Errorf("no pokemon by species [%v] or name [%v]", serr, nerr)
	} else if serr == nil && nerr == nil {
		// both name & species yielded pokemon
		if sdex.Number != ndex.Number {
			return "", fmt.Errorf(
				"name & species indicate different pokemon %s (%s) %s (%s)",
				b.Species, sdex.Name, b.Name, ndex.Name,
			)
		}
		dex = sdex
	} else if serr != nil {
		// found by species
		dex = sdex
	} else {
		// found by name
		dex = ndex
	}

	// --- check pokemon ability
	found := false
	for _, ab := range dex.Abilities {
		if b.Ability == "" {
			// if not given, set first ability we see
			// (order of dict iteration is undefined, so this is random)
			b.Ability = ab
		}
		found = data.Strip(b.Ability) == data.Strip(ab)
		if found {
			break
		}
	}
	if !found {
		// ability doesn't exist or this pokemon can't have this ability
		return "", fmt.Errorf("pokemon %s cannot get ability %s", b.Species, b.Ability)
	}

	// --- check pokemon move(s)
	for _, moveID := range b.Moves {
		// check that each move can be found
		_, err := data.MoveDex(moveID)
		if err != nil {
			return "", err
		}
	}

	// --- check nature
	if b.Nature != "" && !ValidNature(b.Nature) {
		return "", fmt.Errorf("no nature found matching %s", b.Nature)
	}

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
			"", // isShiny
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
	), nil
}

// packbool turns a given bool in a blank or the given value
func packbool(b bool, value string) string {
	if b {
		return value
	}
	return ""
}

// Stats holds values for the 6 pokemon stats
type Stats struct {
	HP             int `json:"hp"`
	Attack         int `json:"atk"`
	Defense        int `json:"def"`
	SpecialAttack  int `json:"spa"`
	SpecialDefense int `json:"spd"`
	Speed          int `json:"spe"`
}

// Sum returns the total of all of the stats
func (s *Stats) Sum() int {
	return s.HP + s.Attack + s.Defense + s.SpecialAttack + s.SpecialDefense + s.Speed
}
