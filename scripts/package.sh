#!/bin/bash

#Adapted from https://github.com/tidwall/tile38/blob/1.22.3/scripts/package.sh
set -e
cd $(dirname "${BASH_SOURCE[0]}")/..

PLATFORM="Linux"
GOOS="linux"
GOARCH="amd64"
VERSION=$(git describe --tags --abbrev=0)

echo Packaging $PLATFORM Binary

# Remove previous build directory, if needed.
bdir=tile38-webserver-$VERSION-$GOOS-$GOARCH
rm -rf packages/$bdir && mkdir -p packages/$bdir

# Make the binaries.
GOOS=$GOOS GOARCH=$GOARCH ./scripts/build.sh

# Copy the executable binaries.
mv tile38-webserver packages/$bdir

# Copy documention and license.
cp README.md packages/$bdir
cp LICENSE packages/$bdir
cp scripts/start_server packages/$bdir

# Compress the package.
cd packages
tar -zcf $bdir.tar.gz $bdir

# Remove build directory.
rm -rf $bdir
