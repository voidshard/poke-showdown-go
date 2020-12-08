package sim

import (
	"encoding/json"
	"fmt"
	"github.com/voidshard/poke-showdown-go/pkg/cmd"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
	"math/rand"
	"os"
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
		// Nb. for some reason if we push these in too quickly pokemon-showdown
		// will not reply with a "sideupdate" for one+ of the teams leaving us
		// unsure what to do.
		// We sleep here to give the process time to absorb what we've said before
		// pushing in more instructions.
		// It only seems to be a problem with this action.
		time.Sleep(time.Millisecond * 100)
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
	update     *activeCmd
	events     *activeCmd
	state      chan *structs.BattleState
	incomplete chan *structs.BattleState
	spec       *structs.BattleSpec
	errors     chan error
}

func NewSimV3(binary string, spec *structs.BattleSpec) (Simulation, error) {
	if spec.Seed == 0 {
		spec.Seed = rng.Int()
	}

	update := startCmd(binary, []string{"simulate-battle"})
	err := update.start(spec)
	if err != nil {
		return nil, err
	}

	events := startCmd(binary, []string{"simulate-battle", "-S"})
	err = events.start(spec)
	if err != nil {
		return nil, err
	}

	sv3 := &SimV3{
		update: update,
		events: events,
		state:  make(chan *structs.BattleState),
		spec:   spec,
		errors: make(chan error),
	}
	go sv3.readPump()

	return sv3, nil
}

func (s *SimV3) readPump() {
	state := structs.NewBattleState()
	messages := []string{}
	messagesReady := false

	for {
		select {
		case msg := <-s.events.stdout:
			msgs, err := parseMessages(msg)
			if err != nil {
				s.events.queueError(err)
				s.errors <- err
			}
			for _, msg := range msgs {
				messages = append(messages, msg)
				if strings.Contains(msg, "|turn|") || strings.Contains(msg, "|win|") {
					messagesReady = true
				}
			}

			if messagesReady && (state.Winner != "" || len(state.Field) == len(s.spec.Players)) {
				state.Messages = messages
				s.state <- state
				state = structs.NewBattleState()
				messages = []string{}
				messagesReady = false
			}
		case msg := <-s.update.stdout:
			err := parseState(msg, state)
			if err != nil {
				s.update.queueError(err)
				s.errors <- err
			}

			if messagesReady && (state.Winner != "" || len(state.Field) == len(s.spec.Players)) {
				state.Messages = messages
				s.state <- state
				state = structs.NewBattleState()
				messages = []string{}
				messagesReady = false
			}
		}
	}
}

func parseState(raw string, state *structs.BattleState) error {
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
				return err
			}
			state.Field[player] = update
		} else if strings.HasPrefix(lines[i], "|error|") {
			if strings.Contains(lines[i], "Can't choose for Team Preview") {
				continue
			}
			return fmt.Errorf(strings.Replace(lines[i], "|error|", "", 1))
		} else if strings.HasPrefix(lines[i], "|win|") {
			bits := strings.Split(lines[i], "|")
			state.Winner = bits[len(bits)-1]
		}
	}
	return nil
}

func parseMessages(raw string) ([]string, error) {
	msgs := []string{}
	lines := strings.Split(strings.Trim(strings.TrimSpace(raw), "\x00"), "\n")
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "|request|") {
			continue
		} else if strings.HasPrefix(lines[i], "|error|") {
			if strings.Contains(lines[i], "Can't choose for Team Preview") {
				continue
			}
			return msgs, fmt.Errorf(strings.Replace(lines[i], "|error|", "", 1))
		} else if strings.HasPrefix(lines[i], "|") {
			if strings.Count(lines[i], "|") > 1 {
				msgs = append(msgs, lines[i])
			}
		}
	}
	return msgs, nil
}

// Write writes one player decision(s) for their pokemon (1 or more) to
// the simulator.
func (s *SimV3) Write(act *structs.Action) error {
	msg := act.Pack()

	s.update.stdin <- msg
	time.Sleep(time.Millisecond * 100)
	err := s.update.latestErrors()
	if err != nil {
		return err
	}

	s.events.stdin <- msg
	time.Sleep(time.Millisecond * 100)

	return s.events.latestErrors()
}

func (s *SimV3) Read() <-chan *structs.BattleState {
	return s.state
}

func (s *SimV3) Stop() {
	s.update.stop()
	s.events.stop()
}

func (s *SimV3) Errors() <-chan error {
	return s.errors
}
