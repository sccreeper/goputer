#!/bin/bash

rm -rf ./build

mkdir ./build

go build -buildmode=plugin -o ./build/gp32.so ./main.go