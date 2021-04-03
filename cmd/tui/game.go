package main

import (
	"fmt"
	"strings"

	"github.com/voidshard/poke-showdown-go/pkg/event"
	"github.com/voidshard/poke-showdown-go/pkg/field"
	"github.com/voidshard/poke-showdown-go/pkg/sim"

	"github.com/rivo/tview"
)

type game struct {
	proc  sim.SimulatorStream
	field field.Watcher

	p1a *tview.TextView
	p1b *tview.TextView
	p2a *tview.TextView
	p2b *tview.TextView

	p1Opts *tview.List
	p2Opts *tview.List
	p1Tree *ChoiceTree
	p2Tree *ChoiceTree

	output *tview.TextView
}

func (g *game) run() {
	for u := range g.proc.Updates() {
		g.field.Update(u)

		if u.Event != nil {
			fmt.Fprintf(g.output, "%s\n", u.Event.String())
			if u.Event.Type == event.Win {
				return
			}
		} else if u.Error != nil {
			fmt.Fprintf(g.output, "[err] %s %v\n", u.Error.Error(), sim.IsInvalidChoice(u.Error) || sim.IsUnavailableChoice(u.Error))
			if sim.IsInvalidChoice(u.Error) || sim.IsUnavailableChoice(u.Error) {
				g.p1Tree.Redo()
				g.p2Tree.Redo()
			}
		} else if u.Side != nil {
			fmt.Fprintf(g.output, "[update] %s\n", u.Side.Player)

			if u.Side.Player == "p1" {
				g.updateSide(u.Side, g.p1a, g.p1b)
				g.p1Tree.Update(u.Side)
			} else {
				g.updateSide(u.Side, g.p2a, g.p2b)
				g.p2Tree.Update(u.Side)
			}
		}
	}
}

func (g *game) stop() {
	g.proc.Stop()
}

func clearList(l *tview.List) {
	for {
		if l.GetItemCount() > 0 {
			l.RemoveItem(0)
		} else {
			return
		}
	}
}

func (g *game) foeSide(player string) []*sim.Pokemon {
	switch player {
	case "p1":
		return g.p2Tree.Active()
	case "p2":
		return g.p1Tree.Active()
	}
	return nil
}

func (g *game) updateSide(side *sim.Side, views ...*tview.TextView) {
	for i, slot := range side.Field {
		v := views[i]
		if slot.Ident == "" {
			v.SetText("")
			continue
		}

		pkm := side.Pokemon[i]
		v.SetText(fmt.Sprintf(
			"%s %s\n%s\n%s\nAbility: %s\nItem: %s",
			slot.ID,
			pkm.Details,
			pkm.Condition,
			strings.Join(pkm.Dex.Types, "/"),
			pkm.Ability,
			pkm.Item,
		))
	}
}

func newGame(spec *sim.BattleSpec) (*game, error) {
	proc, err := sim.NewSimulatorStream(spec)
	if err != nil {
		return nil, err
	}

	isd := sim.IsDoubles(spec.Format)
	p1Opts := tview.NewList().SetWrapAround(true)
	p2Opts := tview.NewList().SetWrapAround(true)

	g := &game{
		proc:   proc,
		field:  field.NewWatcher(),
		p1a:    tview.NewTextView().SetWordWrap(true),
		p1b:    tview.NewTextView().SetWordWrap(true),
		p2a:    tview.NewTextView().SetWordWrap(true),
		p2b:    tview.NewTextView().SetWordWrap(true),
		p1Opts: p1Opts,
		p2Opts: p2Opts,
		output: tview.NewTextView().SetWordWrap(true).SetScrollable(true),
	}

	g.p1Tree = NewChoiceTree(g, p1Opts, isd)
	g.p2Tree = NewChoiceTree(g, p2Opts, isd)

	return g, nil
}
