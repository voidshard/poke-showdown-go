package pokedata

// Stats represent the 6 main pokemon stats
type Stats struct {
	HP  int `json:"hp"`
	Atk int `json:"atk"`
	Def int `json:"def"`
	Spa int `json:"spa"`
	Spd int `json:"spd"`
	Spe int `json:"spe"`
}
