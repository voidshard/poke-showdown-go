package main

import (
	"fmt"
	"strings"

	"github.com/voidshard/poke-showdown-go/pkg/sim"

	"github.com/rivo/tview"
)

type ChoiceTree struct {
	game *game

	view *tview.List

	side        *sim.Side
	forceSwitch bool
	act         *sim.Action

	isDoubles bool
}

func NewChoiceTree(g *game, view *tview.List, isDoubles bool) *ChoiceTree {
	return &ChoiceTree{
		game:      g,
		view:      view,
		isDoubles: isDoubles,
	}
}

func (c *ChoiceTree) Redo() {
	c.Update(c.side)
}

func (c *ChoiceTree) Active() []*sim.Pokemon {
	onfield := []*sim.Pokemon{}
	for _, p := range c.side.Pokemon {
		if p.Active {
			onfield = append(onfield, p)
		}
	}
	return onfield
}

func (c *ChoiceTree) Update(side *sim.Side) {
	clearList(c.view)
	c.side = side

	if c.side.Wait {
		// we don't need to do anything
		return
	}

	c.act = &sim.Action{Player: side.Player, Specs: []*sim.ActionSpec{}}

	// determine if this update is asking for us to decide
	// our turn actions (attack / switch) or is a mid-turn
	// force switch (roar, pokemon faint) etc.
	c.forceSwitch = false
	for _, slot := range side.Field {
		if slot.Switch {
			c.forceSwitch = true
			break
		}
	}

	c.createActionSpec(side.Field[0])
}

func (c *ChoiceTree) canSwitchIn() []*sim.Pokemon {
	mates := []*sim.Pokemon{}
	for _, p := range c.side.Pokemon {
		// TODO: edge case around double switching
		if p.Active || p.Status.IsFainted {
			continue
		}
		mates = append(mates, p)
	}
	return mates
}

func (c *ChoiceTree) specAdded() {
	clearList(c.view)
	if len(c.act.Specs) >= len(c.side.Field) {
		c.game.proc.Write(c.act)
	} else if len(c.side.Field) < 2 && c.isDoubles {
		c.act.Specs = append(c.act.Specs, &sim.ActionSpec{Type: sim.ActionPass})
		c.specAdded()
	} else {
		next := c.side.Field[len(c.act.Specs)]
		c.createActionSpec(next)
	}
}

func (c *ChoiceTree) btnSwitch(i int, p *sim.Pokemon) {
	c.view.AddItem(
		p.Details,
		fmt.Sprintf("%s %s", strings.Join(p.Dex.Types, "/"), p.Condition),
		[]rune(fmt.Sprintf("%d", i))[0],
		func() {
			c.act.Specs = append(
				c.act.Specs,
				&sim.ActionSpec{
					Type: sim.ActionSwitch,
					ID:   fmt.Sprintf("%d", p.Index+1),
				},
			)
			c.specAdded()
		},
	)
}

func otherPlayer(p string) string {
	if p == "p1" {
		return "p2"
	}
	return "p1"
}

func toSlot(me string, index int) string {
	player := me
	if index > 0 {
		player = otherPlayer(me)
	}

	switch index {
	case -2, 2:
		return fmt.Sprintf("%sb", player)
	case -3, 3:
		return fmt.Sprintf("%sc", player)
	}
	return fmt.Sprintf("%sa", player)
}

func (c *ChoiceTree) chooseTarget(m *sim.Move, slot *sim.Slot, spec *sim.ActionSpec) {
	clearList(c.view)
	for i, tnum := range sim.ValidDoublesTargets(m.Target) {
		t := tnum

		slot := toSlot(c.side.Player, t)
		pkm := c.game.field.WhoIs(slot)
		if pkm == nil {
			continue
		}

		c.view.AddItem(
			fmt.Sprintf("%s", pkm.Details),
			fmt.Sprintf("[%s] %s", slot, pkm.Condition),
			[]rune(fmt.Sprintf("%d", i))[0],
			func() {
				spec.Target = t
				c.act.Specs = append(c.act.Specs, spec)
				c.specAdded()
			},
		)
	}
	c.view.AddItem("..", "", '.', func() {
		clearList(c.view)
		c.createActionSpec(slot)
	})
}

