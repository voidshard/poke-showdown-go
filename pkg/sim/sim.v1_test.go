package sim

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
	"os"
	"sync"
	"testing"
)

func TestLatestErrorsNone(t *testing.T) {
	sv1 := &SimV1{unreadErrors: []error{}, unreadErrorsLock: &sync.Mutex{}}

	err := sv1.latestErrors()

	assert.Nil(t, err)
}

func TestLatestErrors(t *testing.T) {
	sv1 := &SimV1{
		unreadErrors:     []error{fmt.Errorf("one"), fmt.Errorf("two")},
		unreadErrorsLock: &sync.Mutex{},
	}

	err := sv1.latestErrors()

	assert.NotNil(t, err)
	assert.Equal(t, 0, len(sv1.unreadErrors))
}

func TestCollateEvents(t *testing.T) {
	out := make(chan string)
	sv1 := &SimV1{
		stdout:           out,
		unreadErrorsLock: &sync.Mutex{},
		unreadErrors:     []error{},
		state:            make(chan *structs.BattleState),
		messages:         make(chan string),
		players:          []string{"p1", "p2"},
	}

	go func() {
		for _, msg := range []string{
			"|error|foo occured\n",
			`|-event|something happened`,
			"|sideupdate\np1\n|request|{\"active\":[{\"moves\":[{\"move\":\"Will-O-Wisp\",\"id\":\"willowisp\",\"pp\":24,\"maxpp\":24,\"target\":\"normal\",\"disabled\":false},{\"move\":\"Nasty Plot\",\"id\":\"nastyplot\",\"pp\":32,\"maxpp\":32,\"target\":\"self\",\"disabled\":false},{\"move\":\"Fire Blast\",\"id\":\"fireblast\",\"pp\":8,\"maxpp\":8,\"target\":\"normal\",\"disabled\":false},{\"move\":\"Solar Beam\",\"id\":\"solarbeam\",\"pp\":16,\"maxpp\":16,\"target\":\"normal\",\"disabled\":false}],\"canDynamax\":true,\"maxMoves\":{\"maxMoves\":[{\"move\":\"maxguard\",\"target\":\"self\"},{\"move\":\"maxguard\",\"target\":\"self\"},{\"move\":\"maxflare\",\"target\":\"adjacentFoe\"},{\"move\":\"maxovergrowth\",\"target\":\"adjacentFoe\"}]}}],\"side\":{\"name\":\"p1\",\"id\":\"p1\",\"pokemon\":[{\"ident\":\"p1: Ninetales\",\"details\":\"Ninetales, L50, M\",\"condition\":\"159/159\",\"active\":true,\"stats\":{\"atk\":107,\"def\":106,\"spa\":112,\"spd\":131,\"spe\":131},\"moves\":[\"willowisp\",\"nastyplot\",\"fireblast\",\"solarbeam\"],\"baseAbility\":\"drought\",\"item\":\"heavydutyboots\",\"pokeball\":\"pokeball\",\"ability\":\"drought\"}]}}",
			"|sideupdate\np2\n|request|{\"active\":[{\"moves\":[{\"move\":\"Protect\",\"id\":\"protect\",\"pp\":16,\"maxpp\":16,\"target\":\"self\",\"disabled\":false},{\"move\":\"Foul Play\",\"id\":\"foulplay\",\"pp\":24,\"maxpp\":24,\"target\":\"normal\",\"disabled\":false},{\"move\":\"Wish\",\"id\":\"wish\",\"pp\":16,\"maxpp\":16,\"target\":\"self\",\"disabled\":false},{\"move\":\"Toxic\",\"id\":\"toxic\",\"pp\":16,\"maxpp\":16,\"target\":\"normal\",\"disabled\":false}],\"canDynamax\":true,\"maxMoves\":{\"maxMoves\":[{\"move\":\"maxguard\",\"target\":\"self\"},{\"move\":\"maxdarkness\",\"target\":\"adjacentFoe\"},{\"move\":\"maxguard\",\"target\":\"self\"},{\"move\":\"maxguard\",\"target\":\"self\"}]}}],\"side\":{\"name\":\"p2\",\"id\":\"p2\",\"pokemon\":[{\"ident\":\"p2: Umbreon\",\"details\":\"Umbreon, L50, F\",\"condition\":\"181/181\",\"active\":true,\"stats\":{\"atk\":96,\"def\":141,\"spa\":91,\"spd\":161,\"spe\":96},\"moves\":[\"protect\",\"foulplay\",\"wish\",\"toxic\"],\"baseAbility\":\"synchronize\",\"item\":\"leftovers\",\"pokeball\":\"pokeball\",\"ability\":\"synchronize\"}]}}",
		} {
			out <- msg
		}
		close(out)
	}()

	go sv1.collateEvents()

	st := <-sv1.state
	assert.NotNil(t, st)
	assert.Equal(t, st.Field["p1"].Pokemon[0].Species(), "Ninetales")
	assert.Equal(t, st.Field["p2"].Pokemon[0].Species(), "Umbreon")

	assert.Equal(t, []string{"|-event|something happened"}, st.Messages)
	assert.Equal(t, 1, len(sv1.unreadErrors))
}

