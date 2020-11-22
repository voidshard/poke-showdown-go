package structs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHP(t *testing.T) {
	cases := []struct {
		Condition string
		Current   int
		Max       int
	}{
		{"150/150", 150, 150},
		{"0 fnt", 0, -1},
		{"100/150", 100, 150},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			cur, max, err := (&Pokemon{Condition: tt.Condition}).HP()

			assert.Nil(t, err)
			assert.Equal(t, tt.Current, cur)
			assert.Equal(t, tt.Max, max)
		})
	}
}

func TestPokemonStatus(t *testing.T) {
	cases := []struct {
		Condition string

		Asleep    bool
		Burned    bool
		Paralyzed bool
		Fainted   bool
		Poisoned  bool
		Toxiced   bool
	}{
		{"150/150", false, false, false, false, false, false},
		{"150/150 slp", true, false, false, false, false, false},
		{"150/150 brn", false, true, false, false, false, false},
		{"150/150 par", false, false, true, false, false, false},
		{"0 fnt", false, false, false, true, false, false},
		{"150/150 psn", false, false, false, false, true, false},
		{"150/150 tox", false, false, false, false, true, true},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			pkm := &Pokemon{Condition: tt.Condition}

			assert.Equal(t, tt.Asleep, pkm.IsAsleep())
			assert.Equal(t, tt.Burned, pkm.IsBurned())
			assert.Equal(t, tt.Paralyzed, pkm.IsParalyzed())
			assert.Equal(t, tt.Fainted, pkm.IsFainted())
			assert.Equal(t, tt.Poisoned, pkm.IsPoisoned())
			assert.Equal(t, tt.Toxiced, pkm.IsToxiced())

			if tt.Toxiced {
				assert.True(t, pkm.IsPoisoned())
			}
		})
	}
}

func TestPokemonSpecies(t *testing.T) {
	cases := []struct {
		Given  string
		Expect string
	}{
		{"Umbreon, L5, F", "Umbreon"},
		{"Whatever, L10, M", "Whatever"},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := (&Pokemon{Details: tt.Given}).Species()

			assert.Equal(t, tt.Expect, result)
		})
	}
}

func TestPokemonLevel(t *testing.T) {
	cases := []struct {
		Given  string
		Expect int
	}{
		{"Umbreon, L5, F", 5},
		{"Whatever, L10, M", 10},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result, err := (&Pokemon{Details: tt.Given}).Level()

			assert.Nil(t, err)
			assert.Equal(t, tt.Expect, result)
		})
	}
}

func TestDecodeUpdate(t *testing.T) {
	in := []byte("{\"active\":[{\"moves\":[{\"move\":\"Nasty Plot\",\"id\":\"nastyplot\",\"pp\":32,\"maxpp\":32,\"target\":\"self\",\"disabled\":false},{\"move\":\"Dark Pulse\",\"id\":\"darkpulse\",\"pp\":24,\"maxpp\":24,\"target\":\"any\",\"disabled\":false},{\"move\":\"Sludge Bomb\",\"id\":\"sludgebomb\",\"pp\":16,\"maxpp\":16,\"target\":\"normal\",\"disabled\":false},{\"move\":\"Flamethrower\",\"id\":\"flamethrower\",\"pp\":24,\"maxpp\":24,\"target\":\"normal\",\"disabled\":false}],\"canDynamax\":true,\"maxMoves\":{\"maxMoves\":[{\"move\":\"maxguard\",\"target\":\"self\"},{\"move\":\"maxdarkness\",\"target\":\"adjacentFoe\"},{\"move\":\"maxooze\",\"target\":\"adjacentFoe\"},{\"move\":\"maxflare\",\"target\":\"adjacentFoe\"}]}}],\"side\":{\"name\":\"p2\",\"id\":\"p2\",\"pokemon\":[{\"ident\":\"p2: Zoroark\",\"details\":\"Zoroark, L5, F\",\"condition\":\"23/23\",\"active\":true,\"stats\":{\"atk\":18,\"def\":13,\"spa\":19,\"spd\":13,\"spe\":18},\"moves\":[\"nastyplot\",\"darkpulse\",\"sludgebomb\",\"flamethrower\"],\"baseAbility\":\"illusion\",\"item\":\"lifeorb\",\"pokeball\":\"pokeball\",\"ability\":\"illusion\"},{\"ident\":\"p2: Umbreon\",\"details\":\"Umbreon, L5, F\",\"condition\":\"27/27\",\"active\":false,\"stats\":{\"atk\":14,\"def\":18,\"spa\":13,\"spd\":20,\"spe\":14},\"moves\":[\"protect\",\"foulplay\",\"wish\",\"toxic\"],\"baseAbility\":\"synchronize\",\"item\":\"leftovers\",\"pokeball\":\"pokeball\",\"ability\":\"synchronize\"}]}}")

	result, err := DecodeUpdate(in)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, 2, len(result.Pokemon))
	assert.Equal(t, "p2", result.Player)
	assert.Equal(t, []int{0}, result.Active)
	assert.Equal(t, "Zoroark", result.Pokemon[0].Species())
	assert.NotNil(t, result.Pokemon[0].Options)
	assert.Equal(t, "Umbreon", result.Pokemon[1].Species())
}
