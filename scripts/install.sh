#!/usr/bin/sh

# creating working directory
mkdir -p $HOME/.darknode
cd $HOME/.darknode
curl -s 'https://darknode.republicprotocol.com/darknode.zip' > darknode.zip
unzip darknode.zip

ostype="$(uname -s)"
cputype="$(uname -m)"

# Download terraform
if [ $ostype = 'Linux' -a $cputype = 'x86_64' ]; then
    TERRAFORM_URL='https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip'
    curl -s 'https://darknode.republicprotocol.com/darknode_linux_amd64' > ./bin/darknode
elif [ $ostype = 'Darwin' -a $cputype = 'x86_64' ]; then
    TERRAFORM_URL='https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_darwin_amd64.zip'
    curl -s 'https://darknode.republicprotocol.com/darknode_darwin_amd64' > ./bin/darknode
else
   echo 'unsupported OS type'
   cd ..
   rm -rf .darknode
   exit 1
fi

chmod +x bin/darknode

curl -s "$TERRAFORM_URL" > terraform.zip

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
if ! [ -x "$(command -v darknode)" ]; then
  if test -n $BASH_VERSION  &&  [ -f "$HOME/.bash_profile" ] ; then
    echo 'export PATH=$PATH:$HOME/.darknode/bin' >> ~/.bash_profile
    source ~/.bash_profile
  elif test -n $ZSH_VERSION && [ -f "$HOME/.zprofile" ] ; then
    echo 'export PATH=$PATH:$HOME/.darknode/bin' >> ~/.zprofile
    source ~/.zprofile
  elif [ -f '~/.profile' ]; then
    echo 'export PATH=$PATH:$HOME/.darknode/bin' >> ~/.profile
    source ~/.profile
  fi

  echo ''
  echo 'If you are using a custom shell, make sure you update your PATH.'
  echo '$ export PATH=$PATH:$HOME/.darknode/bin'
fi

echo ''
echo 'Done! Restart terminal and run the command below to begin.'
echo 'darknode up --help'