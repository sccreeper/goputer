#!/bin/bash

rm -rf ./build

mkdir ./build
mkdir ./build/goputerpy

go build -buildmode=plugin -ldflags "-s -w" -o ./build/goputerpy.so ./cli_interface.go
go build -buildmode=c-shared -ldflags "-s -w" -o ./build/goputerpy/bindings.so ./goputerpy/bindings.go
rm ./build/goputerpy/bindings.h

find ./goputerpy -name "*.py" -exec cp -prv "{}" "./build/goputerpy" ";"
cp ./main.py ./build/main.py
cp -r ./rendering ./build/
