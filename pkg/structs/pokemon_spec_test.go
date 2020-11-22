package structs

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

	result := PackTeam(team)
	assert.Equal(t, packedTeam, result)
}
