package sim

import (
	"encoding/json"
	"fmt"
	"github.com/voidshard/poke-showdown-go/pkg/cmd"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
	"os"
	"strings"
	"syscall"
	"time"
)

type SimV1 struct {
	// internal channels
	stdout <-chan string
	stderr <-chan string
	errors <-chan error
	stdin  chan string
	ctrl   chan os.Signal

	players []string

	// channels we supply the caller
	allErrs  chan error
	state    chan *structs.BattleState
	messages chan string

	// used in order to check the validity of actions
	latest *structs.BattleState
}

func NewSimV1(binary string, spec *structs.BattleSpec) (Simulation, error) {
	stdin := make(chan string)
	ctrl := make(chan os.Signal)
	state := make(chan *structs.BattleState)
	messages := make(chan string)
	allErrs := make(chan error)

	stdout, stderr, errs := cmd.Run(
		binary,
		[]string{"simulate-battle"},
		stdin,
		ctrl,
	)

	sim := &SimV1{
		stdout:   stdout,
		stderr:   stderr,
		stdin:    stdin,
		ctrl:     ctrl,
		errors:   errs,
		allErrs:  allErrs,
		state:    state,
		messages: messages,
		players:  []string{},
	}

	go func() {
		// pushes errors from command back to caller
		// (We do this so we can add errors at this level too)
		for err := range errs {
			allErrs <- err
		}
	}()

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

	for _, order := range orders {
		stdin <- order
		// Nb. for some reason if we push these in too quickly pokemon-showdown
		// will not reply with a "sideupdate" for one+ of the teams leaving us
		// unsure what to do.
		// We sleep here to give the process time to absorb what we've said before
		// pushing in more instructions.
		// It only seems to be a problem with this action.
		time.Sleep(time.Millisecond * 300)
	}

	return sim, nil
}

//
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
			s.allErrs <- err
			continue
		}
		for _, msg := range msgs {
			s.messages <- msg
		}
		if state.Winner != "" || len(state.Field) == len(s.players) {
			s.latest = state
			s.state <- state
			state = structs.NewBattleState()
		}
	}
}

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

			update := &structs.Update{}
			err := json.Unmarshal([]byte(encoded), update)
			if err != nil {
				return msgs, err
			}
			state.Field[player] = update
		} else if strings.HasPrefix(lines[i], "|error|") {
			state.Error = lines[i]
		} else if strings.HasPrefix(lines[i], "|win|") {
			bits := strings.Split(lines[i], "|")
			state.Winner = bits[len(bits)-1]
		} else if strings.HasPrefix(lines[i], "|") {
			msgs = append(msgs, lines[i])
		}
	}

	return msgs, nil
}

func (s *SimV1) Close() {
	defer close(s.stdin)
	defer close(s.ctrl)
	defer close(s.state)
	defer close(s.allErrs)
	s.ctrl <- syscall.SIGINT
}

func (s *SimV1) Read() <-chan *structs.BattleState {
	return s.state
}

func (s *SimV1) Messages() <-chan string {
	return s.messages
}

func (s *SimV1) Write(act *structs.Action) error {
	// TODO: check validty of action
	msg := act.Pack()
	s.stdin <- msg
	return nil
}

func (s *SimV1) Errors() <-chan error {
	return s.allErrs
}
