package sim

import (
	"encoding/json"
	"fmt"
	"github.com/voidshard/poke-showdown-go/pkg/cmd"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	// indirection allows for easier testing
	runCommand = cmd.Run

	// random number generator
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type activeCmd struct {
	stdout           <-chan string
	stderr           <-chan string
	errors           <-chan error
	stdin            chan string
	ctrl             chan os.Signal
	unreadErrors     []error
	unreadErrorsLock *sync.Mutex
}

// stop kills the cmd
func (a *activeCmd) stop() {
	defer close(a.stdin)
	defer close(a.ctrl)
	a.ctrl <- syscall.SIGINT
}

func (a *activeCmd) start(spec *structs.BattleSpec) error {
	// pushes in initial stdin to kick off battle
	// #1 push in the battle format
	data, err := json.Marshal(map[string]interface{}{
		"seed":     []int{spec.Seed, spec.Seed, spec.Seed, spec.Seed},
		"formatid": spec.Format,
	})
	if err != nil {
		return err
	}
	a.stdin <- fmt.Sprintf(">start %s\n", string(data))

	// #2 for each player we need to announce them & their team in packed format
	orders := []string{}
	for num, team := range spec.Players {
		player := fmt.Sprintf("p%d", num+1)

		pteam := ""
		if team != nil {
			pteam = structs.PackTeam(team)
		}

		data, err = json.Marshal(map[string]interface{}{
			"name": player,
			"team": pteam,
		})
		if err != nil {
			return err
		}

		// specify player's team
		a.stdin <- fmt.Sprintf(">player %s %s\n", player, string(data))

		if team == nil {
			continue
		}

		// specify the players team order (we give the order in which we got them).
		members := []string{}
		for i := range team {
			members = append(members, fmt.Sprintf("%d", i+1))
		}
		orders = append(orders, fmt.Sprintf(">%s team %s\n", player, strings.Join(members, ",")))
	}

	// #3 we now need to tell the simulator what order the player's team should
	// be in (ie, battle order).
	for _, order := range orders {
		a.stdin <- order
	}

	return a.latestErrors()
}

// latestErrors returns errors from our errors buffer to the user (if any).
func (a *activeCmd) latestErrors() error {
	a.unreadErrorsLock.Lock()
	defer a.unreadErrorsLock.Unlock()

	if len(a.unreadErrors) > 0 {
		root := a.unreadErrors[0]

		if len(a.unreadErrors) > 1 {
			for _, err := range a.unreadErrors[1:] {
				root = fmt.Errorf("%s: %w", root.Error(), err)
			}
		}

		a.unreadErrors = []error{}
		return root
	}

	return nil
}

func startCmd(cmd string, args []string) *activeCmd {
	stdin := make(chan string)
	ctrl := make(chan os.Signal)
	stdout, stderr, errs := runCommand(cmd, args, stdin, ctrl)

	sim := &activeCmd{
		stdout:           stdout,
		stderr:           stderr,
		errors:           errs,
		stdin:            stdin,
		ctrl:             ctrl,
		unreadErrors:     []error{},
		unreadErrorsLock: &sync.Mutex{},
	}

	go func() {
		// roll up errors into our buffer
		for err := range errs {
			sim.queueError(err)
		}
	}()

	return sim
}

func (a *activeCmd) queueError(err error) {
	a.unreadErrorsLock.Lock()
	defer a.unreadErrorsLock.Unlock()
	a.unreadErrors = append(a.unreadErrors, err)
}

type SimV3 struct {
	process *activeCmd
	spec    *structs.BattleSpec
	errors  chan error
	current *structs.BattleState
}

func NewSimV3(binary string, spec *structs.BattleSpec) (Simulation, error) {
	if spec.Seed == 0 {
		spec.Seed = rng.Int()
	}

	process := startCmd(binary, []string{"simulate-battle"})
	err := process.start(spec)
	if err != nil {
		return nil, err
	}

	sv3 := &SimV3{
		process: process,
		spec:    spec,
		errors:  make(chan error),
	}

	_, err = sv3.readState()
	return sv3, err
}

func (s *SimV3) State() *structs.BattleState {
	return s.current
}

func (s *SimV3) readState() (*structs.BattleState, error) {
	state := structs.NewBattleState()
	events := []*structs.Event{}
	turnOver := false

	for msg := range s.process.stdout {
		msgs, err := parseStdout(msg, state)
		if err != nil {
			return nil, err
		}
		for _, msg := range msgs {
			evt := structs.ParseEvent(msg)
			if evt != nil {
				events = append(events, evt)
			}
		}

		winner := state.Winner != ""
		turnOver = state.Turn > -1 || winner
		complete := len(state.Field) == len(s.spec.Players)

		forceSwitch := false
		if complete {
			for _, f := range state.Field {
				for _, pos := range f.ForceSwitch {
					if pos {
						forceSwitch = true
					}
				}
			}
		}

		if (turnOver || forceSwitch) && (winner || complete) {
			state.Events = events
			if state.Turn == -1 {
				state.Turn = s.current.Turn
			}

			s.current = state
			return state, nil
		}
	}

	return nil, nil
}

func parseStdout(raw string, state *structs.BattleState) ([]string, error) {
	// requires showdown version 0.11.4+ to fix a bug where messages are not returned
	msgs := []string{}
	lines := strings.Split(strings.Trim(strings.TrimSpace(raw), "\x00"), "\n")
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "|request|") {
			player := lines[i-1] // player name preceeds |request| line
			bits := strings.Split(lines[i], "|")

			encoded := bits[len(bits)-1]
			if strings.Contains(encoded, "\"teamPreview\":true") {
				// Specifically ignore preview messages.
				continue
			}

			// the simulator is asking a player to make a choice
			update, err := structs.DecodeUpdate([]byte(encoded))
			if err != nil {
				return msgs, err
			}
			state.Field[player] = update
		} else if strings.HasPrefix(lines[i], "|switch|") || strings.HasPrefix(lines[i], "|-damage|") || strings.HasPrefix(lines[i], "|-heal|") {
			// Nb. the simulator returns duplicates of some messages
			// - in particular switch, damage and heal messages
			// can be replayed with both actual health & percentage
			// (/100) values. This is super annoying and it's easy
			// to figure out a percentage given the health values.
			// So here we strip out the duplicates ..
			bitsthis := strings.Split(lines[i], "|")

			// if it's not the same event type as previous
			if !strings.HasPrefix(lines[i-1], "|"+bitsthis[1]) {
				msgs = append(msgs, lines[i])
				continue
			}

			bitsprev := strings.Split(lines[i-1], "|")
			// if it is the same type as previous & about the same pokemon
			if bitsthis[2] == bitsprev[2] {
				continue
			}
		} else if strings.HasPrefix(lines[i], "|error|") {
			if strings.Contains(lines[i], "Can't choose for Team Preview") {
				// we always give a team layout. If it's not valid
				// for the format then whatever - it does no harm.
				continue
			}
			return msgs, fmt.Errorf(strings.Replace(lines[i], "|error|", "", 1))
		} else if strings.HasPrefix(lines[i], "|win|") {
			bits := strings.Split(lines[i], "|")
			state.Winner = bits[len(bits)-1]
		} else if strings.HasPrefix(lines[i], "|turn|") {
			bits := strings.Split(lines[i], "|")
			i, err := strconv.ParseInt(bits[len(bits)-1], 10, 64)
			if err != nil {
				return msgs, err
			}
			state.Turn = int(i)
		} else if strings.HasPrefix(lines[i], "|split|") {
			continue
		} else if strings.HasPrefix(lines[i], "|start") {
			continue
		} else if strings.HasPrefix(lines[i], "|poke|") {
			continue
		} else if strings.HasPrefix(lines[i], "|t:|") {
			continue
		} else if strings.HasPrefix(lines[i], "|player|") {
			continue
		} else if strings.HasPrefix(lines[i], "|teamsize|") {
			continue
		} else if strings.HasPrefix(lines[i], "|gametype|") {
			continue
		} else if strings.HasPrefix(lines[i], "|gen|") {
			continue
		} else if strings.HasPrefix(lines[i], "|tier|") {
			continue
		} else if strings.HasPrefix(lines[i], "|rule|") {
			continue
		} else if strings.HasPrefix(lines[i], "|") {
			if strings.Count(lines[i], "|") > 1 {
				msgs = append(msgs, lines[i])
			}
		}
	}
	return msgs, nil
}

// Turn supplies all decisions needed to move to the next turn
func (s *SimV3) Turn(decisions []*structs.Action) (*structs.BattleState, error) {
	for _, decision := range decisions {
		err := s.write(decision)
		if err != nil {
			return nil, err
		}
	}

	return s.readState()
}

// Write writes one player decision(s) for their pokemon (1 or more) to
// the simulator.
func (s *SimV3) write(act *structs.Action) error {
	msg := act.Pack()
	s.process.stdin <- msg
	return s.process.latestErrors()
}

// Stop simulation & kill subprocess(es)
func (s *SimV3) Stop() {
	s.process.stop()
}
