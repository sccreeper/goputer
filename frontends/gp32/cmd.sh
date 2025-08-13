#!/bin/bash

rm -rf ./build

mkdir ./build

go build -ldflags "-s -w" -o ./build/gp32 ./main.go