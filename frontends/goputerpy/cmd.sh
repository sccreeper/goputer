#!/bin/bash

rm -rf ./build

mkdir ./build
mkdir ./build/goputerpy

ext=".so"
if [ "$GOOS" = "windows" ]; then
    ext=".dll"
fi

go build -buildmode=c-shared -o "./build/bindings${ext}" ./goputerpy/bindings.go
rm ./build/bindings.h

pyinstaller --onefile --name "goputerpy$( [ $GOOS = windows ] && echo .exe )" \
    --paths ../../.venv/lib64/python3.13/site-packages \
    --paths ../../.venv/lib/python3.13/site-packages \
    --distpath ./build/ \
    --workpath ./build/temp \
    --noconfirm main.py

rm -rf ./build/temp
