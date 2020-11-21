package main

import (
	"os"
)

const (
	// env var to pokemon-showdown nodejs binary
	binEnv = "POKE_SIM_SHOWDOWN_BIN"

	// default binary used if not set
	binDefault = "pokemon-showdown"
)

func main() {
	showdown := os.Getenv(envBin)
	if showdown == "" {
		showdown = binDefault
	}

}
