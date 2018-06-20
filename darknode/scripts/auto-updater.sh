#!/bin/bash

maxdelay=$((3*60*60))  # 3 hours
mindelay=$((1*60*60))  # 1 hour

# mkdir /home/ubuntu/.darknode/ui

while true
do
  randomdelay=$(($RANDOM%maxdelay)) # $RANDOM is a value between 0 and 32767 (9 hrs)
  delay=$((mindelay + randomdelay))
  sleep $((delay)) &&
    echo "Checking for darknode updates..."
    timestamp=$(date +%Y-%m-%d-%H-%M-%S) &&
    # Install darknode
    export GOBIN=/home/ubuntu/go/bin &&
    mkdir -p /home/ubuntu/go/src/github.com/republicprotocol &&
    cd /home/ubuntu/go/src/github.com/republicprotocol &&
    cd republic-go &&
    git pull &&
    cd cmd/darknode &&
    go install &&
    # mv $GOBIN/darknode $GOBIN/dark-node &&
    cd /home/ubuntu &&
    # go get -u github.com/republicprotocol/republic-go/cmd/dark-node &&
    # cp -r /home/ubuntu/go/src/github.com/republicprotocol/republic-go/dark-node-ui/* /home/ubuntu/.darknode/ui &&
    sudo systemctl restart darknode.service &&
    echo $timestamp > $HOME/.darknode/update.log
done