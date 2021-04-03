package event

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Win                = "win"
	Turn               = "turn"
	Move               = "move"
	Switch             = "switch"
	Drag               = "drag" // a forced "switch" essentially
	DetailsChange      = "detailschange"
	FormeChange        = "-formechange"
	Replace            = "replace"
	Swap               = "swap"
	Cant               = "cant"
	Faint              = "faint"
	Fail               = "-fail"
	Block              = "-block"
	NoTarget           = "-notarget"
	Miss               = "-miss"
	Damage             = "-damage"
	Heal               = "-heal"
	SetHP              = "-sethp"
	Status             = "-status"
	CureStatus         = "-curestatus"
	CureTeam           = "-cureteam"
	Boost              = "-boost"
	Unboost            = "-unboost"
	SetBoost           = "-setboost"
	SwapBoost          = "-swapboost"
	InvertBoost        = "-invertboost"
	ClearBoost         = "-clearboost"
	ClearAllBoost      = "-clearallboost"
	ClearPositiveBoost = "-clearpositiveboost"
	ClearNegativeBoost = "-clearnegativeboost"
	CopyBoost          = "-copyboost"
	Weather            = "-weather"
	FieldStart         = "-fieldstart"
	FieldEnd           = "-fieldend"
	SideStart          = "-sidestart"
	SideEnd            = "-sideend"
	Start              = "-start"
	End                = "-end"
	Crit               = "-crit"
	SuperEffective     = "-supereffective"
	Resisted           = "-resisted"
	Immune             = "-immune"
	Item               = "-item"
	EndItem            = "-enditem"
	Ability            = "-ability"
	EndAbility         = "-endability"
	Transform          = "-transform"
	Mega               = "-mega"
	Primal             = "-primal"
	Burst              = "-burst"
	ZPower             = "-zpower"
	ZBroken            = "-zbroken"
	Activate           = "-activate"
	Hint               = "-hint"
	Center             = "-center"
	Message            = "-message"
	Combine            = "-combine"
	Waiting            = "-waiting"
	Prepare            = "-prepare"
	MustRecharge       = "-mustrecharge"
	HitCount           = "-hitcount"
	SingleMove         = "-singlemove"
	SingleTurn         = "-singleturn"
)

// all known events
var events = map[string]bool{
	Win:                true,
	Turn:               true,
	Move:               true,
	Switch:             true,
	Drag:               true,
	DetailsChange:      true,
	FormeChange:        true,
	Replace:            true,
	Swap:               true,
	Cant:               true,
	Faint:              true,
	Fail:               true,
	Block:              true,
	NoTarget:           true,
	Miss:               true,
	Damage:             true,
	Heal:               true,
	SetHP:              true,
	Status:             true,
	CureStatus:         true,
	CureTeam:           true,
	Boost:              true,
	Unboost:            true,
	SetBoost:           true,
	SwapBoost:          true,
	InvertBoost:        true,
	ClearBoost:         true,
	ClearAllBoost:      true,
	ClearPositiveBoost: true,
	ClearNegativeBoost: true,
	CopyBoost:          true,
	Weather:            true,
	FieldStart:         true,
	FieldEnd:           true,
	SideStart:          true,
	SideEnd:            true,
	Start:              true,
	End:                true,
	Crit:               true,
	SuperEffective:     true,
	Resisted:           true,
	Immune:             true,
	Item:               true,
	EndItem:            true,
	Ability:            true,
	EndAbility:         true,
	Transform:          true,
	Mega:               true,
	Primal:             true,
	Burst:              true,
	ZPower:             true,
	ZBroken:            true,
	Activate:           true,
	Hint:               true,
	Center:             true,
	Message:            true,
	Combine:            true,
	Waiting:            true,
	Prepare:            true,
	MustRecharge:       true,
	HitCount:           true,
	SingleMove:         true,
	SingleTurn:         true,
}

// Event represents some event that has happened during a turn
type Event struct {
	Type      string
	Name      string // name of ability, pokemon, item, move, stat (if known)
	Magnitude int    // magnitude of event, if applicable

	Subject *Subject
	Targets []*Subject

	Metadata map[string]string
}

// String returns a string representation of this event for debugging
func (e *Event) String() string {
	targets := []string{}
	for _, t := range e.Targets {
		targets = append(targets, t.String())
	}
	sub := ""
	if e.Subject != nil {
		sub = e.Subject.String()
	}
	return fmt.Sprintf(
		"%s %s %d %s %v %v",
		e.Type,
		e.Name,
		e.Magnitude,
		sub,
		targets,
		e.Metadata,
	)
}

// extractMetadata returns key value tags from an event line
func extractMetadata(line string) map[string]string {
	// showdown adds metadata with "[key] some value"
	// .. the value is optional.

	in := line
	found := map[string]string{}
	for {
		idx := strings.Index(in, "[")
		if idx == -1 {
			return found
		}
		in = in[idx:]

		idx = strings.Index(in, "]")
		name := in[1:idx]

		in = in[idx+1:]
		idx = strings.Index(in, "|")
		if idx == -1 {
			idx = strings.Index(in, "[")
		}
		if idx != -1 {
			found[name] = strings.TrimSpace(in[0:idx])
		} else {
			found[name] = strings.TrimSpace(in)
		}
	}

	return found
}

