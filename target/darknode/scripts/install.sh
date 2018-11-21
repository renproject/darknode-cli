#!/bin/sh

# Print commands before executing
set -x

# Install services
mkdir -p $HOME/.config/systemd/user
mv ./provisions/darknode-updater.service $HOME/.config/systemd/user/darknode-updater.service
mv ./provisions/darknode.service $HOME/.config/systemd/user/darknode.service

# Configure darknode and the updater
mkdir ./.darknode/
mv ./darknode-config.json ./.darknode/config.json
mv ./scripts/updater.sh ./.darknode/bin/update.sh

rm -rf ./scripts/
rm -rf ./services/

# Start services
systemctl --user enable darknode-updater.service
systemctl --user enable darknode.service
systemctl --user start darknode-updater.service
systemctl --user start darknode.service
