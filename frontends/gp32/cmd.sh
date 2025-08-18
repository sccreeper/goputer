#!/bin/bash

rm -rf ./build

mkdir ./build

go build -ldflags "-s -w" -o "./build/gp32$( [ $GOOS = windows ] && echo .exe )" ./main.go