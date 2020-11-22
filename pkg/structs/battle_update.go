package structs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Update represents the current status & options of one side in a battle
type Update struct {
	// Player id
	Player string `json:"player"`

	// True if this player's side doesn't need to make a decision for the battle
	// to progress (ie. a pokemon has been switched / fainted)
	Wait bool `json:"wait"`

	// true if pokemon are forced to switch
	ForceSwitch []bool `json:"forceSwitch"`

	// Active is the list of pokemon indexes that are currently active.
	// Ie. 0, 2 means the first and third pokemon in the team are active.
	Active []int

	// Pokemon is the list of pokemon on this player's team.
	// Order here is important (for "switch" instructions)
	Pokemon []*Pokemon `json:"pokemon"`
}

// DecodeUpdate turns a raw simulater update into our finished side update
func DecodeUpdate(data []byte) (*Update, error) {
	raw := &rawUpdate{}
	err := json.Unmarshal(data, raw)
	if err != nil {
		return nil, err
	}

	update := &Update{
		Wait:        raw.Wait,
		ForceSwitch: raw.ForceSwitch,
		Pokemon:     raw.Team.Pokemon,
		Player:      raw.Team.Player,
		Active:      []int{},
	}
	for index, pkm := range update.Pokemon {
		if !pkm.Active {
			continue
		}

		if raw.Active == nil || !pkm.Active {
			continue
		}

		key := pkm.moveHash()
		var rawOpts *activeData
		for _, i := range raw.Active { // is a list 1-2 pokemon long
			if key == i.moveHash() {
				update.Active = append(update.Active, index)
				rawOpts = i
				break
			}
		}
		if rawOpts == nil {
			continue
		}

		pkm.Options = &Options{
			CanMegaEvolve: rawOpts.CanMegaEvolve,
			CanDynamax:    rawOpts.CanDynamax,
			CanZMove:      len(rawOpts.ZMove) > 0,
			Moves:         rawOpts.Move,
			ZMoves:        rawOpts.ZMove,
			DMoves:        rawOpts.Dynamax.Moves,
		}
	}

	return update, nil
}

// Move struct expands on simply the move name to include relevant in battle info.
type Move struct {
	// Id is the name string, but caseless & run together
	Id string `json:"id"`

	// Human readable name
	Name string `json:"move"`

	// Target details what pokemon can be targeted in doubles+ battles
	Target string `json:"target"`

	// General house keeping information
	PP       int  `json:"pp"`
	MaxPP    int  `json:"maxpp"`
	Disabled bool `json:"disabled"`
}

// Options (only set on active pokemon) is information detailing what a pokemon
// can do next turn.
type Options struct {
	// Booleans detailing if these options are available
	CanMegaEvolve bool
	CanDynamax    bool
	CanZMove      bool

	// Moves including PP, MaxPP, Target data
	Moves []*Move

	// ZMoves available only if z-crystal is held.
	// Nb. some ZMoves may be nil (implies that the given move is not z-move valid).
	// Ie. [Move1, nil, Move3, nil]
	ZMoves []*Move

	// Dynamax moves, if available.
	DMoves []*Move
}

// Pokemon is a pokemon taking part in a battle
type Pokemon struct {
	// Unique `player: name` string for this pokemon
	Ident string `json:"ident"`

	// True if the pokemon is active in the battle currently.
	Active bool `json:"active"`

	// Options are things an active pokemon can do.
	// Only set if Active is True.
	Options *Options `json:"options"`

	// Stats represents pokemon core stats
	Stats *StatValues `json:"stats"`

	// Moves the pokemon can do (names only)
	Moves []string `json:"moves"`

	// BaseAbility is the pokemon ability naturally
	BaseAbility string `json:"baseAbility"`

	// Ability is the pokemon's current ability (usually same as BaseAbility
	// outside of special circumstances)
	Ability string `json:"ability"`

	// Item is what the pokemon is holding in battle
	Item string `json:"item"`

	// Pokeball the pokemon has been caught in
	Pokeball string `json:"pokeball"`

	// From observation condition is a string like
	// 30/130
	// 0 fnt
	// 130/240 slp
	// ie. HP/MaxHap status1 status2
	Condition string `json:"condition"`

	// Name, Level, Gender run together with ','
	// ie. Chesnaught, L82, M
	Details string `json:"details"`
}

