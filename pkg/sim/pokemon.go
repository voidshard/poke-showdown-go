package sim

import (
	"log"

	"github.com/voidshard/poke-showdown-go/pkg/internal/structs"
	data "github.com/voidshard/poke-showdown-go/pkg/pokedata"
)

// Pokemon is the final struct that amasses all the data we can manage to find.
type Pokemon struct {
	ID string

	// Unique 'player: name' string for pokemon
	Ident string

	// Index on the team (changes based on team order)
	Index int

	// If pokemon is currently on the field
	Active bool

	// Raw pokemon stats
	Stats *data.Stats

	// Pokemon moves
	Moves []*Move

	// The pokemon's current ability
	Ability string

	// Item the pokemon is holding
	Item string

	// Pokeon status field(s)
	Status *structs.Status

	// Name, Level, Gender run together with ','
	// Nb. Level is not given if it is 100.
	// ie. Chesnaught, L82, M
	Details string

	// From observation condition is a string like
	// 30/130
	// 0 fnt
	// 130/240 slp
	// ie. HP/MaxHap status1 status2
	Condition string

	// Species is the pokemons major species
	Species string

	// Level is the pokemons current level
	Level int

	// Shiny indicates bling
	Shiny bool

	// Gender indicates male / femaleness.
	// Nb. not all pokemon have a gender
	Gender string

	// Pokedex data
	Dex *data.PokeDexItem
}

// toPokemon builds a final Pokemon from a parsed pokemon struct.
// We also add pokedex & movedex entries to flesh out the data.
func toPokemon(in *structs.Pokemon) *Pokemon {
	moves := []*Move{}
	for _, id := range in.Moves {
		m, err := NewMove(id)
		if err != nil {
			log.Printf("failed to find move '%s'\n", id)
			continue
		}
		moves = append(moves, m)
	}

	dex, err := data.PokeDex(in.Species)
	if err != nil {
		log.Println("failed to find pokemon", in.Species)
	}

	return &Pokemon{
		Ident:     in.Ident,
		Index:     in.Index,
		Active:    in.Active,
		Stats:     in.Stats,
		Moves:     moves,
		Ability:   in.Ability,
		Item:      in.Item,
		Status:    in.Status,
		Details:   in.Details,
		Condition: in.Condition,
		Species:   in.Species,
		Level:     in.Level,
		Shiny:     in.Shiny,
		Gender:    in.Gender,
		Dex:       dex,
	}
}
