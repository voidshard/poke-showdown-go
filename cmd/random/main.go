package main

import (
	"bufio"
	"fmt"
	"github.com/voidshard/poke-showdown-go/pkg/sim"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
	"os"
	"strconv"
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

			for player, data := range state.Field {
				if data.Wait {
					continue
				}

				for {
					act := askUser(player)
					err := g.battle.Write(act)
					if err == nil {
						break
					} else {
						fmt.Printf("Err: %v\n", err)
					}
				}
			}

		}
	}
	fmt.Println("-End-")
}

func printField(f map[string]*structs.Update) {
	fmt.Println("---")
	for player, update := range f {
		fmt.Printf("Player: %s [wait:%v] [switch:%v]\n", player, update.Wait, update.ForceSwitch)
		active := []*structs.Pokemon{}
		for i, pkm := range update.Team.Pokemon {
			if pkm.Active {
				active = append(active, pkm)
			}
			fmt.Printf("    [%d] %s %s\n", i+1, pkm.Details, pkm.Condition)
		}

		fmt.Println("  [Active]")
		for _, pkm := range active {
			bonus := []string{}
			if pkm.Options.CanDynamax {
				bonus = append(bonus, "[max]")
			}
			if pkm.Options.CanMegaEvolve {
				bonus = append(bonus, "[mega]")
			}
			if pkm.Options.CanZMove {
				bonus = append(bonus, "[zmove]")
			}
			fmt.Printf("    %s %s %s %s %v %s\n", pkm.Details, pkm.Condition, pkm.Item, pkm.Ability, pkm.Moves, strings.Join(bonus, ""))
		}
		fmt.Println("")
	}
}

func askUser(player string) *structs.Action {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(player, " => ")
	text, _ := reader.ReadString('\n')

	act := &structs.Action{Player: player, Specs: []*structs.ActionSpec{}}

	bits := strings.Split(strings.TrimSpace(text), ",")
	for _, b := range bits {
		spec := parseSpec(b)
		act.Specs = append(act.Specs, spec)
	}

	return act
}

func parseSpec(move string) *structs.ActionSpec {
	bits := strings.Split(strings.TrimSpace(move), " ")

	atype := structs.ActionMove
	if bits[0] == "switch" {
		atype = structs.ActionSwitch
	}

	id, err := strconv.ParseInt(bits[1], 10, 64)
	if err != nil {
		panic(err)
	}

	spec := &structs.ActionSpec{
		Type: atype,
		ID:   int(id),
	}

	if len(bits) > 2 {
		id, err = strconv.ParseInt(bits[2], 10, 64)
		if err != nil {
			panic(err)
		}

		spec.Target = int(id)
	}

	if len(bits) > 3 {
		switch bits[3] {
		case "mega":
			spec.Mega = true
		case "zmove":
			spec.ZMove = true
		case "max":
			spec.Max = true
		}
	}

	return spec
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

	game.Run()
}
