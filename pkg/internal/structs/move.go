package structs

// Move is returned from the pokemon simulator for active pokemon
type Move struct {
	ID string `json:"id"`

	PP       int  `json:"pp"`
	Disabled bool `json:"disabled"`

	Name string `json:"move"`
}
