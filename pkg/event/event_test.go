package event

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var dataTestParseEvent = []struct {
	In     string
	Expect *Event
}{
	{
		"|t:|1467", // this is a valid event, we just don't use it
		nil,
	},
	{
		"|what", // this is a garbage event
		nil,
	},
	{
		"comment for you", // a comment string
		nil,
	},
	{
		">start", // a command
		nil,
	},
	{
		"|move|p1a: Lugia|Explosion|p2b: Umbreon|[spread] p1b,p2a,p2b",
		&Event{
			Type:    Move,
			Name:    "Explosion",
			Subject: &Subject{Player: "p1", Position: "a"},
			Targets: []*Subject{
				&Subject{Player: "p2", Position: "b"},
				&Subject{Player: "p1", Position: "b"},
				&Subject{Player: "p2", Position: "a"},
				&Subject{Player: "p2", Position: "b"},
			},
			Metadata: map[string]string{
				"spread": "p1b,p2a,p2b",
			},
		},
	},
	{
		"|move|p1a: Ninetales|Inferno|p2a: Umbreon|[miss]",
		&Event{
			Type:    Move,
			Name:    "Inferno",
			Subject: &Subject{Player: "p1", Position: "a"},
			Targets: []*Subject{
				&Subject{Player: "p2", Position: "a"},
			},
			Metadata: map[string]string{
				"miss": "",
			},
		},
	},
	{
		"|-weather|SunnyDay|[upkeep]",
		&Event{
			Type:     Weather,
			Name:     "SunnyDay",
			Metadata: map[string]string{"upkeep": ""},
		},
	},
	{
		"|-heal|p2a: Umbreon|4/13 brn|[from] item: Leftovers",
		&Event{
			Type:      Heal,
			Name:      "4/13 brn",
			Subject:   &Subject{Player: "p2", Position: "a"},
			Metadata:  map[string]string{"from": "item: Leftovers"},
			Magnitude: 4,
		},
	},
	{
		"|switch|p2a: Umbreon|Umbreon, L5, M|27/27",
		&Event{
			Type:      Switch,
			Name:      "Umbreon, L5, M",
			Subject:   &Subject{Player: "p2", Position: "a"},
			Metadata:  map[string]string{"status": ""},
			Magnitude: 100,
		},
	},
	{
		"|detailschange|p2a: Gallade|Gallade-Mega, L50, M",
		&Event{
			Type:     DetailsChange,
			Name:     "Gallade-Mega, L50, M",
			Subject:  &Subject{Player: "p2", Position: "a"},
			Metadata: map[string]string{},
		},
	},
	{
		"|-mega|p2a: Gallade|Gallade|Galladite",
		&Event{
			Type:     Mega,
			Name:     "Galladite",
			Subject:  &Subject{Player: "p2", Position: "a"},
			Metadata: map[string]string{},
		},
	},
	{
		"|-boost|p2a: Gallade|atk|2",
		&Event{
			Type:      Boost,
			Name:      "atk",
			Magnitude: 2,
			Subject:   &Subject{Player: "p2", Position: "a"},
			Metadata:  map[string]string{},
		},
	},
	{
		"|-hitcount|p1a: Lugia|3", // implies Lugia was hit
		&Event{
			Type:      HitCount,
			Magnitude: 3,
			Metadata:  map[string]string{},
			Subject:   &Subject{Player: "p1", Position: "a"},
		},
	},
	{
		"|-end|p2a: Zoroark|Illusion",
		&Event{
			Type:     End,
			Name:     "Illusion",
			Subject:  &Subject{Player: "p2", Position: "a"},
			Metadata: map[string]string{},
		},
	},
	{
		"|replace|p2a: Zoroark|Zoroark, L5, M",
		&Event{
			Type:     Replace,
			Name:     "Zoroark, L5, M",
			Subject:  &Subject{Player: "p2", Position: "a"},
			Metadata: map[string]string{},
		},
	},
	{
		"|-endability|p2a: Zoroark|Limber|[from] move: Transform",
		&Event{
			Type:     EndAbility,
			Name:     "Limber",
			Subject:  &Subject{Player: "p2", Position: "a"},
			Metadata: map[string]string{"from": "move: Transform"},
		},
	},
	{
		"|-transform|p2a: Zoroark|p1a: Lugia",
		&Event{
			Type:    Transform,
			Name:    "p1a: Lugia",
			Subject: &Subject{Player: "p2", Position: "a"},
			Targets: []*Subject{
				&Subject{Player: "p1", Position: "a"},
			},
			Metadata: map[string]string{},
		},
	},
	{
		"|-fail|p1a: Lugia|heal",
		&Event{
			Type:     Fail,
			Name:     "heal",
			Subject:  &Subject{Player: "p1", Position: "a"},
			Metadata: map[string]string{},
		},
	},
	{
		"|-ability|p2a: Zoroark|Intimidate|boost",
		&Event{
			Type:     Ability,
			Name:     "Intimidate",
			Subject:  &Subject{Player: "p2", Position: "a"},
			Metadata: map[string]string{},
		},
	},
}

