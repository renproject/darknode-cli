#!/usr/bin/sh

# creating working directory
mkdir -p $HOME/.darknode
cd $HOME/.darknode
wget https://darknode.republicprotocol.com/darknode.zip
unzip darknode.zip

ostype='$(uname -s)'
cputype='$(uname -m)'

# Download terraform
if [[ $ostype = 'Linux' -a $cputype = 'x86_64' ]]; then
    TERRAFORM_URL='https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip'
    wget https://darknode.republicprotocol.com/darknode_linux_amd64
    mv darknode_linux_amd64 ./bin/darknode
elif [[ $ostype = 'Darwin' -a $cputype = 'x86_64' ]]; then
    TERRAFORM_URL='https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_darwin_amd64.zip'
    wget https://darknode.republicprotocol.com/darknode_darwin_amd64
    mv darknode_darwin_amd64 ./bin/darknode
else
   echo 'unsupported OS type'
   cd ..
   rm -rf .darknode
   exit 1
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
if ! [ -x '$(command -v darknode)' ]; then
  if test -n $ZSH_VERSION; then
    echo 'export PATH=$PATH:$HOME/.darknode/bin' >> ~/.zshrc
    source ~/.zshrc
  elif test -n $BASH_VERSION; then
    echo 'export PATH=$PATH:$HOME/.darknode/bin' >> ~/.bash_profile
    source ~/.bash_profile
  elif test -n $KSH_VERSION; then
    echo 'export PATH=$PATH:$HOME/.darknode/bin' >> ~/.kshrc
    source ~/.kshrc
  elif test -n $FCEDIT; then
    echo 'export PATH=$PATH:$HOME/.darknode/bin' >> ~/.kshrc
    source ~/.kshrc
  else
    echo 'Seems your are using a custom sh.'
    echo 'Please add /.darknode/bin to the PATH variable'
  fi
fi

echo ''
echo 'Done! Run the command below to begin.'
echo 'darknode up --help'