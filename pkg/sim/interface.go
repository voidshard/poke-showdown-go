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
	// Read returns BattleState (turn updates for each side in battle) as they
	// become available.
	Read() <-chan *structs.BattleState

	// Write commits a player choice to the engine. Note that a choice must
	// be made for each player (unless they explicitly do not need to) before a
	// new turn can happen -- which results in another BattleState
	Write(*structs.Action) error

	// Stop stops the battle and closes running processes
	Stop()
}
