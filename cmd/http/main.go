package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/voidshard/poke-showdown-go/pkg/sim"
	"github.com/voidshard/poke-showdown-go/pkg/structs"
)

const (
	// various constants we return to indicate more precisely what went wrong
	errNoData  = "no-data"
	errBadData = "bad-data"
	errSimTurn = "failed-turn"
	errUnknown = "unknown-error"
	errTimeout = "err-timeout"
)

var (
	// cli flags for the http server
	flagPort = flag.Int("p", 8080, "Specify HTTP port to listen on")
	flagBin  = flag.String("b", "pokemon-showdown", "Specify pokemon showdown binary path")
)

// errMsg packages an error message into a json format, failure
// to report our failure is not an option
func errMsg(name, msg string, err error) []byte {
	data, err := json.Marshal(map[string]string{
		"type": name,
		"msg":  fmt.Sprintf("%s: %v", msg, err),
	})
	if err != nil {
		return []byte(fmt.Sprintf("{\"type\": \"%s\"}", name))
	}
	return data
}

// simulateHandle parses JSON data & passes it to sim.Simulate
func simulateHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMsg(errNoData, "no HTTP body", nil))
		return
	}

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read HTTP body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMsg(errBadData, "failed to read HTTP body", err))
		return
	}

	in := structs.SimData{}
	err = json.Unmarshal(data, &in)
	if err != nil {
		log.Println("failed to parse JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMsg(errBadData, "failed to parse JSON", err))
		return
	}

	// basic sanity checks
	if in.Spec == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMsg(errBadData, "spec required", err))
		return
	}
	if in.Spec.Players == nil || len(in.Spec.Players) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMsg(errBadData, "two players required", err))
		return
	}
	for _, p := range in.Spec.Players {
		if p == nil || len(p) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errMsg(errBadData, "players must have at least one pokemon", err))
			return
		}
	}
	if in.Spec.Format == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMsg(errBadData, "spec battle format required", err))
		return
	}
	if in.Actions == nil || len(in.Actions) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMsg(errBadData, "actions required", err))
		return
	}
	if len(in.Actions) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMsg(errBadData, "at least two actions required", err))
		return
	}
	// end sanity checks

	result, err := sim.Simulate(in.Spec, in.Actions)
	if err != nil {
		log.Println("simulation failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		if errors.Is(err, sim.ErrTimeout) {
			w.Write(errMsg(errTimeout, "timedout running simulation", err))
			return
		}

		if errors.Is(err, sim.ErrFailedTurn) {
			w.Write(errMsg(errSimTurn, "simulator error", err))
			return
		}

		w.Write(errMsg(errUnknown, "unknown error during simulation", err))
		return
	}

	resultdata, err := json.Marshal(result)
	if err != nil {
		log.Println("failed to marshal JSON reply: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg(errBadData, "failed to marshal JSON", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resultdata)
}

// healthHandle returns that we're up & ready to serve traffic
func healthHandle(w http.ResponseWriter, r *http.Request) {
	// TODO: ensure we don't timeout if the pokemon-showdown binary isn't found
	_, err := sim.Simulate(
		&structs.BattleSpec{
			Format: structs.FormatGen8,
			Players: [][]*structs.PokemonSpec{
				[]*structs.PokemonSpec{&structs.PokemonSpec{Name: "pikachu", Moves: []string{"tackle"}}},
				[]*structs.PokemonSpec{&structs.PokemonSpec{Name: "pikachu", Moves: []string{"thundershock"}}},
			},
		},
		[]*structs.Action{
			&structs.Action{Player: "p1", Specs: []*structs.ActionSpec{&structs.ActionSpec{ID: 0}}},
			&structs.Action{Player: "p2", Specs: []*structs.ActionSpec{&structs.ActionSpec{ID: 0}}},
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	flag.Parse()
	http.HandleFunc("/v1/simulate", simulateHandle)
	http.HandleFunc("/_health", healthHandle)
	fmt.Printf("listening on :%d\n", *flagPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *flagPort), nil))
}
