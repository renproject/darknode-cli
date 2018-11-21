#!/bin/sh

ostype="$(uname -s)"
cputype="$(uname -m)"

while true
do
  R=$(($RANDOM%72))
  if test $R -eq 0; then
    echo "Updating darknode..."
    timestamp=$(date +%Y-%m-%d-%H-%M-%S)
    if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
      curl -s 'https://releases.republicprotocol.com/darknode/latest/darknode-linux.zip' > darknode.zip
    elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
      curl -s 'https://releases.republicprotocol.com/darknode/latest/darknode-darwin.zip' > darknode.zip
    else
      echo 'unsupported OS type or architecture'
      exit 1
    fi
    unzip -o darknode.zip
    mv ./darknode ./.darknode/bin/darknode
    rm darknode.zip
    systemctl --user restart darknode.service
    echo $timestamp >> /home/darknode/.darknode/update.log
    echo "Finish updating"
  fi
  sleep 1h
done
