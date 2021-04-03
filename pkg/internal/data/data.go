package data

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var (
	// data read from assets.go
	parsedPokedex map[string]*PokeDexItem
	parsedMovedex map[string]*MoveDexItem

	// ErrNotFound implies we didn't find a match for some ID
	ErrNotFound = fmt.Errorf("not found")

	// remove all non alpha numeric
	invalid = regexp.MustCompile("[^a-zA-Z0-9]+")
)

// strip removes non alpha-num chars and switches to lowercase
// 'Charizard-Mega-X' -> 'charizardmegax'
// In theory this makes Name fields match ID fields, as understood by
// pokemon-showdown.
func Strip(in string) string {
	return strings.ToLower(invalid.ReplaceAllString(in, ""))
}

// PokeDex returns data by a pokemon's string ID (it's name lowercase
// and stripped of symbols).
func PokeDex(in string) (*PokeDexItem, error) {
	id := Strip(in)

	if parsedPokedex == nil {
		raw, err := Asset("pokedex.json")
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, &parsedPokedex)
		if err != nil {
			return nil, err
		}
	}

	result, ok := parsedPokedex[id]
	if !ok {
		return nil, fmt.Errorf("%w pokemon '%s'", ErrNotFound, id)
	}

	return result, nil
}

// Move returns move data from the move dex given it's id
// (name lowercase, symbols removed)
func MoveDex(in string) (*MoveDexItem, error) {
	id := Strip(in)

	if parsedMovedex == nil {
		raw, err := Asset("moves.json")
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, &parsedMovedex)
		if err != nil {
			return nil, err
		}
	}

	result, ok := parsedMovedex[id]
	if !ok {
		return nil, fmt.Errorf("%w move '%s'", ErrNotFound, id)
	}

	result.Accuracy = parseAccuracy(result.UnparsedAccuracy)

	return result, nil
}

// parseAccuracy because accuracy can be either a bool or an int .. irritating.
func parseAccuracy(in interface{}) int {
	switch in.(type) {
	case bool:
		return 1000
	}

	data, err := json.Marshal(in)
	if err != nil {
		return 0
	}

	var value int
	err = json.Unmarshal(data, &value)
	if err != nil {
		return 0
	}

	return value
}