func TestParseMessage(t *testing.T) {
	state := structs.NewBattleState()

	msgs, err := parseMessage(
		"|sideupdate\np1\n|request|{\"active\":[{\"moves\":[{\"move\":\"Will-O-Wisp\",\"id\":\"willowisp\",\"pp\":24,\"maxpp\":24,\"target\":\"normal\",\"disabled\":false},{\"move\":\"Nasty Plot\",\"id\":\"nastyplot\",\"pp\":32,\"maxpp\":32,\"target\":\"self\",\"disabled\":false},{\"move\":\"Fire Blast\",\"id\":\"fireblast\",\"pp\":8,\"maxpp\":8,\"target\":\"normal\",\"disabled\":false},{\"move\":\"Solar Beam\",\"id\":\"solarbeam\",\"pp\":16,\"maxpp\":16,\"target\":\"normal\",\"disabled\":false}],\"canDynamax\":true,\"maxMoves\":{\"maxMoves\":[{\"move\":\"maxguard\",\"target\":\"self\"},{\"move\":\"maxguard\",\"target\":\"self\"},{\"move\":\"maxflare\",\"target\":\"adjacentFoe\"},{\"move\":\"maxovergrowth\",\"target\":\"adjacentFoe\"}]}}],\"side\":{\"name\":\"p1\",\"id\":\"p1\",\"pokemon\":[{\"ident\":\"p1: Ninetales\",\"details\":\"Ninetales, L50, M\",\"condition\":\"159/159\",\"active\":true,\"stats\":{\"atk\":107,\"def\":106,\"spa\":112,\"spd\":131,\"spe\":131},\"moves\":[\"willowisp\",\"nastyplot\",\"fireblast\",\"solarbeam\"],\"baseAbility\":\"drought\",\"item\":\"heavydutyboots\",\"pokeball\":\"pokeball\",\"ability\":\"drought\"}]}}\n|win|p1\n|-event|hithere",
		state,
	)

	field := state.Field["p1"]

	assert.Nil(t, err)
	assert.NotNil(t, field)
	assert.Equal(t, []string{"|-event|hithere"}, msgs)
	assert.Equal(t, "p1", state.Winner)
}

func TestParseMessageError(t *testing.T) {
	state := structs.NewBattleState()

	msgs, err := parseMessage("|error|what\n\np1\n", state)

	assert.Equal(t, 0, len(msgs))
	assert.NotNil(t, err)
}

func TestNewSimV1(t *testing.T) {
	f := fakeRun{}
	runCommand = f.Run

	spec := &structs.BattleSpec{
		Format: structs.FormatGen8,
		Players: map[string][]*structs.PokemonSpec{
			"p1": []*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "Ninetales",
					Item:    "heavydutyboots",
					Ability: "drought",
					Moves:   []string{"willowisp", "nastyplot", "fireblast", "solarbeam"},
					Level:   50,
					Gender:  "M",
				},
			},
			"p2": []*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "Umbreon",
					Item:    "leftovers",
					Ability: "synchronize",
					Moves:   []string{"protect", "foulplay", "wish", "toxic"},
					Level:   50,
					Gender:  "F",
				},
			},
		},
	}

	result, err := NewSimV1("somepath", spec)

	sv1 := result.(*SimV1)

	// check struct was setup
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, sv1.stdout)
	assert.NotNil(t, sv1.stderr)
	assert.NotNil(t, sv1.stdin)
	assert.NotNil(t, sv1.ctrl)
	assert.NotNil(t, sv1.errors)
	assert.NotNil(t, sv1.state)
	assert.NotNil(t, sv1.messages)
	assert.NotNil(t, sv1.unreadErrors)
	assert.NotNil(t, sv1.unreadErrorsLock)

	assert.Contains(t, sv1.players, "p1")
	assert.Contains(t, sv1.players, "p2")
	assert.Equal(t, 2, len(sv1.players))

	// check we wrote the expected number of messages
	assert.Equal(t, 5, len(f.in))
	assert.Equal(t, f.in[0], ">start {\"formatid\":\"[Gen 8] Anything Goes\"}\n")

	teams := f.in[1] + f.in[2]
	assert.Contains(t, teams, ">player p1 {\"name\":\"p1\",\"team\":\"Ninetales||heavydutyboots|drought|willowisp,nastyplot,fireblast,solarbeam||85,85,85,85,85,85|M|31,31,31,31,31,31||50|0,,,\"}\n")
	assert.Contains(t, teams, ">player p2 {\"name\":\"p2\",\"team\":\"Umbreon||leftovers|synchronize|protect,foulplay,wish,toxic||85,85,85,85,85,85|F|31,31,31,31,31,31||50|0,,,\"}\n")

	orders := f.in[3] + f.in[4]
	assert.Contains(t, orders, ">p1 team 1")
	assert.Contains(t, orders, ">p2 team 1")
}

type fakeRun struct {
	cmd  string
	args []string
	in   []string
	sigs []os.Signal
}

func (f *fakeRun) Run(cmd string, args []string, stdin <-chan string, ctrl chan os.Signal) (<-chan string, <-chan string, <-chan error) {
	f.cmd = cmd
	f.args = args

	f.in = []string{}
	f.sigs = []os.Signal{}

	go func() {
		for {
			select {
			case input := <-stdin:
				f.in = append(f.in, input)
			case sig := <-ctrl:
				f.sigs = append(f.sigs, sig)
			}
		}
	}()

	return make(chan string), make(chan string), make(chan error)
}
