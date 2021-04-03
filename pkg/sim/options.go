package sim

import (
	"github.com/voidshard/poke-showdown-go/pkg/internal/structs"
)

// Options includes what a pokemon might want to do next turn
type Options struct {
	// Booleans detailing if these options are available
	CanMegaEvolve bool
	CanDynamax    bool
	CanZMove      bool

	// Moves including PP, MaxPP, Target data
	Moves []*Move

	// ZMoves available only if z-crystal is held.
	// Nb. some ZMoves may be nil (implies that the given move is not z-move valid).
	// Ie. [Move1, nil, Move3, nil]
	ZMoves []*Move

	// Dynamax moves, if available.
	DMoves []*Move
}

// toOptions converts showdowns Active data block to our
// 'options' struct.
func toOptions(in *structs.ActiveData) *Options {
	if in == nil {
		return nil
	}

	opts := &Options{}
	for _, m := range in.Moves {
		mov := toMove(m)
		if mov != nil {
			opts.Moves = append(opts.Moves, mov)
		}
	}
	opts.CanMegaEvolve = in.CanMegaEvolve

	for _, m := range in.ZMoves {
		if m == nil {
			continue
		}

		mov := toMove(m)
		if mov != nil {
			opts.ZMoves = append(opts.ZMoves, mov)
		}
	}
	opts.CanZMove = len(opts.ZMoves) > 0

	for _, m := range in.Dynamax.Moves {
		mov := toMove(m)
		if mov != nil {
			opts.DMoves = append(opts.DMoves, mov)
		}
	}
	opts.CanDynamax = in.CanDynamax

	return opts
}
