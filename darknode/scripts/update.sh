#!/usr/bin/env bash

cd ./go/src/github.com/republicprotocol/republic-go
sudo git stash
sudo git checkout master
sudo git fetch origin master
sudo git reset --hard origin/master
cd cmd/darknode
go install
cd
sudo service darknode restart
