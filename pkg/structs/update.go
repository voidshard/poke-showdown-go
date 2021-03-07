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
		if raw.Active == nil || !pkm.Active {
			continue
		}

		pkm.Data = &DerivedData{
			// helpfully pre-parse fields for user(s)
			Species:     pkm.species(),
			Gender:      pkm.gender(),
			IsAsleep:    pkm.isAsleep(),
			IsBurned:    pkm.isBurned(),
			IsPoisoned:  pkm.isBurned(),
			IsToxiced:   pkm.isToxiced(),
			IsFrozen:    pkm.isFrozen(),
			IsParalyzed: pkm.isParalyzed(),
			IsFainted:   pkm.isFainted(),
		}

		curHP, maxHP, err := pkm.parseHP()
		if err != nil {
			return nil, err
		}
		pkm.Data.HPNow = curHP
		pkm.Data.HPMax = maxHP

		lvl, err := pkm.parseLevel()
		if err != nil {
			return nil, err
		}
		pkm.Data.Level = lvl

		pkm.Data.Shiny = strings.Contains(pkm.Details, ", shiny")

		key := pkm.moveHash()
		var rawOpts *activeData
		for slot, i := range raw.Active { // is a list 1-3 pokemon long
			// try to match active pokemon (about which we are given
			// no identifiers other than their moves) to pokemon in the team
			if key == i.moveHash() {
				update.Active = append(update.Active, index)
				rawOpts = i
				pkm.Slot = slot
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

const (
	TargetAdjAlly       = "adjacentAlly" // helping hand
	TargetAdjAllyOrSelf = "adjacentAllyOrSelf"
	TargetAdjFoe        = "adjacentFoe"
	TargetAll           = "all" // eg. explosion
	TargetAdj           = "allAdjacent"
	TargetAdjFoes       = "allAdjacentFoes"
	TargetAllies        = "allies"   // life dew
	TargetAllySide      = "allySide" // tailwind
	TargetAllyTeam      = "allyTeam" // heal bell
	TargetAny           = "any"
	TargetFoeSide       = "foeSide"      // electroweb
	TargetNormal        = "normal"       // waterfall
	TargetRandomNormal  = "randomNormal" // metronome
	TargetScripted      = "scripted"     // counter
	TargetSelf          = "self"         // recover
)

// Move struct expands on simply the move name to include relevant in battle info.
type Move struct {
	// Id is the name string, but no symbols & lowercase
	Id string `json:"id"`

	// Human readable name
	Name string `json:"move"`

	// https://github.com/smogon/pokemon-showdown/blob/master/sim/dex-moves.ts
	// Target details what pokemon can be targeted in doubles+ battles
	/**
	 * adjacentAlly - Only relevant to Doubles or Triples, the move only targets an ally of the user.
	 * adjacentAllyOrSelf - The move can target the user or its ally.
	 * adjacentFoe - The move can target a foe, but not (in Triples) a distant foe.
	 * all - The move targets the field or all Pokémon at once.
	 * allAdjacent - The move is a spread move that also hits the user's ally.
	 * allAdjacentFoes - The move is a spread move.
	 * allies - The move affects all active Pokémon on the user's team.
	 * allySide - The move adds a side condition on the user's side.
	 * allyTeam - The move affects all unfainted Pokémon on the user's team.
	 * any - The move can hit any other active Pokémon, not just those adjacent.
	 * foeSide - The move adds a side condition on the foe's side.
	 * normal - The move can hit one adjacent Pokémon of your choice.
	 * randomNormal - The move targets an adjacent foe at random.
	 * scripted - The move targets the foe that damaged the user.
	 * self - The move affects the user of the move.
	 */
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

	// Represents the pokemons position on the field.
	// This is only relevant if the pokemon is `Active` otherwise
	// it (naturally) doesn't have a place on the field.
	Slot int `json:"slot"`

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

	// Data is information we parse out & attach for ease of use
	Data *DerivedData `json:"data"`

	// Name, Level, Gender run together with ','
	// Nb. Level is not given if it is 100.
	// ie. Chesnaught, L82, M
	Details string `json:"details"`

	// From observation condition is a string like
	// 30/130
	// 0 fnt
	// 130/240 slp
	// ie. HP/MaxHap status1 status2
	Condition string `json:"condition"`
}

type DerivedData struct {
	// Species is the pokemons major species
	Species string

	// Level is the pokemons current level
	Level int

	// Shiny indicates bling
	Shiny bool

	// Gender indicates male / femaleness.
	// Nb. not all pokemon have a gender
	Gender string

	// HPNow is the pokemons current HP
	HPNow int

	// HPMax is the pokemons full HP (if known)
	HPMax int

	// various status effects
	IsAsleep    bool
	IsBurned    bool
	IsPoisoned  bool
	IsToxiced   bool
	IsFrozen    bool
	IsParalyzed bool
	IsFainted   bool
}

func (p *Pokemon) gender() string {
	for _, chunk := range strings.Split(p.Details, ", ")[1:] {
		switch chunk {
		case "F", "M", "N":
			return chunk
		default:
			continue
		}
	}
	return ""
}

func (p *Pokemon) parseLevel() (int, error) {
	bits := strings.Split(p.Details, ", ")
	if len(bits) <= 1 {
		return -1, fmt.Errorf("unable to parse level: %s", p.Details)
	}

	for _, chunk := range bits[1:] {
		if strings.HasPrefix(chunk, "L") {
			lvl, err := strconv.ParseInt(chunk[1:], 10, 64)
			return int(lvl), err
		}
	}

	return 100, nil
}

func (p *Pokemon) species() string {
	bits := strings.Split(p.Details, ", ")
	return bits[0]
}

// IsAsleep returns if the pokemon is asleep
func (p *Pokemon) isAsleep() bool {
	return strings.Contains(p.Condition, " slp")
}

// IsFrozen returns if the pokemon is frozen
func (p *Pokemon) isFrozen() bool {
	return strings.Contains(p.Condition, " frz")
}

// IsBurned returns if the pokemon is burned
func (p *Pokemon) isBurned() bool {
	return strings.Contains(p.Condition, " brn")
}

// IsParalyzed returns if the pokemon is burned
func (p *Pokemon) isParalyzed() bool {
	return strings.Contains(p.Condition, " par")
}

// IsFainted returns if the pokemon has fainted
func (p *Pokemon) isFainted() bool {
	return strings.Contains(p.Condition, " fnt")
}

// IsPoisoned returns if the pokemon is poisoned
// (either standard poison or toxic)
func (p *Pokemon) isPoisoned() bool {
	return strings.Contains(p.Condition, " psn") || p.isToxiced()
}

// IsToxiced returns if the pokemon has been "badly" poisoned
func (p *Pokemon) isToxiced() bool {
	return strings.Contains(p.Condition, " tox")
}

// HP returns the pokemon
// - current hp
// - max hp (if known)
// Or returns an error.
// Note that if the pokemon has fainted we no longer know the max HP :(
func (p *Pokemon) parseHP() (int, int, error) {
	now, max, _, err := parseCondition(p.Condition)
	return now, max, err
}

func parseCondition(condition string) (int, int, string, error) {
	if strings.Contains(condition, "fnt") {
		return 0, -1, "fnt", nil
	}

	bits := strings.SplitN(condition, " ", 2)

	if bits[0] == "0" {
		// fainted
		return 0, -1, "fnt", nil
	}

	hpstats := strings.Split(bits[0], "/")
	if len(hpstats) != 2 {
		return -1, -1, "", fmt.Errorf("unable to read HP stats: %s [from %s]", bits[0], condition)
	}

	cur, err := strconv.ParseInt(hpstats[0], 10, 64)
	if err != nil {
		return -1, -1, "", err
	}
	max, err := strconv.ParseInt(hpstats[1], 10, 64)

	st := ""
	if len(bits) > 1 {
		st = bits[1]
	}

	return int(cur), int(max), st, err
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
