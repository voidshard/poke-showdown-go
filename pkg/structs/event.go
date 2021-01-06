package structs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	EventMove               = "move"
	EventSwitch             = "switch"
	EventDrag               = "drag" // a forced "switch" essentially
	EventDetailsChange      = "detailschange"
	EventFormeChange        = "-formechange"
	EventReplace            = "replace"
	EventSwap               = "swap"
	EventCant               = "cant"
	EventFaint              = "faint"
	EventFail               = "-fail"
	EventBlock              = "-block"
	EventNoTarget           = "-notarget"
	EventMiss               = "-miss"
	EventDamage             = "-damage"
	EventHeal               = "-heal"
	EventSetHP              = "-sethp"
	EventStatus             = "-status"
	EventCureStatus         = "-curestatus"
	EventCureTeam           = "-cureteam"
	EventBoost              = "-boost"
	EventUnboost            = "-unboost"
	EventSetBoost           = "-setboost"
	EventSwapBoost          = "-swapboost"
	EventInvertBoost        = "-invertboost"
	EventClearBoost         = "-clearboost"
	EventClearAllBoost      = "-clearallboost"
	EventClearPositiveBoost = "-clearpositiveboost"
	EventClearNegativeBoost = "-clearnegativeboost"
	EventCopyBoost          = "-copyboost"
	EventWeather            = "-weather"
	EventFieldStart         = "-fieldstart"
	EventFieldEnd           = "-fieldend"
	EventSideStart          = "-sidestart"
	EventSideEnd            = "-sideend"
	EventStart              = "-start"
	EventEnd                = "-end"
	EventCrit               = "-crit"
	EventSuperEffective     = "-supereffective"
	EventResisted           = "-resisted"
	EventImmune             = "-immune"
	EventItem               = "-item"
	EventEndItem            = "-enditem"
	EventAbility            = "-ability"
	EventEndAbility         = "-endability"
	EventTransform          = "-transform"
	EventMega               = "-mega"
	EventPrimal             = "-primal"
	EventBurst              = "-burst"
	EventZPower             = "-zpower"
	EventZBroken            = "-zbroken"
	EventActivate           = "-activate"
	EventHint               = "-hint"
	EventCenter             = "-center"
	EventMessage            = "-message"
	EventCombine            = "-combine"
	EventWaiting            = "-waiting"
	EventPrepare            = "-prepare"
	EventMustRecharge       = "-mustrecharge"
	EventHitCount           = "-hitcount"
	EventSingleMove         = "-singlemove"
	EventSingleTurn         = "-singleturn"
)

