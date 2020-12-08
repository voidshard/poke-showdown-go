# poke-showdown-go
Golang wrapper for [pokemon-showdown battle simulator](https://github.com/smogon/pokemon-showdown/blob/master/sim/README.md).


### What

This library provides a Golang interface for interacting with the pokemon-showdown battle simulator. In short:
- A battle can be started with a given pokemon team and/or a random team in doubles or singles formats (Gen8)
- Decision(s) of players (human or otherwise) can be written to the simulator interface
- Once all players have decided what to do a new battle state is returned


### Warnings

- This is still under development and the API / struct names may change in the future (suggestions & PRs welcome).
- There's very little validation logic in here at the moment, so be careful setting legal values (ie. IVs are between 0-31, EVs 0-255, level 1-100). What happens when you set illegal values is currently undefined.
- This isn't intended to hold data about all pokemon species, items, abilities etc. For that you can use pokemon-showdown to [generate teams](https://github.com/smogon/pokemon-showdown/blob/master/COMMANDLINE.md) or check out [one](https://github.com/veekun/pokedex) of [these](https://github.com/PokeAPI).


### Installing

In addition to the library itself you'll need to install [the engine](https://github.com/smogon/pokemon-showdown/blob/master/sim/README.md). Anecdotally this seems to require a fairly recent version of [NodeJS](https://nodejs.org/en/download/) to be available (this lib was built using v14.15.0). 


### Example

Here we set up a BattleSpec that outlines
- the format of the battle (Gen8 singles)
- the player(s) (p1 & p2)
- a team of one pokemon for each player (Ninetales vs Umbreon)

Note that most fields that are not set on a PokemonSpec will be set to standard values by the simulator itself. Have a look at the [spec](https://github.com/voidshard/poke-showdown-go/blob/main/pkg/structs/pokemon_spec.go) to see what values you can configure.

```golang
import (
	"github.com/voidshard/poke-showdown-go/pkg/sim"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
)

var spec = &structs.BattleSpec{
	Format: structs.FormatGen8,
	Players: map[string][]*structs.PokemonSpec{
		"p1": []*structs.PokemonSpec{
			&structs.PokemonSpec{
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
		"p2": []*structs.PokemonSpec{
			&structs.PokemonSpec{
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

Once we have a spec we can start a new battle with 
```golang
battle, _ := sim.Start("pokemon-showdown", spec)
defer battle.Stop()
```
The string "pokemon-showdown" here is the path to the binary of the same name.

At this point we can begin reading the [current state](https://github.com/voidshard/poke-showdown-go/blob/main/pkg/structs/battle_state.go) of the battle from the simulator. This will tell us 
- the winning player (if any)
- the state of each players field & team
- bonus messages
```golang
for state := battle.Read() {
    // do something
}
```

Each player gets an update including
- the player id
- whether the player needs to make a decision (default) or 'wait'
- which pokemon positions need to switch out
- the currently active pokemon
- the players' team of pokemon

Each [pokemon in the players team](https://github.com/voidshard/poke-showdown-go/blob/main/pkg/structs/battle_update.go#L119) includes all of the general information you might expect (status, HP, stats, ability, items, moves etc). Pokemon that are 'active' include extra information detailing extra move information & options available to them (Dynamax, Megaevolution, Zmoves etc).

For each player that needs to make a choice an [action](https://github.com/voidshard/poke-showdown-go/blob/main/pkg/structs/action.go) is expected that includes information on what to do for each active pokemon.
```golang
act := &structs.Action{
    Player: "p1",
    Specs: []*structs.ActionSpec{
        &structs.ActionSpec{Type: structs.ActionMove, ID: 0},
    },
}

_ = battle.Write(act)
```
This tells the simulator that player p1 wishes their pokemon to use move 0 (that is, the 1st available move). In the same way we could indicate that we wish to switch out to pokemon in position #5 in the team (ie. the 6th team member).
```golang
        &structs.ActionSpec{Type: structs.ActionSwitch, ID: 5},
```
In doubles battles we need to supply two of these specs, in order, that indicate what each active pokemon should do. Additionally in doubles battles we may be required to give a move [target](https://github.com/voidshard/poke-showdown-go/blob/main/pkg/structs/action.go#L47), perhaps better explained [here](https://github.com/smogon/pokemon-showdown/blob/master/sim/SIM-PROTOCOL.md).

Once every player that needs to has called [write](https://github.com/voidshard/poke-showdown-go/blob/main/pkg/sim/interface.go#L21) the engine will output the next turns data, publishing the current state of the field. The cycle continues until a victor is announced.


### TODO

- parse battle update messages ("it is sunny", "pokemon used move") into event structs
- fix tests after some internal reworking
