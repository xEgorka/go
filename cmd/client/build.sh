#!/bin/bash

PLATFORMS=("darwin-arm64" "linux-amd64" "windows-amd64")

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%%-*}
    GOARCH=${PLATFORM##*-}
    echo "Building for $GOOS/$GOARCH..."

    export GOOS
    export GOARCH

    BIN_EXT=""
    if [[ "$GOOS" == "windows" ]]; then
        BIN_EXT=".exe"
    fi

    go build -ldflags="-s -w" -o "bin/gophkeeper-$GOOS-$GOARCH$BIN_EXT"
done

echo "Build completed."
