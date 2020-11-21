package structs

import (
	"fmt"
	"strings"
)

// ActionType is a user choice (move / switch)
type ActionType string

const (
	// ActionMove indicates the pokemon wishes to use an attack / status move
	ActionMove ActionType = "move"

	// ActionSwitch indicates the user wishes to switch out a pokemon to another
	// team member
	ActionSwitch ActionType = "switch"
)

// Action is a list of ActionSpecs (in order of the user pokemon)
type Action struct {
	// Player making the choices
	Player string `json:"player"`

	// Specs is the details of the decision for each pokemon the player
	// has in the battle
	Specs []*ActionSpec `json:"spec"`
}

// ActionSpec is the desired action details for a single pokemon
type ActionSpec struct {
	// move or switch
	Type ActionType `json:"type"`

	// move or pokemon slot ID
	// Technically pokemon-showdown allows a move name
	// or pokemon name, but this is much harder for us to enforce &
	// check the validity of.
	ID int `json:"id"`

	// Indicates target pokemon
	// To specify this on a non-targeting or self targeting move is invalid.
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

// Pack represents this action as a showdown simulator compliant string
func (a *Action) Pack() string {
	lines := []string{}
	for _, spec := range a.Specs {
		switch spec.Type {
		case ActionMove:
			lines = append(lines, fmt.Sprintf("move %s", packMove(spec)))
		case ActionSwitch:
			// showdown switch indexes start from 1
			lines = append(lines, fmt.Sprintf("switch %d", spec.ID+1))
		}
	}
	return fmt.Sprintf(">%s %s\n", a.Player, strings.Join(lines, ","))
}

// packMove packs a piece of an action into a showdown style movespec
func packMove(a *ActionSpec) string {
	// no move is required
	if a.Pass {
		return "pass"
	}

	target := ""
	if a.Target != 0 {
		target = fmt.Sprintf(" %d", a.Target)
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
		"%d%s%s",
		a.ID+1, // showdown move indexes start from 1
		target,
		bonus,
	)
}
