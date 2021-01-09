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

func Run(spec *structs.BattleSpec) {
	battle, err := sim.Start("pokemon-showdown", spec)
	if err != nil {
		panic(err)
	}
	defer battle.Stop()

	state := battle.State()
	for {
		fmt.Printf("----- TURN %d ------\n", state.Turn)
		fmt.Printf("events: %d\n", len(state.Events))
		for _, e := range state.Events {
			fmt.Println(e.String())
		}

		if state.Winner != "" {
			fmt.Println("Winner:", state.Winner)
			return
		}

		fmt.Println("")
		printField(state.Field)

		for {
			decisions := []*structs.Action{}

			for player, data := range state.Field {
				if data.Wait {
					continue
				}

				act := askUser(player)
				decisions = append(decisions, act)
			}

			newstate, err := battle.Turn(decisions)
			if err == nil {
				state = newstate
				break
			}
			fmt.Printf("Error writing choice: %v\n", err)
		}

	}
	fmt.Println("-End-")
}

func printField(f map[string]*structs.Update) {
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
		Players: [][]*structs.PokemonSpec{
			[]*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "Ninetales",
					Item:    "heavydutyboots",
					Ability: "drought",
					Moves:   []string{"inferno", "willowisp", "nastyplot", "solarbeam"},
					Level:   1,
					Gender:  "M",
				},
			},
			[]*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "Umbreon",
					Item:    "leftovers",
					Ability: "synchronize",
					Moves:   []string{"wish", "foulplay", "protect", "toxic"},
					Level:   99,
					Gender:  "F",
				},
			},
		},
	}

	spec1v1Mega = &structs.BattleSpec{
		Format: structs.FormatGen8,
		Players: [][]*structs.PokemonSpec{
			[]*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "zoroark",
					Item:    "lifeorb",
					Ability: "illusion",
					Moves:   []string{"nastyplot", "darkpulse", "sludgebomb", "flamethrower"},
					Level:   50,
				},
			},
			[]*structs.PokemonSpec{
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
		Players: [][]*structs.PokemonSpec{
			[]*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "Lugia",
					Item:    "heavydutyboots",
					Ability: "multiscale",
					//Moves:   []string{"aeroblast", "psyshock", "calmmind", "roost"},
					Moves: []string{"tackle", "recover", "calmmind", "roost"},
					Level: 5,
				},
				&structs.PokemonSpec{
					Name:    "Ninetales",
					Item:    "heavydutyboots",
					Ability: "drought",
					Moves:   []string{"willowisp", "nastyplot", "fireblast", "solarbeam"},
					Level:   5,
				},
			},
			[]*structs.PokemonSpec{
				&structs.PokemonSpec{
					Name:    "zoroark",
					Item:    "lifeorb",
					Ability: "intimidate",
					//Moves:   []string{"nastyplot", "darkpulse", "sludgebomb", "flamethrower"},
					Moves: []string{"transform", "darkpulse", "sludgebomb", "flamethrower"},
					Level: 5,
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
	Run(spec2v2Singles)
}
