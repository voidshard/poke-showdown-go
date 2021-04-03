package sim

import (
	"strings"
)

const (
	// All pokemon natures
	/*
	   Adamant	Attack	Sp. Atk
	   Bashful	Sp. Atk	Sp. Atk
	   Bold	        Defense	Attack
	   Brave	Attack	Speed
	   Calm	        Sp. Def	Attack
	   Careful	Sp. Def	Sp. Atk
	   Docile	Defense	Defense
	   Gentle	Sp. Def	Defense
	   Hardy	Attack	Attack
	   Hasty	Speed	Defense
	   Impish	Defense	Sp. Atk
	   Jolly	Speed	Sp. Atk
	   Lax	        Defense	Sp. Def
	   Lonely	Attack	Defense
	   Mild 	Sp. Atk	Defense
	   Modest	Sp. Atk	Attack
	   Naive	Speed	Sp. Def
	   Naughty	Attack	Sp. Def
	   Quiet	Sp. Atk	Speed
	   Quirky	Sp. Def	Sp. Def
	   Rash	        Sp. Atk	Sp. Def
	   Relaxed	Defense	Speed
	   Sassy	Sp. Def	Speed
	   Serious	Speed	Speed
	   Timid	Speed	Attack
	*/
	NatureAdamant = "adamant"
	NatureBashful = "bashful"
	NatureBold    = "bold"
	NatureBrave   = "brave"
	NatureCalm    = "calm"
	NatureCareful = "careful"
	NatureDocile  = "docile"
	NatureGentle  = "gentle"
	NatureHardy   = "hardy"
	NatureHasty   = "hasty"
	NatureImpish  = "impish"
	NatureJolly   = "jolly"
	NatureLax     = "lax"
	NatureLonely  = "lonely"
	NatureMild    = "mild"
	NatureModest  = "modest"
	NatureNaive   = "naive"
	NatureNaughty = "naughty"
	NautreQuiet   = "quiet"
	NatureQuirky  = "quirky"
	NatureRash    = "rash"
	NatureRelaxed = "relaxed"
	NatureSassy   = "sassy"
	NatureSerious = "serious"
	NatureTimid   = "timid"
)

var (
	natures = []string{
		NatureAdamant,
		NatureBashful,
		NatureBold,
		NatureBrave,
		NatureCalm,
		NatureCareful,
		NatureDocile,
		NatureGentle,
		NatureHardy,
		NatureHasty,
		NatureImpish,
		NatureJolly,
		NatureLax,
		NatureLonely,
		NatureMild,
		NatureModest,
		NatureNaive,
		NatureNaughty,
		NautreQuiet,
		NatureQuirky,
		NatureRash,
		NatureRelaxed,
		NatureSassy,
		NatureSerious,
		NatureTimid,
	}
)

// Natures returns all valid natures
func Natures() []string {
	return natures
}

// ValidNature returns if the given nature is valid
func ValidNature(nature string) bool {
	in := strings.ToLower(nature)
	for _, n := range natures {
		if in == n {
			return true
		}
	}
	return false
}
