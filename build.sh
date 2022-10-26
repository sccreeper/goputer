#!/usr/bin/bash

go build -ldflags="-X main.Commit=$(git rev-parse HEAD)" -o ./goputer ./cmd/govmcmd/main.go