package structs

import ()

// BattleState is the current state of the entire battle (for all players).
type BattleState struct {
	// Field is a map of player -> Update that details the current state
	// of the given players field.
	Field map[string]*Update

	// Winning player if event type is 'end'
	// Obviously implies battle is over
	Winner string

	Messages []string
	Events   []*Event
}

// NewBattleState returns a new empty BattleState
func NewBattleState() *BattleState {
	return &BattleState{
		Field:    map[string]*Update{},
		Messages: []string{},
		Events:   []*Event{},
	}
}

// Event represents some event that has happened during a turn
type Event struct {
	Type string
}
