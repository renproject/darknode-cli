#!/bin/sh

while true
do
  R=$(($RANDOM%72))
  if test $R -eq 0; then
    echo "Updating system..."
    sudo DEBIAN_FRONTEND=noninteractive apt-get -y update
    sudo DEBIAN_FRONTEND=noninteractive apt-get -y upgrade
    sudo DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade
    sudo DEBIAN_FRONTEND=noninteractive apt-get -y auto-remove
    echo "Updating darknode..."
    timestamp=$(date +%Y-%m-%d-%H-%M-%S)
    export GOBIN=/home/ubuntu/go/bin
    export GOPATH=/home/ubuntu/go
    mkdir -p /home/ubuntu/go/src/github.com/republicprotocol
    cd /home/ubuntu/go/src/github.com/republicprotocol/republic-go
    sudo git reset --hard HEAD
    sudo git clean -f -d
    sudo git pull
    cd cmd/darknode
    go install
    sudo systemctl restart darknode.service
    echo $timestamp >> /home/ubuntu/.darknode/update.log
    echo "Finish updating"
  fi
  sleep 1h
done
