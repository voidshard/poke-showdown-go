package sim

// SimulatorStream wraps pokemon-showdown simulate battle
type SimulatorStream interface {
	// Write a player decision into the simulator
	Write(*Action) error

	// Updates is all events, errors and requests (for input) from
	// the simulator process
	Updates() <-chan *Update

	// Stop closes everything
	Stop()
}
