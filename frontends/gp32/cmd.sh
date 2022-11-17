#!/bin/bash

rm -rf ./build

mkdir ./build

go build -buildmode=plugin -ldflags "-s -w" -o ./build/gp32.so ./main.go