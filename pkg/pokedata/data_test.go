package pokedata

import (
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	umbreon   = new(PokeDexItem)
	charizard = new(PokeDexItem)
	venusaur  = new(PokeDexItem)

	thunder    = new(MoveDexItem)
	protect    = new(MoveDexItem)
	shellsmash = new(MoveDexItem)
	flail      = new(MoveDexItem)
	flatter    = new(MoveDexItem)
)

func init() {
	dataThunder := []byte(`{"num":87,"accuracy":70,"basePower":110,"category":"Special","name":"Thunder","pp":10,"priority":0,"flags":{"protect":1,"mirror":1},"secondary":{"chance":30,"status":"par"},"target":"normal","type":"Electric","contestType":"Cool","desc":"Has a 30% chance to paralyze the target. This move can hit a target using Bounce, Fly, or Sky Drop, or is under the effect of Sky Drop. If the weather is Primordial Sea or Rain Dance, this move does not check accuracy. If the weather is Desolate Land or Sunny Day, this move's accuracy is 50%. If this move is used against a Pokemon holding Utility Umbrella, this move's accuracy remains at 70%.","shortDesc":"30% chance to paralyze. Can't miss in rain."}`)
	dataFlail := []byte(`{"num":175,"accuracy":100,"basePower":0,"category":"Physical","name":"Flail","pp":15,"priority":0,"flags":{"contact":1,"protect":1,"mirror":1},"secondary":null,"target":"normal","type":"Normal","zMove":{"basePower":160},"maxMove":{"basePower":130},"contestType":"Cute","desc":"The power of this move is 20 if X is 33 to 48, 40 if X is 17 to 32, 80 if X is 10 to 16, 100 if X is 5 to 9, 150 if X is 2 to 4, and 200 if X is 0 or 1, where X is equal to (user's current HP * 48 / user's maximum HP), rounded down.","shortDesc":"More power the less HP the user has left."}`)
	dataFlatter := []byte(`{"num":260,"accuracy":100,"basePower":0,"category":"Status","name":"Flatter","pp":15,"priority":0,"flags":{"protect":1,"reflectable":1,"mirror":1,"mystery":1},"volatileStatus":"confusion","boosts":{"spa":1},"secondary":null,"target":"normal","type":"Dark","zMove":{"boost":{"spd":1}},"contestType":"Clever","desc":"Raises the target's Special Attack by 1 stage and confuses it.","shortDesc":"Raises the target's Sp. Atk by 1 and confuses it."}`)
	dataShellsmash := []byte(`{"num":504,"accuracy":true,"basePower":0,"category":"Status","name":"Shell Smash","pp":15,"priority":0,"flags":{"snatch":1},"boosts":{"def":-1,"spd":-1,"atk":2,"spa":2,"spe":2},"secondary":null,"target":"self","type":"Normal","zMove":{"effect":"clearnegativeboost"},"contestType":"Tough","desc":"Lowers the user's Defense and Special Defense by 1 stage. Raises the user's Attack, Special Attack, and Speed by 2 stages.","shortDesc":"Lowers Def, SpD by 1; raises Atk, SpA, Spe by 2."}`)
	dataProtect := []byte(`{"num":182,"accuracy":true,"basePower":0,"category":"Status","name":"Protect","pp":10,"priority":4,"flags":{},"stallingMove":true,"volatileStatus":"protect","condition":{"duration":1,"onTryHitPriority":3},"secondary":null,"target":"self","type":"Normal","zMove":{"effect":"clearnegativeboost"},"contestType":"Cute","desc":"The user is protected from most attacks made by other Pokemon during this turn. This move has a 1/X chance of being successful, where X starts at 1 and triples each time this move is successfully used. X resets to 1 if this move fails, if the user's last move used is not Baneful Bunker, Detect, Endure, King's Shield, Obstruct, Protect, Quick Guard, Spiky Shield, or Wide Guard, or if it was one of those moves and the user's protection was broken. Fails if the user moves last this turn.","shortDesc":"Prevents moves from affecting the user this turn."}`)

	err := json.Unmarshal(dataProtect, protect)
	if err != nil {
		panic(err)
	}
	protect.Accuracy = 1000

	err = json.Unmarshal(dataThunder, thunder)
	if err != nil {
		panic(err)
	}
	thunder.Accuracy = 70

	err = json.Unmarshal(dataShellsmash, shellsmash)
	if err != nil {
		panic(err)
	}
	shellsmash.Accuracy = 1000

	err = json.Unmarshal(dataFlail, flail)
	if err != nil {
		panic(err)
	}
	flail.Accuracy = 100

	err = json.Unmarshal(dataFlatter, flatter)
	if err != nil {
		panic(err)
	}
	flatter.Accuracy = 100

	dataUmbreon := []byte(`{"num":197,"name":"Umbreon","types":["Dark"],"genderRatio":{"M":0.875,"F":0.125},"baseStats":{"hp":95,"atk":65,"def":110,"spa":60,"spd":130,"spe":65},"abilities":{"0":"Synchronize","H":"Inner Focus"},"heightm":1,"weightkg":27,"color":"Black","prevo":"Eevee","evoType":"levelFriendship","evoCondition":"at night","eggGroups":["Field"],"tier":"RU"}`)
	dataCharmegay := []byte(`{"num":6,"name":"Charizard-Mega-Y","baseSpecies":"Charizard","forme":"Mega-Y","types":["Fire","Flying"],"genderRatio":{"M":0.875,"F":0.125},"baseStats":{"hp":78,"atk":104,"def":78,"spa":159,"spd":115,"spe":100},"abilities":{"0":"Drought"},"heightm":1.7,"weightkg":100.5,"color":"Red","eggGroups":["Monster","Dragon"],"requiredItem":"Charizardite Y","tier":"Illegal","isNonstandard":"Past"}`)
	dataVengmax := []byte(`{"num":3,"name":"Venusaur-Gmax","baseSpecies":"Venusaur","forme":"Gmax","types":["Grass","Poison"],"genderRatio":{"M":0.875,"F":0.125},"baseStats":{"hp":80,"atk":82,"def":83,"spa":100,"spd":100,"spe":80},"abilities":{"0":"Overgrow","H":"Chlorophyll"},"heightm":2,"weightkg":0,"color":"Green","eggGroups":["Monster","Grass"],"changesFrom":"Venusaur","tier":"AG","isNonstandard":"Gigantamax"}`)

	err = json.Unmarshal(dataUmbreon, umbreon)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(dataCharmegay, charizard)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(dataVengmax, venusaur)
	if err != nil {
		panic(err)
	}
}

