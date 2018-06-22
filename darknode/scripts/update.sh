#!/usr/bin/env bash

cd ./go/src/github.com/republicprotocol/republic-go
sudo git stash
sudo git checkout nightly
sudo git fetch origin nightly
sudo git reset --hard origin/nightly
cd cmd/darknode
go install
cd
sudo service darknode restart
