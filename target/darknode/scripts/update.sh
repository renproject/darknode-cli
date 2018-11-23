#!/bin/sh

curl -s 'https://releases.republicprotocol.com/darknode/latest/darknode.tar.gz' > darknode.tar.gz
tar -xzvf darknode.tar.gz
mv ./darknode ./.darknode/bin/darknode
rm darknode.tar.gz
systemctl --user restart darknode.service
