#!/bin/bash

rm -rf ./build
mkdir ./build

echo Building WASM...

GOOS=js GOARCH=wasm go build -ldflags "-s -w" -o ./static/main.wasm main.go

echo Copying files...

cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./static/
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
git rev-parse HEAD > ./static/ver
date >> ./static/ver

cp ./run.sh ./build/
cp ./pages_deploy.sh ./build/

echo Building JS...

env $(cat .env) npx parcel build index.html --dist-dir ./build/dist/
