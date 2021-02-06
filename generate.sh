#!/bin/bash

#
# Generates json marshal/unmarshal functions for structs
#
# tools:
#   (easyjson) go get -u github.com/mailru/easyjson/...
#

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

easyjson -all $DIR/pkg/structs/action.go
easyjson -all $DIR/pkg/structs/battle_spec.go
easyjson -all $DIR/pkg/structs/battle_state.go
easyjson -all $DIR/pkg/structs/event.go
easyjson -all $DIR/pkg/structs/pokemon_spec.go
easyjson -all $DIR/pkg/structs/update.go
easyjson -all $DIR/pkg/structs/simulate.go

