package main

import (
	"github.com/voidshard/poke-showdown-go/pkg/sim"
	"github.com/voidshard/poke-showdown-go/pkg/structs"

	"fmt"
)

func main() {
	spec := &structs.BattleSpec{
		Format: "gen7randombattle",
		Players: map[string][]*structs.PokemonSpec{
			"p1": nil,
			"p2": nil,
		},
	}

	battle, err := sim.Start("pokemon-showdown", spec)
	defer battle.Close()

	fmt.Sprintln(battle, err)

	go func() {
		for err := range battle.Errors() {
			if err != nil {
				panic(err)
			}
		}
	}()

	for e := range battle.Read() {
		fmt.Println(e.Type, e.Player, e.Update)
	}
}
