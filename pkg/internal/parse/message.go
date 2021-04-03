package parse

import (
	"github.com/voidshard/poke-showdown-go/pkg/event"
	"github.com/voidshard/poke-showdown-go/pkg/internal/structs"
)

var messageNumber int

// Message is some output from the PS process.
// We wrap these together because order is important in the interleaving
// of events, updates and errors.
type Message struct {
	Num    int
	Event  *event.Event
	Update *structs.Update
	Error  error
}

// message makes a new message for the given item
func message(in interface{}) *Message {
	m := &Message{Num: messageNumber}
	switch in.(type) {
	case *structs.Update:
		m.Update = in.(*structs.Update)
	case *event.Event:
		m.Event = in.(*event.Event)
	case error:
		m.Error = in.(error)
	}
	messageNumber++
	return m
}
