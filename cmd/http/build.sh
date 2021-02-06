#! /bin/sh

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd $DIR

go build -o app main.go 
docker build -t poke-showdown-go .

cd -
