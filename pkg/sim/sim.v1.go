package sim

import (
	"encoding/json"
	"github.com/voidshard/poke-showdown-go/pkg/cmd"
	"os/signal"
)

type SimV1 struct {
	stdout chan string
	stderr chan string
	stdin  chan string
	ctrl   chan os.Signal
	errors chan string
}

func NewSimV1(binary string, spec *BattleSpec) (Simulation, error) {
	stdin := make(chan string)
	ctrl := make(chan os.Signal)

	stdout, stderr, errs := cmd.Run(
		binary,
		[]string{"simulate-battle"},
		stdin,
		ctrl,
	)

	sim := &SimV1{
		stdout: stdout,
		stderr: stderr,
		stdin:  stdin,
		ctrl:   ctrl,
		errors: errs,
	}

	go func() { // push in initial cmds to kick off battle
		// #1 push in the battle format
		stdin <- fmt.Sprintf(">start %s\n", sim.encodeJSON(map[string]string{"formatid": spec.Format}))

		// #2 for each player we need to announce them & their team in packed format
		for player, team := range spec.Players {
			stdin <- fmt.Sprintf(
				">player %s %s\n",
				player,
				sim.encodeJSON(map[string]interface{}{"name": player, "team": team}),
			)
		}
	}()

	return sim, nil
}

func (s *SimV1) encodeJSON(msg map[string]interface{}) string {
	data, err := json.Marshal(msg)
	if err != nil {
		s.errors <- err
		return ""
	}
	return string(data)
}

func (s *SimV1) Read() <-chan *Event {

}

func (s *SimV1) Write(act *Action) error {
}
