#!/bin/sh

# Print commands before executing
set -x

# Install services
mkdir -p $HOME/.config/systemd/user
mv ./services/darknode-updater.service $HOME/.config/systemd/user/darknode-updater.service
mv ./services/darknode.service $HOME/.config/systemd/user/darknode.service

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
mkdir -p ./.darknode/bin
mv ./darknode ./.darknode/bin/darknode
mv ./darknode-config.json ./.darknode/config.json
mv ./scripts/updater.sh ./.darknode/bin/update.sh

rm -rf ./scripts/
rm -rf ./services/
rm darknode.zip

# Start services
systemctl --user enable darknode-updater.service
systemctl --user enable darknode.service
systemctl --user start darknode-updater.service
systemctl --user start darknode.service
