package pokeutils

import (
	"encoding/json"
	"fmt"
	"github.com/voidshard/poke-showdown-go/pkg/internal/cmd"
	"github.com/voidshard/poke-showdown-go/pkg/sim"
	"os"
	"os/exec"
	"time"
)

// Option is some option for RandomTeam
type Option func(*opts)

// opts is our function inputs
type opts struct {
	Binary string
}

// buildOpts turns options into an opts struct
func buildOpts(in []Option) *opts {
	cfg := &opts{Binary: "pokemon-showdown"}
	for _, o := range in {
		o(cfg)
	}
	return cfg
}

// Binary allows one to set the path to pokemon-showdown
func Binary(path string) Option {
	return func(o *opts) {
		o.Binary = path
	}
}

// RandomTeam uses pokemon-showdown to generate a random team of pokemon.
func RandomTeam(in ...Option) ([]*sim.PokemonSpec, error) {
	cfg := buildOpts(in)
	team, err := randomTeam(cfg.Binary)
	return team, err
}

// randomTeam uses pokemon-showdown to generate a random team of pokemon.
func randomTeam(binary string) ([]*sim.PokemonSpec, error) {
	out, err := exec.Command(binary, "generate-team").Output()
	if err != nil {
		return nil, err
	}

	stdin := make(chan string)
	cntrl := make(chan os.Signal)

	sout, serr, errs := cmd.Run(
		binary, []string{"unpack-team"}, stdin, cntrl, cmd.Seperator("\n"),
	)
	go func() {
		stdin <- string(out) + "\n\n"
	}()

	timer := time.NewTimer(time.Second)

	select {
	case data := <-sout:
		pkm := []*sim.PokemonSpec{}
		err = json.Unmarshal([]byte(data), &pkm)
		return pkm, err
	case data := <-serr:
		return nil, fmt.Errorf("failed to unpack team: %s", string(data))
	case err := <-errs:
		return nil, err
	case <-timer.C:
		return nil, fmt.Errorf("time out unpacking team")
	}

	return nil, fmt.Errorf("failed to generate team")
}
