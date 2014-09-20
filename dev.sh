#!/bin/sh

# Install dirs auto
mkdir -p src
mkdir -p pkg
mkdir -p bin

export GOPATH=$GOPATH:"$PWD":"$PWD"/vendor
export PATH=$PATH:"$PWD"/bin