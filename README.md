# poke-showdown-go
Golang wrapper for [pokemon-showdown battle simulator](https://github.com/smogon/pokemon-showdown/blob/master/sim/README.md).


### What

This library provides a Golang interface for interacting with the pokemon-showdown battle simulator. In short:
- A battle can be started with a given pokemon team and/or a random team in doubles or singles formats (Gen8)
- Decision(s) of players (human or otherwise) can be written to the simulator interface
- Once all players have decided what to do a new battle state is returned


### Installing

In addition to the library itself you'll need to install [the engine](https://github.com/smogon/pokemon-showdown/blob/master/sim/README.md). Anecdotally this seems to require a fairly recent version of [NodeJS](https://nodejs.org/en/download/) to be available (this lib was built using v14.15.0). 

You'll need to make sure the Pokemon-Showdown JS lib is version [0.11.4](https://github.com/smogon/pokemon-showdown/commit/baaeb1e23bd1d59e7690568e36da510bfa540d03)+. 


### Example

A simple spec for a 1v1 battle

```golang
import (
	"github.com/voidshard/poke-showdown-go/pkg/sim"
)

var spec = &sim.BattleSpec{
	Format: sim.FormatGen8,
	Players: [][]*sim.PokemonSpec{
		[]*sim.PokemonSpec{
			&sim.PokemonSpec{
				Name:    "Ninetales",
				Item:    "heavydutyboots",
				Ability: "drought",
				Moves:   []string{
                                    "willowisp", 
                                    "nastyplot", 
                                    "fireblast", 
                                    "solarbeam",
                                },
				Level:   50,
			},
		},
		[]*sim.PokemonSpec{
			&sim.PokemonSpec{
				Name:    "Umbreon",
				Item:    "leftovers",
				Ability: "synchronize",
				Moves:   []string{
                                    "protect", 
                                    "foulplay",
                                    "wish",
                                    "toxic",
                                },
				Level:   50,
			},
		},
	},
}
```
Before a battle is started we do check some things like 
- there are two players with at least one pokemon each (two each in doubles)
- pokemon are found in showdown's [pokedex](https://play.pokemonshowdown.com/data/pokedex.json)
- moves are found in showdown's [moves](https://play.pokemonshowdown.com/data/moves.json)
- given pokemon have the listed ability
We also clamp down on IVs (0-31) and EVs (0-225) & a few other basic things.


Once we have a spec we can start a new battle with 
```golang
battle, _ := sim.NewSimulatorStream(spec)
defer battle.Stop()
```
This kicks off a battle as an event stream. Events, errors and requests for input are sent to us as they occur by the simulator as 'Updates' (more specifically each Update is either an error, event or side update). 

```
for update := range battle.Updates() {
    // do something
}
```

Asynchronously we can push player decisions into the simulator as either
- Move (a pokemon should use some technique)
- Switch (we want to switch pokemon)
- Pass (the given pokemon should do nothing)
```golang
battle.Write(&structs.Action{
       Player: "p1",
       Specs: []*structs.ActionSpec{
           &structs.ActionSpec{Type: structs.ActionMove, ID: "tackle"},
       },
})
```
Note that 'Specs' in the Action struct is a list, so in doubles two specs are expected per player per decision.


You can find a trivial demo terminal UI application in cmd/tui.


### Events

Events are parsed from [pokemon-showdown](https://github.com/smogon/pokemon-showdown/blob/master/sim/SIM-PROTOCOL.md) in to a standard Golang struct including the event
- type
- name
- magnitude (if applicable)
- subject (field slot the event refers to)
- targets (additional field slots this refers to)
- metadata

The notion of slots is taken from showdown and is used to represent a player (p1, p2 etc) & field position (a, b, c). In singles there are two slots (p1a, p2a) and in doubles four (p1a, p1b, p2a, p2b).

