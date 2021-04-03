package sim

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPack(t *testing.T) {
	cases := []struct {
		Given  *Action
		Expect string
	}{
		{
			&Action{
				Player: "foo",
				Specs: []*ActionSpec{
					&ActionSpec{
						Type:   ActionMove,
						ID:     "1",
						Target: 2,
						Mega:   true,
					},
					&ActionSpec{
						Type:   ActionMove,
						ID:     "4",
						Target: -1,
					},
				},
			},
			">foo move 1 2 mega,move 4 -1\n",
		},
		{
			&Action{
				Player: "foo",
				Specs: []*ActionSpec{
					&ActionSpec{
						Type: ActionSwitch,
						ID:   "3",
					},
				},
			},
			">foo switch 3\n",
		},
		{
			&Action{
				Player: "foo",
				Specs: []*ActionSpec{
					&ActionSpec{
						Type:   ActionMove,
						ID:     "4",
						Target: -1,
					},
				},
			},
			">foo move 4 -1\n",
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := tt.Given.Pack()

			assert.Equal(t, tt.Expect, result)
		})
	}
}

func TestPackMove(t *testing.T) {
	cases := []struct {
		Given  *ActionSpec
		Expect string
	}{
		{&ActionSpec{ID: "2", Target: 3, Mega: true}, "2 3 mega"},
		{&ActionSpec{ID: "3", Target: -2, Max: true}, "3 -2 max"},
		{&ActionSpec{ID: "2", Target: -1}, "2 -1"},
		{&ActionSpec{ID: "5", Target: -1}, "5 -1"},
		{&ActionSpec{ID: "1"}, "1"},
		{&ActionSpec{ID: "1", Mega: true}, "1 mega"},
		{&ActionSpec{ID: "1", Max: true}, "1 max"},
		{&ActionSpec{ID: "1", ZMove: true}, "1 zmove"},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			out := packMove(tt.Given)

			assert.Equal(t, tt.Expect, out)
		})
	}
}
