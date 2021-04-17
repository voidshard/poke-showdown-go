package sim

import (
	"log"

	"github.com/voidshard/poke-showdown-go/pkg/event"
	"github.com/voidshard/poke-showdown-go/pkg/internal/parse"
)

// stream treats the showdown simulator as an bidirectional event stream
// (as it is intended).
type stream struct {
	proc  *parse.Process
	out   chan *Update
	idmap map[string]map[string]string
	spec  *BattleSpec
}

// Write some battle instruction to the simulator
func (s *stream) Write(in *Action) error {
	return s.proc.Write(in.Pack())
}

// Updates returns an event channel for outgoing updates.
// These are compressed into one chan because they're strictly in order.
func (s *stream) Updates() <-chan *Update {
	return s.out
}

// Stop closes the running process(es) & our update chan.
func (s *stream) Stop() {
	defer func() {
		// If the chan is closed (process died/exited) then we're ok with
		// that (we're in Stop() after all) - we don't want to panic here.
		if r := recover(); r != nil {
			log.Println("squashed panic during stop", r)
		}
	}()
	s.proc.Stop()
	close(s.out)
}

// NewSimulatorStream starts a new event stream & showdown process for the given
// battle spec.
func NewSimulatorStream(spec *BattleSpec) (SimulatorStream, error) {
	err := spec.validate()
	if err != nil {
		return nil, err
	}

	counts := []int{}
	teams := []string{}
	for _, t := range spec.Players {
		pstring, err := PackTeam(t)
		if err != nil {
			return nil, err
		}
		teams = append(teams, pstring)
		counts = append(counts, len(t))
	}

	proc, err := parse.NewProcess(&parse.Config{
		Seed:      spec.Seed,
		Format:    string(spec.Format),
		Teams:     teams,
		TeamSizes: counts,
	})
	if err != nil {
		proc.Stop()
		return nil, err
	}

	ss := &stream{
		proc:  proc,
		out:   make(chan *Update),
		idmap: map[string]map[string]string{},
		spec:  spec,
	}
	go func() {
		for m := range proc.Messages() {
			delta := toUpdate(m)
			ss.fillIDs(delta)
			ss.out <- delta

			if m.Event != nil {
				if m.Event.Type == event.Win {
					ss.Stop()
					return
				}
			}
		}
	}()

	return ss, nil
}

// fillIDs is where we match showdown returned pokemon to IDs
// that we accept in PokemonSpec
func (s *stream) fillIDs(u *Update) {
	if u.Side == nil {
		return
	}

	pokes, ok := s.idmap[u.Side.Player]
	if !ok {
		// if this is the first side update message then
		// the order in which we gave pokemon (spec) is
		// the order if which they're returned
		pokes = map[string]string{}
		pidx := playerIndex(u.Side.Player)

		for i, p := range u.Side.Pokemon {
			givenID := s.spec.Players[pidx][i].ID
			pokes[p.Ident] = givenID
			p.ID = givenID
		}

		s.idmap[u.Side.Player] = pokes
	} else {
		for _, p := range u.Side.Pokemon {
			givenID, _ := pokes[p.Ident]
			p.ID = givenID
		}
	}
}

// playerIndex returns a players spec index from their ID
func playerIndex(name string) int {
	switch name {
	case "p2":
		return 1
	default:
		return 0
	}
}

// toUpdate parses the given message to an update (we do this
// to have a hard break between internal & external structs).
func toUpdate(m *parse.Message) *Update {
	u := &Update{Number: m.Num}
	if m.Event != nil {
		u.Event = m.Event
	}
	if m.Error != nil {
		u.Error = m.Error
	}
	if m.Update != nil {
		u.Side = toSide(m.Update)
	}
	return u
}
