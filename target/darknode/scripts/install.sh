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

get_latest_release linux
unzip darknode.zip

mkdir -p ./.darknode
mv ./darknode-config.json ./.darknode/config.json

cd darknode
chmod +x install.sh
./install.sh
cd

rm -rf darknode
rm darknode.zip



