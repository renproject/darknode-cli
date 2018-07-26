#!/bin/bash

maxdelay=$((3*60*60))  # 3 hours
mindelay=$((1*60*60))  # 1 hour


# mkdir $HOME/.darknode/ui

while true
do
  randomdelay=$(($RANDOM%maxdelay)) # $RANDOM is a value between 0 and 32767 (9 hrs)
  delay=$((mindelay + randomdelay))
  sleep $((delay)) &&
    echo "Checking for darknode updates..."
    timestamp=$(date +%Y-%m-%d-%H-%M-%S) &&
    # Install darknode
    export GOBIN=$HOME/go/bin &&
    export GOPATH=$HOME/go &&
    cd $HOME/go/src/github.com/republicprotocol/republic-go &&
    sudo git reset --hard HEAD &&
    sudo git clean -f -d &&
    sudo git pull &&
    cd cmd/darknode &&
    go install &&
    sudo systemctl restart darknode.service &&
    echo $timestamp >> $HOME/.darknode/update.log &&
    echo "Finish updating"
done




