#!/usr/bin/env bash
export GOPATH=$(pwd)

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


