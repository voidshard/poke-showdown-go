package structs

// SimData is the input expected for a valid simulate call
type SimData struct {
	Spec    *BattleSpec `json:"spec"`
	Actions []*Action   `json:"actions"`
}
