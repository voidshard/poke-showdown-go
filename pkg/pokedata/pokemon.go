package pokedata

// PokeDexItem is data parsed from Showdown data files
// https://play.pokemonshowdown.com/data/pokedex.json
type PokeDexItem struct {
	Number       int                `json:"num"`
	Name         string             `json:"name"`
	Types        []string           `json:"types"`
	GenderRatio  map[string]float64 `json:"genderRatio"`
	Abilities    map[string]string  `json:"abilities"`
	HeightMeters float64            `json:"heightm"`
	WeightKG     float64            `json:"weightkg"`
	Color        string             `json:"color"`
	Evolutions   []string           `json:"evos"`
	EggGroups    []string           `json:"eggGroups"`
	Tier         string             `json:"tier"`

	Stats Stats `json:"baseStats"`

	PreEvolution       string `json:"prevo"`
	EvolutionLevel     int    `json:"evoLevel"`
	EvolutionType      string `json:"evoType"`
	EvolutionItem      string `json:"evoItem"`
	EvolutionCondition string `json:"evoCondition"`
	IsNonstandard      string `json:"isNonstandard"`

	// used for mega evos
	OtherFormes  []string `json:"otherFormes"`
	FormeOrder   []string `json:"formeOrder"`
	RequiredItem string   `json:"requiredItem"`
	BaseSpecies  string   `json:"baseSpecies"`

	// used for arceus .. as far as I can tell
	Forme         string   `json:"forme"`
	Gender        string   `json:"gender"`
	RequiredItems []string `json:"requiredItems"`
}
