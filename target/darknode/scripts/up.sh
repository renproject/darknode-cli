#!/bin/sh

# Print commands before executing
set -x

# Do until not locked - will enter infinite loop if update fails
until sudo apt update; do sleep 2; done

# Update the system-level updates
sudo DEBIAN_FRONTEND=noninteractive apt-get -y update
sudo DEBIAN_FRONTEND=noninteractive apt-get -y upgrade
sudo DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade
sudo DEBIAN_FRONTEND=noninteractive apt-get -y auto-remove

# Install services
sudo mv ./provisions/darknode-updater.service /etc/systemd/system/darknode-updater.service
sudo mv ./provisions/darknode.service /etc/systemd/system/darknode.service
sudo mv ./scripts/updater.sh

# Install golang
wget https://dl.google.com/go/go1.10.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz
rm go1.10.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> $HOME/.profile
sudo ln -s /usr/local/go/bin/go /usr/bin/go

# Setup UFW
sudo apt-get install ufw
sudo ufw allow 22/tcp     # ssh
sudo ufw allow 18514/tcp  # republicprotocol
sudo ufw allow 18515/tcp  # status page
sudo ufw --force enable

# Configure darknode and the updater
mkdir ./.darknode/
mv ./darknode-config.json ./.darknode/config.json
mv ./scripts/updater.sh ./.darknode/updater.sh

# Install dep
mkdir -p $HOME/go/bin
export GOBIN=$HOME/go/bin
export GOPATH=$HOME/go
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Install darknode
until sudo apt install -y gcc; do sleep 2; done
mkdir -p $HOME/go/src/github.com/republicprotocol
cd $HOME/go/src/github.com/republicprotocol
git clone https://github.com/republicprotocol/republic-go.git
cd republic-go/cmd/darknode
$GOBIN/dep ensure
go install
cd $HOME

# Will fail if there are any files still in ./provisions/
rm -rf ./provisions/
rm -rf ./scripts/

# Start services
sudo systemctl daemon-reload
sudo systemctl enable darknode-updater.service
sudo systemctl enable darknode.service 
sudo systemctl start darknode-updater.service
sudo systemctl start darknode.service
