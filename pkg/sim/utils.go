package sim

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/voidshard/poke-showdown-go/pkg/structs"
)

var (
	// ErrFailedTurn is wrapped around errors returned from simulate
	ErrFailedTurn = fmt.Errorf("failed to apply turn")

	// ErrTimeout is used if the simulator blocks for too long
	ErrTimeout = fmt.Errorf("timeout (check input?)")

	// EnvPokemonShowdown is the env var we check for a pokemon-showdown.
	// Defaults to `pokemon-showdown` if not given.
	EnvPokemonShowdown = "PATH_POKEMON_SHOWDOWN"

	// EnvSimulationTimeout dictates how long a sim should last
	// before we declare "timeout"
	// This is a Golang duration for ParseDuration see golang.org/pkg/time/#ParseDuration
	// Defaults to 2s if not given or unable to parse.
	EnvSimulationTimeout = "SIMULATION_TIMEOUT"
)

// timeoutDuration returns how long the simulation should last for
// before raising a timeout error.
func timeoutDuration() time.Duration {
	timeout := os.Getenv(EnvSimulationTimeout)
	if timeout == "" {
		return time.Second * 2
	}

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		log.Printf("err parsing duration '%s' using 2s instead: %v\n", timeout, err)
		return time.Second * 2
	}
	return duration
}

// Simulate runs pokemon-showdown, passing in the given set of actions as turn(s)
// and returning the final battle state.
func Simulate(spec *structs.BattleSpec, actions []*structs.Action) (*structs.BattleState, error) {
	binary := os.Getenv(EnvPokemonShowdown)
	if binary == "" {
		binary = "pokemon-showdown"
	}

	// prep the simulator
	battle, err := Start(binary, spec)
	if err != nil {
		return nil, err
	}
	defer battle.Stop()

	// bucket actions into turns
	turns := toTurns(actions)

	// prep a timeout in case actions are incomplete
	done := make(chan error)
	tick := time.NewTicker(timeoutDuration())
	defer tick.Stop()

	// push the turn(s) into the simulator
	go func() {
		for _, acts := range turns {
			if acts == nil {
				continue
			}

			_, err := battle.Turn(acts)
			if err != nil {
				done <- fmt.Errorf("%w: %v", ErrFailedTurn, err)
				return
			}
		}

		done <- nil
	}()

	select {
	case err = <-done:
		return battle.State(), err
	case <-tick.C:
		// probably implies the simulator was waiting for input,
		// which would imply the action(s) were out of order or
		// simply incomplete
		return nil, ErrTimeout
	}

	// we shouldn't reach here but .. anyway
	return battle.State(), nil
}

// toTurns breaks a list of actions that are in order into valid "turn" choices.
func toTurns(actions []*structs.Action) [][]*structs.Action {
	// strictly speaking, this doesn't precisely know what "turn" an action
	// is done on, since things like switches can happen mid turn or be proper
	// turns in their own right (eg. knockouts vs actively switching).
	// What we really care about though is that each of these pseudo "turns"
	// makes the simulator write out the next state .. in the end
	// we'll feed in all turns and get to the desired end state (if possible)
	counts := map[string]int{}
	result := [][]*structs.Action{}

	for _, act := range actions {
		playerTurn := counts[act.Player]

		if len(result) <= playerTurn {
			result = append(result, []*structs.Action{})
		}

		turnActs := result[playerTurn]
		turnActs = append(turnActs, act)
		result[playerTurn] = turnActs

		counts[act.Player] = playerTurn + 1
	}

	return result
}
