#!/bin/sh

# Print commands before executing
set -x

# Install services
mkdir -p $HOME/.config/systemd/user
mv ./services/darknode-updater.service $HOME/.config/systemd/user/darknode-updater.service
mv ./services/darknode.service $HOME/.config/systemd/user/darknode.service

curl -s 'https://releases.republicprotocol.com/darknode/latest/darknode.tar.gz' > darknode.tar.gz
tar -xzvf darknode.tar.gz
mkdir -p ./.darknode/bin
mv ./darknode ./.darknode/bin/darknode
mv ./darknode-config.json ./.darknode/config.json
mv ./scripts/updater.sh ./.darknode/bin/update.sh

rm -rf ./scripts/
rm -rf ./services/
rm darknode.tar.gz

# Start services
loginctl enable-linger darknode
systemctl --user enable darknode.service
systemctl --user enable darknode-updater.service
systemctl --user start darknode.service
systemctl --user start darknode-updater.service