func TestParseEvent(t *testing.T) {
	for _, tt := range dataTestParseEvent {
		t.Run(tt.In, func(t *testing.T) {
			result := Parse(tt.In)

			assert.Equal(t, tt.Expect, result)
		})
	}
}

func TestSubjectString(t *testing.T) {
	result := (&Subject{Player: "p1", Position: "a"}).String()

	assert.Equal(t, "p1a", result)
}

func TestParseInt(t *testing.T) {
	valid := parseint("42")
	assert.Equal(t, 42, valid)

	invalid := parseint("zip")
	assert.Equal(t, 0, invalid)
}

var dataTestExtractSubjects = []struct {
	In     string
	Expect []*Subject
}{
	{
		"p1a",
		[]*Subject{&Subject{Player: "p1", Position: "a"}},
	},
	{
		"|blah|BLAH: what|p1a: foobar|[zap]",
		[]*Subject{&Subject{Player: "p1", Position: "a"}},
	},
	{
		"|move|p1a: Lugia|Explosion|p2b: Umbreon|[spread] p1b,p2a,p2b",
		[]*Subject{
			&Subject{Player: "p1", Position: "a"},
			&Subject{Player: "p2", Position: "b"},
			&Subject{Player: "p1", Position: "b"},
			&Subject{Player: "p2", Position: "a"},
			&Subject{Player: "p2", Position: "b"},
		},
	},
}

func TestExtractSubjects(t *testing.T) {
	for _, tt := range dataTestExtractSubjects {
		t.Run(tt.In, func(t *testing.T) {
			result := ParseSubjects(tt.In)

			assert.Equal(t, tt.Expect, result)
		})
	}
}

var dataTestExtractMetadata = []struct {
	In     string
	Expect map[string]string
}{
	{
		"|move|p1a: Ninetales|Inferno|p2a: Umbreon|[miss]",
		map[string]string{"miss": ""},
	},
	{
		"|-weather|SunnyDay|[upkeep]",
		map[string]string{"upkeep": ""},
	},
	{
		"|-heal|p2a: Umbreon|4/13 brn|[from] item: Leftovers",
		map[string]string{"from": "item: Leftovers"},
	},
	{
		"|-damage|p2a: Umbreon|3/13 brn|[from] brn",
		map[string]string{"from": "brn"},
	},
	{
		"|-weather|SunnyDay|[from] ability: Drought|[of] p1a: Ninetales",
		map[string]string{
			"from": "ability: Drought",
			"of":   "p1a: Ninetales",
		},
	},
	{
		// I don't know if showdown does this, but ..
		"[hi][there]",
		map[string]string{"hi": "", "there": ""},
	},
}

func TestExtractMetadata(t *testing.T) {
	for _, tt := range dataTestExtractMetadata {
		t.Run(tt.In, func(t *testing.T) {
			result := extractMetadata(tt.In)

			assert.Equal(t, tt.Expect, result)
		})
	}
}
