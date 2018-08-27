#!/bin/sh

# Create directory for build files
mkdir -p build

# Compress darknode directory
cd target/darknode
zip -r ../../build/darknode.zip ./*
cd ../..

# Copy install script to build folder
cp scripts/install.sh build/install.sh
cp scripts/update.sh build/update.sh

# Generate binaries
docker-machine create default
eval $(docker-machine env default) # Setup the environment for the Docker client
go get github.com/karalabe/xgo
xgo --targets=darwin/amd64,linux/amd64 ./cmd
mv cmd-darwin-10.6-amd64 build/darknode_darwin_amd64
mv cmd-linux-amd64 build/darknode_linux_amd64
docker-machine rm -f default