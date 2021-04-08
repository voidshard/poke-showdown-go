#!/bin/sh

#
# Generates an assets.go file in pkg/pokedata from pokemon-showdown 
# - pokedex.json
# - moves.json
#
# Data is embedded in to the library with go-bindata
# - see: https://github.com/go-bindata/go-bindata
#

set -eu

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

# fetch files
mkdir ${DIR}/assets
wget "https://play.pokemonshowdown.com/data/pokedex.json" -O ${DIR}/assets/pokedex.json
wget "https://play.pokemonshowdown.com/data/moves.json" -O ${DIR}/assets/moves.json

# generate assets.go
go-bindata -prefix "${DIR}/assets" -o ${DIR}/pkg/pokedata/assets.go ${DIR}/assets/
sed -i 's/package main/package pokedata/g' ${DIR}/pkg/pokedata/assets.go

# clean up
rm -v ${DIR}/assets/pokedex.json
rm -v ${DIR}/assets/moves.json
rmdir ${DIR}/assets
