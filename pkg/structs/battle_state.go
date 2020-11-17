package structs

import ()

// BattleState is the current state of the entire battle (for all players).
type BattleState struct {
	// update represents some change to the current battle
	Field map[string]*Update

	// an error was encountered
	Error string

	// Winning player if event type is 'end'
	// Obviously implies battle is over
	Winner string

	// Extra messages for the caller
	Messages []string
}

// NewBattleState returns a new empty BattleState
func NewBattleState() *BattleState {
	return &BattleState{
		Field:    map[string]*Update{},
		Messages: []string{},
	}
}
