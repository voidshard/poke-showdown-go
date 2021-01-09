package main

import (
	"fmt"
	"github.com/voidshard/poke-showdown-go/pkg/pokeutils"
)

func main() {
	team, err := pokeutils.RandomTeam()
	if err != nil {
		panic(err)
	}
	fmt.Println(team)
}
