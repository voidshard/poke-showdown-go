package parse

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/voidshard/poke-showdown-go/pkg/event"
	"github.com/voidshard/poke-showdown-go/pkg/internal/cmd"
	"github.com/voidshard/poke-showdown-go/pkg/internal/structs"
)

var (
	// random number generator
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	// indirection allows for easier testing
	runCommand = cmd.Run

	// ErrInvalidChoice indicates that the given choice is not possible
	// ie. switching to an already active pokemon
	ErrInvalidChoice = fmt.Errorf("invalid choice see github.com/smogon/pokemon-showdown/blob/master/sim/SIM-PROTOCOL.md")

	// ErrUnavailableChoice indicates that the given (usually valid) choice
	// cannot be done for some reason (ie. move is disabled)
	ErrUnavailableChoice = fmt.Errorf("unavailable choice https://github.com/smogon/pokemon-showdown/blob/master/sim/SIM-PROTOCOL.md")
)

// Process wraps an active pokemon showdown process
type Process struct {
	// underlying process (mostly a holder for it's stdout, strerr, stdin)
	raw *activeCmd

	// output messages from simulator
	messages chan *Message
}

// Config holds settings relevant to kicking off a showdown battle
type Config struct {
	Binary    string
	Seed      int
	Format    string
	Teams     []string
	TeamSizes []int
}

// defaults sets unset fields
func (c *Config) defaults() {
	if c.Seed == 0 {
		c.Seed = rng.Int()
	}
	if c.Binary == "" {
		c.Binary = "pokemon-showdown"
	}
	if c.Format == "" {
		c.Format = "[Gen 8] Anything Goes"
	}
}

// NewProcess starts a new battle simulation
func NewProcess(cfg *Config) (*Process, error) {
	cfg.defaults()

	stdin := make(chan string)
	ctrl := make(chan os.Signal)
	stdout, stderr, errs := runCommand(
		cfg.Binary,
		[]string{"simulate-battle"},
		stdin,
		ctrl,
	)

	raw := &activeCmd{
		stdout: stdout,
		stderr: stderr,
		errors: errs,
		stdin:  stdin,
		ctrl:   ctrl,
	}

	proc := &Process{
		raw:      raw,
		messages: make(chan *Message),
	}

	go func() {
		// push errors from the raw handler into our error chan
		for err := range errs {
			log.Printf("err: %v\n", err)
			proc.messages <- message(err)
		}
	}()

	go func() {
		// read messages from raw stdout and parse into showdown structs
		// results are pushed into process chans
		for msg := range proc.raw.stdout {
			log.Printf("msg: %s\n", msg)
			proc.parseStdout(msg)
		}
	}()

	err := raw.start(cfg.Seed, cfg.Format, cfg.TeamSizes, cfg.Teams)
	return proc, err
}

// parseStdout reads messages from the showdown stdout and parses them into events,
// errors or side updates (as applicable). Not all lines that are printed are parsed;
// some are unimportant, diagnostic info (stuff we already know) or simply not useful
// as events.
// Eg. team preview messages, duplicate health/damage updates (with percentages that
// we can calculate anyways if needed), server time, info about the format, players
// or format rules.
func (s *Process) parseStdout(raw string) {
	// requires showdown version 0.11.4+ to fix a bug where messages are not returned
	lines := strings.Split(strings.Trim(strings.TrimSpace(raw), "\x00"), "\n")
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "|request|") {
			bits := strings.Split(lines[i], "|")

			encoded := bits[len(bits)-1]
			if strings.Contains(encoded, "\"teamPreview\":true") {
				// Specifically ignore preview messages.
				continue
			}

			// the simulator is asking a player to make a choice
			update, err := structs.DecodeUpdate([]byte(encoded))
			if err != nil {
				s.messages <- message(err)
				continue
			}
			s.messages <- message(update)
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
				evt := event.Parse(lines[i])
				if evt != nil {
					s.messages <- message(evt)
				}
				continue
			}

			bitsprev := strings.Split(lines[i-1], "|")
			// if it is the same type as previous & about the same pokemon
			if bitsthis[2] == bitsprev[2] {
				continue
			}
		} else if strings.HasPrefix(lines[i], "|error|") {
			if strings.Contains(lines[i], "Can't choose for Team Preview") {
				// we totally ignore team preview messages, so we also
				// ignore this error which occurs if we give a team
				// ordering when the particular format doesn't allow
				// us to give one.
				// That is, we always give an ordering; if the server
				// accepts it then great, if not then it does no harm.
				continue
			}

			err := fmt.Errorf(strings.Replace(lines[i], "|error|", "", 1))

			// there are two specific sub errors we're really interested in
			// since they both indicate that a user explicitly did something
			// wrong.
			if strings.Contains(lines[i], "[Invalid choice]") {
				err = fmt.Errorf("%w %v", ErrInvalidChoice, err)
			} else if strings.Contains(lines[i], "[Unavailable choice]") {
				err = fmt.Errorf("%w %v", ErrUnavailableChoice, err)
			}

			s.messages <- message(err)
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
				// if we got this far, then it's probably a noteworthy event
				evt := event.Parse(lines[i])
				if evt != nil {
					s.messages <- message(evt)
				}
			}
		}
	}
}

// Messages are parsed messages (in order) from the showdown simulator
func (s *Process) Messages() <-chan *Message {
	return s.messages
}

// Write writes one player decision(s) for their pokemon (1 or more) to
// the simulator.
func (s *Process) Write(in string) error {
	log.Printf(in)
	s.raw.stdin <- in
	return nil
}

// Stop simulation & kill subprocess(es)
func (s *Process) Stop() {
	s.raw.stop()
}