// Parse returns an event (if possible) from a string
// github.com/smogon/pokemon-showdown/blob/master/sim/SIM-PROTOCOL.md
func Parse(line string) *Event {
	bits := strings.Split(line, "|")
	if len(bits) == 1 {
		// not an event
		return nil
	}

	_, ok := events[bits[1]]
	if !ok {
		// not an event we care about
		return nil
	}

	e := &Event{
		Metadata: extractMetadata(line),
		Type:     bits[1],
	}

	subjects := ParseSubjects(line)
	if len(subjects) > 0 {
		// the first is always the source or principal subject
		// ie. the pokemon getting healed / damaged, the source
		// of a move / ability, or generally "who we're talking about"
		e.Subject = subjects[0]
	}
	if len(subjects) > 1 {
		// other subject(s) are (generally) targets of the event
		// or are otherwise related. Context is important.
		// nb. this can include the principal subject too -
		// explosion for example is used by one pokemon but
		// affects everyone (even the user)
		e.Targets = subjects[1:]
	}

	// because we always parse out subjects "POKEMON" & [tag] information above we
	// only need to handle event specific information (setting "Name" and
	// additional information / magnitude)
	switch e.Type {
	case Switch, Drag:
		//|switch|POKEMON|DETAILS|HP STATUS or |drag|POKEMON|DETAILS|HP STATUS
		//|switch|p1a: Ninetales|Ninetales, L5, M|24/24
		e.Name = bits[3]
		now, max, status, err := parseCondition(bits[4])
		if err == nil && max > 0 {
			// a percentage represented as an int 0-100
			e.Magnitude = int((float64(now) / float64(max)) * 100)
		}
		e.Metadata["status"] = status
	case Move, SingleTurn, SingleMove, Status, CureStatus, Start, End, Item, EndItem, Ability, Transform, Prepare, FormeChange, DetailsChange, Replace, EndAbility:
		//|-singleturn|p2a: Umbreon|Protect
		//|-curestatus|POKEMON|STATUS
		//|move|p1a: Ninetales|Inferno|p2a: Umbreon|[miss]
		//|-start|p2a: Cinccino|Disable|Rock Blast|[from] ability: Cursed Body|[of] p1b: Froslass
		e.Name = bits[3]
	case Win, Activate, Weather, FieldStart, FieldEnd:
		//|win|USER
		//|-fieldstart|CONDITION
		//|-weather|WEATHER
		e.Name = bits[2]
	case Turn:
		// |turn|NUMBER
		e.Magnitude = parseint(bits[2])
	case Damage, Heal, SetHP:
		//|-heal|p2a: Umbreon|100/100 brn|[from] item: Leftovers
		//|-damage|p2a: Umbreon|327/348 brn|[from] brn
		now, _, _, _ := parseCondition(bits[3])
		e.Magnitude = now
		e.Name = bits[3]
	case Swap:
		//|swap|POKEMON|POSITION
		e.Name = bits[2]
		e.Magnitude = parseint(bits[3])
	case Cant:
		//|cant|POKEMON|REASON or |cant|POKEMON|REASON|MOVE
		if len(bits) >= 5 {
			e.Name = bits[4]
		}
		e.Metadata["reason"] = bits[3]
	case Boost, Unboost, SetBoost:
		//|-unboost|p1a: Lugia|spd|1
		//|-boost|p1a: Lugia|atk|2
		e.Name = bits[3]
		e.Magnitude = parseint(bits[4])
	case Fail:
		//|-fail|p2a: Umbreon # <-- we don't need to worry about this one here
		//|-fail|p2a: Umbreon|ACTION
		if len(bits) >= 4 {
			e.Name = bits[3]
		}
	case ClearPositiveBoost, SwapBoost, Mega:
		//|-mega|p2a: Gallade|Gallade|Galladite
		//|-swapboost|SOURCE|TARGET|STATS
		//|-clearpositiveboost|TARGET|POKEMON|EFFECT
		e.Name = bits[4]
	case SideStart, SideEnd:
		//|-sidestart|SIDE|CONDITION
		e.Name = bits[3]
		e.Metadata["side"] = bits[2]
	case Burst:
		//|-burst|POKEMON|SPECIES|ITEM
		e.Name = bits[4]
		e.Metadata["species"] = bits[3]
	case HitCount:
		//|-hitcount|POKEMON|NUM
		e.Magnitude = parseint(bits[3])
	}

	return e
}

// parseCondition parses a showdown style pokemon 'condition' string
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

// parseint parses the int.
// If not possible we return the default value (0).
func parseint(s string) int {
	i, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return int(i)
}
