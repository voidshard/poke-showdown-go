package sim

import (
	"github.com/stretchr/testify/assert"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
	"testing"
)

func TestSimulateSwitch(t *testing.T) {
	spec := &structs.BattleSpec{
		Format: structs.FormatGen8,
		Players: [][]*structs.PokemonSpec{
			[]*structs.PokemonSpec{ // p1
				&structs.PokemonSpec{
					Name:    "Umbreon",
					Species: "umbreon",
					Level:   50,
					Moves:   []string{"wish", "toxic", "protect", "bite"},
					Ability: "synchronize",
				},
			},
			[]*structs.PokemonSpec{ // p2
				&structs.PokemonSpec{
					Name:    "Pikachu",
					Species: "pikachu",
					Level:   50,
					Moves:   []string{"tackle"},
					Ability: "voltabsorb", // seems legit
				},
				&structs.PokemonSpec{
					Name:    "Ninetales",
					Species: "ninetales",
					Level:   1,
					Moves:   []string{"solarbeam", "flamethrower", "willowisp", "sunnyday"},
					Ability: "flashfire",
				},
			},
		},
		Seed: 4910,
	}

	actions := []*structs.Action{
		&structs.Action{
			Player: "p1",
			Specs:  []*structs.ActionSpec{&structs.ActionSpec{Type: structs.ActionMove, ID: 3}},
		}, // player 1 uses bite
		&structs.Action{
			Player: "p2",
			Specs:  []*structs.ActionSpec{&structs.ActionSpec{Type: structs.ActionSwitch, ID: 1}},
		}, // player 2 switches in ninetales
		&structs.Action{
			Player: "p2",
			Specs:  []*structs.ActionSpec{&structs.ActionSpec{Type: structs.ActionSwitch, ID: 1}},
		}, // player 2 switches in pikachu (ninetales was KO'd)

		&structs.Action{
			Player: "p1",
			Specs:  []*structs.ActionSpec{&structs.ActionSpec{Type: structs.ActionMove, ID: 1}},
		}, // player 1 uses toxic
		&structs.Action{
			Player: "p2",
			Specs:  []*structs.ActionSpec{&structs.ActionSpec{Type: structs.ActionMove, ID: 0}},
		}, // player 2 uses tackle
	}

	result, err := Simulate(spec, actions)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, result.Turn)
	assert.Equal(t, "Tackle", result.Events[len(result.Events)-5].Name)
	assert.Equal(t, "Toxic", result.Events[len(result.Events)-3].Name)
}

func TestSimulate(t *testing.T) {
	spec := &structs.BattleSpec{
		Format: structs.FormatGen8,
		Players: [][]*structs.PokemonSpec{
			[]*structs.PokemonSpec{ // p1
				&structs.PokemonSpec{
					Name:    "Umbreon",
					Species: "umbreon",
					Level:   50,
					Moves:   []string{"wish", "toxic", "protect", "bite"},
					Ability: "synchronize",
				},
			},
			[]*structs.PokemonSpec{ // p2
				&structs.PokemonSpec{
					Name:    "Ninetales",
					Species: "ninetales",
					Level:   50,
					Moves:   []string{"solarbeam", "flamethrower", "willowisp", "sunnyday"},
					Ability: "flashfire",
				},
			},
		},
		Seed: 12345678,
	}

	actions := []*structs.Action{
		// turn 1
		&structs.Action{
			Player: "p1",
			Specs:  []*structs.ActionSpec{&structs.ActionSpec{Type: structs.ActionMove, ID: 1}},
		}, // player 1 uses toxic
		&structs.Action{
			Player: "p2",
			Specs:  []*structs.ActionSpec{&structs.ActionSpec{Type: structs.ActionMove, ID: 0}},
		}, // player 2 uses solarbeam

		// turn 2
		&structs.Action{
			Player: "p1",
			Specs:  []*structs.ActionSpec{&structs.ActionSpec{Type: structs.ActionMove, ID: 2}},
		}, // player 1 blocks with protect
		&structs.Action{
			Player: "p2",
			Specs:  []*structs.ActionSpec{&structs.ActionSpec{Type: structs.ActionMove, ID: 0}},
		}, // player 1 continues solarbeam

		// turn 3
		&structs.Action{
			Player: "p1",
			Specs:  []*structs.ActionSpec{&structs.ActionSpec{Type: structs.ActionMove, ID: 3}},
		}, // player 1 uses bite
		&structs.Action{
			Player: "p2",
			Specs:  []*structs.ActionSpec{&structs.ActionSpec{Type: structs.ActionMove, ID: 1}},
		}, // player 2 uses flamethrower
	}

	result, err := Simulate(spec, actions)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 4, result.Turn)

	p1field := result.Field["p1"]
	assert.Equal(t, 138, p1field.Pokemon[0].Data.HPNow)
	assert.Equal(t, 181, p1field.Pokemon[0].Data.HPMax)

	p2field := result.Field["p2"]
	assert.Equal(t, 71, p2field.Pokemon[0].Data.HPNow)
	assert.Equal(t, 159, p2field.Pokemon[0].Data.HPMax)
	assert.Equal(t, true, p2field.Pokemon[0].Data.IsToxiced)
}
