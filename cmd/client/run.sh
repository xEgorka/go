#!/usr/bin/env bash

main () {

    PROJECT=$( cd "$(dirname "$0")" ; pwd -P )
    cd $PROJECT
    go mod tidy
    if ! command -v mockgen &> /dev/null
    then
        go install github.com/golang/mock/mockgen@v1.6.0
    fi
    mockgen -destination=internal/app/mocks/mock_auth.go -package=mocks github.com/xEgorka/project3/cmd/client/internal/app/auth Auth
    mockgen -destination=internal/app/mocks/mock_storage.go -package=mocks github.com/xEgorka/project3/cmd/client/internal/app/storage Storage
    cd "$(dirname "$0")"
    go run main.go
    exit 0
}

main
