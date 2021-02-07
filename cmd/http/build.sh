#! /bin/sh

#
# Script builds docker image for a simple stateless HTTP server that wraps our sim logic.
#

set -eu

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd $DIR

# compile the app
echo "compiling .."
go build -o app main.go 

# build the image
echo "building image .."
docker build -t poke-showdown-go .

# remove the app
echo "cleaning up .."
rm -v app

echo "done"
cd -
