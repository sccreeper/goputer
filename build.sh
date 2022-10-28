#!/usr/bin/bash

go build -ldflags="-X main.Commit=$(git rev-parse HEAD)" -o ./goputer ./cmd/goputer/main.go