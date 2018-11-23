#!/bin/sh

while true
do
  R=$(($RANDOM%72))
  if test $R -eq 0; then
    echo "Updating darknode..."
    timestamp=$(date +%Y-%m-%d-%H-%M-%S)
    curl -s 'https://releases.republicprotocol.com/darknode/latest/darknode.tar.gz' > darknode.tar.gz
    tar -xzvf darknode.tar.gz
    mv ./darknode ./.darknode/bin/darknode
    rm darknode.tar.gz
    systemctl --user restart darknode.service
    echo $timestamp >> /home/darknode/.darknode/update.log
    echo "Finish updating"
  fi
  sleep 1h
done
