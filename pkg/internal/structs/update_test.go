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
			cur, max, err := (&Pokemon{Condition: tt.Condition}).parseHP()

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
		Frozen    bool
	}{
		{"150/150", false, false, false, false, false, false, false},
		{"150/150 slp", true, false, false, false, false, false, false},
		{"150/150 brn", false, true, false, false, false, false, false},
		{"150/150 par", false, false, true, false, false, false, false},
		{"0 fnt", false, false, false, true, false, false, false},
		{"150/150 psn", false, false, false, false, true, false, false},
		{"150/150 tox", false, false, false, false, true, true, false},
		{"150/150 frz", false, false, false, false, false, false, true},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			pkm := &Pokemon{Condition: tt.Condition}

			assert.Equal(t, tt.Asleep, pkm.isAsleep())
			assert.Equal(t, tt.Burned, pkm.isBurned())
			assert.Equal(t, tt.Paralyzed, pkm.isParalyzed())
			assert.Equal(t, tt.Fainted, pkm.isFainted())
			assert.Equal(t, tt.Poisoned, pkm.isPoisoned())
			assert.Equal(t, tt.Toxiced, pkm.isToxiced())
			assert.Equal(t, tt.Frozen, pkm.isFrozen())

			if tt.Toxiced {
				assert.True(t, pkm.isPoisoned())
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
		{"Umbreon, F", "Umbreon"},
		{"Whatever, L10, M", "Whatever"},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := (&Pokemon{Details: tt.Given}).species()

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
		{"Umbreon, F", 100}, // don't even get me started. WTF
		{"Whatever, L10, M", 10},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result, err := (&Pokemon{Details: tt.Given}).parseLevel()

			assert.Nil(t, err)
			assert.Equal(t, tt.Expect, result)
		})
	}
}

func TestDecodeUpdateSwitch(t *testing.T) {
	in := []byte("{\"wait\":true,\"side\":{\"name\":\"p1\",\"id\":\"p1\",\"pokemon\":[{\"ident\":\"p1: Cursola\",\"details\":\"Cursola, L88, M\",\"condition\":\"106/249 psn\",\"active\":true,\"stats\":{\"atk\":172,\"def\":138,\"spa\":305,\"spd\":279,\"spe\":103},\"moves\":[\"shadowball\",\"hydropump\",\"icebeam\",\"earthpower\"],\"baseAbility\":\"perishbody\",\"item\":\"choicespecs\",\"pokeball\":\"pokeball\",\"ability\":\"perishbody\"},{\"ident\":\"p1: Aurorus\",\"details\":\"Aurorus, L86, M\",\"condition\":\"352/352\",\"active\":false,\"stats\":{\"atk\":137,\"def\":173,\"spa\":220,\"spd\":207,\"spe\":149},\"moves\":[\"ancientpower\",\"stealthrock\",\"thunderwave\",\"blizzard\"],\"baseAbility\":\"snowwarning\",\"item\":\"heavydutyboots\",\"pokeball\":\"pokeball\",\"ability\":\"snowwarning\"},{\"ident\":\"p1: Tapu Lele\",\"details\":\"Tapu Lele, L82\",\"condition\":\"249/249\",\"active\":false,\"stats\":{\"atk\":144,\"def\":170,\"spa\":260,\"spd\":236,\"spe\":203},\"moves\":[\"psyshock\",\"moonblast\",\"focusblast\",\"calmmind\"],\"baseAbility\":\"psychicsurge\",\"item\":\"lifeorb\",\"pokeball\":\"pokeball\",\"ability\":\"psychicsurge\"},{\"ident\":\"p1: Reuniclus\",\"details\":\"Reuniclus, L84, M\",\"condition\":\"322/322\",\"active\":false,\"stats\":{\"atk\":114,\"def\":174,\"spa\":258,\"spd\":191,\"spe\":55},\"moves\":[\"focusblast\",\"psychic\",\"shadowball\",\"trickroom\"],\"baseAbility\":\"magicguard\",\"item\":\"lifeorb\",\"pokeball\":\"pokeball\",\"ability\":\"magicguard\"},{\"ident\":\"p1: Buzzwole\",\"details\":\"Buzzwole, L80\",\"condition\":\"302/302\",\"active\":false,\"stats\":{\"atk\":269,\"def\":269,\"spa\":131,\"spd\":131,\"spe\":173},\"moves\":[\"stoneedge\",\"dualwingbeat\",\"leechlife\",\"closecombat\"],\"baseAbility\":\"beastboost\",\"item\":\"choiceband\",\"pokeball\":\"pokeball\",\"ability\":\"beastboost\"},{\"ident\":\"p1: Naganadel\",\"details\":\"Naganadel, L76\",\"condition\":\"236/236\",\"active\":false,\"stats\":{\"atk\":115,\"def\":155,\"spa\":237,\"spd\":155,\"spe\":228},\"moves\":[\"sludgewave\",\"nastyplot\",\"flamethrower\",\"dracometeor\"],\"baseAbility\":\"beastboost\",\"item\":\"lifeorb\",\"pokeball\":\"pokeball\",\"ability\":\"beastboost\"}]}}")

	result, err := DecodeUpdate(in)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, 0, len(result.Active))
}

