package sim

import (
	"github.com/voidshard/poke-showdown-go/pkg/event"
)

// Update represents some update to the game
type Update struct {
	Number int
	Side   *Side
	Event  *event.Event
	Error  error
}
