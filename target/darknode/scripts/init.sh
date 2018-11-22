#!/bin/sh

# create new user and enable ssh login
sudo adduser darknode --gecos \",,,\" --disabled-password
sudo usermod -aG systemd-journal darknode
sudo rsync --archive --chown=darknode:darknode ~/.ssh /home/darknode

until sudo apt update; do sleep 2; done

# Update the system-level updates
sudo DEBIAN_FRONTEND=noninteractive apt-get -y update
sudo DEBIAN_FRONTEND=noninteractive apt-get -y upgrade
sudo DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade
sudo DEBIAN_FRONTEND=noninteractive apt-get -y auto-remove

# Setup UFW
sudo apt-get install ufw
sudo ufw limit 22/tcp     # ssh
sudo ufw allow 18514/tcp  # republicprotocol
sudo ufw limit 18515/tcp  # status page
sudo ufw --force enable
