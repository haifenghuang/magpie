#!/usr/bin/env bash
export GOPATH=$(pwd)

# for newer go version's build error: "package XXX is not in GOROOT"
export GO111MODULE=off

olddir=$(pwd)
# build-wasm
cd wasm/cmd/wasm
echo "Building wasm..."
GOOS=js GOARCH=wasm go build -o ../../assets/magpie.wasm

# run server
cd ../server
echo "Running server..."
echo "    Now open the browser, and type 'http://localhost:9090'"
go run main.go


