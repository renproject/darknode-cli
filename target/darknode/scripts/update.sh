#!/bin/sh

# Print commands before executing
set -x

# Install the darknode
ostype="$(uname -s)"
cputype="$(uname -m)"

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

# Start services
systemctl --user restart darknode-updater.service
systemctl --user restart darknode.service
