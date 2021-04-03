package parse

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"
)

// activeCmd is a holder for a commands stdout, err, in etc
type activeCmd struct {
	stdout <-chan string
	stderr <-chan string
	errors <-chan error
	stdin  chan string
	ctrl   chan os.Signal
}

// stop kills the cmd
func (a *activeCmd) stop() {
	defer close(a.stdin)
	defer close(a.ctrl)
	a.ctrl <- syscall.SIGINT
}

func (a *activeCmd) start(seed int, format string, teamsizes []int, teams []string) error {
	// pushes in initial stdin to kick off battle
	// #1 push in the battle format
	data, err := json.Marshal(map[string]interface{}{
		"seed":     []int{seed, seed, seed, seed},
		"formatid": format,
	})
	if err != nil {
		return err
	}
	a.stdin <- fmt.Sprintf(">start %s\n", string(data))

	// #2 for each player we need to announce them & their team in packed format
	orders := []string{}
	for num, pteam := range teams {
		player := fmt.Sprintf("p%d", num+1)

		data, err = json.Marshal(map[string]interface{}{
			"name": player,
			"team": pteam,
		})
		if err != nil {
			return err
		}

		// specify player's team
		a.stdin <- fmt.Sprintf(">player %s %s\n", player, string(data))

		if pteam == "" {
			continue
		}

		// specify the players team order (we give the order in which we got them).
		members := []string{}
		for i := 0; i < teamsizes[num]; i++ {
			members = append(members, fmt.Sprintf("%d", i+1))
		}
		orders = append(orders, fmt.Sprintf(">%s team %s\n", player, strings.Join(members, ",")))
	}

	// #3 we now need to tell the simulator what order the player's team should
	// be in (ie, battle order).
	for _, order := range orders {
		a.stdin <- order
	}

	return nil
}
