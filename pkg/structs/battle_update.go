package structs

// Update is essentially the state of one side of the battle currently, and indicates
// that the engine is waiting for a choice (from one or more players) before
// it can proceed to the next round.
type Update struct {
	// battle update but this player doesn't need to make a choice
	// (ie. a pokemon has been switched / fainted)
	Wait bool `json:"wait"`

	// true if pokemon are forced to switch
	ForceSwitch []bool `json:"forceSwitch"`

	// active pokemon movesets (includes additional pp/disabled data)
	Active []PokemonOptions `json:"active"`

	// the team of pokemon
	Team Team `json:"side"`
}

// PokemonOptions is a list of additional move data in order to facilitate a user
// choice.
type PokemonOptions struct {
	Move []Move `json:"moves"`
}

// Move struct expands on simply the move name to include relevant in battle info.
type Move struct {
	Name     string `json:"move"`
	Id       string `json:"id"`
	PP       int    `json:"pp"`
	MaxPP    int    `json:"maxpp"`
	Target   string `json:"target"`
	Disabled bool   `json:"disabled"`
}

// Team represents an entire pokemon team.
// Nb. order here is important.
type Team struct {
	Name    string    `json:"name"`
	Player  string    `json:"player"`
	Pokemon []Pokemon `json:"pokemon"`
}

// Pokemon is a pokemon taking part in a battle
type Pokemon struct {
	Ident string `json:"ident"`

	Active      bool       `json:"active"`
	Stats       StatValues `json:"stats"`
	Moves       []string   `json:"moves"`
	BaseAbility string     `json:"baseAbility"`
	Ability     string     `json:"ability"`
	Item        string     `json:"item"`
	Pokeball    string     `json:"pokeball"`

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
