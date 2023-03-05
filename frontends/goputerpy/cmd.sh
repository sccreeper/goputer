#!/bin/bash

rm -rf ./build

mkdir ./build
mkdir ./build/goputerpy

go build -buildmode=plugin -o ./build/goputerpy.so ./cli_interface.go
go build -buildmode=c-shared -o ./build/bindings.so ./goputerpy/bindings.go
rm ./build/bindings.h


pyinstaller -F -n goputerpy --paths ../../.venv/lib/python3.10/site-packages --distpath ./build/ --workpath ./build/temp -y main.py
rm -rf ./build/temp

# find ./goputerpy -name "*.py" -exec cp -prv "{}" "./build/goputerpy" ";"
# cp ./main.py ./build/main.py
# cp -r ./rendering ./build/
