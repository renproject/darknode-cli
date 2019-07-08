#!/bin/sh

if [ -d "$HOME/.darknode" ] && [ -d "$HOME/.darknode/darknodes" ]; then
    cd $HOME/.darknode
    curl -s 'https://releases.republicprotocol.com/darknode-cli/resources.zip' > resources.zip
    unzip -o resources.zip
else
    echo "cannot find the darknode-cli"
    echo "please install darknode-cli first"
    exit 1
fi

# Check the os type and cpu architecture
ostype="$(uname -s)"
cputype="$(uname -m)"

# Download the latest binary darknode
if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
    TERRAFORM_URL='https://releases.hashicorp.com/terraform/0.12.3/terraform_0.12.3_linux_amd64.zip'
    curl -s 'https://releases.republicprotocol.com/darknode-cli/darknode_linux_amd64' > ./bin/darknode
elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
    TERRAFORM_URL='https://releases.hashicorp.com/terraform/0.12.3/terraform_0.12.3_darwin_amd64.zip'
    curl -s 'https://releases.republicprotocol.com/darknode-cli/darknode_darwin_amd64' > ./bin/darknode
else
   echo 'unsupported OS type or architecture'
   exit 1
fi

# download terraform
curl -s "$TERRAFORM_URL" > terraform.zip
unzip terraform
chmod +x terraform
mv terraform bin/terraform
rm terraform.zip

# update darknode binary
chmod +x bin/darknode
rm resources.zip

echo ''
echo 'Done! Your darknode-cli has been updated.'