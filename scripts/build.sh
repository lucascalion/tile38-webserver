#!/bin/bash

# Adapted from https://github.com/tidwall/tile38/blob/1.22.3/scripts/build.sh

set -e
cd $(dirname "${BASH_SOURCE[0]}")/..

# Check the Go installation
if [ "$(which go)" == "" ]; then
	echo "error: Go is not installed. Please download and follow installation"\
		 "instructions at https://golang.org/dl to continue."
	exit 1
fi

# Hardcode some values to the core package.
if [ -d ".git" ]; then
	VERSION=$(git describe --tags --abbrev=0)
	GITSHA=$(git rev-parse --short HEAD)
fi

# Set final Go environment options
LDFLAGS="$LDFLAGS -extldflags '-static'"
export CGO_ENABLED=0

if [ "$NOMODULES" != "1" ]; then
	export GO111MODULE=on
	export GOFLAGS=-mod=vendor
	go mod vendor
fi

# Build and store objects into original directory.
go build -ldflags "$LDFLAGS" -o tile38-webserver *.go

