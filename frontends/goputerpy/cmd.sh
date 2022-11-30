#!/bin/bash

rm -rf ./build

mkdir ./build

go build -buildmode=plugin -ldflags "-s -w" -o ./build/goputerpy.so ./cli_interface.go
go build -buildmode=c-shared -ldflags "-s -w" -o ./build/bindings.so ./goputerpy/bindings.go

cp -r ./goputerpy/* ./build/
