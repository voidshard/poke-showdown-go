package sim

import (
	"github.com/stretchr/testify/assert"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
	"testing"
)

var dataTestParseStdout = `sideupdate
p1
|request|{"teamPreview":true,"maxTeamSize":6,"side":{"name":"p1","id":"p1","pokemon":[{"ident":"p1: Lugia","details":"Lugia, L5","condition":"28/28","active":true,"stats":{"atk":16,"def":20,"spa":16,"spd":23,"spe":18},"moves":["tackle","recover","calmmind","roost"],"baseAbility":"multiscale","item":"heavydutyboots","pokeball":"pokeball","ability":"multiscale"},{"ident":"p1: Ninetales","details":"Ninetales, L5, M","condition":"24/24","active":false,"stats":{"atk":15,"def":15,"spa":15,"spd":17,"spe":17},"moves":["willowisp","nastyplot","fireblast","solarbeam"],"baseAbility":"drought","item":"heavydutyboots","pokeball":"pokeball","ability":"drought"}]}}
sideupdate
p2
|request|{"teamPreview":true,"maxTeamSize":6,"side":{"name":"p2","id":"p2","pokemon":[{"ident":"p2: Zoroark","details":"Zoroark, L5, M","condition":"23/23","active":true,"stats":{"atk":18,"def":13,"spa":19,"spd":13,"spe":18},"moves":["transform","darkpulse","sludgebomb","flamethrower"],"baseAbility":"intimidate","item":"lifeorb","pokeball":"pokeball","ability":"intimidate"},{"ident":"p2: Umbreon","details":"Umbreon, L5, M","condition":"27/27","active":false,"stats":{"atk":14,"def":18,"spa":13,"spd":20,"spe":14},"moves":["protect","foulplay","wish","toxic"],"baseAbility":"synchronize","item":"leftovers","pokeball":"pokeball","ability":"synchronize"}]}}
sideupdate
p1
|request|{"active":[{"moves":[{"move":"Tackle","id":"tackle","pp":56,"maxpp":56,"target":"normal","disabled":false},{"move":"Recover","id":"recover","pp":16,"maxpp":16,"target":"self","disabled":false},{"move":"Calm Mind","id":"calmmind","pp":32,"maxpp":32,"target":"self","disabled":false},{"move":"Roost","id":"roost","pp":16,"maxpp":16,"target":"self","disabled":false}],"canDynamax":true,"maxMoves":{"maxMoves":[{"move":"maxstrike","target":"adjacentFoe"},{"move":"maxguard","target":"self"},{"move":"maxguard","target":"self"},{"move":"maxguard","target":"self"}]}}],"side":{"name":"p1","id":"p1","pokemon":[{"ident":"p1: Lugia","details":"Lugia, L5","condition":"28/28","active":true,"stats":{"atk":16,"def":20,"spa":16,"spd":23,"spe":18},"moves":["tackle","recover","calmmind","roost"],"baseAbility":"multiscale","item":"heavydutyboots","pokeball":"pokeball","ability":"multiscale"},{"ident":"p1: Ninetales","details":"Ninetales, L5, M","condition":"24/24","active":false,"stats":{"atk":15,"def":15,"spa":15,"spd":17,"spe":17},"moves":["willowisp","nastyplot","fireblast","solarbeam"],"baseAbility":"drought","item":"heavydutyboots","pokeball":"pokeball","ability":"drought"}]}}
sideupdate
p2
|request|{"active":[{"moves":[{"move":"Transform","id":"transform","pp":16,"maxpp":16,"target":"normal","disabled":false},{"move":"Dark Pulse","id":"darkpulse","pp":24,"maxpp":24,"target":"any","disabled":false},{"move":"Sludge Bomb","id":"sludgebomb","pp":16,"maxpp":16,"target":"normal","disabled":false},{"move":"Flamethrower","id":"flamethrower","pp":24,"maxpp":24,"target":"normal","disabled":false}],"canDynamax":true,"maxMoves":{"maxMoves":[{"move":"maxguard","target":"self"},{"move":"maxdarkness","target":"adjacentFoe"},{"move":"maxooze","target":"adjacentFoe"},{"move":"maxflare","target":"adjacentFoe"}]}}],"side":{"name":"p2","id":"p2","pokemon":[{"ident":"p2: Zoroark","details":"Zoroark, L5, M","condition":"23/23","active":true,"stats":{"atk":18,"def":13,"spa":19,"spd":13,"spe":18},"moves":["transform","darkpulse","sludgebomb","flamethrower"],"baseAbility":"intimidate","item":"lifeorb","pokeball":"pokeball","ability":"intimidate"},{"ident":"p2: Umbreon","details":"Umbreon, L5, M","condition":"27/27","active":false,"stats":{"atk":14,"def":18,"spa":13,"spd":20,"spe":14},"moves":["protect","foulplay","wish","toxic"],"baseAbility":"synchronize","item":"leftovers","pokeball":"pokeball","ability":"synchronize"}]}}
update
|t:|1609958280
|player|p1|p1||
|player|p2|p2||
|teamsize|p1|2
|teamsize|p2|2
|gametype|singles
|gen|8
|tier|[Gen 8] Anything Goes
|clearpoke
|poke|p1|Lugia, L5|
|poke|p1|Ninetales, L5, M|
|poke|p2|Zoroark, L5, M|
|poke|p2|Umbreon, L5, M|
|rule|HP Percentage Mod: HP is shown in percentages
|rule|Endless Battle Clause: Forcing endless battles is banned
|teampreview
|
|t:|1609958280
|start
|split|p1
|switch|p1a: Lugia|Lugia, L5|28/28
|switch|p1a: Lugia|Lugia, L5|100/100
|split|p2
|switch|p2a: Zoroark|Zoroark, L5, M|23/23
|switch|p2a: Zoroark|Zoroark, L5, M|100/100
|-ability|p2a: Zoroark|Intimidate|boost
|-unboost|p1a: Lugia|atk|1
|turn|1`

func TestParseStdout(t *testing.T) {
	st := structs.NewBattleState()

	msgs, err := parseStdout(dataTestParseStdout, st)

	assert.Nil(t, err)
	assert.Equal(t, 4, len(msgs))
	assert.Equal(t, 2, len(st.Field))

	assert.Equal(t, 1, st.Turn)
	assert.Equal(t, "", st.Winner)

	p1, ok := st.Field["p1"]
	assert.True(t, ok)

	p2, ok := st.Field["p2"]
	assert.True(t, ok)

	assert.Equal(t, "p1: Lugia", p1.Pokemon[0].Ident)
	assert.Equal(t, "p1: Ninetales", p1.Pokemon[1].Ident)
	assert.Equal(t, "p2: Zoroark", p2.Pokemon[0].Ident)
	assert.Equal(t, "p2: Umbreon", p2.Pokemon[1].Ident)
}
