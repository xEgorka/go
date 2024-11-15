#!/usr/bin/env bash

main () {

    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo -ne "docker start " && docker start postgres && echo
        trap "echo -ne 'docker stop ' && docker stop postgres" SIGHUP SIGINT
    fi

    PROJECT=$( cd "$(dirname "$0")/../.." ; pwd -P )
    cd $PROJECT
    go mod tidy
    if ! command -v mockgen &> /dev/null
    then
        go install github.com/golang/mock/mockgen@v1.6.0
    fi
    mockgen -destination=internal/app/mocks/mock_storage.go -package=mocks github.com/xEgorka/project3/internal/app/storage Storage
    mockgen -destination=internal/app/mocks/mock_auth.go -package=mocks github.com/xEgorka/project3/internal/app/auth Auth
    
    cd "$(dirname "$0")"
    go run -ldflags "-X github.com/xEgorka/project3/internal/app/server.buildVersion=v0.1 -X 'github.com/xEgorka/project3/internal/app/server.buildDate=$(date +'%Y-%m-%d')' -X github.com/xEgorka/project3/internal/app/server.buildCommit=$(git rev-parse --short HEAD)" main.go -d $CONNINFO
    exit 0
}

main
