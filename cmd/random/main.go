package main

import (
	"bufio"
	"fmt"
	"github.com/voidshard/poke-showdown-go/pkg/sim"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
	"os"
	"strings"
)

type Game struct {
	// battle simulation
	battle sim.Simulation

	numPlayers int
}

func (g *Game) Run() {
	for {
		select {
		case msg := <-g.battle.Messages():
			fmt.Println("msg:", msg)
		case state := <-g.battle.Read():
			if state.Error != "" {
				fmt.Println("Error:", state.Error)
			}

			if state.Winner != "" {
				fmt.Println("Winner:", state.Winner)
				return
			}

			printField(state.Field)

			acts := []*structs.Action{}
			for player, data := range state.Field {
				if data.Wait {
					continue
				}
				act := askUser(player)
				acts = append(acts, act)
			}

			for _, act := range acts {
				g.battle.Write(act)
			}

		}
	}
	fmt.Println("-End-")
}

func printField(f map[string]*structs.Update) {
	for player, update := range f {
		fmt.Printf("Player: %s [wait:%v] [switch:%v]\n", player, update.Wait, update.ForceSwitch)
		active := []structs.Pokemon{}
		for i, pkm := range update.Team.Pokemon {
			if pkm.Active {
				active = append(active, pkm)
			}
			fmt.Printf("    [%d] %s %s\n", i+1, pkm.Details, pkm.Condition)
		}

		fmt.Println("  [Active]")
		for _, pkm := range active {
			fmt.Printf("    %s %s %s %s %v\n", pkm.Details, pkm.Condition, pkm.Item, pkm.Ability, pkm.Moves)
		}
	}
}

func askUser(player string) *structs.Action {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(player, " => ")
	text, _ := reader.ReadString('\n')

	bits := strings.Split(text, " ")
	atype := structs.ActionMove
	if bits[0] == "switch" {
		atype = structs.ActionSwitch
	}

	act := &structs.Action{
		Player: player,
		Specs: []structs.ActionSpec{
			structs.ActionSpec{
				Type:  atype,
				Value: bits[1],
			},
		},
	}

	if len(bits) > 2 {
		switch bits[2] {
		case "mega":
			act.Specs[0].Mega = true
		case "zmove":
			act.Specs[0].ZMove = true
		case "max":
			act.Specs[0].Max = true
		}
	}

	return act
}

/*
{"name":"Ninetales","species":"Ninetales","item":"heavydutyboots","ability":"drought","moves":["willowisp","nastyplot","fireblast","solarbeam"],"nature":"","evs":{"hp":85,"atk":0,"def":85,"spa":85,"spd":85,"spe":85},"ivs":{"hp":31,"atk":0,"def":31,"spa":31,"spd":31,"spe":31},"level":86}

{"name":"Umbreon","species":"Umbreon","item":"leftovers","ability":"synchronize","moves":["protect","foulplay","wish","toxic"],"nature":"","evs":{"hp":85,"atk":85,"def":85,"spa":85,"spd":85,"spe":85},"level":84}
*/

var (
	spec1v1 = &structs.BattleSpec{
		Format: structs.FormatGen8,
		Players: map[string][]*structs.PokemonSpec{
			"p1": []*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "Ninetales",
					Item:    "heavydutyboots",
					Ability: "drought",
					Moves:   []string{"willowisp", "nastyplot", "fireblast", "solarbeam"},
					Level:   50,
				},
			},
			"p2": []*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "Umbreon",
					Item:    "leftovers",
					Ability: "synchronize",
					Moves:   []string{"protect", "foulplay", "wish", "toxic"},
					Level:   50,
				},
			},
		},
	}

	/*
		{"name":"Kommo-o","species":"Kommo-o","item":"throatspray","ability":"soundproof","moves":["closecombat","clangingscales","clangoroussoul","poisonjab"],"nature":"","evs":{"hp":85,"atk":85,"def":85,"spa":85,"spd":85,"spe":85},"level":80}
		{"name":"Lugia","species":"Lugia","item":"heavydutyboots","ability":"multiscale","moves":["aeroblast","psyshock","calmmind","roost"],"nature":"","evs":{"hp":85,"atk":0,"def":85,"spa":85,"spd":85,"spe":85},"gender":"N","ivs":{"hp":31,"atk":0,"def":31,"spa":31,"spd":31,"spe":31},"level":72}
	*/
	spec2v2 = &structs.BattleSpec{
		Format: structs.FormatGen8,
		Players: map[string][]*structs.PokemonSpec{
			"p1": []*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "Lugia",
					Item:    "heavydutyboots",
					Ability: "multiscale",
					Moves:   []string{"aeroblast", "psyshock", "calmmind", "roost"},
					Level:   5,
				},
				&structs.PokemonSpec{
					Name:    "Ninetales",
					Item:    "heavydutyboots",
					Ability: "drought",
					Moves:   []string{"willowisp", "nastyplot", "fireblast", "solarbeam"},
					Level:   5,
				},
			},
			"p2": []*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "Kommo-o",
					Item:    "throatspray",
					Ability: "soundproof",
					Moves:   []string{"closecombat", "clangingscales", "clangoroussoul", "poisonjab"},
					Level:   5,
				},
				&structs.PokemonSpec{
					Name:    "Umbreon",
					Item:    "leftovers",
					Ability: "synchronize",
					Moves:   []string{"protect", "foulplay", "wish", "toxic"},
					Level:   5,
				},
			},
		},
	}
)

func main() {
	launch(spec2v2)
}

func launch(spec *structs.BattleSpec) {
	// this is a quick hack to test run pokemon battles
	battle, err := sim.Start("pokemon-showdown", spec)
	if err != nil {
		panic(err)
	}
	defer battle.Close()

	game := &Game{
		battle:     battle,
		numPlayers: len(spec.Players),
	}

	go func() {
		for err := range battle.Errors() {
			if err != nil {
				panic(err)
			}
		}
	}()

	game.Run()
}
