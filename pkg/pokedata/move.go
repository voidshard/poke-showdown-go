package pokedata

// MoveDexItem is data parsed from Showdown data files
// https://play.pokemonshowdown.com/data/moves.json
type MoveDexItem struct {
	Number           int `json:"num"`
	Accuracy         int
	Power            int                    `json:"basePower"`
	Category         string                 `json:"category"`
	Priority         int                    `json:"priority"`
	Flags            map[string]int         `json:"flags"`
	Secondary        Secondary              `json:"secondary"`
	Type             string                 `json:"type"`
	Description      string                 `json:"desc"`
	VolatileStatus   string                 `json:"volatileStatus"`
	ZMove            map[string]interface{} `json:"zMove"`
	Condition        map[string]interface{} `json:"condition"`
	StallingMove     bool                   `json:"stallingMove"`
	Boosts           Stats                  `json:"boosts"`
	UnparsedAccuracy interface{}            `json:"accuracy"`
	MaxPP            int                    `json:"pp"`
	Name             string                 `json:"name"`
	Target           string                 `json:"target"`
}

// secondary move side effect data
type Secondary struct {
	Chance int    `json:"chance"`
	Status string `json:"status"`
}
