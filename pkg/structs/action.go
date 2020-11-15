package structs

import (
	"fmt"
	"strings"
)

type ActionType string

const (
	ActionMove   ActionType = "move"
	ActionSwitch ActionType = "switch"
)

type Action struct {
	Player string `json:"player"`

	Specs []ActionSpec `json:"spec"`
}

type ActionSpec struct {
	// move or switch
	Type ActionType `json:"type"`

	// move or pokemon slot ID
	// Technically pokemon-showdown allows a move name
	// or pokemon name, but this is much harder for us to enforce &
	// check the validity of.
	ID int `json:"id"`

	// Indicates target pokemon
	Target int `json:"target"`
	// Singles:
	// Not used
	//
	// Doubles:
	// +2 +1
	// -1 -2
	//
	// Triples:
	// +3 +2 +1
	// -1 -2 -3

	// Used in doubles / triples if a slot doesn't need
	// to move (ie. a pokemon fainted or something)
	Pass bool `json:"pass"`

	// if type is `move` add some transformation first
	Mega  bool `json:"mega"`
	ZMove bool `json:"zmove"`
	Max   bool `json:"max"`
}

func (a *Action) Pack() string {
	lines := []string{}
	for _, spec := range a.Specs {
		switch spec.Type {
		case ActionMove:
			lines = append(lines, fmt.Sprintf("move %s", packMove(&spec)))
		case ActionSwitch:
			lines = append(lines, fmt.Sprintf("switch %d", spec.ID))
		}
	}
	return fmt.Sprintf(">%s %s\n", a.Player, strings.Join(lines, ","))
}

func packMove(a *ActionSpec) string {
	// no move is required
	if a.Pass {
		return "pass"
	}

	// bonus options we can append to the end
	bonus := ""
	if a.Mega {
		bonus = " mega"
	} else if a.ZMove {
		bonus = " zmove"
	} else if a.Max {
		bonus = " max"
	}

	return fmt.Sprintf(
		"%d%s",
		a.ID,
		bonus,
	)
}
