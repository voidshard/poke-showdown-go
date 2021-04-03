package sim

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsDoubles(t *testing.T) {
	assert.True(t, IsDoubles(FormatGen8Doubles))
	assert.True(t, !IsDoubles(FormatGen8))
}

var dataTestValidate = []struct {
	Name string
	In   *BattleSpec
	Err  bool
}{
	{
		"no-data",
		&BattleSpec{},
		true,
	},
	{
		"no-players",
		&BattleSpec{Players: [][]*PokemonSpec{}},
		true,
	},
	{
		"not-enough-players",
		&BattleSpec{Players: [][]*PokemonSpec{
			[]*PokemonSpec{&PokemonSpec{}, &PokemonSpec{}},
		}},
		true,
	},
	{
		"not-enough-pokemon",
		&BattleSpec{Players: [][]*PokemonSpec{
			[]*PokemonSpec{&PokemonSpec{}, &PokemonSpec{}},
			[]*PokemonSpec{},
		}},
		true,
	},
	{
		"too-many-players",
		&BattleSpec{Players: [][]*PokemonSpec{
			[]*PokemonSpec{&PokemonSpec{}, &PokemonSpec{}},
			[]*PokemonSpec{&PokemonSpec{}},
			[]*PokemonSpec{&PokemonSpec{}},
		}},
		true,
	},
	{
		"valid-for-singles",
		&BattleSpec{Players: [][]*PokemonSpec{
			[]*PokemonSpec{&PokemonSpec{}, &PokemonSpec{}},
			[]*PokemonSpec{&PokemonSpec{}},
		}},
		false,
	},
	{
		"invalid-for-doubles",
		&BattleSpec{
			Format: FormatGen8Doubles,
			Players: [][]*PokemonSpec{
				[]*PokemonSpec{&PokemonSpec{}, &PokemonSpec{}},
				[]*PokemonSpec{&PokemonSpec{}},
			}},
		true,
	},
	{
		"valid-for-doubles",
		&BattleSpec{
			Format: FormatGen8Doubles,
			Players: [][]*PokemonSpec{
				[]*PokemonSpec{&PokemonSpec{}, &PokemonSpec{}},
				[]*PokemonSpec{&PokemonSpec{}, &PokemonSpec{}},
			}},
		false,
	},
	{
		"full-team-pokemon",
		&BattleSpec{
			Players: [][]*PokemonSpec{
				[]*PokemonSpec{
					&PokemonSpec{}, &PokemonSpec{},
					&PokemonSpec{}, &PokemonSpec{},
					&PokemonSpec{}, &PokemonSpec{},
				},
				[]*PokemonSpec{
					&PokemonSpec{}, &PokemonSpec{},
					&PokemonSpec{}, &PokemonSpec{},
					&PokemonSpec{}, &PokemonSpec{},
				},
			}},
		false,
	},
	{
		"too-many-pokemon",
		&BattleSpec{
			Players: [][]*PokemonSpec{
				[]*PokemonSpec{
					&PokemonSpec{}, &PokemonSpec{},
					&PokemonSpec{}, &PokemonSpec{},
					&PokemonSpec{}, &PokemonSpec{},
					&PokemonSpec{},
				},
				[]*PokemonSpec{
					&PokemonSpec{}, &PokemonSpec{},
					&PokemonSpec{}, &PokemonSpec{},
					&PokemonSpec{}, &PokemonSpec{},
				},
			}},
		true,
	},
}

func TestValidate(t *testing.T) {
	for _, tt := range dataTestValidate {
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.In.validate()

			if tt.Err {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