func (c *ChoiceTree) btnMove(i int, m *sim.Move, slot *sim.Slot, spec *sim.ActionSpec) {
	needsTarget := c.isDoubles && sim.RequiresDoublesTarget(m.Target)

	c.view.AddItem(
		fmt.Sprintf("%s [%d %d %s]", m.Name, m.Power, m.PP, m.Type),
		m.Description,
		[]rune(fmt.Sprintf("%d", i))[0],
		func() {
			if needsTarget {
				c.chooseTarget(m, slot, spec)
				return
			}
			c.act.Specs = append(c.act.Specs, spec)
			c.specAdded()
		},
	)
}

func (c *ChoiceTree) createActionSpec(slot *sim.Slot) {
	/* Essentially we think of each root option as a new menu giving us 5 menus
	- switch
	- fight (use vanilla move)
	- zmove
	- dynamax (use dynamax move, dynamax if not already)
	- megaevolve (megaevolve if not already, use vanilla move)

	There are quite a few edge cases to consider
	- we must pass (in doubles if only the other pokemon slot must switch)
	- we cannot switch (arena trap, outrage)
	- we must switch (faint, roar, uturn)
	- we cannot perform some moves (disabled, out of pp)
	- whether we need to choose a target or not
	*/

	pkm := c.side.Pokemon[slot.Index]
	mates := c.canSwitchIn()
	canSwitch := len(mates) > 0 && !slot.Trapped

	if (c.forceSwitch && !slot.Switch) || (pkm.Status.IsFainted && !canSwitch) {
		// this update is asking us to make a switch for slots on our
		// field, but this slot doesn't need to (so it should pass)
		c.act.Specs = append(c.act.Specs, &sim.ActionSpec{Type: sim.ActionPass})
		c.specAdded()
		return
	}

	if canSwitch {
		c.view.AddItem("switch", "", 's', func() {
			clearList(c.view)
			for i, pkm := range mates {
				p := pkm
				c.btnSwitch(i, p)
			}
			c.view.AddItem("..", "", '.', func() {
				clearList(c.view)
				c.createActionSpec(slot)
			})
		})

	}
	if slot.Switch {
		return
	}

	if slot.Options == nil {
		if c.view.GetItemCount() == 0 { // if you have no options then you pass
			c.act.Specs = append(c.act.Specs, &sim.ActionSpec{Type: sim.ActionPass})
			c.specAdded()
		}
		return
	}

	if slot.Options.CanMegaEvolve {
		c.view.AddItem("mega-evolve", "", 'm', func() {
			clearList(c.view)
			for i, move := range slot.Options.Moves {
				m := move
				if (m.Disabled || m.PP == 0) && !slot.Trapped {
					continue
				}
				c.btnMove(
					i,
					m,
					slot,
					&sim.ActionSpec{
						Type: sim.ActionMove,
						ID:   m.ID,
						Mega: true,
					},
				)
			}
			c.view.AddItem("..", "", '.', func() {
				clearList(c.view)
				c.createActionSpec(slot)
			})
		})
	}

	if slot.Options.CanZMove {
		c.view.AddItem("z-move", "", 'z', func() {
			clearList(c.view)
			for i, move := range slot.Options.DMoves {
				m := move
				c.btnMove(
					i,
					m,
					slot,
					&sim.ActionSpec{
						Type:  sim.ActionMove,
						ID:    m.ID,
						ZMove: true,
					},
				)
			}
			c.view.AddItem("..", "", '.', func() {
				clearList(c.view)
				c.createActionSpec(slot)
			})
		})
	}

	if slot.Options.CanDynamax {
		c.view.AddItem("dynamax", "", 'd', func() {
			clearList(c.view)
			for i, move := range slot.Options.DMoves {
				m := move
				c.btnMove(
					i,
					m,
					slot,
					&sim.ActionSpec{
						Type: sim.ActionMove,
						ID:   m.ID,
						Max:  true,
					},
				)
			}
			c.view.AddItem("..", "", '.', func() {
				clearList(c.view)
				c.createActionSpec(slot)
			})
		})
	}

	c.view.AddItem("fight", "", 'f', func() {
		clearList(c.view)
		for i, move := range slot.Options.Moves {
			m := move
			if (m.Disabled || m.PP == 0) && !slot.Trapped {
				continue
			}
			c.btnMove(
				i,
				m,
				slot,
				&sim.ActionSpec{Type: sim.ActionMove, ID: m.ID},
			)
		}
		c.view.AddItem("..", "", '.', func() {
			clearList(c.view)
			c.createActionSpec(slot)
		})
	})
}
