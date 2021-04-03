package event

import (
	"fmt"
	"regexp"
	"strings"
)

// Subject (pokemon) of an event (source, target etc)
type Subject struct {
	Player   string
	Position string
}

// String returns the [player][position] string used by pokemon showdown
// to reference a slot eg. p2a (first slot of the second player).
func (s *Subject) String() string {
	return fmt.Sprintf("%s%s", s.Player, s.Position)
}

// ParseSubjects returns all subjects in a given line
func ParseSubjects(line string) []*Subject {
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
			},
		)
	}

	return subs
}