// all known events
var events = map[string]bool{
	EventMove:               true,
	EventSwitch:             true,
	EventDrag:               true,
	EventDetailsChange:      true,
	EventFormeChange:        true,
	EventReplace:            true,
	EventSwap:               true,
	EventCant:               true,
	EventFaint:              true,
	EventFail:               true,
	EventBlock:              true,
	EventNoTarget:           true,
	EventMiss:               true,
	EventDamage:             true,
	EventHeal:               true,
	EventSetHP:              true,
	EventStatus:             true,
	EventCureStatus:         true,
	EventCureTeam:           true,
	EventBoost:              true,
	EventUnboost:            true,
	EventSetBoost:           true,
	EventSwapBoost:          true,
	EventInvertBoost:        true,
	EventClearBoost:         true,
	EventClearAllBoost:      true,
	EventClearPositiveBoost: true,
	EventClearNegativeBoost: true,
	EventCopyBoost:          true,
	EventWeather:            true,
	EventFieldStart:         true,
	EventFieldEnd:           true,
	EventSideStart:          true,
	EventSideEnd:            true,
	EventStart:              true,
	EventEnd:                true,
	EventCrit:               true,
	EventSuperEffective:     true,
	EventResisted:           true,
	EventImmune:             true,
	EventItem:               true,
	EventEndItem:            true,
	EventAbility:            true,
	EventEndAbility:         true,
	EventTransform:          true,
	EventMega:               true,
	EventPrimal:             true,
	EventBurst:              true,
	EventZPower:             true,
	EventZBroken:            true,
	EventActivate:           true,
	EventHint:               true,
	EventCenter:             true,
	EventMessage:            true,
	EventCombine:            true,
	EventWaiting:            true,
	EventPrepare:            true,
	EventMustRecharge:       true,
	EventHitCount:           true,
	EventSingleMove:         true,
	EventSingleTurn:         true,
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

func extractSubjects(line string) []*Subject {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	fields := strings.Fields(re.ReplaceAllString(line, " "))

	fe := regexp.MustCompile(`p[1-4][a-c]`) // ie: p1a p2c ...
	subs := []*Subject{}
	for _, field := range fields {
		if !fe.MatchString(field) {
			continue
		}

		subs = append(
			subs,
			&Subject{
				Player:   field[0:2],
				Position: field[2:3],
			})
	}

	return subs
}

func ParseEvent(line string) *Event {
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

	subjects := extractSubjects(line)
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
	case EventMove, EventSingleTurn, EventSingleMove, EventStatus, EventCureStatus, EventStart, EventEnd, EventItem, EventEndItem, EventAbility, EventTransform, EventPrepare, EventSwitch, EventDrag, EventFormeChange, EventDetailsChange, EventReplace, EventEndAbility:
		//|-singleturn|p2a: Umbreon|Protect
		//|switch|p1a: Ninetales|Ninetales, L5, M|24/24
		//|-curestatus|POKEMON|STATUS
		//|move|p1a: Ninetales|Inferno|p2a: Umbreon|[miss]
		e.Name = bits[3]
	case EventActivate, EventWeather, EventFieldStart, EventFieldEnd:
		//|-fieldstart|CONDITION
		//|-weather|WEATHER
		e.Name = bits[2]
	case EventDamage, EventHeal, EventSetHP:
		//|-heal|p2a: Umbreon|100/100 brn|[from] item: Leftovers
		//|-damage|p2a: Umbreon|327/348 brn|[from] brn
		now, _, _, _ := parseCondition(bits[3])
		e.Magnitude = now
		e.Name = bits[3]
	case EventSwap:
		//|swap|POKEMON|POSITION
		e.Name = bits[2]
		e.Magnitude = parseint(bits[3])
	case EventCant:
		//|cant|POKEMON|REASON or |cant|POKEMON|REASON|MOVE
		if len(bits) >= 5 {
			e.Name = bits[4]
		}
		e.Metadata["reason"] = bits[3]
	case EventBoost, EventUnboost, EventSetBoost:
		//|-unboost|p1a: Lugia|spd|1
		//|-boost|p1a: Lugia|atk|2
		e.Name = bits[3]
		e.Magnitude = parseint(bits[4])
	case EventFail:
		//|-fail|p2a: Umbreon # <-- we don't need to worry about this one here
		//|-fail|p2a: Umbreon|ACTION
		if len(bits) >= 4 {
			e.Name = bits[3]
		}
	case EventClearPositiveBoost, EventSwapBoost, EventMega:
		//|-mega|p2a: Gallade|Gallade|Galladite
		//|-swapboost|SOURCE|TARGET|STATS
		//|-clearpositiveboost|TARGET|POKEMON|EFFECT
		e.Name = bits[4]
	case EventSideStart, EventSideEnd:
		//|-sidestart|SIDE|CONDITION
		e.Name = bits[3]
		e.Metadata["side"] = bits[2]
	case EventBurst:
		//|-burst|POKEMON|SPECIES|ITEM
		e.Name = bits[4]
		e.Metadata["species"] = bits[3]
	case EventHitCount:
		//|-hitcount|POKEMON|NUM
		e.Magnitude = parseint(bits[3])
	}

	return e
}

// Subject (pokemon) of an event (source, target etc)
type Subject struct {
	Player   string
	Position string
}

func (s *Subject) String() string {
	return fmt.Sprintf("%s%s", s.Player, s.Position)
}

func (s *Subject) ActiveIndex() int {
	switch s.Position {
	case "b":
		return 1
	case "c":
		return 2
	}
	return 0
}

func parseint(s string) int {
	i, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return int(i)
}
