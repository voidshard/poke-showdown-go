package sim

import (
	"github.com/voidshard/poke-showdown-go/pkg/structs"
)

func Start(cmd string, spec *structs.BattleSpec) (Simulation, error) {
	return NewSimV1(cmd, spec)
}

type Simulation interface {
	Read() <-chan *structs.BattleState
	Messages() <-chan string
	Errors() <-chan error
	Write(*structs.Action) error
	Close()
}
