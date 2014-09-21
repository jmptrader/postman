#!/bin/sh

export GOPATH=$GOPATH:"$PWD":"$PWD"/vendor

cd vendor
# install dependences
go get github.com/ugorji/go/codec
