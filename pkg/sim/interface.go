package sim

import (
	"github.com/voidshard/pkg/structs"
)

type Simulator interface {
	Start(string, *BattleSpec) (Simulation, error)
}

type Simulation interface {
	Read() chan *Event
	Write(*Action) error
}

type Action struct {
	Player string
}

type Event struct {
}

type BattleSpec struct {
	Format string

	Players map[string]*structs.Battlemon
}
