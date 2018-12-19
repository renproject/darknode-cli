#!/bin/sh

get_latest_release() {
  curl -s https://api.github.com/repos/republicprotocol/republic-go/releases/latest \
    | grep "browser_download_url.*darknode-$1.zip" \
    | cut -d : -f 2,3 \
    | tr -d \" \
    | wget -qi -
  mv darknode-$1.zip darknode.zip
}

# Print commands before executing
set -x

# Install services
mkdir -p $HOME/.config/systemd/user
mkdir -p ./.darknode/bin

mv ./services/darknode-updater.service $HOME/.config/systemd/user/darknode-updater.service
mv ./services/darknode.service $HOME/.config/systemd/user/darknode.service

get_latest_release

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

