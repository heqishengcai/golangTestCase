#!/usr/bin/env bash

set -e

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
OLDGOBIN="$GOBIN"
export GOPATH="$CURDIR"
export GOBIN="$CURDIR/bin/"
echo 'GOPATH:' $GOPATH
echo 'GOBIN:' $GOBIN
#go build -o mygo -race -work -v -ldflags "-s" src/api/local/testCase.go
go build -o mygo -race -work -v -ldflags "-s" src/api/canRun/testCase.go

if [ ! -d ./bin ]; then
	mkdir bin
fi

if [ -e ./mygo ]; then
   mv mygo ./bin/
fi

export GOPATH="$OLDGOPATH"
export GOBIN="$OLDGOBIN"

echo 'build finished'