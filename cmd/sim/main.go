package main

import (
	"os"
)

const (
	// env var to pokemon-showdown nodejs binary
	envBin = "POKE_SIM_SHOWDOWN_BIN"
)

func main() {
	showdown := os.Getenv(envBin)
	if showdown == "" {
		showdown = "./node_modules/pokemon-showdown/pokemon-showdown"
	}

}
