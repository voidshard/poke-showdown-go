package main

import (
	"bufio"
	"fmt"
	"github.com/voidshard/poke-showdown-go/pkg/sim"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
	"os"
	"strings"
)

const human = "p1"

type Game struct {
	// battle simulation
	battle sim.Simulation

	numPlayers int
}

func (g *Game) Run() {
	for {
		field := map[string]*structs.BattleUpdate{}
		for evt := range g.battle.Read() {
			if evt.Type == structs.EventSideUpdate {
				field[evt.Player] = evt.Update
				if len(field) == g.numPlayers {
					break
				}
			}
			fmt.Println(evt, field)
		}

		printField(field)

		acts := decideActions(field)
		for _, act := range acts {
			g.battle.Write(act)
		}
	}
}

func printField(field map[string]*structs.BattleUpdate) {
	fmt.Println("---")
	for player, side := range field {
		active := []structs.Pokemon{}

		for i, pkm := range side.Team.Pokemon {
			fmt.Printf("[%d] %s\n", i+1, pkm.Details)
			if pkm.Active {
				active = append(active, pkm)
			}
		}

		for _, pkm := range active {
			fmt.Println(player, pkm.Details, pkm.Condition)
			fmt.Println(pkm.Item, pkm.Ability)
			fmt.Println(pkm.Moves)
			fmt.Println()
		}
	}
}

func decideActions(field map[string]*structs.BattleUpdate) []*structs.Action {
	acts := []*structs.Action{}
	for player, data := range field {
		if data.Wait {
			continue
		}
		if len(data.ForceSwitch) > 0 {
			fmt.Println("[switch]")
			for i, pkm := range data.Team.Pokemon {
				fmt.Printf("\t[%d] %s [%s %s %s] %v\n", i+1, pkm.Details, pkm.Condition, pkm.Item, pkm.Ability, pkm.Moves)
			}
		}
		action := askPlayer(player)
		acts = append(acts, action)
	}
	return acts
}

func askPlayer(player string) *structs.Action {
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

func main() {
	// this is a quick hack to test run pokemon battles

	spec := &structs.BattleSpec{
		Format: "gen7randombattle",
		Players: map[string][]*structs.PokemonSpec{
			human: nil,
			"p2":  nil,
		},
	}

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
