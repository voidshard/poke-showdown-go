package sim

import (
	"github.com/voidshard/poke-showdown-go/pkg/structs"
)

// Start kicks off a new simulation with the given battle spec
func Start(cmd string, spec *structs.BattleSpec) (Simulation, error) {
	return NewSimV3(cmd, spec)
}

// Simulation represents an interface to the battle-simulator
type Simulation interface {
	// Turn triggers a new turn. An action is required for each player in the
	// battle simulation.
	Turn([]*structs.Action) (*structs.BattleState, error)

	// State returns the last successfully read state
	State() *structs.BattleState

	// Stop the battle and closes running processes
	Stop()
}
