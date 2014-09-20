#!/bin/sh

# Install dirs auto
mkdir -p src
mkdir -p pkg
mkdir -p bin
mkdir -p _tmp

export GOPATH=$GOPATH:"$PWD":"$PWD"/vendor
export PATH=$PATH:"$PWD"/bin
export POSTMAN_DB_DIR="$PWD"/_tmp