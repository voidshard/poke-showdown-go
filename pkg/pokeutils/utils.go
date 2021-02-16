package pokeutils

import (
	"encoding/json"
	"fmt"
	"github.com/voidshard/poke-showdown-go/pkg/cmd"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
	"os"
	"os/exec"
	"time"
)

type Option func(*opts)

type opts struct {
	Binary string
}

func buildOpts(in []Option) *opts {
	cfg := &opts{Binary: "pokemon-showdown"}
	for _, o := range in {
		o(cfg)
	}
	return cfg
}

func Binary(path string) Option {
	return func(o *opts) {
		o.Binary = path
	}
}

func RandomTeam(in ...Option) ([]*structs.PokemonSpec, error) {
	cfg := buildOpts(in)
	team, err := randomTeam(cfg.Binary)
	return team, err
}

func randomTeam(binary string) ([]*structs.PokemonSpec, error) {
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
		pkm := []*structs.PokemonSpec{}
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