func (p *Pokemon) Level() (int, error) {
	bits := strings.Split(p.Details, ", ")
	if len(bits) != 3 {
		return -1, fmt.Errorf("unable to parse level: %s", p.Details)
	}
	lvl, err := strconv.ParseInt(bits[1][1:], 10, 64)
	return int(lvl), err
}

func (p *Pokemon) Species() string {
	bits := strings.Split(p.Details, ", ")
	return bits[0]
}

// IsAsleep returns if the pokemon is asleep
func (p *Pokemon) IsAsleep() bool {
	return strings.Contains(p.Condition, " slp")
}

// IsBurned returns if the pokemon is burned
func (p *Pokemon) IsBurned() bool {
	return strings.Contains(p.Condition, " brn")
}

// IsParalyzed returns if the pokemon is burned
func (p *Pokemon) IsParalyzed() bool {
	return strings.Contains(p.Condition, " par")
}

// IsFainted returns if the pokemon has fainted
func (p *Pokemon) IsFainted() bool {
	return strings.Contains(p.Condition, " fnt")
}

// IsPoisoned returns if the pokemon is poisoned
// (either standard poison or toxic)
func (p *Pokemon) IsPoisoned() bool {
	return strings.Contains(p.Condition, " psn") || p.IsToxiced()
}

// IsToxiced returns if the pokemon has been "badly" poisoned
func (p *Pokemon) IsToxiced() bool {
	return strings.Contains(p.Condition, " tox")
}

// HP returns the pokemon
// - current hp
// - max hp (if known)
// Or returns an error.
// Note that if the pokemon has fainted we no longer know the max HP :(
func (p *Pokemon) HP() (int, int, error) {
	bits := strings.SplitN(p.Condition, " ", 2)

	if bits[0] == "0" {
		// fainted
		return 0, -1, nil
	}

	hpstats := strings.Split(bits[0], "/")
	if len(hpstats) != 2 {
		return -1, -1, fmt.Errorf("unable to read HP stats: %s [from %s]", bits[0], p.Condition)
	}

	cur, err := strconv.ParseInt(hpstats[0], 10, 64)
	if err != nil {
		return -1, -1, err
	}
	max, err := strconv.ParseInt(hpstats[1], 10, 64)
	return int(cur), int(max), err
}

// ----------------------
// func & structs below here are used to parse & understand simulator data only

// moveHash generates a string from the pokemon's move IDs (we use this to
// compare vs Active.Moves)
func (p *Pokemon) moveHash() string {
	return strings.Join(p.Moves, "#")
}

//rawUpdate is the native format of showdown `sideupdate`
type rawUpdate struct {
	// battle update but this player doesn't need to make a choice
	// (ie. a pokemon has been switched / fainted)
	Wait bool `json:"wait"`

	// true if pokemon are forced to switch
	ForceSwitch []bool `json:"forceSwitch"`

	// active pokemon movesets (includes additional pp/disabled data)
	Active []*activeData `json:"active"`

	// the team of pokemon
	Team rawTeam `json:"side"`
}

// activeData is the native format of showdown `active` (Options)
type activeData struct {
	Ident string `json:"Ident"`

	Move  []*Move `json:"moves"`
	ZMove []*Move `json:"canZMove"`

	CanDynamax    bool `json:"canDynamax"`
	CanMegaEvolve bool `json:"canMegaEvo"`

	Dynamax struct {
		Moves []*Move `json:"maxMoves"`
	} `json:"maxMoves"`
}

// moveHash creates a unique string from the pokemons moves
func (a *activeData) moveHash() string {
	moveNames := []string{}
	for _, m := range a.Move {
		moveNames = append(moveNames, m.Id)
	}
	return strings.Join(moveNames, "#")
}

// rawTeam represents an entire pokemon team.
// Nb. order here is important.
type rawTeam struct {
	Player  string     `json:"name"`
	Pokemon []*Pokemon `json:"pokemon"`
}
