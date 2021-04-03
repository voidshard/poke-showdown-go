package structs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/voidshard/poke-showdown-go/pkg/internal/data"
)

// Pokemon is a pokemon taking part in a battle
type Pokemon struct {
	// Unique `player: name` string for this pokemon
	Ident string `json:"ident"`

	Index int

	// True if the pokemon is active in the battle currently.
	Active bool `json:"active"`

	// Stats represents pokemon core stats
	Stats *data.Stats `json:"stats"`

	// Moves the pokemon can do (names only)
	Moves []string `json:"moves"`

	// BaseAbility is the pokemon ability naturally
	BaseAbility string `json:"baseAbility"`

	// Ability is the pokemon's current ability (usually same as BaseAbility
	// outside of special circumstances)
	Ability string `json:"ability"`

	// Item is what the pokemon is holding in battle
	Item string `json:"item"`

	// Pokeball the pokemon has been caught in
	Pokeball string `json:"pokeball"`

	// Status is information we parse out & attach for ease of use
	Status *Status `json:"status"`

	// Name, Level, Gender run together with ','
	// Nb. Level is not given if it is 100.
	// ie. Chesnaught, L82, M
	Details string `json:"details"`

	// From observation condition is a string like
	// 30/130
	// 0 fnt
	// 130/240 slp
	// ie. HP/MaxHap status1 status2
	Condition string `json:"condition"`

	// Species is the pokemons major species
	Species string

	// Level is the pokemons current level
	Level int

	// Shiny indicates bling
	Shiny bool

	// Gender indicates male / femaleness.
	// Nb. not all pokemon have a gender
	Gender string
}

type Status struct {
	// HPNow is the pokemons current HP
	HPNow int

	// HPMax is the pokemons full HP (if known)
	HPMax int

	// various status effects
	IsAsleep    bool
	IsBurned    bool
	IsPoisoned  bool
	IsToxiced   bool
	IsFrozen    bool
	IsParalyzed bool
	IsFainted   bool
}

func (p *Pokemon) gender() string {
	for _, chunk := range strings.Split(p.Details, ", ")[1:] {
		switch chunk {
		case "F", "M", "N":
			return chunk
		default:
			continue
		}
	}
	return ""
}

func (p *Pokemon) parseLevel() (int, error) {
	bits := strings.Split(p.Details, ", ")
	if len(bits) <= 1 {
		return -1, fmt.Errorf("unable to parse level: %s", p.Details)
	}

	for _, chunk := range bits[1:] {
		if strings.HasPrefix(chunk, "L") {
			lvl, err := strconv.ParseInt(chunk[1:], 10, 64)
			return int(lvl), err
		}
	}

	return 100, nil
}

func (p *Pokemon) species() string {
	bits := strings.Split(p.Details, ", ")
	return bits[0]
}

// IsAsleep returns if the pokemon is asleep
func (p *Pokemon) isAsleep() bool {
	return strings.Contains(p.Condition, " slp")
}

// IsFrozen returns if the pokemon is frozen
func (p *Pokemon) isFrozen() bool {
	return strings.Contains(p.Condition, " frz")
}

// IsBurned returns if the pokemon is burned
func (p *Pokemon) isBurned() bool {
	return strings.Contains(p.Condition, " brn")
}

// IsParalyzed returns if the pokemon is burned
func (p *Pokemon) isParalyzed() bool {
	return strings.Contains(p.Condition, " par")
}

// IsFainted returns if the pokemon has fainted
func (p *Pokemon) isFainted() bool {
	return strings.Contains(p.Condition, " fnt")
}

// IsPoisoned returns if the pokemon is poisoned
// (either standard poison or toxic)
func (p *Pokemon) isPoisoned() bool {
	return strings.Contains(p.Condition, " psn") || p.isToxiced()
}

// IsToxiced returns if the pokemon has been "badly" poisoned
func (p *Pokemon) isToxiced() bool {
	return strings.Contains(p.Condition, " tox")
}

// HP returns the pokemon
// - current hp
// - max hp (if known)
// Or returns an error.
// Note that if the pokemon has fainted we no longer know the max HP :(
func (p *Pokemon) parseHP() (int, int, error) {
	now, max, _, err := parseCondition(p.Condition)
	return now, max, err
}

// parseCondition parses a showdown style pokemon 'condition' string
func parseCondition(condition string) (int, int, string, error) {
	if strings.Contains(condition, "fnt") {
		return 0, -1, "fnt", nil
	}

	bits := strings.SplitN(condition, " ", 2)

	if bits[0] == "0" {
		// fainted
		return 0, -1, "fnt", nil
	}

	hpstats := strings.Split(bits[0], "/")
	if len(hpstats) != 2 {
		return -1, -1, "", fmt.Errorf("unable to read HP stats: %s [from %s]", bits[0], condition)
	}

	cur, err := strconv.ParseInt(hpstats[0], 10, 64)
	if err != nil {
		return -1, -1, "", err
	}
	max, err := strconv.ParseInt(hpstats[1], 10, 64)

	st := ""
	if len(bits) > 1 {
		st = bits[1]
	}

	return int(cur), int(max), st, err
}
