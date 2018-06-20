#!/usr/bin/env bash

# creating working directory
cd ~
wget https://s3.amazonaws.com/darknode/darknode.zip
unzip darknode.zip
cd ~/.darknode

# Download terraform
if [[ "$OSTYPE" == "linux-gnu" ]]; then
        TERRAFORM_URL="https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip?_ga=2.244744645.1962621656.1520813979-1356604133.1518061974"
elif [[ "$OSTYPE" == "darwin"* ]]; then
        TERRAFORM_URL="https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_darwin_amd64.zip?_ga=2.125313226.130847896.1520825932-912852027.1520206290"
fi

wget $TERRAFORM_URL

# mv terraform* terraform.zip
mv terraform* terraform.zip

# unzip darknode
unzip terraform

# chmod +x darknode-setup
chmod +x terraform

# rm darknode.zip
rm terraform.zip
./gen-config
./darknode-setup