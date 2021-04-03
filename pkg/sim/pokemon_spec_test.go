package sim

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnforceLimits(t *testing.T) {
	p := PokemonSpec{
		Moves:            []string{"a", "b", "c", "d", "e", "f"},
		EffortValues:     &Stats{255, 255, 255, 0, 0, 0},
		IndividualValues: &Stats{39, -15, 31, 31, 31, 31},
		Gender:           "X",
		Level:            300,
		Happiness:        -150,
	}

	p.enforceLimits()

	assert.Equal(t, []string{"a", "b", "c", "d"}, p.Moves)
	assert.Equal(t, &Stats{85, 85, 85, 85, 85, 85}, p.EffortValues)
	assert.Equal(t, &Stats{31, 0, 31, 31, 31, 31}, p.IndividualValues)
	assert.Equal(t, "", p.Gender)
	assert.Equal(t, 100, p.Level)
	assert.Equal(t, 0, p.Happiness)
}

func TestPackBool(t *testing.T) {
	defVal := "value"

	cases := []struct {
		Given  bool
		Expect string
	}{
		{true, defVal},
		{false, ""},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := packbool(tt.Given, defVal)

			assert.Equal(t, tt.Expect, result)
		})
	}
}

const packedTeam = `Musharna||leftovers|synchronize|calmmind,moonlight,moonblast,psychic||85,0,85,85,85,85||31,0,31,31,31,31||88|0,,,]Lycanroc|lycanrocdusk|lifeorb|toughclaws|swordsdance,closecombat,psychicfangs,stoneedge||85,85,85,85,85,85||31,31,31,31,31,31||82|0,,,]Skuntank||lifeorb|aftermath|suckerpunch,fireblast,toxic,crunch||85,85,85,85,85,85||31,31,31,31,31,31||86|0,,,]Terrakion||choiceband|justified|earthquake,stoneedge,quickattack,closecombat||85,85,85,85,85,85|N|31,31,31,31,31,31||82|0,,,]Persian|persianalola|lifeorb|furcoat|nastyplot,powergem,thunderbolt,darkpulse||85,0,85,85,85,85||31,0,31,31,31,31||86|0,,,]Genesect||lifeorb|download|shiftgear,icebeam,thunderbolt,ironhead||85,85,85,85,85,85|N|31,31,31,31,31,31||76|0,,,`

const unpackedTeam = `[{"name":"Musharna","species":"Musharna","item":"leftovers","ability":"synchronize","moves":["calmmind","moonlight","moonblast","psychic"],"nature":"","evs":{"hp":85,"atk":0,"def":85,"spa":85,"spd":85,"spe":85},"ivs":{"hp":31,"atk":0,"def":31,"spa":31,"spd":31,"spe":31},"level":88},{"name":"Lycanroc","species":"lycanrocdusk","item":"lifeorb","ability":"toughclaws","moves":["swordsdance","closecombat","psychicfangs","stoneedge"],"nature":"","evs":{"hp":85,"atk":85,"def":85,"spa":85,"spd":85,"spe":85},"level":82},{"name":"Skuntank","species":"Skuntank","item":"lifeorb","ability":"aftermath","moves":["suckerpunch","fireblast","toxic","crunch"],"nature":"","evs":{"hp":85,"atk":85,"def":85,"spa":85,"spd":85,"spe":85},"level":86},{"name":"Terrakion","species":"Terrakion","item":"choiceband","ability":"justified","moves":["earthquake","stoneedge","quickattack","closecombat"],"nature":"","evs":{"hp":85,"atk":85,"def":85,"spa":85,"spd":85,"spe":85},"gender":"N","level":82},{"name":"Persian","species":"persianalola","item":"lifeorb","ability":"furcoat","moves":["nastyplot","powergem","thunderbolt","darkpulse"],"nature":"","evs":{"hp":85,"atk":0,"def":85,"spa":85,"spd":85,"spe":85},"ivs":{"hp":31,"atk":0,"def":31,"spa":31,"spd":31,"spe":31},"level":86},{"name":"Genesect","species":"Genesect","item":"lifeorb","ability":"download","moves":["shiftgear","icebeam","thunderbolt","ironhead"],"nature":"","evs":{"hp":85,"atk":85,"def":85,"spa":85,"spd":85,"spe":85},"gender":"N","level":76}]`

func TestPackTeam(t *testing.T) {
	team := []*PokemonSpec{}
	err := json.Unmarshal([]byte(unpackedTeam), &team)

	assert.Nil(t, err)
	if err != nil {
		return
	}

	result, err := PackTeam(team)

	assert.Nil(t, err)
	assert.Equal(t, packedTeam, result)
}

func TestClamp(t *testing.T) {
	assert.Equal(t, 130, clamp(0, 200, 130))
	assert.Equal(t, 0, clamp(0, 200, -1))
	assert.Equal(t, 200, clamp(0, 200, 940))
}

func TestSum(t *testing.T) {
	given := &Stats{10, 20, 30, 40, 50, 60}

	result := given.Sum()

	assert.Equal(t, 210, result)
}