var dataTestMoveDex = []struct {
	Name   string
	Expect *MoveDexItem
}{
	{
		"protect",
		protect,
	},
	{
		"Protect",
		protect,
	},
	{
		"flail",
		flail,
	},
	{
		"flatter",
		flatter,
	},
	{
		"shellsmash",
		shellsmash,
	},
	{
		"Thunder",
		thunder,
	},
	{
		"thunder",
		thunder,
	},
}

var dataTestParseAccuracy = []struct {
	Name   string
	In     interface{}
	Expect int
}{
	{"bool", true, 1000},
	{"int", 54, 54},
}

func TestParseAccuracy(t *testing.T) {
	for _, tt := range dataTestParseAccuracy {
		t.Run(tt.Name, func(t *testing.T) {
			result := parseAccuracy(tt.In)

			assert.Equal(t, tt.Expect, result)
		})
	}
}

func TestMoveDex(t *testing.T) {
	for _, tt := range dataTestMoveDex {
		t.Run(tt.Name, func(t *testing.T) {
			result, err := MoveDex(tt.Name)

			assert.Nil(t, err)
			assert.Equal(t, tt.Expect, result)
		})
	}
}

var dataTestPokeDex = []struct {
	Name   string
	Expect *PokeDexItem
}{
	{
		"umbreon",
		umbreon,
	},
	{
		"Umbreon",
		umbreon,
	},
	{
		"charizardmegay",
		charizard,
	},
	{
		"Charizard-Mega-Y",
		charizard,
	},
	{
		"venusaurgmax",
		venusaur,
	},
	{
		"Venusaur-Gmax",
		venusaur,
	},
}

func TestPokeDex(t *testing.T) {
	for _, tt := range dataTestPokeDex {
		t.Run(tt.Name, func(t *testing.T) {
			result, err := PokeDex(tt.Name)

			assert.Nil(t, err)
			assert.Equal(t, tt.Expect, result)
		})
	}
}
