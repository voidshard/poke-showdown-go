package sim

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

	// ActionPass is used if, for example, a slot cannot be used.
	// Ie. in doubles where you have one pokemon left.
	ActionPass ActionType = "pass"
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
	ID string `json:"id"`

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
		case ActionPass:
			lines = append(lines, "pass")
		case ActionMove, "": // default to move if not given
			lines = append(lines, fmt.Sprintf("move %s", packMove(spec)))
		case ActionSwitch:
			lines = append(lines, fmt.Sprintf("switch %s", spec.ID))
		}
	}
	return fmt.Sprintf(">%s %s\n", a.Player, strings.Join(lines, ","))
}

// packMove packs a piece of an action into a showdown style movespec
func packMove(a *ActionSpec) string {
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
		"%s%s%s",
		a.ID,
		target,
		bonus,
	)
}
