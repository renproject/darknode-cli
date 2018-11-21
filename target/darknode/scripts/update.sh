#!/bin/sh

while true
do
  R=$(($RANDOM%72))
  if test $R -eq 0; then
    echo "Updating darknode..."
    timestamp=$(date +%Y-%m-%d-%H-%M-%S)
    export GOBIN=/home/darknode/go/bin
    export GOPATH=/home/darknode/go
    mkdir -p /home/ubuntu/go/src/github.com/republicprotocol
    cd /home/ubuntu/go/src/github.com/republicprotocol/republic-go
    git reset --hard HEAD
    git clean -f -d
    git pull
    cd cmd/darknode
    go install
    systemctl --user restart darknode.service
    echo $timestamp >> /home/darknode/.darknode/update.log
    echo "Finish updating"
  fi
  sleep 1h
done
