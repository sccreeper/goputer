#!/usr/bin/bash

go build  -ldflags="-X main.Commit=$(git rev-parse HEAD)" -o ./out ./cmd/govmcmd/main.go