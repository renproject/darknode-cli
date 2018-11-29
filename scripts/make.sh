#!/bin/sh

cd cmd
xgo -go 1.10 --targets=darwin/amd64 .
xgo -go 1.10 --targets=linux/amd64 .
mv cmd-linux-amd64 ../target/bin/darknode_linux_amd64
mv cmd-darwin-10.6-amd64 ../target/bin/darknode_darwin_amd64
