package sim

import (
	"encoding/json"
	"fmt"
	"github.com/voidshard/poke-showdown-go/pkg/cmd"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	// indirection allows for easier testing
	runCommand = cmd.Run
)

// SimV1 is our struct for maintaining communication & parsing of an active battle
// simulation process.
type SimV1 struct {
	// internal channels
	stdout <-chan string
	stderr <-chan string
	errors <-chan error
	stdin  chan string
	ctrl   chan os.Signal

	players []string

	// so that we can return errors synchronously from our async read / writing we
	// buffer errors here after calling Write() to return them.
	unreadErrors     []error
	unreadErrorsLock *sync.Mutex

	// channels we supply the caller
	state    chan *structs.BattleState
	messages chan string

	// used in order to check the validity of actions
	latest *structs.BattleState
}

// NewSimV1 creates and launches a new simulation with the given Spec setup.
func NewSimV1(binary string, spec *structs.BattleSpec) (Simulation, error) {
	stdin := make(chan string)
	ctrl := make(chan os.Signal)
	state := make(chan *structs.BattleState)
	messages := make(chan string)

	stdout, stderr, errs := runCommand(
		binary,
		[]string{"simulate-battle"},
		stdin,
		ctrl,
	)

	sim := &SimV1{
		stdout:           stdout,
		stderr:           stderr,
		stdin:            stdin,
		ctrl:             ctrl,
		errors:           errs,
		state:            state,
		messages:         messages,
		players:          []string{},
		unreadErrors:     []error{},
		unreadErrorsLock: &sync.Mutex{},
	}

	go func() {
		// roll up errors into our buffer
		for err := range errs {
			sim.unreadErrorsLock.Lock()
			sim.unreadErrors = append(sim.unreadErrors, err)
			sim.unreadErrorsLock.Unlock()
		}
	}()

	// collateEvents & parse messages the simulator is writing back to us
	go sim.collateEvents()

	// pushes in initial stdin to kick off battle
	// #1 push in the battle format
	data, err := json.Marshal(map[string]interface{}{"formatid": spec.Format})
	if err != nil {
		return nil, err
	}
	stdin <- fmt.Sprintf(">start %s\n", string(data))

	// #2 for each player we need to announce them & their team in packed format
	orders := []string{}
	for player, team := range spec.Players {
		pteam := ""
		if team != nil {
			pteam = structs.PackTeam(team)
		}

		data, err = json.Marshal(map[string]interface{}{
			"name": player,
			"team": pteam,
		})
		if err != nil {
			return nil, err
		}

		// specify player's team
		stdin <- fmt.Sprintf(">player %s %s\n", player, string(data))
		sim.players = append(sim.players, player)

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
		stdin <- order
		// Nb. for some reason if we push these in too quickly pokemon-showdown
		// will not reply with a "sideupdate" for one+ of the teams leaving us
		// unsure what to do.
		// We sleep here to give the process time to absorb what we've said before
		// pushing in more instructions.
		// It only seems to be a problem with this action.
		time.Sleep(time.Millisecond * 100)
	}

	return sim, sim.latestErrors()
}

// collateEvents reads and parses relevant messages together to discover messages &
// battle updates.
func (s *SimV1) collateEvents() {
	// pokemon-showdown streams events as messages, this is cool, but really
	// we want to amass all these into "here is the state of play now"
	// so we'll collect events as they arrive & determine when we have a full
	// "turn" (ie. we need a new choice(s) from player(s) in order to move
	// to the next turn)

	state := structs.NewBattleState()
	for message := range s.stdout {
		msgs, err := parseMessage(message, state)
		if err != nil {
			s.unreadErrorsLock.Lock()
			s.unreadErrors = append(s.unreadErrors, err)
			s.unreadErrorsLock.Unlock()
			continue
		}
		for _, msg := range msgs {
			state.Messages = append(state.Messages, msg)
		}
		if state.Winner != "" || len(state.Field) == len(s.players) {
			s.latest = state
			s.state <- state
			state = structs.NewBattleState()
		}
	}
}

// parseMessage parses some message (a multi line string) into message types we
// care about.
func parseMessage(raw string, state *structs.BattleState) ([]string, error) {
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

			update, err := structs.DecodeUpdate([]byte(encoded))
			if err != nil {
				return msgs, err
			}
			state.Field[player] = update
		} else if strings.HasPrefix(lines[i], "|error|") {
			if strings.Contains(lines[i], "Can't choose for Team Preview") {
				// message indicates we gave a team preview for
				// a non-preview battle. Has no effect.
				continue
			}
			return msgs, fmt.Errorf(strings.Replace(lines[i], "|error|", "", 1))
		} else if strings.HasPrefix(lines[i], "|win|") {
			bits := strings.Split(lines[i], "|")
			state.Winner = bits[len(bits)-1]
		} else if strings.HasPrefix(lines[i], "|") {
			if strings.Count(lines[i], "|") > 1 {
				msgs = append(msgs, lines[i])
			}
		}
	}

	return msgs, nil
}

// Close stops the simulator and kills all our channels.
func (s *SimV1) Stop() {
	defer close(s.stdin)
	defer close(s.ctrl)
	defer close(s.state)
	s.ctrl <- syscall.SIGINT
}

// Read returns a read-only channel for getting new BattleState structs as they are
// written out.
func (s *SimV1) Read() <-chan *structs.BattleState {
	return s.state
}

// latestErrors returns errors from our errors buffer to the user (if any).
func (s *SimV1) latestErrors() error {
	s.unreadErrorsLock.Lock()
	defer s.unreadErrorsLock.Unlock()

	if len(s.unreadErrors) > 0 {
		root := s.unreadErrors[0]

		if len(s.unreadErrors) > 1 {
			for _, err := range s.unreadErrors[1:] {
				root = fmt.Errorf("%s: %w", root.Error(), err)
			}
		}

		s.unreadErrors = []error{}
		return root
	}

	return nil
}

// Write writes one player decision(s) for their pokemon (1 or more) to
// the simulator.
func (s *SimV1) Write(act *structs.Action) error {
	msg := act.Pack()
	s.stdin <- msg
	time.Sleep(time.Millisecond * 100)
	return s.latestErrors()
}