func TestDecodeUpdate(t *testing.T) {
	in := []byte("{\"active\":[{\"moves\":[{\"move\":\"Nasty Plot\",\"id\":\"nastyplot\",\"pp\":32,\"maxpp\":32,\"target\":\"self\",\"disabled\":false},{\"move\":\"Dark Pulse\",\"id\":\"darkpulse\",\"pp\":24,\"maxpp\":24,\"target\":\"any\",\"disabled\":false},{\"move\":\"Sludge Bomb\",\"id\":\"sludgebomb\",\"pp\":16,\"maxpp\":16,\"target\":\"normal\",\"disabled\":false},{\"move\":\"Flamethrower\",\"id\":\"flamethrower\",\"pp\":24,\"maxpp\":24,\"target\":\"normal\",\"disabled\":false}],\"canDynamax\":true,\"maxMoves\":{\"maxMoves\":[{\"move\":\"maxguard\",\"target\":\"self\"},{\"move\":\"maxdarkness\",\"target\":\"adjacentFoe\"},{\"move\":\"maxooze\",\"target\":\"adjacentFoe\"},{\"move\":\"maxflare\",\"target\":\"adjacentFoe\"}]}}],\"side\":{\"name\":\"p2\",\"id\":\"p2\",\"pokemon\":[{\"ident\":\"p2: Zoroark\",\"details\":\"Zoroark, L5, F\",\"condition\":\"23/23\",\"active\":true,\"stats\":{\"atk\":18,\"def\":13,\"spa\":19,\"spd\":13,\"spe\":18},\"moves\":[\"nastyplot\",\"darkpulse\",\"sludgebomb\",\"flamethrower\"],\"baseAbility\":\"illusion\",\"item\":\"lifeorb\",\"pokeball\":\"pokeball\",\"ability\":\"illusion\"},{\"ident\":\"p2: Umbreon\",\"details\":\"Umbreon, L5, F\",\"condition\":\"27/27\",\"active\":false,\"stats\":{\"atk\":14,\"def\":18,\"spa\":13,\"spd\":20,\"spe\":14},\"moves\":[\"protect\",\"foulplay\",\"wish\",\"toxic\"],\"baseAbility\":\"synchronize\",\"item\":\"leftovers\",\"pokeball\":\"pokeball\",\"ability\":\"synchronize\"}]}}")

	result, err := DecodeUpdate(in)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, 2, len(result.Team.Pokemon))
	assert.Equal(t, "p2", result.Team.Player)
	assert.NotNil(t, result.Active)
	assert.Equal(t, 1, len(result.Active))
	assert.Equal(t, 4, len(result.Active[0].Moves))
	assert.Equal(t, "Nasty Plot", result.Active[0].Moves[0].Name)
	assert.Equal(t, "nastyplot", result.Active[0].Moves[0].ID)
	assert.Equal(t, "Zoroark", result.Team.Pokemon[0].species())
	assert.NotNil(t, result.Team.Pokemon[0].Status)
	assert.Equal(t, "Umbreon", result.Team.Pokemon[1].species())
	assert.Equal(t, "foulplay", result.Team.Pokemon[1].Moves[1])
}
