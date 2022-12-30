#!/bin/bash

rm -rf ./build
mkdir ./build

GOOS=js GOARCH=wasm go build -ldflags "-s -w" -o ./static/main.wasm main.go
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./static/
git rev-parse HEAD > ./static/ver
date >> ./static/ver

cp ./run.sh ./build/
cp ./pages_deploy.sh ./build/
env $(cat .env) npx parcel build index.html --dist-dir ./build/dist/
