package parse

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var dataTestParseStdout = `p1
|request|{"teamPreview":true,"maxTeamSize":6,"side":{"name":"p1","id":"p1","pokemon":[{"ident":"p1: Pincurchin","details":"Pincurchin, L88, M","condition":"228/228","active":true,"stats":{"atk":228,"def":217,"spa":210,"spd":200,"spe":77},"moves":["suckerpunch","risingvoltage","spikes","scald"],"baseAbility":"electricsurge","item":"focussash","pokeball":"pokeball","ability":"electricsurge"}]}}
2021/04/02 08:46:00 msg: sideupdate
p2
|request|{"teamPreview":true,"maxTeamSize":6,"side":{"name":"p2","id":"p2","pokemon":[{"ident":"p2: Liepard","details":"Liepard, L88, M","condition":"256/256","active":true,"stats":{"atk":205,"def":138,"spa":205,"spd":138,"spe":237},"moves":["uturn","knockoff","copycat","encore"],"baseAbility":"prankster","item":"focussash","pokeball":"pokeball","ability":"prankster"}]}}
2021/04/02 08:46:00 msg: sideupdate
p1
|request|{"active":[{"moves":[{"move":"Sucker Punch","id":"suckerpunch","pp":8,"maxpp":8,"target":"normal","disabled":false},{"move":"Rising Voltage","id":"risingvoltage","pp":32,"maxpp":32,"target":"normal","disabled":false},{"move":"Spikes","id":"spikes","pp":32,"maxpp":32,"target":"foeSide","disabled":false},{"move":"Scald","id":"scald","pp":24,"maxpp":24,"target":"normal","disabled":false}],"canDynamax":true,"maxMoves":{"maxMoves":[{"move":"maxdarkness","target":"adjacentFoe"},{"move":"maxlightning","target":"adjacentFoe"},{"move":"maxguard","target":"self"},{"move":"maxgeyser","target":"adjacentFoe"}]}}],"side":{"name":"p1","id":"p1","pokemon":[{"ident":"p1: Pincurchin","details":"Pincurchin, L88, M","condition":"228/228","active":true,"stats":{"atk":228,"def":217,"spa":210,"spd":200,"spe":77},"moves":["suckerpunch","risingvoltage","spikes","scald"],"baseAbility":"electricsurge","item":"focussash","pokeball":"pokeball","ability":"electricsurge"}]}}
2021/04/02 08:46:00 msg: sideupdate
p2
|request|{"active":[{"moves":[{"move":"U-turn","id":"uturn","pp":32,"maxpp":32,"target":"normal","disabled":false},{"move":"Knock Off","id":"knockoff","pp":32,"maxpp":32,"target":"normal","disabled":false},{"move":"Copycat","id":"copycat","pp":32,"maxpp":32,"target":"self","disabled":false},{"move":"Encore","id":"encore","pp":8,"maxpp":8,"target":"normal","disabled":false}],"canDynamax":true,"maxMoves":{"maxMoves":[{"move":"maxflutterby","target":"adjacentFoe"},{"move":"maxdarkness","target":"adjacentFoe"},{"move":"maxguard","target":"self"},{"move":"maxguard","target":"self"}]}}],"side":{"name":"p2","id":"p2","pokemon":[{"ident":"p2: Liepard","details":"Liepard, L88, M","condition":"256/256","active":true,"stats":{"atk":205,"def":138,"spa":205,"spd":138,"spe":237},"moves":["uturn","knockoff","copycat","encore"],"baseAbility":"prankster","item":"focussash","pokeball":"pokeball","ability":"prankster"}]}}
2021/04/02 08:46:00 msg: update
|t:|1617349560
|player|p1|p1||
|player|p2|p2||
|teamsize|p1|1
|teamsize|p2|1
|gametype|singles
|gen|8
|tier|[Gen 8] Anything Goes
|clearpoke
|poke|p1|Pincurchin, L88, M|
|poke|p2|Liepard, L88, M|
|rule|HP Percentage Mod: HP is shown in percentages
|rule|Endless Battle Clause: Forcing endless battles is banned
|teampreview
|
|t:|1617349560
|start
|split|p1
|switch|p1a: Pincurchin|Pincurchin, L88, M|228/228
|switch|p1a: Pincurchin|Pincurchin, L88, M|100/100
|split|p2
|switch|p2a: Liepard|Liepard, L88, M|256/256
|switch|p2a: Liepard|Liepard, L88, M|100/100
|-fieldstart|move: Electric Terrain|[from] ability: Electric Surge|[of] p1a: Pincurchin
|turn|1`

func TestParseStdout(t *testing.T) {
	proc := &Process{messages: make(chan *Message)}
	msgs := []*Message{}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for m := range proc.messages {
			msgs = append(msgs, m)
		}
	}()

	proc.parseStdout(dataTestParseStdout)

	close(proc.messages)
	wg.Wait()

	assert.Equal(t, 6, len(msgs))
	assert.NotNil(t, msgs[0].Update)
	assert.NotNil(t, msgs[1].Update)
	assert.NotNil(t, msgs[2].Event)
	assert.NotNil(t, msgs[3].Event)
	assert.NotNil(t, msgs[4].Event)
	assert.NotNil(t, msgs[5].Event)
	assert.Equal(t, 0, msgs[0].Num)
	assert.Equal(t, 1, msgs[1].Num)
	assert.Equal(t, 2, msgs[2].Num)
	assert.Equal(t, 3, msgs[3].Num)
	assert.Equal(t, 4, msgs[4].Num)
	assert.Equal(t, 5, msgs[5].Num)
	assert.Equal(t, "switch", msgs[2].Event.Type)
	assert.Equal(t, "switch", msgs[3].Event.Type)
	assert.Equal(t, "-fieldstart", msgs[4].Event.Type)
	assert.Equal(t, "turn", msgs[5].Event.Type)
}
