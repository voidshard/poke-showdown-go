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
		fmt.Printf("Player: %s [wait:%v] [switch:%v] [active:%v]\n", player, update.Wait, update.ForceSwitch, update.Active)
		active := []*structs.Pokemon{}
		for i, pkm := range update.Pokemon {
			if pkm.Active {
				active = append(active, pkm)
			}
			fmt.Printf("    [%d] %s %s\n", i, pkm.Details, pkm.Condition)
		}

		fmt.Println("  [Active]")
		for _, pkm := range active {
			bonus := []string{}
			if pkm.Options != nil {
				if pkm.Options.CanDynamax {
					bonus = append(bonus, "[max]")
				}
				if pkm.Options.CanMegaEvolve {
					bonus = append(bonus, "[mega]")
				}
				if pkm.Options.CanZMove {
					bonus = append(bonus, "[zmove]")
				}
			}
			fmt.Printf("    %s %s %s %s %v %s\n", pkm.Details, pkm.Condition, pkm.Item, pkm.Ability, pkm.Moves, strings.Join(bonus, ""))
			fmt.Printf(
				"        Atk:%d Def:%d Spa:%d Sde:%d Spe:%d\n",
				pkm.Stats.Attack,
				pkm.Stats.Defense,
				pkm.Stats.SpecialAttack,
				pkm.Stats.SpecialDefense,
				pkm.Stats.Speed,
			)
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

var (
	spec1v1 = &structs.BattleSpec{
		Format: structs.FormatGen8,
		Players: map[string][]*structs.PokemonSpec{
			"p1": []*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "Ninetales",
					Item:    "heavydutyboots",
					Ability: "drought",
					//Moves:   []string{"willowisp", "nastyplot", "fireblast", "solarbeam"},
					Moves:  []string{"willowisp", "spore", "substitute", "poisonpowder"},
					Level:  50,
					Gender: "M",
				},
			},
			"p2": []*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "Umbreon",
					Item:    "leftovers",
					Ability: "synchronize",
					Moves:   []string{"protect", "foulplay", "wish", "toxic"},
					//Moves:  []string{"spore", "attract", "wish", "toxic"},
					Level:  50,
					Gender: "F",
				},
			},
		},
	}

	spec1v1Mega = &structs.BattleSpec{
		Format: structs.FormatGen8,
		Players: map[string][]*structs.PokemonSpec{
			"p1": []*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "zoroark",
					Item:    "lifeorb",
					Ability: "illusion",
					Moves:   []string{"nastyplot", "darkpulse", "sludgebomb", "flamethrower"},
					Level:   50,
				},
			},
			"p2": []*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "gallade",
					Item:    "galladite",
					Ability: "justified",
					Moves:   []string{"swordsdance", "closecombat", "zenheadbutt", "knockoff"},
					Level:   50,
				},
			},
		},
	}

	spec2v2 = &structs.BattleSpec{
		Format: structs.FormatGen8Doubles,
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
					Name:    "zoroark",
					Item:    "lifeorb",
					Ability: "illusion",
					Moves:   []string{"nastyplot", "darkpulse", "sludgebomb", "flamethrower"},
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

	spec2v2Singles = &structs.BattleSpec{
		Format:  structs.FormatGen8,
		Players: spec2v2.Players,
	}
)

func main() {
	launch(spec1v1)
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
