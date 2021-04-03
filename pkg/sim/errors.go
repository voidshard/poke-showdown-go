package sim

import (
	"errors"

	"github.com/voidshard/poke-showdown-go/pkg/internal/parse"
)

// IsInvalidChoice returns if the given error is a ErrInvalidChoice
func IsInvalidChoice(err error) bool {
	return errors.Is(err, parse.ErrInvalidChoice)
}

// IsUnavailableChoice returns if the given error is ErrUnavailableChoice
func IsUnavailableChoice(err error) bool {
	return errors.Is(err, parse.ErrUnavailableChoice)
}
