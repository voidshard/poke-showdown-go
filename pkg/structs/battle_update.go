package structs

import (
	"encoding/json"
	"strings"
)

type rawUpdate struct {
	// battle update but this player doesn't need to make a choice
	// (ie. a pokemon has been switched / fainted)
	Wait bool `json:"wait"`

	// true if pokemon are forced to switch
	ForceSwitch []bool `json:"forceSwitch"`

	// active pokemon movesets (includes additional pp/disabled data)
	Active []*activeData `json:"active"`

	// the team of pokemon
	Team Team `json:"side"`
}

type activeData struct {
	Ident string `json:"Ident"`

	Move  []*Move `json:"moves"`
	ZMove []*Move `json:"canZMove"`

	CanDynamax    bool `json:"canDynamax"`
	CanMegaEvolve bool `json:"canMegaEvo"`

	Dynamax struct {
		Moves []*Move `json:"maxMoves"`
	} `json:"maxMoves"`
}

func (a *activeData) moveHash() string {
	moveNames := []string{}
	for _, m := range a.Move {
		moveNames = append(moveNames, m.Id)
	}
	return strings.Join(moveNames, "#")
}

type Update struct {
	// battle update but this player doesn't need to make a choice
	// (ie. a pokemon has been switched / fainted)
	Wait bool `json:"wait"`

	// true if pokemon are forced to switch
	ForceSwitch []bool `json:"forceSwitch"`

	// the team of pokemon
	Team Team `json:"team"`
}

func DecodeUpdate(data []byte) (*Update, error) {
	raw := &rawUpdate{}
	err := json.Unmarshal(data, raw)
	if err != nil {
		return nil, err
	}

	// we literally have no way to identify the pokemon this is talking
	// about aside from the moves it gives us.
	// Luckily the moves are very unlikely to be both the same & in the same
	// order across two pokemon ..
	extra := map[string]*activeData{}
	if raw.Active != nil {
		for _, i := range raw.Active {
			extra[i.moveHash()] = i
		}
	}

	update := &Update{Wait: raw.Wait, ForceSwitch: raw.ForceSwitch, Team: raw.Team}
	for _, pkm := range update.Team.Pokemon {
		if !pkm.Active {
			continue
		}

		rawOpts, ok := extra[pkm.moveHash()]
		if !ok {
			continue
		}

		pkm.Options = &Options{
			CanMegaEvolve: rawOpts.CanMegaEvolve,
			CanDynamax:    rawOpts.CanDynamax,
			CanZMove:      len(rawOpts.ZMove) > 0,
			Moves:         rawOpts.Move,
			ZMoves:        rawOpts.ZMove,
			DMoves:        rawOpts.Dynamax.Moves,
		}
	}

	return update, nil
}

// PokemonOptions is a list of additional move data in order to facilitate a user
// choice.

// Move struct expands on simply the move name to include relevant in battle info.
type Move struct {
	Name   string `json:"move"`
	Target string `json:"target"`

	Id       string `json:"id"`
	PP       int    `json:"pp"`
	MaxPP    int    `json:"maxpp"`
	Disabled bool   `json:"disabled"`
}

// Team represents an entire pokemon team.
// Nb. order here is important.
type Team struct {
	Name    string     `json:"name"`
	Player  string     `json:"player"`
	Pokemon []*Pokemon `json:"pokemon"`
}

type Options struct {
	CanMegaEvolve bool
	CanDynamax    bool
	CanZMove      bool

	Moves  []*Move
	ZMoves []*Move
	DMoves []*Move
}

// Pokemon is a pokemon taking part in a battle
type Pokemon struct {
	Ident string `json:"ident"`

	Active  bool     `json:"active"`
	Options *Options `json:"options"`

	Stats       *StatValues `json:"stats"`
	Moves       []string    `json:"moves"`
	BaseAbility string      `json:"baseAbility"`
	Ability     string      `json:"ability"`
	Item        string      `json:"item"`
	Pokeball    string      `json:"pokeball"`

	// From observation condition is a string like
	// 30/130
	// 0 fnt
	// 130/240 slp
	// ie. HP/MaxHap status1 status2
	Condition string `json:"condition"`

	// Name, Level, Gender run together with ','
	// ie. Chesnaught, L82, M
	Details string `json:"details"`
}

func (p *Pokemon) moveHash() string {
	return strings.Join(p.Moves, "#")
}
