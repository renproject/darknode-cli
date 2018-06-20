#!/usr/bin/env bash

cd ./go/src/github.com/republicprotocol/republic-go
sudo git checkout master
sudo git pull
cd cmd/darknode
go install
cd
sudo service darknode restart
