#!/bin/bash

rm -rf ./build
mkdir ./build

GOOS=js GOARCH=wasm go build -ldflags "-s -w" -o ./static/main.wasm main.go
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./static/
cp ./run.sh ./build/

npx parcel build index.html --dist-dir ./build/dist/
