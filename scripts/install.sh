#!/usr/bin/env bash

# creating working directory
mkdir -p $HOME/.darknode
cd $HOME/.darknode
wget https://darknode.republicprotocol.com/darknode.zip
unzip darknode.zip

# Download terraform
if [[ "$OSTYPE" == "linux-gnu" ]]; then
        TERRAFORM_URL="https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip"
        wget https://darknode.republicprotocol.com/darknode_linux_amd64
        mv darknode_linux_amd64 ./bin/darknode
elif [[ "$OSTYPE" == "darwin"* ]]; then
        TERRAFORM_URL="https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_darwin_amd64.zip"
        wget https://darknode.republicprotocol.com/darknode_darwin_amd64
        mv darknode_darwin_amd64 ./bin/darknode
fi
chmod +x bin/darknode

wget $TERRAFORM_URL

# unzip darknode
mv terraform* terraform.zip
unzip terraform

# chmod +x darknode-setup
chmod +x terraform
mv terraform bin/terraform

# rm darknode.zip
rm darknode.zip
rm terraform.zip

# make sure the binary is installed in the path
echo 'export PATH=$PATH:$HOME/.darknode/bin' >> $HOME/.profile

echo ''
echo 'Done!'