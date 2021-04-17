/* main package here is only to test / demo the library.
Code in here should not be relied on and is at best simply a rough sketch of how one might utilise the library.
*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/voidshard/poke-showdown-go/pkg/pokeutils"
	"github.com/voidshard/poke-showdown-go/pkg/sim"

	"github.com/rivo/tview"
)

func init() {
	f, err := os.OpenFile("tui.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

func main() {
	//format := sim.FormatGen8Doubles
	format := sim.FormatGen8

	teama, err := pokeutils.RandomTeam()
	if err != nil {
		panic(err)
	}

	teamb, err := pokeutils.RandomTeam()
	if err != nil {
		panic(err)
	}

	// set some IDs for the lols
	for i, p := range teama {
		p.ID = fmt.Sprintf("teama-%d", i)
	}
	for i, p := range teamb {
		p.ID = fmt.Sprintf("teamb-%d", i)
	}

	game, err := newGame(&sim.BattleSpec{
		Format:  format,
		Players: [][]*sim.PokemonSpec{teama, teamb},
	})
	if err != nil {
		panic(err)
	}
	defer game.stop()

	grid := tview.NewGrid().SetColumns(0, 0, 0, 0, 0).SetRows(-2, -2, 0).SetBorders(true)
	grid.AddItem(game.p1Opts, 0, 0, 5, 1, 0, 0, false)
	grid.AddItem(game.p2Opts, 0, 4, 5, 1, 0, 0, false)
	grid.AddItem(game.p1a, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(game.p1b, 1, 2, 1, 1, 0, 0, false)
	grid.AddItem(game.p2a, 0, 3, 1, 1, 0, 0, false)
	grid.AddItem(game.p2b, 0, 2, 1, 1, 0, 0, false)
	grid.AddItem(game.output, 3, 1, 2, 3, 0, 0, false)

	done := make(chan error)

	go func() {
		game.run()
		done <- nil
	}()

	go func() {
		done <- tview.NewApplication().SetRoot(grid, true).EnableMouse(true).Run()
	}()

	err = <-done
	if err != nil {
		panic(err)
	}
}
