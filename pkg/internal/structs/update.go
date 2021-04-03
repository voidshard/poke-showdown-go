package structs

import (
	"encoding/json"
	"strings"
)

// DecodeUpdate turns a raw simulater update into our finished side update
func DecodeUpdate(data []byte) (*Update, error) {
	raw := &Update{}
	err := json.Unmarshal(data, raw)
	if err != nil {
		return nil, err
	}

	for i, pkm := range raw.Team.Pokemon {
		curHP, maxHP, err := pkm.parseHP()
		if err != nil {
			return nil, err
		}

		pkm.Index = i

		pkm.Status = &Status{
			// helpfully pre-parse fields for user(s)
			HPNow:       curHP,
			HPMax:       maxHP,
			IsAsleep:    pkm.isAsleep(),
			IsBurned:    pkm.isBurned(),
			IsPoisoned:  pkm.isBurned(),
			IsToxiced:   pkm.isToxiced(),
			IsFrozen:    pkm.isFrozen(),
			IsParalyzed: pkm.isParalyzed(),
			IsFainted:   pkm.isFainted(),
		}
		pkm.Species = pkm.species()
		pkm.Gender = pkm.gender()

		lvl, err := pkm.parseLevel()
		if err != nil {
			return nil, err
		}
		pkm.Level = lvl

		pkm.Shiny = strings.Contains(pkm.Details, ", shiny")
	}

	return raw, nil
}

// Update is the native format of showdown `sideupdate`
// Nb. we change the name of this to 'Side' in outer layers
type Update struct {
	// battle update but this player doesn't need to make a choice
	// (ie. a pokemon has been switched / fainted)
	Wait bool `json:"wait"`

	// true if pokemon are forced to switch
	ForceSwitch []bool `json:"forceSwitch"`

	// active pokemon movesets (includes additional pp/disabled data)
	Active []*ActiveData `json:"active"`

	// the team of pokemon
	Team Team `json:"side"`
}

// ActiveData is the native format of showdown `active` (Options)
type ActiveData struct {
	Moves  []*Move `json:"moves"`
	ZMoves []*Move `json:"canZMove"`

	CanDynamax    bool `json:"canDynamax"`
	CanMegaEvolve bool `json:"canMegaEvo"`

	Dynamax struct {
		Moves []*Move `json:"maxMoves"`
	} `json:"maxMoves"`

	Trapped bool `json:"trapped"`
}

// Team represents an entire pokemon team.
// Nb. order here is important.
type Team struct {
	Player  string     `json:"name"`
	Pokemon []*Pokemon `json:"pokemon"`
}
