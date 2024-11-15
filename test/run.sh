#!/usr/bin/env bash
set -eu

PROJECT=$( cd "$(dirname "$0")/.." ; pwd -P )

cd $PROJECT
go test ./... -coverprofile ./test/coverage.out
go tool cover -func ./test/coverage.out
go tool cover -html=./test/coverage.out -o ./test/coverage.html
